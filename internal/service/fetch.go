package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func setCommonHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "open.assembly.go.kr")
	req.Header.Set("User-Agent", "Mozilla/5.0")
}

func fetchHTTP(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setCommonHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code error: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, nil
}

func fetchCCLAndInfSeq(ctx context.Context, id string) (ccl string, infSeq string, error error) {
	if id == "" {
		return "", "", fmt.Errorf("empty id")
	}
	for _, r := range id {
		// disallow control chars, spaces and characters that would change the URL path or query
		if r <= 32 || r == '/' || r == '?' || r == '#' || r == '\\' || r == '%' || r == '&' || r == ':' {
			return "", "", fmt.Errorf("unsafe id value: %q", id)
		}
	}
	url := fmt.Sprintf("https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s", id)
	body, err := fetchHTTP(ctx, url, nil)
	if err != nil {
		return "", "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	ccl, exists := doc.
		Find("#metaInfo table tbody img").
		Attr("alt")
	if !exists {
		return "", "", fmt.Errorf("CCL image not found for ID %s", id)
	}
	infSeq, exists = doc.
		Find("#openapi-search-form input[name='infSeq']").
		Attr("value")
	return ccl, infSeq, nil
}
