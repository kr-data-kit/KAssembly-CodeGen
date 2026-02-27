package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultPageSize = 500
)

type SummaryResponse struct {
	Total   int       `json:"total"`
	Pages   int       `json:"pages"`
	APIList []Summary `json:"data"`
	Count   int       `json:"count"`
	Page    int       `json:"page"`
	Rows    int       `json:"rows"`
}

type Summary struct {
	No               int    `json:"ROW_NUM"`   // 공공데이터 순서
	Tag              string `json:"opentyTag"` // 공개유형
	ID               string `json:"infaId"`    // 공공데이터ID
	Title            string `json:"infaNm"`    // 공공데이터명
	CategoryID       string `json:"cateId"`    // 분류체계ID
	CategoryName     string `json:"cateNm"`    // 분류체계명
	OrganizationCode string `json:"orgCd"`     // 제공기관코드
	OrganizationName string `json:"orgNm"`     // 제공기관명
	Description      string `json:"infaExp"`   // 공공데이터설명
	OpenDate         string `json:"openYmd"`   // 공개일자 (YYYY-MM-DD)
	ServiceTypesRaw  string `json:"openSrv"`   // 서비스유형 (<유형 코드>-<순서>,... 형태의 문자열)
}

const (
	SummaryURL = "https://open.assembly.go.kr/portal/infs/list/selectInfsListPaging.do"
)

// TODO : Summary에 데이터 표를 제공하는지도 포함시킬 수 있음
// 서비스 유형을 bool 필드로 분리해서 ServiceSummary에 포함시키는 것도 고려해볼 수 있음

func FetchSummary(ctx context.Context) ([]Summary, error) {
	check, err := fetchSummaryResponse(ctx, 1, defaultPageSize)
	if err != nil {
		return nil, err
	}

	if check.Total <= defaultPageSize {
		return check.APIList, nil
	}

	var allItems []Summary
	allItems = append(allItems, check.APIList...)

	for page := 2; page <= check.Pages; page++ {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("summary fetch cancelled: %w", ctx.Err())
		}

		pageResult, err := fetchSummaryResponse(ctx, page, defaultPageSize)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, pageResult.APIList...)
	}
	return allItems, nil
}

func fetchSummaryResponse(ctx context.Context, page int, rows int) (*SummaryResponse, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	bodyData := url.Values{}
	bodyData.Set("page", fmt.Sprintf("%d", page))
	bodyData.Set("rows", fmt.Sprintf("%d", rows))

	req, err := http.NewRequestWithContext(ctx, "POST", SummaryURL, strings.NewReader(bodyData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code error: %d", resp.StatusCode)
	}

	var apiResp SummaryResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)

	return &apiResp, nil
}
