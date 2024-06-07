package branch

import (
	"errors"

	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

func GetGeneratorBranchReport(format string) (generators.IReport, error) {

	switch format {
	case generators.PDF:
		zap.Log.Info("Excel reporting generation")
		return NewBranchReport(), nil

	case generators.EXCEL:
		zap.Log.Info("Excel reporting generation")
		return nil, nil

	default:
		zap.Log.Info("File format not implemented:", format)
		return nil, errors.New("File format not implemented: " + format)
	}
}
