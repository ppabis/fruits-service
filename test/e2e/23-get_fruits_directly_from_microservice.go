package e2e

import (
	"encoding/json"
	"fmt"
)

func GetFruitsDirectlyFromMicroservice(url string) error {
	resp, err := httpClient.Get(url + "/fruit/2")
	if err != nil {
		return err
	}

	var record struct {
		Username string `json:"username"`
		Fruit    string `json:"fruit"`
	}
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		return err
	}

	if record.Username != "bob" {
		return fmt.Errorf("expected bob, got %s", record.Username)
	}

	if record.Fruit != "banana" {
		return fmt.Errorf("expected banana, got %s", record.Fruit)
	}

	return nil
}
