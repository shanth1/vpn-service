package qr

import (
	"context"
	"fmt"

	"github.com/skip2/go-qrcode"
)

type infra struct{}

func NewQRInfra() *infra {
	return &infra{}
}

func (*infra) Generate(ctx context.Context, fileName string, data []byte) error {
	dataString := string(data)
	outputFilename := fileName + ".png"

	err := qrcode.WriteFile(dataString, qrcode.Medium, 256, outputFilename)
	if err != nil {
		return fmt.Errorf("qrcode.WriteFile: %w", err)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
