package outlier

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

// EnableHttpmock so that PR data is loaded from file rather than fetched from the real GitHub API
func EnableHttpmock() {
	httpmock.Activate()
	httpmock.Reset()
	responder := func(req *http.Request) (*http.Response, error) {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		var data []byte
		switch {
		case strings.Contains(string(reqBody), `"pullRequestsStates":["MERGED"]`):
			data, err = ioutil.ReadFile("outlier/testdata/merged-prs-response.json")
			if err != nil {
				panic(err)
			}
		default:
			panic("unknown graphql query")
		}
		response := httpmock.NewStringResponse(200, string(data))
		response.Header.Set("Content-Type", "application/json")
		return response, nil
	}
	httpmock.RegisterResponder(http.MethodPost, "https://api.github.com/graphql", responder)
}
