package adapters_test

import (
	"context"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/adapters"
	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type FakeHttpClient struct {
	Returns adapters.DeepSeekResponse
}

// Post implements [HttpClient].
func (f *FakeHttpClient) Post(url string, req any, resp any) error {
	if r, ok := resp.(*adapters.DeepSeekResponse); ok {
		r.Choices.Text = f.Returns.Choices.Text
	}
	return nil
}

func NewFakeHttpClient() *FakeHttpClient {
	return &FakeHttpClient{
		Returns: adapters.DeepSeekResponse{
			Choices: struct {
				Text string "json:\"text\""
			}{
				Text: "func test():",
			},
		},
	}
}

func TestDeepSeek(t *testing.T) {
	t.Run("send to callback", func(t *testing.T) {
		httpClient := NewFakeHttpClient()
		deepSeek := adapters.NewDeepSeek("deepseek-coder", 100, httpClient)
		executed := false
		req := ports.CompletionRequest{
			Prompt:      "func ",
			Temperature: 1.0,
			TopP:        2,
			Stop:        []string{"EOF"},
		}

		err := deepSeek.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
			executed = true
			return nil
		})

		if err != nil {
			t.Fatalf("error to execute completion")
		}

		if !executed {
			t.Fatal("Deepseek not call function completion")
		}
	})

	t.Run("send request to deepseek", func(t *testing.T) {
		httpClient := NewFakeHttpClient()
		deepSeek := adapters.NewDeepSeek("deepseek-coder", 100, httpClient)
		req := ports.CompletionRequest{
			Prompt:      "func ",
			Temperature: 1.0,
			TopP:        2,
			Stop:        []string{"EOF"},
		}

		err := deepSeek.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
			if resp.Response != "func test():" {
				t.Fatalf("unexpected response \"%s\"", resp.Response)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("invalid error: %s", err.Error())
		}
	})
}
