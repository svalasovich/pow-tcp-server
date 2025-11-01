package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/svalasovich/golang-template/internal/config"
)

func TestServer_ListenAndServe(t *testing.T) {
	// given
	// Configure server
	cfg := config.HTTPServer{
		Port: 9000,
	}
	testSubject := NewServer(cfg)

	// Configure Handler
	pattern := "/hello"
	response := []byte(gofakeit.Word())
	testSubject.Router().Get(pattern, func(writer http.ResponseWriter, _ *http.Request) {
		_, err := writer.Write(response)
		require.NoError(t, err)
	})

	// when
	err := testSubject.ListenAndServe()

	// then
	// Check response
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d%s", cfg.Port, pattern), nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, response, body)

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, testSubject.Shutdown(ctx))

	// Check shutdown
	require.NoError(t, err)
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d", cfg.Port), nil)
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	require.Error(t, err)
}
