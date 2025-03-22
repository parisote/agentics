package agentics

import (
	"context"
	"errors"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	anthropics_option "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	openai_option "github.com/openai/openai-go/option"
)

type ModelType string

const (
	OpenAI    ModelType = "openai"
	Anthropic ModelType = "anthropic"
)

type ModelClient struct {
	provider ModelProvider
}

type ModelProvider interface {
	Execute(ctx context.Context, prompt string, messages []string) (*ModelResponse, error)
}

type ModelResponse struct {
	Content string
}

func (r *ModelResponse) GetContent() string {
	return r.Content
}

type OpenAIProvider struct {
	Client openai.Client
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		Client: openai.NewClient(
			openai_option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		),
	}
}

func (p *OpenAIProvider) Execute(ctx context.Context, prompt string, messages []string) (*ModelResponse, error) {
	chatCompletion, err := p.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(messages[0]),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		return nil, errors.New("error executing openai model")
	}

	return &ModelResponse{
		Content: chatCompletion.Choices[0].Message.Content,
	}, nil
}

type AnthropicProvider struct {
	Client *anthropic.Client
}

func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		Client: anthropic.NewClient(
			anthropics_option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
		),
	}
}

func (p *AnthropicProvider) Execute(ctx context.Context, prompt string, messages []string) (*ModelResponse, error) {
	return nil, nil
}

func NewModelClient(modelType ModelType) *ModelClient {
	var provider ModelProvider

	switch modelType {
	case OpenAI:
		provider = NewOpenAIProvider()
	case Anthropic:
		provider = NewAnthropicProvider()
	}

	return &ModelClient{
		provider: provider,
	}
}
