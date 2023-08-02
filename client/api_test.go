package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/appcenter/model"
)

func TestCreateAPIWithClientParams(t *testing.T) {
	const authToken = "MYTOKEN"

	requestCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount < 3 { // Simulate error
			requestCount++
			w.WriteHeader(502)
			return
		}

		tokenValues := r.Header.Values("x-api-token")
		if len(tokenValues) != 1 {
			w.WriteHeader(401)
			return
		}
		if tokenValues[0] != authToken {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(201)
	}))
	defer ts.Close()

	api := CreateAPIWithClientParams(authToken)
	api.baseURL = ts.URL

	err := api.AddTesterToRelease("", 1, model.ReleaseOptions{})

	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}
}
