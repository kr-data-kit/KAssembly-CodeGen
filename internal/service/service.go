package service

import (
	"context"
	"fmt"
)

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

type CreateServiceOptions struct {
	CreateServiceEn bool
}

func CreateService(
	ctx context.Context,
	summary ServiceSummary,
	options ...CreateServiceOptions,
) (*Service, error) {
	spec, err := FetchServiceSpec(ctx, summary.DocURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service spec: %w", err)
	}

	service := &Service{
		ID:          summary.ID,
		Title:       summary.Title,
		Description: summary.Description,
		URL:         summary.ServiceURL,

		StructName: getStructName(spec.ResponseKey),
		AlterStructNames: []string{
			summary.ID,

			// TODO: add more alternatives if needed
		},

		Endpoint:    spec.Endpoint,
		ResponseKey: spec.ResponseKey,
		Params:      spec.Variables,
		Cols:        spec.Columns,

		CCL:                  summary.License,
		CommercialUseAllowed: getCommercialUseAllowed(summary.License),
		AttributionRequired:  getAttributionRequired(summary.License),
	}

	if len(options) > 0 {
		option := options[0]
		if option.CreateServiceEn {
			serviceNameEn := GetServiceNameEn(service.ID)
			if serviceNameEn != "" {
				service.AlterStructNames = append(service.AlterStructNames, serviceNameEn)
			}
		}
	}

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
