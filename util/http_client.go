package util

import (
	"context"
	"github.com/go-resty/resty/v2"
)

func ApiClientGet(ctx context.Context, url, endpoint string, headers map[string]string) (*resty.Response, error) {
	headers["Content-Type"] = "application/json"

	client := resty.New()
	client.SetHostURL(url)
	resp, err := client.R().SetContext(ctx).
		SetHeaders(
			headers,
		).Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
