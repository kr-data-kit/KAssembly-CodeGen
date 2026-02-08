package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type OpenSrvApiResponse struct {
	OPENSRVAPI []OpenSrvApiItem `json:"OPENSRVAPI"`
	RESULT     *Result          `json:"RESULT,omitempty"`
}

type OpenSrvApiItem struct {
	Head []HeadItem       `json:"head,omitempty"`
	Row  []ServiceSummary `json:"row,omitempty"`
}

type HeadItem struct {
	ListTotalCount *int    `json:"list_total_count,omitempty"`
	Result         *Result `json:"RESULT,omitempty"`
}

type Result struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

type ServiceSummary struct {
	ID           string `json:"INF_ID"`    // 공공데이터ID
	Title        string `json:"INF_NM"`    // 공공데이터명
	Description  string `json:"INF_EXP"`   // 공공데이터설명
	Category     string `json:"CATE_NM"`   // 분류체계
	OpenDate     string `json:"OPEN_DTTM"` // 공개일자 (YYYY-MM-DD)
	Organization string `json:"ORG_NM"`    // 제공기관
	LastModified string `json:"LOAD_DTTM"` // 최종수정일자 (YYYY-MM-DD)
	Source       string `json:"SRC_EXP"`   // 원본시스템
	DocURL       string `json:"DDC_URL"`   // 명세서URL
	ServiceURL   string `json:"SRV_URL"`   // 서비스URL
	License      string `json:"CCL_NM"`    // 이용허락조건
	UpdateCycle  string `json:"LOAD_NM"`   // 공개주기
	UpdatePeriod string `json:"LOAD_CONT"` // 공개시기
}

func buildSummaryURL(apiKey string, pageIndex int, pageSize int) string {
	params := url.Values{}
	params.Add("KEY", apiKey)
	params.Add("Type", "json")
	params.Add("pIndex", strconv.Itoa(pageIndex))
	params.Add("pSize", strconv.Itoa(pageSize))
	return opensrvapiUrl + "?" + params.Encode()
}

const (
	opensrvapiUrl   = "https://open.assembly.go.kr/portal/openapi/OPENSRVAPI"
	defaultPageSize = 500
)

func FetchServiceSummaries(ctx context.Context, apiKey string) ([]ServiceSummary, error) {
	allSummaries := []ServiceSummary{}

	for pageIndex := 1; ; pageIndex++ {
		url := buildSummaryURL(apiKey, pageIndex, defaultPageSize)
		body, err := fetchHTTP(ctx, url)
		if err != nil {
			return nil, err
		}
		var resp OpenSrvApiResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		if resp.RESULT != nil {
			result := *resp.RESULT
			if result.Code == "INFO-200" {
				break
			}
			return nil, fmt.Errorf("API error: %s - %s", result.Code, result.Message)
		}
		if len(resp.OPENSRVAPI) < 2 || len(resp.OPENSRVAPI[0].Head) < 2 {
			return nil, fmt.Errorf("invalid response structure")
		}
		summaries := resp.OPENSRVAPI[1].Row
		if len(summaries) == 0 {
			break
		}
		allSummaries = append(allSummaries, summaries...)
	}

	return allSummaries, nil
}
