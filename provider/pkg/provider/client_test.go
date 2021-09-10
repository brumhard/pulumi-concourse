package provider

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const concourseBaseURL = "http://localhost:8080"

var newBasicAuthClientTests = []struct {
	name      string
	username  string
	password  string
	expectErr bool
}{
	{name: "correct creds", username: "test", password: "test", expectErr: false},
	{name: "wrong password", username: "test", password: "lel", expectErr: true},
	{name: "no username", username: "", password: "test", expectErr: true},
}

func Test_newPasswordGrantHTTPClient(t *testing.T) {
	for _, test := range newBasicAuthClientTests {
		t.Run(test.name, func(t *testing.T) {
			client, err := newPasswordGrantHTTPClient(concourseBaseURL, test.username, test.password)
			if test.expectErr {
				require.Error(t, err)
				return
			}

			// test correct creds can query authenticated endpoint
			req, err := http.NewRequest(http.MethodGet, concourseBaseURL+"/api/v1/workers", nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
