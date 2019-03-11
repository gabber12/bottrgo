package classifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type MachineBox struct {
	HostPort string
}
type response struct {
	Sentences []sentence
	success   bool
}
type sentence struct {
	Text      string
	Sentiment float64
	Entities  []Entity
}

//Classify Returns
func (m *MachineBox) Classify(text string) (*Classification, error) {
	TEXT_CLASSIFICATION_API := fmt.Sprintf("%s/textbox/check", m.HostPort)
	resp, err := http.PostForm(TEXT_CLASSIFICATION_API, url.Values{"text": []string{text}})
	if err != nil {
		logrus.Errorf("Couldnot reach Machinebox server %v", err)
		return nil, fmt.Errorf("Could not reach MachineBox Server %v", err)
	}
	sen := response{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read body %v", err)
	}

	json.Unmarshal(body, &sen)
	sentiment := 0.0
	var entities []Entity
	for _, v := range sen.Sentences {
		sentiment = v.Sentiment + sentiment
		entities = append(entities, v.Entities...)

	}
	sentimentScore := sentiment / float64(len(sen.Sentences))

	return &Classification{SentimentScore: sentimentScore, Entities: entities}, nil
}
