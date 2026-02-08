package command

import (
	"context"
	"fmt"
	"openassemblybind/internal/service"
)

func ListCommand(
	key string,
	method string,
) error {
	// TODO: add ctx config
	ctx := context.Background()

	data, err := service.FetchServiceSummaries(
		ctx,
		key,
	)
	if err != nil {
		return err
	}

	switch method {
	case "detailed":
		return nil
	case "simple":
		return listSimple(ctx, data)
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
}

func listSimple(
	ctx context.Context,
	data []service.ServiceSummary,
) error {
	for _, d := range data {
		spec, err := service.FetchServiceSpec(ctx, d.DocURL)
		if err != nil {
			return err
		}
		fmt.Printf("%-20s : %s\n", spec.ResponseKey, spec.Title)
	}
	return nil
}
