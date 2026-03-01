package service

import (
	"context"
	"fmt"
	"strings"
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

	ProvidesAPI  bool
	APIInfSeq    string
	ProvidesData bool
	DataInfSeq   string

	CCL                  string
	CommercialUseAllowed bool
	AttributionRequired  bool
}

func getStructName(ResponseKey string) string {
	if checkRandomName(ResponseKey) {
		return fmt.Sprintf("%s%s", string(ResponseKey[0]-32), ResponseKey[1:])
	}
	// TODO : 나중에 고도화 필요
	return ResponseKey
}

func getCommercialUseAllowed(License string) bool {
	return strings.Contains(License, "상업용 금지")
}

func getAttributionRequired(License string) bool {
	return strings.Contains(License, "출처표시")
}

func checkRandomName(name string) bool {
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			return false
		}
	}
	return true
}

type ServiceResult struct {
	Service *Service
	Error   error
}

func GenerateServices(ctx context.Context) (chan *ServiceResult, error) {
	returnChan := make(chan *ServiceResult)

	summaries, err := FetchSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch summaries: %w", err)
	}
	go func() {
		defer close(returnChan)

		for _, item := range summaries {
			// Check if context is cancelled
			if ctx.Err() != nil {
				returnChan <- &ServiceResult{
					Service: nil,
					Error:   fmt.Errorf("service generation cancelled: %w", ctx.Err()),
				}
				return
			}

			service := &Service{
				ID:          item.ID,
				Title:       item.Title,
				Description: item.Description,

				AlterStructNames: []string{
					item.ID,
				},

				ProvidesAPI:  false,
				ProvidesData: false,
			}

			provides := strings.Split(item.ServiceTypesRaw, ",")
			for _, item := range provides {
				parts := strings.Split(item, "-")
				if len(parts) != 2 {
					continue
				}
				switch parts[0] {
				case "A":
					service.ProvidesAPI = true
					service.APIInfSeq = parts[1]
				case "S":
					service.ProvidesData = true
					service.DataInfSeq = parts[1]
				}
			}

			if !service.ProvidesAPI {
				continue
			}

			query, err := FetchQueryData(ctx, item.ID)
			if err != nil {
				returnChan <- &ServiceResult{
					Service: nil,
					Error:   fmt.Errorf("failed to fetch query data for %s: %w", item.ID, err),
				}
				continue
			}

			spec, err := FetchServiceSpec(ctx, item.ID, service.APIInfSeq)
			if err != nil {
				returnChan <- &ServiceResult{
					Service: nil,
					Error:   fmt.Errorf("failed to fetch service spec for %s: %w", item.ID, err),
				}
				continue
			}

			service.URL = fmt.Sprintf("https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s", item.ID)

			service.StructName = getStructName(spec.ResponseKey)

			service.Endpoint = spec.Endpoint
			service.ResponseKey = spec.ResponseKey

			service.Params = spec.Variables
			service.Cols = spec.Columns

			service.CCL = query.CCL
			service.CommercialUseAllowed = getCommercialUseAllowed(query.CCL)
			service.AttributionRequired = getAttributionRequired(query.CCL)

			returnChan <- &ServiceResult{
				Service: service,
				Error:   nil,
			}

			extra, err := CheckExtra(ctx, service)
			if extra != nil {
				returnChan <- &ServiceResult{
					Service: extra,
					Error:   err,
				}
			}
		}

	}()

	return returnChan, nil
}
