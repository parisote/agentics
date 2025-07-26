package agentics

import (
	"context"
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
	Execute(ctx context.Context, prompt string, messages []Message, tools []ToolInterface) (*ModelResponse, error)
	ExecuteWithFollowUp(ctx context.Context, prompt string, messages []Message, tools []ToolInterface, toolResult string) (*ModelResponse, error)
	GetModel() string
	SetModel(model string)
}

type ModelResponse struct {
	IsToolCall bool
	ToolCalls  []ToolCall
	Content    string
	Params     []byte //bytes
}

type ToolCall struct {
	Name       string
	Arguments  string
	ToolCallID string
}

func (r *ModelResponse) GetContent() string {
	return r.Content
}

type OpenAIProvider struct {
	Client openai.Client
	Model  string
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		Client: openai.NewClient(
			openai_option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		),
		Model: "gpt-4o",
	}
}

func (p *OpenAIProvider) GetModel() string {
	return p.Model
}

func (p *OpenAIProvider) SetModel(model string) {
	p.Model = model
}

func (p *OpenAIProvider) ExecuteWithFollowUp(ctx context.Context, prompt string, messages []Message, tools []ToolInterface, toolResult string) (*ModelResponse, error) {
	openAIMessages := p.toOpenAIMessages(messages)
	newMessages := []openai.ChatCompletionMessageParamUnion{}
	newMessages = append(newMessages, openai.SystemMessage(prompt))
	newMessages = append(newMessages, openAIMessages...)
	newMessages = append(newMessages, openai.UserMessage(toolResult))

	chatCompletion, err := p.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: newMessages,
		Model:    openai.ChatModel(p.Model),
		Tools:    p.getOpenAITools(tools),
	})
	if err != nil {
		return nil, err
	}

	params, _ := chatCompletion.Choices[0].Message.ToParam().MarshalJSON()

	return &ModelResponse{
		IsToolCall: chatCompletion.Choices[0].FinishReason == "tool_calls",
		ToolCalls:  p.getToolCalls(chatCompletion.Choices[0].Message.ToolCalls),
		Content:    chatCompletion.Choices[0].Message.Content,
		Params:     params,
	}, nil
}

func (p *OpenAIProvider) Execute(ctx context.Context, prompt string, messages []Message, tools []ToolInterface) (*ModelResponse, error) {
	openAIMessages := p.toOpenAIMessages(messages)
	newMessages := []openai.ChatCompletionMessageParamUnion{}
	newMessages = append(newMessages, openai.SystemMessage(prompt))
	newMessages = append(newMessages, openAIMessages...)
	chatCompletion, err := p.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: newMessages,
		Model:    openai.ChatModel(p.Model),
		Tools:    p.getOpenAITools(tools),
	})
	if err != nil {
		return nil, err
	}

	params, _ := chatCompletion.Choices[0].Message.ToParam().MarshalJSON()

	return &ModelResponse{
		IsToolCall: chatCompletion.Choices[0].FinishReason == "tool_calls",
		ToolCalls:  p.getToolCalls(chatCompletion.Choices[0].Message.ToolCalls),
		Content:    chatCompletion.Choices[0].Message.Content,
		Params:     params,
	}, nil
}

func (p *OpenAIProvider) getToolCalls(toolCalls []openai.ChatCompletionMessageToolCall) []ToolCall {
	result := []ToolCall{}

	for _, toolCall := range toolCalls {
		result = append(result, ToolCall{
			Name:       toolCall.Function.Name,
			Arguments:  toolCall.Function.Arguments,
			ToolCallID: toolCall.ID,
		})
	}

	return result
}

func (p *OpenAIProvider) getOpenAITools(tools []ToolInterface) []openai.ChatCompletionToolParam {
	result := []openai.ChatCompletionToolParam{}
	for _, tool := range tools {
		properties := make(map[string]interface{})

		for _, param := range tool.GetParameters() {
			properties[param.Name] = map[string]string{
				"type": param.Type,
			}
		}

		result = append(result, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.GetName(),
				Description: openai.String(tool.GetDescription()),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": properties,
				},
			},
		})
	}

	return result
}

func (p *OpenAIProvider) toOpenAIMessages(m []Message) []openai.ChatCompletionMessageParamUnion {
	result := []openai.ChatCompletionMessageParamUnion{}

	for i, message := range m {
		switch message.Role {
		case "tool":
			if i == 0 || m[i-1].Role != "assistant" {
				result = append(result, openai.AssistantMessage(""))
			}
			result = append(result, openai.ToolMessage(message.Content, message.ToolCallID))
		case "user":
			result = append(result, openai.UserMessage(message.Content))
		case "assistant":
			if len(message.ToolCalls) > 0 {
				result = append(result, openai.AssistantMessage(message.Content))
			} else {
				result = append(result, openai.AssistantMessage(message.Content))
			}
		case "assistant_tool":
			continue
		case "system":
			result = append(result, openai.SystemMessage(message.Content))
		}
	}

	return result
}

type AnthropicProvider struct {
	Client *anthropic.Client
	Model  string
}

func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		Client: anthropic.NewClient(
			anthropics_option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
		),
		Model: "claude-3-sonnet-20240229",
	}
}

func (p *AnthropicProvider) GetModel() string {
	return p.Model
}

func (p *AnthropicProvider) SetModel(model string) {
	p.Model = model
}

func (p *AnthropicProvider) Execute(ctx context.Context, prompt string, messages []Message, tools []ToolInterface) (*ModelResponse, error) {
	// TODO: Implementar para Anthropic
	return nil, nil
}

func (p *AnthropicProvider) ExecuteWithFollowUp(ctx context.Context, prompt string, messages []Message, tools []ToolInterface, toolResult string) (*ModelResponse, error) {
	// TODO: Implementar para Anthropic
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

func NewModelClientWithModel(modelType ModelType, model string) *ModelClient {
	client := NewModelClient(modelType)
	client.provider.SetModel(model)
	return client
}
