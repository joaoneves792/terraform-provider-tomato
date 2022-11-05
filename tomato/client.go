package tomato

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//  "github.com/hashicorp/terraform-plugin-log/tflog"
)

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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: tr},
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
func (c *Client) applyChange(service, nvram_entries string) (string, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return "", fmt.Errorf("define username and password")
	}

	rb := fmt.Sprintf("_ajax=1&_service=%s&%s&_http_id=TID944e75807cf4c9c3", service, nvram_entries)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tomato.cgi?", c.HostURL), bytes.NewBuffer([]byte(rb)))
	if err != nil {
		return "", err
	}

	b, err := c.doRequest(req)
	if err != nil {
		return string(b), err
	}

	return string(b), nil

}

// retrieve NVRAM
func (c *Client) getNVRAM() (map[string]string, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cfg/tomato.cfg?_http_id=TID944e75807cf4c9c3", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	b, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	pad := make([]byte, 1)
	pad[0] = 0x00

	//	tflog.Info(ctx, "NVRAM 8 byte Header: "+hex.EncodeToString(b[0:8]))

	length := int(binary.BigEndian.Uint32(append(pad, b[4:7]...)))
	magic := b[7]
	//	tflog.Info(ctx, "NVRAM dump length: "+strconv.Itoa(length))

	for i := 8; i < length; i++ {
		if b[i] > (0xfd - 0x1) {
			b[i] = 0x0
		} else {
			b[i] = 0xff + magic - b[i]
		}
	}

	n := make(map[string]string)

	se := 8
	for se < length+8 {
		nb := bytes.IndexByte(b[se:length+8], 0x00)
		cfg := string(b[se : se+nb])
		se = se + nb + 1
		eq := bytes.Index([]byte(cfg), []byte("="))
		if eq == -1 {
			break
		}
		n[string(cfg[0:eq])] = string(cfg[eq+1:])
	}

	return n, nil

}
