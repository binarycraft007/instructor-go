package googleai

import "github.com/google/generative-ai-go/genai"

type ChatRequest struct {
	Model   *genai.GenerativeModel
	Session *genai.ChatSession
	Parts   []genai.Part
}
