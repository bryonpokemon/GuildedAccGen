package guilded

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func New(p string) *GuildeadClient {
	g := &GuildeadClient{}
	proxy, _ := url.Parse(p)
	// never returns an error
	jar, _ := cookiejar.New(nil)
	// i don't even know if they do tls checks, but I like using this lol
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				0x1302, 0x1303, 0x1301, 0xC02C, 0xC030, 0xC02B, 0xC02F, 0xCCA9,
				0xCCA8, 0x009F, 0x009E, 0xCCAA, 0xC0AF, 0xC0AD, 0xC0AE, 0xC0AC,
				0xC024, 0xC028, 0xC023, 0xC027, 0xC00A, 0xC014, 0xC009, 0xC013,
				0xC0A3, 0xC09F, 0xC0A2, 0xC09E, 0x006B, 0x0067, 0x0039, 0x0033,
				0x009D, 0x009C, 0xC0A1, 0xC09D, 0xC0A0, 0xC09C, 0x003D, 0x003C,
				0x0035, 0x002F, 0x00FF,
			},
			InsecureSkipVerify: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveID(0x001D),
				tls.CurveID(0x0017),
				tls.CurveID(0x0018),
				tls.CurveID(0x0019),
				tls.CurveID(0x0100),
				tls.CurveID(0x0101),
			},
		},
		Proxy: http.ProxyURL(proxy),
	}
	g.Client = &http.Client{
		Timeout:   20 * time.Second,
		Transport: tr,
		Jar:       jar,
	}

	err := g.getCookie()
	if err != nil {
		log.Println("error getting cookies", err)
	}
	g.ClienID = uuid.New().String()
	g.DeviceID = randHexString(64)
	return g
}

func (g *GuildeadClient) CreateAccount() error {
	g.Username = randHexString(10)
	g.Email = fmt.Sprintf("%s+%s@gmail.com", g.EmailBase, randHexString(5))
	g.Password = randHexString(11)

	payload := &RegisterPayload{
		ExtraInfo: ExtraInfo{
			Platform: "desktop",
		},
		Name:     g.Username,
		Email:    g.Email,
		Password: g.Password,
		FullName: g.Username,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://www.guilded.gg/api/users?type=email", bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Content-Length":      []string{fmt.Sprint(len(b))},
		"Sec-Ch-Ua":           []string{`" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
		"Sec-Ch-Ua-Mobile":    []string{`?0`},
		"User-Agent":          []string{`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36`},
		"Content-Type":        []string{`application/json`},
		"Accept":              []string{`application/json, text/javascript, */*; q=0.01`},
		"X-Requested-With":    []string{`XMLHttpRequest`},
		"Sec-Ch-Ua-Platform":  []string{`"macOS"`},
		"Origin":              []string{`https://www.guilded.gg`},
		"Sec-Fetch-Site":      []string{`same-origin`},
		"Sec-Fetch-Mode":      []string{`cors`},
		"Sec-Fetch-Dest":      []string{`empty`},
		"Referer":             []string{"https://www.guilded.gg/"},
		"guilded-client-id":   []string{g.ClienID},
		"guilded-device-id":   []string{g.DeviceID},
		"guilded-device-type": []string{"desktop"},
		"Accept-Language":     []string{`fr-FR,fr;q=0.9`},
	}
	res, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: wrong status code %d, %s", res.StatusCode, string(body))
	}

	return nil
}

func (g *GuildeadClient) Login() error {

	b := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, g.Email, g.Password)
	req, err := http.NewRequest("POST", "https://www.guilded.gg/api/login", strings.NewReader(b))
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Content-Length":     []string{fmt.Sprint(len(b))},
		"Sec-Ch-Ua":          []string{`" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
		"Sec-Ch-Ua-Mobile":   []string{`?0`},
		"User-Agent":         []string{`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36`},
		"Content-Type":       []string{`application/json`},
		"Accept":             []string{`application/json, text/javascript, */*; q=0.01`},
		"X-Requested-With":   []string{`XMLHttpRequest`},
		"Sec-Ch-Ua-Platform": []string{`"macOS"`},
		"Origin":             []string{`https://www.guilded.gg`},
		"Sec-Fetch-Site":     []string{`same-origin`},
		"Sec-Fetch-Mode":     []string{`cors`},
		"Sec-Fetch-Dest":     []string{`empty`},
		"Referer":            []string{"https://www.guilded.gg/"},
		"guilded-client-id":  []string{g.ClienID},
		"Accept-Language":    []string{`fr-FR,fr;q=0.9`},
	}
	res, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: wrong status code %d, %s", res.StatusCode, string(body))
	}

	return nil
}
func (g *GuildeadClient) ConsumeInvite(invite string) error {

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.guilded.gg/api/invites/%s", invite), nil)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		//"Content-Length":     []string{fmt.Sprint(len(b))},
		"Sec-Ch-Ua":          []string{`" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
		"Sec-Ch-Ua-Mobile":   []string{`?0`},
		"User-Agent":         []string{`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36`},
		"Content-Type":       []string{`application/json`},
		"Accept":             []string{`application/json, text/javascript, */*; q=0.01`},
		"X-Requested-With":   []string{`XMLHttpRequest`},
		"Sec-Ch-Ua-Platform": []string{`"macOS"`},
		"Origin":             []string{`https://www.guilded.gg`},
		"Sec-Fetch-Site":     []string{`same-origin`},
		"Sec-Fetch-Mode":     []string{`cors`},
		"Sec-Fetch-Dest":     []string{`empty`},
		"Referer":            []string{"https://www.guilded.gg/i/" + invite},
		"guilded-client-id":  []string{g.ClienID},
		"Accept-Language":    []string{`fr-FR,fr;q=0.9`},
	}
	res, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: wrong status code %d, %s", res.StatusCode, string(body))
	}

	return nil
}

//Populate the cookie jar with some cookies, idk if this is needed, remove if you need perf
func (g *GuildeadClient) getCookie() error {
	req, err := http.NewRequest("GET", "https://guilded.gg", nil)
	if err != nil {
		return fmt.Errorf("error getting cookies: %s", err)
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error getting cookies: %s", err)
	}
	defer res.Body.Close()
	urlObj, _ := url.Parse("https://guilded.gg")
	g.Client.Jar.SetCookies(urlObj, res.Cookies())
	return nil

}
func (g *GuildeadClient) GetToken() string {
	URL, _ := url.Parse("https://guilded.gg")
	cookies := g.Client.Jar.Cookies(URL)

	var hmac string
	for _, cookie := range cookies {
		if cookie.Name == "hmac_signed_session" {
			hmac = cookie.Value
		}
	}
	return hmac
}
func randHexString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
