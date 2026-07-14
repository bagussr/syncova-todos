package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	domain "syncova-todo/domain/base"
	models "syncova-todo/domain/models"
	"time"
)

type AuthClient struct {
	baseUrl    string
	httpClient *http.Client
	apiKey     string
}

func NewAuthClient(baseURL string, timeout time.Duration) *AuthClient {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &AuthClient{
		baseUrl: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func NewAuthClientWithAPIKey(baseURL string, timeout time.Duration, apiKey string) *AuthClient {
	client := NewAuthClient(baseURL, timeout)
	client.apiKey = apiKey
	return client
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (bool, error) {

	if c.apiKey == "" {
		return false, domain.ErrInternal
	}

	req, err := http.NewRequest(http.MethodGet, c.baseUrl+"/auth/verify_token", nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return false, domain.ErrUnauthorized
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		return false, domain.ErrUnauthorized
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Auth service returned non-OK status:", res.StatusCode)
		return false, domain.ErrUnauthorized
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false, domain.ErrUnauthorized
	}

	if len(body) == 0 {
		fmt.Println("Empty response body from auth service")
		return false, domain.ErrUnauthorized
	}

	var authResponse models.AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		fmt.Println("Error decoding auth response:", err)
		return false, domain.ErrBadRequest
	}

	return true, nil
}
