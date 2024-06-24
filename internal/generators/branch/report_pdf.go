package branch

import (
	"errors"
	"fmt"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/CesarDelgadoM/generator-reports/internal/utils"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/go-pdf/fpdf"
)

const (
	restaurant_type = "restaurant"
	branch_type     = "branch"

	// File variables
	marginX = 10.0
	marginY = 20.0
	gapY    = 2.0
)

type BranchReportPdf struct {
	config   *config.PDF
	pdf      *fpdf.Fpdf
	filename string
	count    int
}

func NewBranchReportPdf(config *config.PDF) generators.IReport {
	return &BranchReportPdf{
		config: config,
	}
}

func (br *BranchReportPdf) GenerateReport(msg *consumer.Message) error {

	switch msg.Type {

	case restaurant_type:
		restaurant := models.UnmarshalRestaurant(msg.Data)
		br.writeRestaurantData(restaurant)

		return nil

	case branch_type:
		branches := models.UnmarshalBranches(msg.Data)
		br.writeBranchesData(branches)

		return nil

	default:
		zap.Log.Warn("Type data not implemented: ", msg.Type)
		return errors.New("type data not implemented")
	}
}

func (br *BranchReportPdf) writeRestaurantData(restaurant *models.Restaurant) {

	font := br.config.Font
	title := br.config.Title
	br.filename = restaurant.Name

	br.pdf = fpdf.New("P", "mm", "A4", "")

	br.pdf.SetMargins(marginX, marginY, marginX)
	br.pdf.AddPage()
	pageW, _ := br.pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX

	br.pdf.SetFont(font, "B", 16)
	_, lineHeight := br.pdf.GetFontSize()
	currentY := br.pdf.GetY() + lineHeight + gapY
	br.pdf.SetXY(marginX, currentY)
	br.pdf.Cell(40, 10, restaurant.Name)

	leftY := br.pdf.GetY() + lineHeight + gapY
	br.pdf.SetFont(font, "B", 30)
	_, lineHeight = br.pdf.GetFontSize()
	br.pdf.SetXY(100, currentY-lineHeight)
	br.pdf.Cell(100, 40, title)

	newY := leftY
	if (br.pdf.GetY() + gapY) > newY {
		newY = br.pdf.GetY() + gapY
	}

	newY += 10.0

	br.pdf.SetXY(marginX, newY)
	br.pdf.SetFont(font, "", 12)
	_, lineHeight = br.pdf.GetFontSize()
	lineBreak := lineHeight + float64(2)

	br.pdf.SetFontStyle("B")
	br.pdf.Cell(safeAreaW/2, lineHeight, "Restaurant Information:")
	br.pdf.Line(marginX, br.pdf.GetY()+lineHeight, 70, br.pdf.GetY()+lineHeight)
	br.pdf.Ln(lineBreak)
	br.pdf.Ln(lineBreak)

	br.pdf.SetFontStyle("I")
	br.pdf.Cell(safeAreaW/2, lineHeight, restaurant.Location)
	br.pdf.Ln(lineBreak)

	br.pdf.SetFontStyle("I")
	br.pdf.Cell(safeAreaW/2, lineHeight, restaurant.Country)
	br.pdf.Ln(lineBreak)

	br.pdf.SetFontStyle("I")
	br.pdf.Cell(safeAreaW/2, lineHeight, restaurant.Founder)
	br.pdf.Ln(lineBreak)

	br.pdf.SetFontStyle("I")
	br.pdf.Cell(safeAreaW/2, lineHeight, restaurant.Fundation)
	br.pdf.Ln(lineBreak)

	br.pdf.SetFontStyle("I")
	br.pdf.Cell(safeAreaW/2, lineHeight, restaurant.Headquarter)
	br.pdf.Ln(lineBreak)
	br.pdf.Ln(lineBreak)
	br.pdf.Ln(lineBreak)
}

func (br *BranchReportPdf) writeBranchesData(branches *[]models.Branch) {

	pageW, _ := br.pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX
	_, lineHeight := br.pdf.GetFontSize()
	lineBreak := lineHeight + float64(2)

	br.pdf.Line(marginX, br.pdf.GetY()+lineHeight, 200, br.pdf.GetY()+lineHeight)
	br.pdf.Ln(lineBreak)

	for _, b := range *branches {

		br.count = br.count + 1

		br.pdf.SetFontStyle("B")
		br.pdf.SetX((safeAreaW / 2) - float64((len(b.Name) / 2)))
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("%d. %s", br.count, b.Name))
		br.pdf.Ln(lineBreak)
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("I")
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Manager: %s", b.Manager))
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("I")
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("City: %s", b.City))
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("I")
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Address: %s", b.Address))
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("I")
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Phone: %s", b.Phone))
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("B")
		br.pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Score: %d", b.Score))
		br.pdf.Ln(lineBreak)
		br.pdf.Ln(lineBreak)
		br.pdf.Ln(lineBreak)

		br.pdf.SetFontStyle("B")
		br.pdf.Cell(safeAreaW/2, lineHeight, "Employees")
		br.pdf.Ln(lineBreak)

		lineHt := 10.0
		const colNumber = 4
		header := [colNumber]string{"No", "Name", "Years", "Sales"}
		colWidth := [colNumber]float64{10.0, 95.0, 35.0, 50.0}

		// Headers
		br.pdf.SetFontStyle("B")
		br.pdf.SetFillColor(200, 200, 200)
		for colJ := 0; colJ < colNumber; colJ++ {
			br.pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "CM", true, 0, "")
		}
		br.pdf.Ln(lineBreak)
		br.pdf.Ln(lineBreak)
		br.pdf.Ln(lineBreak)
	}
}

func (br *BranchReportPdf) CloseReport() (string, error) {
	file := br.filename + "-" + utils.TimestampID() + br.config.Suffix
	zap.Log.Info("Closing file")

	return file, br.pdf.OutputFileAndClose(br.config.Path + file)
}

func (br *BranchReportPdf) DeleteReport() error {
	return nil
}
