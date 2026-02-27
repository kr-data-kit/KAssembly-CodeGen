package command

import (
	"context"
	"fmt"
	"openassemblybinder/internal/service"
)

func ListCommand(
	method string,
) error {
	// TODO: add ctx config
	ctx := context.Background()

	data, err := service.FetchSummary(ctx)
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
	data []service.Summary,
) error {
	// TODO : add list
	return nil
}
