package adapters

import (
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

func TestOllamaPromptGeneration(t *testing.T) {
	tests := []struct {
		name        string
		templateStr string
		req         ports.CompletionRequest
		want        string
		wantErr     bool
	}{
		{
			name:        "Default template",
			templateStr: "<PRE> {{.Prefix}} <SUF> {{.Suffix}} <MID>",
			req: ports.CompletionRequest{
				Prompt: "var x = ",
				Suffix: "\nfunc main() {}",
			},
			want: "<PRE> var x =  <SUF> \nfunc main() {} <MID>",
		},
		{
			name:        "Template with Prompt alias",
			templateStr: "<START>{{.Prompt}}<END>",
			req: ports.CompletionRequest{
				Prompt: "hello",
				Suffix: "world",
			},
			want: "<START>hello<END>",
		},
		{
			name:        "Template with Prediction alias",
			templateStr: "<START>{{.Prediction}}<END>",
			req: ports.CompletionRequest{
				Prompt: "hello",
				Suffix: "world",
			},
			want: "<START>hello<END>",
		},
		{
			name:        "Template with missing variable",
			templateStr: "{{.Prefix}} {{.Missing}}",
			req: ports.CompletionRequest{
				Prompt: "hello",
			},
			want: "hello <no value>", // text/template default behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewOllama("model", 10, "system", tt.templateStr)
			if err != nil {
				t.Fatalf("NewOllama failed: %v", err)
			}

			ollama, ok := provider.(*Ollama)
			if !ok {
				t.Fatalf("Provider is not *Ollama")
			}

			got, err := ollama.generatePrompt(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generatePrompt() = %q, want %q", got, tt.want)
			}
		})
	}
}