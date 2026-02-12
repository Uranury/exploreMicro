package http_pack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Uranury/exploreMicro/service2/internal/models"
	"net/http"
)

type UserClient interface {
	Get(ctx context.Context, userID uint) (*models.User, error)
	Patch(ctx context.Context, userID uint, newBalance float64) (*models.User, error)
}

type userClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPUserClient(baseURL string, client *http.Client) UserClient {
	return &userClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (c *userClient) Get(ctx context.Context, userID uint) (*models.User, error) {
	url := fmt.Sprintf("%s/users?id=%d", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *userClient) Patch(ctx context.Context, userID uint, newBalance float64) (*models.User, error) {
	url := fmt.Sprintf("%s/users?id=%d", c.baseURL, userID)

	body := map[string]float64{
		"balance": newBalance,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to patch user %d, status: %d", userID, resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
