package agentics

import (
	"context"
	"encoding/json"
	"fmt"
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
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		Client: openai.NewClient(
			openai_option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		),
	}
}

func (p *OpenAIProvider) Execute(ctx context.Context, prompt string, messages []Message, tools []ToolInterface) (*ModelResponse, error) {
	openAIMessages := p.toOpenAIMessages(messages)
	newMessages := []openai.ChatCompletionMessageParamUnion{}
	newMessages = append(newMessages, openai.SystemMessage(prompt))
	newMessages = append(newMessages, openAIMessages...)
	chatCompletion, err := p.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: newMessages,
		Model:    openai.ChatModelGPT4o,
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

	for _, message := range m {
		switch message.Role {
		case "tool":
			result = append(result, openai.ToolMessage(message.Content, message.ToolCallID))
		case "user":
			result = append(result, openai.UserMessage(message.Content))
		case "assistant":
			result = append(result, openai.AssistantMessage(message.Content))
		case "assistant_tool":
			var obj openai.ChatCompletionMessageParamUnion
			fmt.Println("mensaje1 ", message.Content)
			if err := json.Unmarshal([]byte(message.Content), &obj); err != nil {
				fmt.Println("ERROR")
			}
			if obj.OfTool == nil {
				panic("")
			}
			fmt.Println("mensaje ", obj)
			result = append(result, obj)
		case "system":
			result = append(result, openai.SystemMessage(message.Content))
		}
	}

	return result
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

func (p *AnthropicProvider) Execute(ctx context.Context, prompt string, messages []Message, tools []ToolInterface) (*ModelResponse, error) {
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
