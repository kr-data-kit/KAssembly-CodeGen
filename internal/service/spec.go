package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
)

type ServiceSpecResponse struct {
	Data ServiceSpec `json:"data"`
}

type ServiceSpec struct {
	Title        string      `json:"infNm"`
	ServiceCode  string      `json:"srvCd"`
	Endpoint     string      `json:"apiEp"`
	ResponseKey  string      `json:"apiRes"`
	TrafficLimit json.Number `json:"apiTrf"`
	Variables    []Variable  `json:"variables"`
	Columns      []Column    `json:"columns"`
	Urls         []Url       `json:"urls"`
	Filters      []Filter    `json:"filters"`
	Messages     []Message   `json:"messages"`
}

type Variable struct {
	ID          string  `json:"colId"`
	Type        string  `json:"reqType"` // "STRING", "INT", etc.
	RequiredRaw *string `json:"reqNeed"`
	Name        string  `json:"colNm"`
	Description string  `json:"colExp"`
	Example     *string `json:"smpColExp"`
}

type Column struct {
	ID          string  `json:"colId"`
	Name        string  `json:"colNm"`
	Unit        *string `json:"unitNm"`
	Description string  `json:"colExp"`
}

type Url struct {
	Name        string `json:"uriNm"`
	Endpoint    string `json:"apiEp"`
	ResponseKey string `json:"apiRes"`
	Uri         string `json:"uri"`
}

type Filter struct {
	ID          string  `json:"colId"`
	Name        string  `json:"colNm"`
	RequiredRaw *string `json:"reqNeed"`
	FiltCode    *string `json:"filtCode"`
}

type Message struct {
	Tag         string `json:"msgTag"`
	Code        string `json:"msgCd"`
	Description string `json:"msgExp"`
}

const (
	specUrl = "https://open.assembly.go.kr/portal/data/openapi/selectOpenApiMeta.do"
)

func buildSpecURL(DocURL string) string {
	u, err := url.Parse(DocURL)
	if err != nil {
		return ""
	}
	q := u.Query()
	infId := q.Get("infId")
	infSeq := q.Get("infSeq")
	if infId == "" {
		// fallback to last path segment
		seg := filepath.Base(u.Path)
		if seg != "" && seg != "/" {
			infId = seg
		}
	}
	if infSeq == "" {
		infSeq = "1"
	}

	params := url.Values{}
	params.Add("infId", infId)
	params.Add("infSeq", infSeq)
	return specUrl + "?" + params.Encode()
}

func FetchServiceSpec(ctx context.Context, DocURL string) (*ServiceSpec, error) {
	url := buildSpecURL(DocURL)
	if url == "" {
		return nil, fmt.Errorf("invalid DocURL: %s", DocURL)
	}

	body, err := fetchHTTP(ctx, url)
	if err != nil {
		return nil, err
	}
	var resp ServiceSpecResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &resp.Data, nil
}
