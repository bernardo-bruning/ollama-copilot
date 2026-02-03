package adapters_test

import (
	"context"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/adapters"
	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type FakeDeepSeekHttpClient struct {
	Returns     adapters.DeepSeekResponse
	LastRequest any
}

// Post implements [HttpClient].
func (f *FakeDeepSeekHttpClient) Post(url string, req any, resp any) error {
	f.LastRequest = req
	if r, ok := resp.(*adapters.DeepSeekResponse); ok {
		r.Choices.Text = f.Returns.Choices.Text
	}
	return nil
}

func NewFakeDeepSeekHttpClient() *FakeDeepSeekHttpClient {
	return &FakeDeepSeekHttpClient{
		Returns: adapters.DeepSeekResponse{
			Choices: struct {
				Text string `json:"text"`
			}{
				Text: "func test():",
			},
		},
	}
}

func TestDeepSeek(t *testing.T) {
	t.Run("send to callback", func(t *testing.T) {
		httpClient := NewFakeDeepSeekHttpClient()
		deepSeek := adapters.NewDeepSeek("deepseek-coder", 100, httpClient)
		executed := false
		req := ports.CompletionRequest{
			Prompt:      "func ",
			Suffix:      "}",
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

		deepSeekRequest, ok := httpClient.LastRequest.(adapters.DeepSeekRequest)
		if !ok {
			t.Fatalf("expected DeepSeekRequest, got %T", httpClient.LastRequest)
		}

		if deepSeekRequest.Suffix != "}" {
			t.Errorf("expected suffix \"}\", got \"%s\"", deepSeekRequest.Suffix)
		}
	})

	t.Run("send request to deepseek", func(t *testing.T) {
		httpClient := NewFakeDeepSeekHttpClient()
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
