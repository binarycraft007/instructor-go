package instructor

import (
	"github.com/google/generative-ai-go/genai"
)

type InstructorGoogleAI struct {
	Client *genai.Client

	provider   Provider
	mode       Mode
	maxRetries int
	validate   bool
}

var _ Instructor = &InstructorAnthropic{}

func FromGoogleAI(client *genai.Client, opts ...Options) *InstructorGoogleAI {
	options := mergeOptions(opts...)

	i := &InstructorGoogleAI{
		Client: client,

		provider:   ProviderGoogleAI,
		mode:       *options.Mode,
		maxRetries: *options.MaxRetries,
		validate:   *options.validate,
	}
	return i
}

func (i *InstructorGoogleAI) MaxRetries() int {
	return i.maxRetries
}

func (i *InstructorGoogleAI) Mode() string {
	return i.mode
}

func (i *InstructorGoogleAI) Provider() string {
	return i.provider
}
func (i *InstructorGoogleAI) Validate() bool {
	return i.validate
}
