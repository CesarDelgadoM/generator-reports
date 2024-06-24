package branch

import (
	"errors"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
)

func GetGeneratorBranchReport(format string, config *config.Config) (generators.IReport, error) {

	switch format {
	case generators.PDF:
		return NewBranchReportPdf(config.Branch.Pdf), nil

	case generators.EXCEL:
		return nil, nil

	default:
		return nil, errors.New("File format not implemented: " + format)
	}
}
