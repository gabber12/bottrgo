package classifier

type Service interface {
	Classify(text string) (cls *Classification, err error)
}

type Classification struct {
	SentimentScore float64
	Entities       []Entity
}

type Entity struct {
	Text string
	Type string
}
