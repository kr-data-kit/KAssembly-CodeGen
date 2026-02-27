package service

import (
	"context"
	"fmt"
	"regexp"
)

/*
장황한 설명: 버전이 다른 API를 자동으로 감지해서 처리하는 로직을 구현해보렸고 했으나, 하드코딩하는 방식으로 넘어갔습니다.
이유는 다음과 같습니다.
1. API 버전 링크가 불완전함 : 해당 버전의 ID를 실제로 조회하면, API가 없거나 아예 접근이 불가능한 경우가 많았습니다.
2. 접근이 가능한 다른 버전의 API 수가 너무 적음 : 2026-02-26 기준으로, 접근이 가능한 API는 의안 접수목록과 의안정보 통합 API 뿐이었습니다.

다음은 다른 버전의 API가 조회되는 경우입니다 (2026-02-26 기준)

# skip
국회도서관 연간보고서 (ID: OO9ZRP000830VN19457) : API가 존재하지 않음
최신외국입법정보 (ID: OK5H9O001029KC12614) : API가 존재하지 않음
국회도서관 행정정보 공표목록 (ID: OYX6MA001159BO19467) : 접근이 불가능함

# extra
의안 접수목록 (ID: OOWY4R001216HX11458) : 접근이 가능하며, 모든 의안을 데이터로 조회할 수 있어 중요함
의안정보 통합 API (ID: OOWY4R001216HX11440) : 접근이 가능하며, 모든 의안을 데이터로 조회할 수 있어 중요함
*/

const (
	ALLBILL   string = "OOWY4R001216HX11440"
	ALLBILL2  string = "OOWY4R001216HX11536"
	BILLRCPV  string = "OOWY4R001216HX11458"
	BILLRCPV2 string = "OOWY4R001216HX11537"
)

// CheckAdditionalIDs checks if there are additional IDs for the given ID by fetching the query data page and extracting all IDs from the links on the page.
// check `extra.go` for more details on why this function exists and how it is used.
func CheckAdditionalIDs(ctx context.Context, id string) ([]string, error) {
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
	re := regexp.MustCompile(`/portal/data/service/selectAPIServicePage\.do/([^"'>\s\?]+)`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	var additionalIDs []string
	for _, match := range matches {
		additionalIDs = append(additionalIDs, match[1])
	}
	return additionalIDs, nil
}

func CheckExtra(ctx context.Context, service *Service) (*Service, error) {
	switch service.ID {
	case ALLBILL2:
		spec, err := FetchServiceSpec(ctx, ALLBILL, "2")
		if err != nil {
			return nil, err
		}

		return &Service{
			ID:          ALLBILL,
			Title:       service.Title + " (Version 1)",
			Description: service.Description + " - Updated Version",
			URL:         fmt.Sprintf("https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s", ALLBILL),
			StructName:  getStructName(spec.ResponseKey),

			AlterStructNames: []string{
				ALLBILL,
			},

			Endpoint:    spec.Endpoint,
			ResponseKey: spec.ResponseKey,

			Params: spec.Variables,
			Cols:   spec.Columns,

			CCL:                  service.CCL,
			CommercialUseAllowed: service.CommercialUseAllowed,
			AttributionRequired:  service.AttributionRequired,
		}, nil
	case BILLRCPV2:

		spec, err := FetchServiceSpec(ctx, BILLRCPV, "2")
		if err != nil {
			return nil, err
		}

		return &Service{
			ID:          BILLRCPV,
			Title:       service.Title + " (Version 1)",
			Description: service.Description + " - Updated Version",
			URL:         fmt.Sprintf("https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s", BILLRCPV),
			StructName:  getStructName(spec.ResponseKey),

			AlterStructNames: []string{
				BILLRCPV,
			},

			Endpoint:    spec.Endpoint,
			ResponseKey: spec.ResponseKey,

			Params: spec.Variables,
			Cols:   spec.Columns,

			CCL:                  service.CCL,
			CommercialUseAllowed: service.CommercialUseAllowed,
			AttributionRequired:  service.AttributionRequired,
		}, nil
	default:
		return nil, nil
	}
}
