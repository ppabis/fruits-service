package e2e

import (
	"fmt"
	"net/http"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// Gets all the fruits from the current monolith
// And compares it to the expected result
func GetFruits(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Prefer-Data", "monolith")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Header.Get("X-Data-Source") != "monolith" {
		return fmt.Errorf("expected X-Data-Source to be monolith, got %s", resp.Header.Get("X-Data-Source"))
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	sel, _ := cascadia.Parse("body ul li")
	matches := cascadia.QueryAll(body, sel)
	if len(matches) != 3 {
		return fmt.Errorf("expected 3 fruits, got %d", len(matches))
	}

	expected := []string{"alice: apple", "bob: banana", "charlie: pineapple"}

	for i, m := range matches {
		if m.FirstChild.Data != expected[i] {
			return fmt.Errorf("expected %x, got %x", expected[i], m.FirstChild.Data)
		}
	}

	return nil

}
