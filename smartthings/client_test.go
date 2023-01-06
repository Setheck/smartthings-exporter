package smartthings

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	testToken := "some-token"

	tests := []struct {
		name       string
		token      string
		httpClient *http.Client
	}{
		{"new client, nil http client", testToken, nil},
		{"new client, with http client", testToken, &http.Client{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			client := NewClient(test.token, test.httpClient)
			assert.Equal(t, testToken, client.token)
			if test.httpClient == nil {
				assert.Equal(t, client.httpClient, http.DefaultClient)
			} else {
				assert.Equal(t, client.httpClient, test.httpClient)
			}
		})
	}
}
