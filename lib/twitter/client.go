package twitter

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
)

// Tweet represents a tweet parsed from the API.
type Tweet struct {
	Data struct {
		ID       string
		Text     string
		Entities struct {
			Hashtags []struct {
				Tag string
			}
			URLs []struct {
				URL      string
				Expanded string `json:"expanded_url"`
			}
		}
	}
}

// A Sampler samples tweets until it is closed.
type Sampler interface {
	Open() (<-chan Tweet, error)
	Close() error
}

func encodeCredentials(consumerKey, consumerSecret string) string {
	credentials := url.UserPassword(consumerKey, consumerSecret).String()
	return base64.StdEncoding.EncodeToString([]byte(credentials))

}


type tokenResp struct {
	TokenType string `json:"token_type"`
	AccessToken string `json:"access_token"`
}


type basicAuthRoundTripper struct {
	username, password string
}

func (rt basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(rt.username, rt.password)
	return http.DefaultTransport.RoundTrip(req)
}

// NewSampler initializes a new Sampler with a a bearer token.
func NewSampler(consumerKey, consumerSecret string) (Sampler, error) {
	client := &http.Client{Transport: basicAuthRoundTripper{consumerKey,
		consumerSecret}}
	data :=url.Values{}
	data.Set("grant_type", "client_credentials")

	resp, err := client.PostForm("https://api.twitter.com/oauth2/token",data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token tokenResp
	err = json.NewDecoder(resp.Body).Decode(&token)
	return &sampler{bearerToken:token.AccessToken}, err

}

type sampler struct {
	bearerToken string
	isOpen bool
	err error
}

type bearerAuthRoundTripper struct {
	token string
}

func (rt bearerAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization","Bearer "+rt.token)
	return http.DefaultTransport.RoundTrip(req)
}

func (c *sampler) Open() (<-chan Tweet, error) {
	client := &http.Client{Transport:bearerAuthRoundTripper{c.bearerToken}}
	resp, err := client.Get("https://api.twitter." +
		"com/labs/1/tweets/stream/sample?format=detailed")
	if err != nil {
		return nil, err
	}

	c.isOpen = true
	decoder := json.NewDecoder(resp.Body)
	stream := make(chan Tweet, 1000)
	go func() {
		for c.isOpen {
			var tweet Tweet
			if c.err = decoder.Decode(&tweet); c.err != nil {
				break
			}
			stream <- tweet
		}
		resp.Body.Close()
		close(stream)
	}()
	return stream, nil
}

func (c *sampler) Close() error {
	c.isOpen = false
	return c.err
}