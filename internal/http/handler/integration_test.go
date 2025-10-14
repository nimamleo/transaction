//go:build integration

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"transaction/internal/account/domain"
	userDomain "transaction/internal/user/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8081"
	apiKey  = "b13cb46cf2b29449d5234f9ac5723189b31b957d9ff5112708963b816ef7e497"
)

type TestClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewTestClient() *TestClient {
	return &TestClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *TestClient) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	return c.client.Do(req)
}

func TestIntegration_HappyPath(t *testing.T) {
	client := NewTestClient()

	user := createUser(t, client)
	fromAccount := createAccount(t, client, user.ID, "USD")
	toAccount := createAccount(t, client, user.ID, "USD")

	depositToAccount(t, client, fromAccount.ID, 10000, "initial-deposit")
	depositToAccount(t, client, toAccount.ID, 5000, "initial-deposit-2")

	transferBetweenAccounts(t, client, fromAccount.ID, toAccount.ID, 2000, "transfer-1")

	checkBalance(t, client, fromAccount.ID, 8000)
	checkBalance(t, client, toAccount.ID, 7000)

	checkTransactionHistory(t, client, fromAccount.ID, 2)
	checkTransactionHistory(t, client, toAccount.ID, 2)
}

func TestIntegration_Idempotency(t *testing.T) {
	client := NewTestClient()

	user := createUser(t, client)
	account := createAccount(t, client, user.ID, "USD")

	depositToAccount(t, client, account.ID, 1000, "idempotent-deposit")
	depositToAccount(t, client, account.ID, 1000, "idempotent-deposit-2")

	checkBalance(t, client, account.ID, 2000)
}

func TestIntegration_InsufficientFunds(t *testing.T) {
	client := NewTestClient()

	user := createUser(t, client)
	fromAccount := createAccount(t, client, user.ID, "USD")
	toAccount := createAccount(t, client, user.ID, "USD")

	depositToAccount(t, client, fromAccount.ID, 1000, "small-deposit")

	resp, err := client.makeRequest("POST", "/api/v1/transfers", map[string]interface{}{
		"from_account_id": fromAccount.ID,
		"to_account_id":   toAccount.ID,
		"amount":          2000,
		"reference":       "insufficient-funds-transfer",
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func createUser(t *testing.T, client *TestClient) *userDomain.User {
	userReq := map[string]string{
		"name":  "Test User",
		"email": fmt.Sprintf("test-%d@example.com", time.Now().UnixNano()),
	}

	resp, err := client.makeRequest("POST", "/api/v1/users", userReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	user := &userDomain.User{
		ID:    response["data"].(map[string]interface{})["id"].(string),
		Name:  userReq["name"],
		Email: userReq["email"],
	}

	resp.Body.Close()
	return user
}

func createAccount(t *testing.T, client *TestClient, userID, currency string) *domain.Account {
	accountReq := map[string]string{
		"user_id":  userID,
		"currency": currency,
	}

	resp, err := client.makeRequest("POST", "/api/v1/accounts", accountReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	account := &domain.Account{
		ID:       response["data"].(map[string]interface{})["id"].(string),
		UserID:   userID,
		Currency: domain.Currency(currency),
		Balance:  0,
	}

	resp.Body.Close()
	return account
}

func depositToAccount(t *testing.T, client *TestClient, accountID string, amount int64, reference string) {
	depositReq := map[string]interface{}{
		"amount":    amount,
		"reference": reference,
	}

	resp, err := client.makeRequest("POST", fmt.Sprintf("/api/v1/accounts/%s/deposit", accountID), depositReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	resp.Body.Close()
}

func transferBetweenAccounts(t *testing.T, client *TestClient, fromAccountID, toAccountID string, amount int64, reference string) {
	transferReq := map[string]interface{}{
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
		"amount":          amount,
		"reference":       reference,
	}

	resp, err := client.makeRequest("POST", "/api/v1/transfers", transferReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	resp.Body.Close()
}

func checkBalance(t *testing.T, client *TestClient, accountID string, expectedBalance int64) {
	resp, err := client.makeRequest("GET", fmt.Sprintf("/api/v1/accounts/%s/balance", accountID), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	balance := int64(response["data"].(map[string]interface{})["balance"].(float64))
	assert.Equal(t, expectedBalance, balance)

	resp.Body.Close()
}

func checkTransactionHistory(t *testing.T, client *TestClient, accountID string, expectedCount int) {
	resp, err := client.makeRequest("GET", fmt.Sprintf("/api/v1/accounts/%s/transactions?limit=10", accountID), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	transactions := response["data"].(map[string]interface{})["transactions"].([]interface{})
	assert.Len(t, transactions, expectedCount)

	resp.Body.Close()
}

func TestIntegration_HealthCheck(t *testing.T) {
	client := NewTestClient()

	resp, err := client.makeRequest("GET", "/health", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])

	resp.Body.Close()
}
