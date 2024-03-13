package exporter

import (
	"api-tester/pkg/api/tester"
	"api-tester/pkg/utilities/goext"
	"fmt"
	"github.com/xuri/excelize/v2"
)

type Exporter struct {
}

func (e *Exporter) ToExcel(testReports []*tester.Report, filename string) error {
	var (
		f     = excelize.NewFile()
		sheet = "Sheet1"
	)

	// 设置标题样式
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "1f7f3b", Bold: true, Family: "Microsoft YaHei", Size: 16},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E6F4EA"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
		Border:    []excelize.Border{{Type: "top", Style: 2, Color: "1f7f3b"}},
	})
	if err := f.SetCellStyle(sheet, "A1", "E1", titleStyle); err != nil {
		return err
	}

	// 设置内容样式
	contentStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "Microsoft YaHei", Size: 15},
	})
	if err := f.SetCellStyle(sheet, "A2", fmt.Sprintf("E%d", len(testReports)+1), contentStyle); err != nil {
		return err
	}

	if err := f.SetColWidth(sheet, "A", "B", 45); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "C", "D", 10); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "E", "E", 20); err != nil {
		return err
	}

	// 设置工作表列名
	_ = f.SetCellValue(sheet, "A1", "接口名称")
	_ = f.SetCellValue(sheet, "B1", "接口路径")
	_ = f.SetCellValue(sheet, "C1", "是否通过")
	_ = f.SetCellValue(sheet, "D1", "状态码")
	_ = f.SetCellValue(sheet, "E1", "耗时(秒)")

	// 遍历测试报告并在工作表中写入数据
	for i, report := range testReports {
		row := i + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), report.ApiName)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), report.ApiPath)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), goext.If(report.IsPassed, "是", "否"))
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), report.Response.StatusCode())
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), report.Elapsed.Seconds())
	}

	if err := f.AddTable(sheet, &excelize.Table{
		Range:             fmt.Sprintf("A1:E%d", len(testReports)+1),
		Name:              "table",
		StyleName:         "TableStyleMedium2",
		ShowFirstColumn:   true,
		ShowLastColumn:    true,
		ShowColumnStripes: true,
	}); err != nil {
		return err
	}

	// 工作表重命名
	_ = f.SetSheetName(sheet, "测试报告")

	// 将 Excel 文件保存到磁盘
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}
