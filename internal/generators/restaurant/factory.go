package restaurant

import (
	"errors"

	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

func GetGeneratorRestaurantReport(format string) (generators.IReport, error) {

	switch format {

	case generators.PDF:
		zap.Log.Info("Excel reporting generation")
		return NewRestaurantReport(), nil

	case generators.EXCEL:
		zap.Log.Info("Excel reporting generation")
		return nil, nil

	default:
		zap.Log.Info("File format not implemented:", format)
		return nil, errors.New("File format not implemented: " + format)
	}
}
