package instructor

type Provider = string

const (
	ProviderOpenAI    Provider = "OpenAI"
	ProviderAnthropic Provider = "Anthropic"
	ProviderCohere    Provider = "Cohere"
	ProviderGoogleAI  Provider = "GoogleAI"
)
