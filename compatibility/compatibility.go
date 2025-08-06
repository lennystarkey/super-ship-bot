package compatibility

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type Result struct {
	User1 Analysis
	User2 Analysis
	Story string
	Style string
}

type Analysis struct {
	Sentiment string
	// Formality float64
	// Favorites []string
}

type tcResp struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

var (
	HfToken = flag.String("hfToken", "", "token")
	// textClassificationApi = "https://router.huggingface.co/hf-inference/models/tabularisai/multilingual-textClassificationApi-analysis"
	textClassificationApi = "https://router.huggingface.co/hf-inference/models/j-hartmann/emotion-english-distilroberta-base"
	textGenerationApi     = "https://router.huggingface.co/v1/chat/completions"
	textGenerationModel   = "openai/gpt-oss-20b:novita"

	styles = []string{"Shakespeare", "a biblical prophet", "a toddler", "a medieval knight", "gen alpha brainrot slang"}
)

func Assess(u1, u2 string) (Result, error) {
	r := Result{}
	m, err := queryTextClassificationApi(u1)
	if err != nil {
		return Result{}, err
	}
	r.User1.Sentiment = highest(m)
	m, err = queryTextClassificationApi(u2)
	if err != nil {
		return Result{}, err
	}
	r.User2.Sentiment = highest(m)

	r.Style = styles[rand.Intn(len(styles))]
	p := fmt.Sprintf("200 characters. dramatically describe the relationship between someone with %v emotions and someone with %v emotions, in the style of %v. Make sure to use complete sentences.", r.User1.Sentiment, r.User2.Sentiment, r.Style)
	s, err := queryTextGenerationApi(p)
	if err != nil {
		return Result{}, err
	}
	r.Story = s

	return r, nil
}

// queries a map of scores from an hf inference
func queryTextClassificationApi(input string) (map[string]float64, error) {
	data, err := json.Marshal(map[string]string{"inputs": input})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	req, err := http.NewRequest("POST", textClassificationApi, bytes.NewBuffer(data))
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

	var apiRsp [][]tcResp
	if err := json.Unmarshal(body, &apiRsp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	m := make(map[string]float64)
	for _, r := range apiRsp[0] {
		m[r.Label] = r.Score
	}
	return m, nil
}

func queryTextGenerationApi(input string) (string, error) {
	data := fmt.Appendf(nil, `{"messages": [{"role": "user","content": "%v"}],"model": "%v","stream": false}`, input, textGenerationModel)
	req, err := http.NewRequest("POST", textGenerationApi, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *HfToken))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var rsp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &rsp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return rsp.Choices[0].Message.Content, nil
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
