package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(`C:\Users\Atai\Desktop\newGoLang\.env`)
}

type AnalyzeResult struct {
	Sentiment string   `json:"sentiment"`
	Tags      []string `json:"tags"`
}

func groqRequest(prompt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")

	body, _ := json.Marshal(map[string]any{
		"model": "llama-3.1-8b-instant",
		"messages": []map[string]any{
			{"role": "user", "content": prompt},
		},
		"temperature":          1,
		"max_completion_tokens": 1024,
		"top_p":                1,
		"stream":               false,
	})

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	choices, ok := result["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response: %v", result)
	}

	text, ok := choices[0].(map[string]any)["message"].(map[string]any)["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in response")
	}

	return text, nil
}

func AnalyzeEntry(content string) (*AnalyzeResult, error) {
	prompt := fmt.Sprintf(`Analyze this journal entry and respond ONLY with valid JSON, no markdown:
{"sentiment": "positive", "tags": ["tag1", "tag2"]}
sentiment must be: positive, neutral, or negative

Entry: "%s"`, content)

	text, err := groqRequest(prompt)
	if err != nil {
		return nil, err
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var analyzed AnalyzeResult
	if err := json.Unmarshal([]byte(text), &analyzed); err != nil {
		return nil, fmt.Errorf("parse error: %w, raw: %s", err, text)
	}

	return &analyzed, nil
}

func GenerateInsight(entries string) (string, error) {
	prompt := fmt.Sprintf(`You are a personal journal assistant. Analyze these journal entries and give a short, warm, personal insight in 2-3 sentences. Focus on patterns, emotions, and one practical suggestion. Plain text only.

Entries:
%s`, entries)

	return groqRequest(prompt)
}