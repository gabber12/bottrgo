package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"tweety/classifier"
	"tweety/twitter"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func getTwitterAuthConfig() twitter.AuthConfig {
	tckey := os.Getenv("TWITTER_CONSUMER_KEY")
	tcsec := os.Getenv("TWITTER_CONSUMER_SECRET")
	tatoken := os.Getenv("TWITTER_ACCESS_TOKEN")
	tatokensecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	return twitter.AuthConfig{
		ConsumerKey:       tckey,
		ConsumerSecret:    tcsec,
		AccessToken:       tatoken,
		AccessTokenSecret: tatokensecret,
	}
}
func classificationText(string string) {

}
func getFilterConfig(textFilter, locationFilter, userFilter *string) twitter.StreamFilter {
	// likeOnTweet := flag.Bool("likeOnTweet", false, "Like on Tweet")
	// retweetOnHashTagMention := flag.Bool("retweetOnTweet", false, "Like on Tweet")
	// commentOnHashTagMention := flag.Bool("commentOnTweet", false, "Like on Tweet")

	textFilterList := strings.Split((*textFilter), "|")
	if len(*textFilter) == 0 {
		textFilterList = []string{}
	}
	locationFilterList := strings.Split((*locationFilter), "|")
	if len(*locationFilter) == 0 {
		locationFilterList = []string{}
	}
	userFilterList := strings.Split((*userFilter), "|")
	if len(*userFilter) == 0 {
		userFilterList = []string{}
	}
	return twitter.StreamFilter{
		TextKeywords: textFilterList,
		Locations:    locationFilterList,
		UserIds:      userFilterList,
	}
}

type Action string

const (
	RETWEET Action = "RETWEET"
)

type Rule interface {
	evaluate(twitter.Tweet, classifier.Classification) bool
	action() Action
}
type FirstStrategy struct{}

func (FirstStrategy) evaluate(tweet twitter.Tweet, cls classifier.Classification) bool {
	return cls.SentimentScore > 0.8
}
func (FirstStrategy) action() Action {
	return RETWEET
}
func handleAction(tweet *twitter.Tweet, cls *classifier.Classification) {
	// logrus.Infof("Tweet Received %s", tweet)
	rule := FirstStrategy{}
	if rule.evaluate(*tweet, *cls) {
		rule.action()
	}
}
func filterExpression(tweet *twitter.Tweet, classification *classifier.Classification) bool {
	return classification.SentimentScore > 0.8
}
func getSentiment(score float64) string {
	if score > 0.5 {
		return "😀"
	}
	return "😟"
}
func main() {
	logrus.Info("Starting Twitter Sevice ...\n")
	classify := flag.Bool("classify", false, "Hash tag list to filter by")
	machineBoxURL := flag.String("mbHost", "http://localhost:8080", "Machine Box Url")

	// botConfig := getConfigByFlag()
	// commentPool := flag.String("commentList", "", "Hash tag list to filter by")

	service := twitter.Service(getTwitterAuthConfig())
	textFilter := flag.String("textFilter", "", "Hash tag list to filter by")
	locationFilter := flag.String("locationFilter", "", "Location list to filter by")
	userFilter := flag.String("userFilter", "", "UserIds list to filter by")
	flag.Parse()

	service.FilteredStream(getFilterConfig(textFilter, locationFilter, userFilter), func(tweet *twitter.Tweet) {
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		fmt.Printf("\n%s\n@%s %s %s\n%s\n", blue("-----------"), cyan(tweet.User.ScreenName), yellow("::"), tweet.FullText, blue("-----------"))

		classification := &classifier.Classification{}
		if *classify {
			var cls classifier.Service
			cls = &classifier.MachineBox{HostPort: *machineBoxURL}
			classification, err := cls.Classify(tweet.FullText)
			if err == nil {
				sentiment := getSentiment(classification.SentimentScore)
				fmt.Printf("%s\n", yellow("Sentiment ", sentiment))
				if !filterExpression(tweet, classification) {
					return
				}
			}
		}

		handleAction(tweet, classification)

	})
	logrus.Info("Service Stopped\n")
}
