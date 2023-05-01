package fruits

import (
	"encoding/json"
	"fmt"

	"monolith/config"
	"net/http"
	"strings"
	"time"
)

func GetFruitsFromMicroservice() (map[string]string, error) {
	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	req, err := http.NewRequest("GET", config.FruitsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body := make([]byte, 1024) // Read at most 1 kilobyte
		resp.Body.Read(body)
		resp.Body.Close()
		error_string := strings.Trim(string(body), "\x00")
		return nil, fmt.Errorf("[GET /] fruits microservice returned %d: %s", resp.StatusCode, error_string)
	}

	jsonFruits := make([]struct {
		Username string `json:"username"`
		Fruit    string `json:"fruit"`
	}, 0, 32)
	err = json.NewDecoder(resp.Body).Decode(&jsonFruits)
	if err != nil {
		return nil, err
	}

	fruits := make(map[string]string)
	for _, f := range jsonFruits {
		fruits[f.Username] = f.Fruit
	}

	return fruits, nil
}
