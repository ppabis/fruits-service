package e2e

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Sets fruits in the following setup
// alice => apple
// bob => banana
// bob => pineapple
// charlie => banana
// charlie => pineapple
// Expected end result: apple, banana, pineapple
func SetFruits(url string) error {
	var values = []struct {
		user  string
		fruit string
	}{
		{"alice", "apple"},
		{"bob", "banana"},
		{"bob", "pineapple"},
		{"charlie", "banana"},
		{"charlie", "pineapple"},
	}

	var currentUser string
	var cookie *http.Cookie
	for _, v := range values {
		if v.user != currentUser {
			cookie = login(url, httpClient, v.user)
			currentUser = v.user
		}

		if cookie == nil {
			return fmt.Errorf("failed to login as %s", v.user)
		}

		err := setFruit(url, v.fruit, httpClient, cookie)
		if err != nil {
			return err
		}

	}

	return nil
}

func login(url string, client *http.Client, user string) *http.Cookie {
	// Login
	form := u.Values{
		"username": {user},
		"password": {Password},
	}.Encode()
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.Jar = jar
	req, _ := http.NewRequest("POST", url+"/login", strings.NewReader(form))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
	cookies := jar.Cookies(req.URL)
	for _, cookie := range cookies {
		if cookie.Name == "session" {
			return cookie
		}
	}
	return nil
}

func setFruit(url string, fruit string, client *http.Client, cookie *http.Cookie) error {
	// Set fruit
	form := u.Values{
		"fruit": {fruit},
	}.Encode()
	req, err := http.NewRequest("PUT", url+"/fruit", strings.NewReader(form))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	_, err = client.Do(req)
	return err
}
