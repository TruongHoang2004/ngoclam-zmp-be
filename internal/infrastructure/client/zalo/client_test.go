package zalo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetUserInfo(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("access_token") != "test_access_token" {
			t.Errorf("expected access_token header to be test_access_token, got %s", r.Header.Get("access_token"))
		}
		if r.Header.Get("code") != "test_token" {
			t.Errorf("expected code header to be test_token, got %s", r.Header.Get("code"))
		}
		if r.Header.Get("secret_key") != "test_secret_key" {
			t.Errorf("expected secret_key header to be test_secret_key, got %s", r.Header.Get("secret_key"))
		}

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "123456789",
			"name": "Test User",
			"picture": {
				"data": {
					"url": "https://example.com/avatar.jpg"
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(nil)
	client.baseURL = server.URL // Inject mock server URL

	userInfo, err := client.GetUserInfo("test_access_token", "test_token", "test_secret_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if userInfo.ID != "123456789" {
		t.Errorf("expected ID to be 123456789, got %s", userInfo.ID)
	}
	if userInfo.Name != "Test User" {
		t.Errorf("expected Name to be Test User, got %s", userInfo.Name)
	}
	if userInfo.Picture.Data.URL != "https://example.com/avatar.jpg" {
		t.Errorf("expected Picture URL to be https://example.com/avatar.jpg, got %s", userInfo.Picture.Data.URL)
	}
}
