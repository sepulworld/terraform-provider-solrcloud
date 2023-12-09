package provider

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:8983"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if username == nil || password == nil {
		return &c, nil
	}

	if username != nil && password != nil {
		basicAuth := base64.StdEncoding.EncodeToString([]byte(*username + ":" + *password))
		c.HTTPClient.Transport = &basicAuthTransport{
			Transport: http.DefaultTransport,
			Username:  *username,
			Password:  *password,
			BasicAuth: "Basic " + basicAuth,
		}
	}

	return &c, nil
}

type basicAuthTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
	BasicAuth string
}

func (bat *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", bat.BasicAuth)
	return bat.Transport.RoundTrip(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return nil, err
	}
	ctx := req.Context()
	tflog.Info(ctx, string(requestDump))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
