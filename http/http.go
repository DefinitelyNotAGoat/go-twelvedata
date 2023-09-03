package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func Get(u *url.URL, client *http.Client) ([]byte, error) {
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code '%d'", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return body, nil
}
