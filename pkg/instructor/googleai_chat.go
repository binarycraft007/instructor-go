package instructor

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/binarycraft007/instructor-go/pkg/instructor/googleai"
)

func (i *InstructorGoogleAI) Chat(
	ctx context.Context,
	request *googleai.ChatRequest,
	response any,
) (*genai.GenerateContentResponse, error) {
	resp, err := chatHandler(i, ctx, request, response)
	if err != nil {
		if resp == nil {
			return &genai.GenerateContentResponse{}, err
		}
		return nilGoogleAIRespWithUsage(resp.(*genai.GenerateContentResponse)), err
	}

	return resp.(*genai.GenerateContentResponse), nil
}

func (i *InstructorGoogleAI) chat(ctx context.Context, request interface{}, schemaIn interface{}) (string, interface{}, error) {
	schema := schemaIn.(*genai.Schema)
	req, ok := request.(*googleai.ChatRequest)
	if !ok {
		return "", nil, fmt.Errorf("invalid request type for %s client", i.Provider())
	}
	req.Model.SetCandidateCount(1)

	switch i.Mode() {
	case ModeToolCall:
		return i.chatToolCall(ctx, req, schema)
	case ModeJSON:
		return i.chatJSON(ctx, req, schema)
	default:
		return "", nil, fmt.Errorf("mode '%s' is not supported for %s", i.Mode(), i.Provider())
	}
}

func (i *InstructorGoogleAI) chatToolCall(ctx context.Context, request *googleai.ChatRequest, schema *genai.Schema) (string, *genai.GenerateContentResponse, error) {
	panic("tool call not implemented googleai")
}

func (i *InstructorGoogleAI) chatJSON(ctx context.Context, request *googleai.ChatRequest, schema *genai.Schema) (string, *genai.GenerateContentResponse, error) {
	request.Model.GenerationConfig.ResponseMIMEType = "application/json"
	request.Model.GenerationConfig.ResponseSchema = schema

	resp, err := request.Session.SendMessage(ctx, request.Parts...)
	if err != nil {
		return "", nil, err
	}

	var respText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			respText += "\n" + string(textPart)
		}
	}

	return respText, resp, nil
}

func (i *InstructorGoogleAI) emptyResponseWithUsageSum(usage *UsageSum) interface{} {
	return &genai.GenerateContentResponse{
		UsageMetadata: &genai.UsageMetadata{
			PromptTokenCount:        int32(usage.InputTokens),
			CachedContentTokenCount: int32(usage.OutputTokens),
			TotalTokenCount:         int32(usage.TotalTokens),
		},
	}
}

func (i *InstructorGoogleAI) emptyResponseWithResponseUsage(response interface{}) interface{} {
	resp, ok := response.(*genai.GenerateContentResponse)
	if !ok || resp == nil {
		return nil
	}

	return &genai.GenerateContentResponse{
		UsageMetadata: resp.UsageMetadata,
	}
}

func (i *InstructorGoogleAI) addUsageSumToResponse(response interface{}, usage *UsageSum) (interface{}, error) {
	resp, ok := response.(*genai.GenerateContentResponse)
	if !ok {
		return response, fmt.Errorf("internal type error: expected *openai.ChatCompletionResponse, got %T", response)
	}

	resp.UsageMetadata.PromptTokenCount += int32(usage.InputTokens)
	resp.UsageMetadata.CandidatesTokenCount += int32(usage.OutputTokens)
	resp.UsageMetadata.TotalTokenCount += int32(usage.TotalTokens)

	return response, nil
}

func (i *InstructorGoogleAI) countUsageFromResponse(response interface{}, usage *UsageSum) *UsageSum {
	resp, ok := response.(*genai.GenerateContentResponse)
	if !ok {
		return usage
	}

	usage.InputTokens += int(resp.UsageMetadata.PromptTokenCount)
	usage.OutputTokens += int(resp.UsageMetadata.CandidatesTokenCount)
	usage.TotalTokens += int(resp.UsageMetadata.TotalTokenCount)

	return usage
}

func nilGoogleAIRespWithUsage(resp *genai.GenerateContentResponse) *genai.GenerateContentResponse {
	if resp == nil {
		return nil
	}

	return &genai.GenerateContentResponse{
		UsageMetadata: resp.UsageMetadata,
	}
}
