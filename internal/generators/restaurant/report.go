package restaurant

import (
	"errors"

	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

const (
	// Formats
	PDF   = "PDF"
	EXCEL = "EXCEL"
)

// Interface to strategies
type IReport interface {
	GenerateReport(msg *consumer.Message)
}

func GetGeneratorRestaurantReport(format string) (IReport, error) {

	switch format {
	case PDF:
		zap.Log.Info("Excel reporting generation")
		return NewRestaurantReport(), nil

	case EXCEL:
		zap.Log.Info("Excel reporting generation")
		return nil, nil

	default:
		zap.Log.Info("File format not implemented:", format)
		return nil, errors.New("File format not implemented: " + format)
	}
}
