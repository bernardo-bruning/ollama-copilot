package adapters_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/adapters"
	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
	"github.com/stretchr/testify/assert"
)

type FakeOpenAIHttpClient struct {
	Returns adapters.OpenAIResponse
}

// Post implements [HttpClient].
func (f *FakeOpenAIHttpClient) Post(url string, req any, resp any) error {
	if r, ok := resp.(*adapters.OpenAIResponse); ok {
		r.Choices = f.Returns.Choices
	} else {
		panic(fmt.Sprintf("expected *adapters.OpenAIResponse but got %T", resp))
	}
	return nil
}

func NewOpenAIFakeHttpClient() *FakeOpenAIHttpClient {
	return &FakeOpenAIHttpClient{
		Returns: adapters.OpenAIResponse{
			Id:      "id",
			Created: 1,
			Choices: []adapters.ChoiceResponse{{Text: "func test():", Index: 0}},
		},
	}
}

func Test(t *testing.T) {
	t.Run("test request", func(t *testing.T) {
		req := adapters.NewOpenAIRequest("gpt-3.5-turbo-instruct", 20, ports.CompletionRequest{
			Prompt: "hello ",
			Suffix: "!",
		})

		assert.Equal(t, req, &adapters.OpenAIRequest{
			Prompt:    "hello ",
			Model:     "gpt-3.5-turbo-instruct",
			MaxTokens: 20,
			Suffix:    func(s string) *string { return &s }("!"),
		})
	})

	t.Run("test response", func(t *testing.T) {
		resp := adapters.NewOpenAIResponse()
		resp.Choices = append(resp.Choices, adapters.ChoiceResponse{Text: "hello"})
		assert.Equal(t, resp.ToPorts(), ports.CompletionResponse{
			Response: "hello",
			Done:     true,
		})
	})

	t.Run("test not apply suffix when model is not gpt-3.5-turbo-instruct", func(t *testing.T) {
		req := adapters.NewOpenAIRequest("gpt-4o-mini-2024-07-18", 20, ports.CompletionRequest{
			Prompt: "hello ",
			Suffix: "!",
		})

		assert.Nil(t, req.Suffix)
	})

	t.Run("test provider", func(t *testing.T) {
		client := NewOpenAIFakeHttpClient()
		provider := adapters.NewOpenAIWithClient("model", 10, client)
		req := ports.CompletionRequest{Prompt: "hello ", Suffix: "!"}
		err := provider.Completion(context.Background(), req, func(r ports.CompletionResponse) error {
			assert.Equal(t, r, ports.CompletionResponse{
				Response: "func test():",
				Done:     true,
			})
			return nil
		})

		assert.NoError(t, err)
	})

}
