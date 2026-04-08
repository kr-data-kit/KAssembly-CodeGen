package endpoint

import (
	"context"
	"fmt"
	"strings"
)

// For Binding
type Endpoint struct {
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

type GenerateResult struct {
	Endpoint *Endpoint
	Error    error
}

func shouldIncludeEndpoint(serviceID string, includeList, excludeList []string) bool {
	if len(excludeList) > 0 {
		for _, excluded := range excludeList {
			if serviceID == excluded {
				return false
			}
		}
	}

	if len(includeList) == 0 {
		return true
	}

	for _, included := range includeList {
		if serviceID == included {
			return true
		}
	}

	return false
}

func isExcluded(serviceID string, excludeList []string) bool {
	for _, excluded := range excludeList {
		if serviceID == excluded {
			return true
		}
	}
	return false
}

func hasRequestedMatch(serviceID, responseKey string, includeList, excludeList []string) bool {
	return shouldIncludeEndpoint(serviceID, includeList, excludeList) || shouldIncludeEndpoint(responseKey, includeList, excludeList)
}

func parseServiceTypes(raw string) (apiProvides bool, apiInfSeq string, dataProvides bool, dataInfSeq string) {
	for _, serviceType := range strings.Split(raw, ",") {
		parts := strings.Split(serviceType, "-")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "A":
			apiProvides = true
			apiInfSeq = parts[1]
		case "S":
			dataProvides = true
			dataInfSeq = parts[1]
		}
	}

	return apiProvides, apiInfSeq, dataProvides, dataInfSeq
}

func shouldEmitExtra(baseIncluded bool, extra *Endpoint, includeList, excludeList []string) bool {
	if extra == nil {
		return false
	}
	if isExcluded(extra.ID, excludeList) || isExcluded(extra.ResponseKey, excludeList) {
		return false
	}
	if baseIncluded {
		return true
	}
	return hasRequestedMatch(extra.ID, extra.ResponseKey, includeList, excludeList)
}

func newEndpointFromSummary(item Summary, spec *ServiceSpec) *Endpoint {
	return &Endpoint{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		URL:         fmt.Sprintf("https://open.assembly.go.kr/portal/data/service/selectAPIServicePage.do/%s", item.ID),

		StructName: getStructName(spec.ResponseKey),
		AlterStructNames: []string{
			item.ID,
		},

		Endpoint:    spec.Endpoint,
		ResponseKey: spec.ResponseKey,
		Params:      spec.Variables,
		Cols:        spec.Columns,

		ProvidesAPI:  true,
		ProvidesData: false,
	}
}

func applyQueryMetadata(target *Endpoint, ccl string) {
	target.CCL = ccl
	target.CommercialUseAllowed = getCommercialUseAllowed(ccl)
	target.AttributionRequired = getAttributionRequired(ccl)
}

func GenerateEndpoints(ctx context.Context, includeList, excludeList []string) (chan *GenerateResult, error) {
	returnChan := make(chan *GenerateResult)

	summaries, err := FetchSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch summaries: %w", err)
	}
	go func() {
		defer close(returnChan)

		for _, item := range summaries {
			if ctx.Err() != nil {
				returnChan <- &GenerateResult{
					Endpoint: nil,
					Error:    fmt.Errorf("service generation cancelled: %w", ctx.Err()),
				}
				return
			}

			apiProvides, apiInfSeq, dataProvides, dataInfSeq := parseServiceTypes(item.ServiceTypesRaw)
			if !apiProvides {
				continue
			}

			spec, err := FetchServiceSpec(ctx, item.ID, apiInfSeq)
			if err != nil {
				returnChan <- &GenerateResult{
					Endpoint: nil,
					Error:    fmt.Errorf("failed to fetch service spec for %s: %w", item.ID, err),
				}
				continue
			}

			enp := newEndpointFromSummary(item, spec)
			enp.ProvidesAPI = apiProvides
			enp.APIInfSeq = apiInfSeq
			enp.ProvidesData = dataProvides
			enp.DataInfSeq = dataInfSeq

			extra, err := CheckExtra(ctx, enp)

			if err != nil {
				returnChan <- &GenerateResult{
					Endpoint: nil,
					Error:    fmt.Errorf("failed to fetch extra service spec for %s: %w", item.ID, err),
				}
				extra = nil
			}

			baseIncluded := hasRequestedMatch(item.ID, spec.ResponseKey, includeList, excludeList)
			extraIncluded := shouldEmitExtra(baseIncluded, extra, includeList, excludeList)

			if !baseIncluded && !extraIncluded {
				continue
			}

			query, err := FetchQueryData(ctx, item.ID)
			if err != nil {
				returnChan <- &GenerateResult{
					Endpoint: nil,
					Error:    fmt.Errorf("failed to fetch query data for %s: %w", item.ID, err),
				}
				continue
			}
			applyQueryMetadata(enp, query.CCL)

			if baseIncluded {
				returnChan <- &GenerateResult{
					Endpoint: enp,
					Error:    nil,
				}
			}

			if extraIncluded {
				applyQueryMetadata(extra, query.CCL)
				returnChan <- &GenerateResult{
					Endpoint: extra,
					Error:    nil,
				}
			}
		}

	}()

	return returnChan, nil
}
