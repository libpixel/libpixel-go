// Package libpixel provides a Client to generate and sign LibPixel URLs.
package libpixel

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"net/url"
)

// The Client to sign and/or generate URLs with.
type Client struct {
	Host   string
	HTTPS  bool
	Secret string
}

// Params for the LibPixel Image API.
//
// For API documentation, see: http://libpixel.com/docs/#image-api
type Params map[string]interface{}

// Sign generates and adds a signature to a URL. A Secret must be provided in
// the Client for this method to work. Returns the signed URL.
func (c *Client) Sign(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return c.sign(u), nil
}

// URL generates a LibPixel URL for a given path and Params. If the Client has a
// Secret defined, the URL will automatically be signed. The Client must define
// a Host for this method to work. Returns the generated URL.
func (c *Client) URL(path string, params Params) (string, error) {
	if path == "" {
		path = "/"
	}

	u := &url.URL{Scheme: "http", Host: c.Host, Path: path}

	if c.HTTPS {
		u.Scheme = "https"
	}

	qs := make(url.Values)
	for k, v := range params {
		qs.Add(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = qs.Encode()

	if c.Secret == "" {
		return u.String(), nil
	}

	return c.sign(u), nil
}

func (c *Client) sign(u *url.URL) string {
	if u.Path == "" {
		u.Path = "/"
	}

	s := u.Path
	if u.RawQuery != "" {
		s += "?" + u.RawQuery
	}

	m := hmac.New(sha1.New, []byte(c.Secret))
	m.Write([]byte(s))
	signature := fmt.Sprintf("%x", m.Sum(nil))

	if u.RawQuery != "" {
		u.RawQuery += "&"
	}
	u.RawQuery += "signature=" + signature

	return u.String()
}
