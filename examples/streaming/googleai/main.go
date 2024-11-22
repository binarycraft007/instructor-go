package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/binarycraft007/instructor-go/pkg/instructor"
	"github.com/binarycraft007/instructor-go/pkg/instructor/googleai"
	"google.golang.org/api/option"
)

type HistoricalFact struct {
	Decade      string `json:"decade"       schema:"description=Decade when the fact occurred"`
	Topic       string `json:"topic"        schema:"description=General category or topic of the fact"`
	Description string `json:"description"  schema:"description=Description or details of the fact"`
}

func (hf HistoricalFact) String() string {
	return fmt.Sprintf(`
Decade:         %s
Topic:          %s
Description:    %s`, hf.Decade, hf.Topic, hf.Description)
}

func main() {
	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		panic(err)
	}
	defer genaiClient.Close()

	model := genaiClient.GenerativeModel("gemini-1.5-flash")
	model.SetMaxOutputTokens(2500)
	cs := model.StartChat()

	client := instructor.FromGoogleAI(
		genaiClient,
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)

	hfStream, err := client.ChatStream(ctx, &googleai.ChatRequest{
		Model:   model,
		Session: cs,
		Parts: []genai.Part{
			genai.Text("Tell me about the history of artificial intelligence up to year 2000"),
		},
	},
		*new(HistoricalFact),
	)
	if err != nil {
		panic(err)
	}

	for instance := range hfStream {
		hf := instance.(*HistoricalFact)
		println(hf.String())
	}
	/*
	   Decade:         1950s
	   Topic:          Birth of AI
	   Description:    The term 'Artificial Intelligence' is coined by John McCarthy at the Dartmouth Conference in 1956, considered the birth of AI as a field. Early research focuses on areas like problem solving, search algorithms, and logic.

	   Decade:         1960s
	   Topic:          Expert Systems and LISP
	   Description:    The language LISP is developed, which becomes widely used in AI applications. Research also leads to the development of expert systems, which emulate human decision-making abilities in specific domains.

	   Decade:         1970s
	   Topic:          AI Winter
	   Description:    AI experiences its first 'winter', a period of reduced funding and interest due to unmet expectations. Despite this, research continues in areas like knowledge representation and natural language processing.

	   Decade:         1980s
	   Topic:          Machine Learning and Neural Networks
	   Description:    The field of machine learning emerges, with a focus on developing algorithms that can learn from data. Neural networks, inspired by the structure of biological brains, gain traction during this decade.

	   Decade:         1990s
	   Topic:          AI in Practice
	   Description:    AI starts to find practical applications in various industries. Speech recognition, image processing, and expert systems are used in fields like healthcare, finance, and manufacturing.
	*/
}

func toPtr[T any](val T) *T {
	return &val
}
