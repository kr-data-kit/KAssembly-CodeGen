package service

import (
	"context"
	"fmt"
)

// For Binding
type Service struct {
	ID          string
	Title       string
	Description string
	URL         string

	StructName       string
	AlterStructNames []string

	Endpoint    string // api endpoint
	ResponseKey string // json response key

	Params []Variable
	Cols   []Column

	CCL                  string
	CommercialUseAllowed bool
	AttributionRequired  bool
}

func CreateService(
	ctx context.Context,
	item Summary,
) (*Service, error) {
	service := &Service{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
	}

	ccl, infSeq, err := fetchCCLAndInfSeq(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	service.URL = fmt.Sprintf(
		"https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s",
		item.ID,
	)

	spec, err := FetchServiceSpec(ctx, item.ID, infSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service spec: %w", err)
	}
	service.StructName = getStructName(spec.ResponseKey)
	service.AlterStructNames = []string{
		item.ID,
	}
	service.Endpoint = spec.Endpoint
	service.ResponseKey = spec.ResponseKey
	service.Params = spec.Variables
	service.Cols = spec.Columns

	service.CCL = ccl
	service.CommercialUseAllowed = getCommercialUseAllowed(ccl)
	service.AttributionRequired = getAttributionRequired(ccl)

	return service, nil
}

func getStructName(ResponseKey string) string {
	if checkRandomName(ResponseKey) {
		return fmt.Sprintf("%s%s", string(ResponseKey[0]-32), ResponseKey[1:])
	}
	// TODO : 나중에 고도화 필요
	return ResponseKey
}

func getCommercialUseAllowed(License string) bool {
	return License != "출처표시 + 상업적 이용금지"
}

func getAttributionRequired(License string) bool {
	return License == "출처표시 + 상업적 이용금지" || License == "출처표시"
}

func checkRandomName(name string) bool {
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			return false
		}
	}
	return true
}
