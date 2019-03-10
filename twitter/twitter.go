// Twitter Manages Twitter tweet fetching etc..
package twitter

import (
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/ChimeraCoder/anaconda"
)

//Twitter client
type Twitter struct {
	Config      AuthConfig
	anacondaAPI *anaconda.TwitterApi
}

//Auth is a config map for Authentication
type AuthConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//BotConfig controls the behaviour of the bot online
type BotConfig struct {
	LikeOnHashTagMention    bool
	RetweetOnHashTagMention bool
	HashTagList             []string
	CommentOnHashTagMention bool
	CommentPoolList         []string
}
type StreamFilter struct {
	TextKeywords []string
	UserIds      []string
	Locations    []string
}
type tweetHandler func(tweet *Tweet)
type Tweet anaconda.Tweet

//Service Creates and Initializes a new Twitter client
func Service(config AuthConfig) *Twitter {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)

	api.SetLogger(&logger{logrus.New()})
	t := Twitter{
		Config:      config,
		anacondaAPI: api,
	}
	return &t
}

// FilterStream Starts listening to tweets.
func (t *Twitter) FilteredStream(streamFilter StreamFilter, tweetHandle tweetHandler) *Twitter {
	t.listen(t.anacondaAPI, streamFilter, tweetHandle)
	return t
}

func (t *Twitter) Stream(tweetHandle tweetHandler) *Twitter {
	t.listen(t.anacondaAPI, StreamFilter{}, tweetHandle)
	return t
}
func (t *Twitter) handle(tweet interface{}, tweetHandle tweetHandler) {
	switch v := tweet.(type) {
	case anaconda.Tweet:
		logrus.Tracef("%-15s: %s\n", v.User.ScreenName, v.Text)
		tw := Tweet(v)
		tweetHandle(&tw)
	default:
		logrus.Infof("Un recognized type found %T", v)
	}
}
func addToFilter(filterMap map[string][]string, key string, array []string) {
	if array != nil && len(array) > 0 {
		filterMap[key] = array
	}
}
func (t *Twitter) listen(api *anaconda.TwitterApi, filter StreamFilter, handler tweetHandler) {
	filterMap := map[string][]string{}

	addToFilter(filterMap, "track", filter.TextKeywords)
	addToFilter(filterMap, "follow", filter.UserIds)
	addToFilter(filterMap, "location", filter.Locations)

	for {
		if len(filterMap) > 0 {
			logrus.Infof("Starting Stream filtering On %v", filterMap)

			tweetStream := api.PublicStreamFilter(filterMap)
			for tw := range tweetStream.C {
				go t.handle(tw, handler)
			}
		} else {
			logrus.Infof("Starting Golbal Stream  %v", filterMap)
			tweetStream := api.PublicStreamSample(url.Values{})
			for tw := range tweetStream.C {
				go t.handle(tw, handler)
			}
		}

		logrus.Warn("Twitter client stream Broke. ReStarting ...")
	}
}

type logger struct{ *logrus.Logger }

func (l *logger) Critical(args ...interface{}) { l.Error(args...) }

func (l *logger) Criticalf(pattern string, args ...interface{}) { l.Errorf(pattern, args...) }

func (l *logger) Notice(args ...interface{}) { l.Info(args...) }

func (l *logger) Noticef(pattern string, args ...interface{}) { l.Infof(pattern, args...) }
