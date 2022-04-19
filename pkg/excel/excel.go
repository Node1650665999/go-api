package excel

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"net/http"
	"os"
	"strconv"
)

//ExcelData 定义了读写excel/csv文件的数据格式
type ExcelData [][]string

//ImportExcel 实现了excel文件的读取
//refer : https://xuri.me/excelize/zh-hans/sheet.html
func ImportExcel(filename io.Reader) (ExcelData,error) {
	xlsx, err := excelize.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	return xlsx.GetRows("Sheet1"), nil
}

//ExportExcel 实现了excel导出
func ExportExcel(filename string,data ExcelData,w http.ResponseWriter)  {
	xlsx := excelize.NewFile()
	for index, rowData :=  range data{
		xlsx.SetSheetRow("Sheet1", "A" + strconv.Itoa(index + 1), &rowData) // SetSheetRow：设置一行数据 SetCellValue：设置一个数据
	}
	// 设置下载的文件名
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	xlsx.Write(w)
	return
}

//SaveExcel 保存excel文件到本地
func SaveExcel(filename string, data ExcelData) error  {
	xlsx := excelize.NewFile()
	for index, rowData :=  range data{
		//以 A1 单元格作为起始坐标按行赋值
		xlsx.SetSheetRow("Sheet1", "A" + strconv.Itoa(index + 1), &rowData) // SetSheetRow：设置一行数据 SetCellValue：设置一个单元格
	}
	return xlsx.SaveAs(filename)
}

//ImportCsv 实现了读取 csv 文件
func ImportCsv(filename io.Reader) (ExcelData, error) {
	reader := csv.NewReader(filename)
	//一次性读完
	return reader.ReadAll()
}

//ExportCsv 实现导出 csv 文件
func ExportCsv(filename string, data ExcelData, w http.ResponseWriter) {
	buf := &bytes.Buffer{}
	buf.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM,避免文件打开乱码

	writer := csv.NewWriter(buf)
	writer.WriteAll(data)

	// 设置下载的文件名
	w.Header().Add("Content-Type", "application/octet-stream")
	//w.Header().Add("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))

	//输出数据
	w.Write(buf.Bytes())
	return
}

//SaveCsv 保存生成的csv文件到本地
func SaveCsv(filename string, data ExcelData) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	return w.WriteAll(data)
}