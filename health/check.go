package health

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Check calls the given service endpoint with a given context and timeout.
// An error will be returned if the connection fails, or the response status
// is not 200 (i.e. StatusOK). A successful check will return only the check message reply.
func Check(ctx context.Context, servicePath string, timeout time.Duration) ([]byte, error) {
	req, err := url.Parse("http://" + servicePath)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, req.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if resp == nil || err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s (%s)", string(body), http.StatusText(resp.StatusCode))
	}

	return body, nil
}

// CheckStatus runs a Check on the given service and returns zero for a healthy service, and one otherwise.
//
//	@param {string} servicePat: service address and path to check e.g. 8080/soh
func CheckStatus(servicePath string, timeout time.Duration) int {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := Check(ctx, servicePath, timeout); err != nil {
		return 1
	}

	return 0
}
