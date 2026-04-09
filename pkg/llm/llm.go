package llm


import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

)

type AnalyzeResult struct {
    Sentiment string   `json:"sentiment"`
    Tags      []string `json:"tags"`
}

func AnalyzeEntry(content string) (*AnalyzeResult, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")

    prompt := fmt.Sprintf(`Analyze this journal entry and respond ONLY with JSON, no other text:
{"sentiment": "positive|neutral|negative", "tags": ["tag1", "tag2"]}

Entry: "%s"`, content)

    body, _ := json.Marshal(map[string]any{
        "contents": []map[string]any{
            {"parts": []map[string]any{{"text": prompt}}},
        },
    })

    url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]any
    json.NewDecoder(resp.Body).Decode(&result)

    // достаём текст ответа
    text := result["candidates"].([]any)[0].(map[string]any)["content"].(map[string]any)["parts"].([]any)[0].(map[string]any)["text"].(string)

    var analyzed AnalyzeResult
    if err := json.Unmarshal([]byte(text), &analyzed); err != nil {
        return nil, fmt.Errorf("can't parse LLM response: %w", err)
    }

    return &analyzed, nil
}