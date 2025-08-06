package compatibility

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Formality float64
	Sentiment float64
	Favorites []string
}

type apiResponse struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

var (
	hfToken           = flag.String("hfToken", "", "token")
	sentimentAnalysis = "https://router.huggingface.co/hf-inference/models/tabularisai/multilingual-sentiment-analysis"
)

func Assess(u1, u2 string) (Result, error) {
	if len(*hfToken) == 0 {
		return Result{}, errors.New("please provide Hugging Face Inference API token")
	}

	return Result{}, nil
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *hfToken))
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
