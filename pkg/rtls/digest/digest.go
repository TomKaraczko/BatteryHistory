package digest

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidChallenge        = errors.New("digest challenge is invalid")
	ErrAlgorithmNotImplemented = errors.New("algorithm not implemented")
)

// Transport is an implementation of http.RoundTripper that takes care of http
// digest authentication.
type transport struct {
	Username  string
	Password  string
	Retry     int
	Transport http.RoundTripper
	mutex     sync.Mutex
}

type challenge struct {
	Realm     string
	Domain    string
	Nonce     string
	Opaque    string
	Stale     string
	Algorithm string
	Qop       string
}

type credentials struct {
	Username   string
	Realm      string
	Nonce      string
	DigestURI  string
	Algorithm  string
	Cnonce     string
	Opaque     string
	MessageQop string
	NonceCount int
	method     string
	password   string
}

// NewTransport creates a new digest transport using the http.DefaultTransport.
func NewTransport(username, password string, retry int) *transport {
	t := &transport{
		Username: username,
		Password: password,
		Retry:    retry,
	}
	t.Transport = http.DefaultTransport
	return t
}

// parseChallenge parses the digest challenge received from the server
func parseChallenge(input string) (*challenge, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "Digest ") {
		return nil, ErrInvalidChallenge
	}

	input = strings.TrimSpace(input[7:])
	parameters := strings.Split(input, ", ")
	c := &challenge{
		Algorithm: "MD5",
	}

	const qs = `"`
	for i := range parameters {
		parameter := strings.SplitN(parameters[i], "=", 2)
		switch parameter[0] {
		case "realm":
			c.Realm = strings.Trim(parameter[1], qs)
		case "domain":
			c.Domain = strings.Trim(parameter[1], qs)
		case "nonce":
			c.Nonce = strings.Trim(parameter[1], qs)
		case "opaque":
			c.Opaque = strings.Trim(parameter[1], qs)
		case "stale":
			c.Stale = strings.Trim(parameter[1], qs)
		case "algorithm":
			c.Algorithm = strings.Trim(parameter[1], qs)
		case "qop":
			c.Qop = strings.Trim(parameter[1], qs)
		default:
			return nil, ErrInvalidChallenge
		}
	}
	return c, nil
}

// hash generates an MD5 hash of the given data
func hash(data string) string {
	hf := md5.New()
	io.WriteString(hf, data)
	return fmt.Sprintf("%x", hf.Sum(nil))
}

// keyedDigest creates an MD5 digest of secret:data
func keyedDigest(secret, data string) string {
	return hash(fmt.Sprintf("%s:%s", secret, data))
}

// credentialsHash generates the hash of the username, realm, and password
func (c *credentials) credentialsHash() string {
	return hash(fmt.Sprintf("%s:%s:%s", c.Username, c.Realm, c.password))
}

// methodHash generates the hash of the HTTP method and the request URI
func (c *credentials) methodHash() string {
	return hash(fmt.Sprintf("%s:%s", c.method, c.DigestURI))
}

// responseHash generates the response hash used for digest authentication
func (c *credentials) responseHash(cnonce string) (string, error) {
	c.NonceCount++
	if c.MessageQop == "auth" {
		if cnonce != "" {
			c.Cnonce = cnonce
		} else {
			b := make([]byte, 8)
			io.ReadFull(rand.Reader, b)
			c.Cnonce = fmt.Sprintf("%x", b)[:16]
		}
		return keyedDigest(c.credentialsHash(), fmt.Sprintf("%s:%08x:%s:%s:%s",
			c.Nonce, c.NonceCount, c.Cnonce, c.MessageQop, c.methodHash())), nil
	} else if c.MessageQop == "" {
		return keyedDigest(c.credentialsHash(), fmt.Sprintf("%s:%s", c.Nonce, c.methodHash())), nil
	}
	return "", ErrAlgorithmNotImplemented
}

// authorize generates the Authorization header value for digest authentication
func (c *credentials) authorize() (string, error) {
	// Note that this is only implemented for MD5 and NOT MD5-sess.
	// MD5-sess is rarely supported and those that do are a big mess.
	if c.Algorithm != "MD5" {
		return "", ErrAlgorithmNotImplemented
	}
	// Note that this is NOT implemented for "qop=auth-int".  Similarly the
	// auth-int server side implementations that do exist are a mess.
	if c.MessageQop != "auth" && c.MessageQop != "" {
		return "", ErrAlgorithmNotImplemented
	}
	resp, err := c.responseHash("")
	if err != nil {
		return "", ErrAlgorithmNotImplemented
	}
	sl := []string{fmt.Sprintf(`username="%s"`, c.Username)}
	sl = append(sl, fmt.Sprintf(`realm="%s"`, c.Realm))
	sl = append(sl, fmt.Sprintf(`nonce="%s"`, c.Nonce))
	sl = append(sl, fmt.Sprintf(`uri="%s"`, c.DigestURI))
	sl = append(sl, fmt.Sprintf(`response="%s"`, resp))
	if c.Algorithm != "" {
		sl = append(sl, fmt.Sprintf(`algorithm="%s"`, c.Algorithm))
	}
	if c.Opaque != "" {
		sl = append(sl, fmt.Sprintf(`opaque="%s"`, c.Opaque))
	}
	if c.MessageQop != "" {
		sl = append(sl, fmt.Sprintf("qop=%s", c.MessageQop))
		sl = append(sl, fmt.Sprintf("nc=%08x", c.NonceCount))
		sl = append(sl, fmt.Sprintf(`cnonce="%s"`, c.Cnonce))
	}
	return fmt.Sprintf("Digest %s", strings.Join(sl, ", ")), nil
}

// newCredentials creates a new credentials instance based on the request and challenge
func (t *transport) newCredentials(req *http.Request, c *challenge) *credentials {
	return &credentials{
		Username:   t.Username,
		Realm:      c.Realm,
		Nonce:      c.Nonce,
		DigestURI:  req.URL.RequestURI(),
		Algorithm:  c.Algorithm,
		Opaque:     c.Opaque,
		MessageQop: c.Qop, // "auth" must be a single value
		NonceCount: 0,
		method:     req.Method,
		password:   t.Password,
	}
}

// RoundTrip performs an HTTP request, handling digest authentication if necessary
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Copy the request so we don't modify the input.
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header)
	for k, s := range req.Header {
		req2.Header[k] = s
	}

	// Make a request to get the 401 that contains the challenge.
	resp, err := t.Transport.RoundTrip(req)
	if err != nil || resp.StatusCode != 401 {
		return resp, err
	}

	challenge, err := parseChallenge(resp.Header.Get("WWW-Authenticate"))
	if err != nil {
		return resp, err
	}

	// Form credentials based on the challenge.
	credentials := t.newCredentials(req2, challenge)
	auth, err := credentials.authorize()
	if err != nil {
		return resp, err
	}

	// We'll no longer use the initial response, so close it
	resp.Body.Close()

	// Make authenticated request.
	req2.Header.Set("Authorization", auth)

	// Retry logic
	retry := t.Retry
	for {
		resp, err = t.Transport.RoundTrip(req2)
		if err == nil && resp.StatusCode != 401 || retry <= 0 {
			break
		}
		retry--
		time.Sleep(500 * time.Millisecond)
	}

	return resp, err
}

// Client returns an HTTP client that uses the digest transport.
func (t *transport) Client() (*http.Client, error) {
	return &http.Client{Transport: t}, nil
}
