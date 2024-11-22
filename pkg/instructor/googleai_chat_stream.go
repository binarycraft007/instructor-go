package instructor

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/binarycraft007/instructor-go/pkg/instructor/googleai"
	"google.golang.org/api/iterator"
)

func (i *InstructorGoogleAI) ChatStream(
	ctx context.Context,
	request *googleai.ChatRequest,
	responseType any,
) (<-chan any, error) {
	stream, err := chatStreamHandler(i, ctx, request, responseType)
	if err != nil {
		return nil, err
	}

	return stream, err
}

func (i *InstructorGoogleAI) chatStream(ctx context.Context, request interface{}, schemaIn interface{}) (<-chan string, error) {
	schema := schemaIn.(*genai.Schema)
	req, ok := request.(*googleai.ChatRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type for %s client", i.Provider())
	}
	req.Model.SetCandidateCount(1)

	switch i.Mode() {
	case ModeJSON:
		return i.chatJSONStream(ctx, req, schema)
	default:
		return nil, fmt.Errorf("mode '%s' is not supported for %s", i.Mode(), i.Provider())
	}
}

func (i *InstructorGoogleAI) chatJSONStream(ctx context.Context, request *googleai.ChatRequest, schema *genai.Schema) (<-chan string, error) {
	request.Model.GenerationConfig.ResponseMIMEType = "application/json"
	request.Model.GenerationConfig.ResponseSchema = schema

	// Create a channel to stream text responses
	ch := make(chan string)

	// Send the request asynchronously
	go func() {
		defer close(ch)

		iter := request.Session.SendMessageStream(ctx, request.Parts...)
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// Send error information through a logger or handle it gracefully
				ch <- "Error: " + err.Error()
				return
			}

			// Extract and stream response content
			for _, part := range resp.Candidates[0].Content.Parts {
				if textPart, ok := part.(genai.Text); ok {
					ch <- string(textPart)
				}
			}
		}
	}()
	return ch, nil
}
