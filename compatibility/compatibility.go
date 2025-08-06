package compatibility

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	User1      Analysis
	User2      Analysis
	Similarity float64
}

type Analysis struct {
	Sentiment string
	// Formality float64
	// Favorites []string
}

type apiResponse struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

var (
	HfToken = flag.String("hfToken", "", "token")
	// sentiment = "https://router.huggingface.co/hf-inference/models/tabularisai/multilingual-sentiment-analysis"
	sentiment = "https://router.huggingface.co/hf-inference/models/j-hartmann/emotion-english-distilroberta-base"
)

func Assess(u1, u2 string) (Result, error) {
	r := Result{}
	m, err := queryApi(sentiment, u1)
	if err != nil {
		return Result{}, err
	}
	r.User1.Sentiment = highest(m)
	m, err = queryApi(sentiment, u2)
	if err != nil {
		return Result{}, err
	}
	r.User2.Sentiment = highest(m)
	return r, nil
}

// queries a map of scores from an hf inference
func queryApi(url, input string) (map[string]float64, error) {
	data, err := json.Marshal(map[string]string{"inputs": input})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *HfToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var apiRsp [][]apiResponse
	if err := json.Unmarshal(body, &apiRsp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	m := make(map[string]float64)
	for _, r := range apiRsp[0] {
		m[r.Label] = r.Score
	}
	return m, nil
}

func highest(m map[string]float64) string {
	if len(m) == 0 {
		return ""
	}

	var mk string
	var mv float64
	first := true

	for k, v := range m {
		if first || v > mv {
			mk = k
			mv = v
			first = false
		}
	}

	return mk
}
