package qr

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

type infra struct{}

func NewInfra() *infra {
	return &infra{}
}

func (*infra) Generate(ctx context.Context, outputDir, fileName string, data []byte) error {
	dataString := string(data)
	filePath := filepath.Join(outputDir, fileName+".png")

	err := qrcode.WriteFile(dataString, qrcode.Medium, 256, filePath)
	if err != nil {
		return fmt.Errorf("qrcode.WriteFile: %w", err)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
