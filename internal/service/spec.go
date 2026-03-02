package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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

func FetchServiceSpec(ctx context.Context, infId string, infSeq string) (*ServiceSpec, error) {
	url := fmt.Sprintf("%s?infId=%s&infSeq=%s", specUrl, url.QueryEscape(infId), url.QueryEscape(infSeq))
	body, err := fetchHTTP(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	var resp ServiceSpecResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &resp.Data, nil
}
