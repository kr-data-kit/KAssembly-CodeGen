package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	queryDataURL = "https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do"
)

type QueryData struct {
	CCL string
}

func FetchQueryData(ctx context.Context, id string) (*QueryData, error) {
	if id == "" {
		return nil, fmt.Errorf("empty id")
	}
	for _, r := range id {
		if r <= 32 || r == '/' || r == '?' || r == '#' || r == '\\' || r == '%' || r == '&' || r == ':' {
			return nil, fmt.Errorf("unsafe id value: %q", id)
		}
	}
	url := fmt.Sprintf("%s/%s", queryDataURL, id)
	body, err := fetchHTTP(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	queryData := &QueryData{}

	queryData.CCL, _ = doc.
		Find("#metaInfo table tbody img").
		Attr("alt")

	return queryData, nil
}
