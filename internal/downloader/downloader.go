// Package downloader общается с iprbookshop: авторизуется по cookies,
// получает токен доступа, ключ шифрования, зашифрованный поток и название книги.
package downloader

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"IprbooksDumper/internal/domain"
)

// userAgent — общий заголовок для всех запросов.
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"

const baseURL = "https://www.iprbookshop.ru"

// Client — авторизованный клиент iprbookshop.
type Client struct {
	http http.Client
}

// New создаёт клиент с cookies авторизации из переданной строки.
func New(cookie string) *Client {
	jar, _ := cookiejar.New(nil)
	if u, err := url.Parse(baseURL); err == nil {
		jar.SetCookies(u, parseCookies(cookie))
	}
	return &Client{http: http.Client{Jar: jar}}
}

// Stream получает зашифрованный поток книги и ключ его расшифровки.
func (c *Client) Stream(bookID int) (stream, key []byte, err error) {
	token, err := c.accessToken(bookID)
	if err != nil {
		return nil, nil, err
	}

	key, err = sessionKeyFromToken(token)
	if err != nil {
		return nil, nil, err
	}

	link := baseURL + "/publications/reader/stream?access_token=" + token
	stream, err = c.get(link, "application/octet-stream")
	if err != nil {
		return nil, nil, err
	}

	return stream, key, nil
}

// Title получает название книги через JSON-API (страница сайта — SPA).
// При недоступности названия возвращает строковое представление ID.
func (c *Client) Title(bookID int) string {
	link := baseURL + "/books/" + strconv.Itoa(bookID)

	body, err := c.get(link, "application/json")
	if err != nil {
		return strconv.Itoa(bookID)
	}

	var resp struct {
		Data struct {
			Book struct {
				Title string `json:"title"`
			} `json:"book"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil || resp.Data.Book.Title == "" {
		return strconv.Itoa(bookID)
	}

	return resp.Data.Book.Title
}

// accessToken получает JWT-токен доступа к книге через reader-API.
func (c *Client) accessToken(bookID int) (string, error) {
	link := baseURL + "/publications/reader/" + strconv.Itoa(bookID) + "/access?publication_type=book"

	body, err := c.get(link, "application/json")
	if err != nil {
		return "", err
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return "", domain.ErrNoAccess
	}

	if !resp.Success || resp.Data.AccessToken == "" {
		return "", domain.ErrNoAccess
	}

	return resp.Data.AccessToken, nil
}

// get делает GET-запрос авторизованным клиентом и возвращает тело ответа.
func (c *Client) get(link, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, link, http.NoBody)
	if err != nil {
		return nil, domain.ErrUnavailable
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, domain.ErrUnavailable
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.ErrUnavailable
	}

	return body, nil
}

// sessionKeyFromToken достаёт session_key из полезной нагрузки JWT и декодирует его в ключ AES.
func sessionKeyFromToken(token string) ([]byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, domain.ErrBadToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, domain.ErrBadToken
	}

	var claims struct {
		SessionKey string `json:"session_key"`
	}
	if uErr := json.Unmarshal(payload, &claims); uErr != nil || claims.SessionKey == "" {
		return nil, domain.ErrBadToken
	}

	key, err := decodeBase64URL(claims.SessionKey)
	if err != nil {
		return nil, domain.ErrBadToken
	}

	return key, nil
}

// decodeBase64URL декодирует base64url (с набивкой и без) в байты.
func decodeBase64URL(s string) ([]byte, error) {
	s = strings.NewReplacer("-", "+", "_", "/").Replace(s)
	if pad := len(s) % 4; pad != 0 {
		s += strings.Repeat("=", 4-pad)
	}
	return base64.StdEncoding.DecodeString(s)
}

// parseCookies парсит строку cookies в массив *http.Cookie.
func parseCookies(cookieStr string) []*http.Cookie {
	var cookies []*http.Cookie

	for _, pair := range strings.Split(cookieStr, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			// Клиентские cookies для исходящих запросов: атрибуты Secure/HttpOnly
			// к ним неприменимы, поэтому предупреждение gosec здесь не актуально.
			cookies = append(cookies, &http.Cookie{ //nolint:gosec
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			})
		}
	}

	return cookies
}
