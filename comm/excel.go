package comm

import (
	"errors"
	"github.com/xuri/excelize/v2"
)

type excel struct {
	File *excelize.File
	TitleRow int
	ValStartRow int
	path string
}

type ExcelCallbackReturn struct {
	ErrFlag bool
	ForType string
	Error error
}

func NewCreateExcel() (excel) {
	model := excel{}

	//path := "E:\\_Project\\Box\\_\\gotest\\wabuy\\_handyDev\\处理Excel\\写\\123.xlsx"
	model.File = excelize.NewFile()
	//model.File.SaveAs(path)
	//f, err := excelize.OpenFile(path, excelize.Options{UnzipXMLSizeLimit: "password"})
	//if err != nil {
	//	return model, err
	//}
	//model.File.SetSheetViewOptions()
	//model.File.SetSheetPrOptions()

	return model
}

func NewOpenExcel(path string) (excel, error) {
	model := excel{path: path, TitleRow: 1, ValStartRow: 1,}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return model, err
	}
	model.File = f
	return model, nil
}

// 设置标题对应多少行
func (model *excel) SetTitleRow(row int) *excel {
	model.TitleRow = row
	return model
}

// 设置从第几行开始处理数据
func (model *excel) SetValStartRow(row int) *excel {
	model.ValStartRow = row - 1
	return model
}

func (model *excel) GetSheetIndex (name string) int {
	return model.File.GetSheetIndex(name)
}

// 处理sheet内容
func (model *excel) HandleSheet(sheetIndex int, title map[string][]string, callback func(info map[string]string) ExcelCallbackReturn) error {
	rows, err := model.File.GetRows(model.File.GetSheetName(sheetIndex))
	if err != nil {
		return err
	}

	key ,exists := getExcelTable(rows, title, model.TitleRow)
	if !exists {
		return errors.New("标题不完整")
	}
	rowsLen := len(rows)
	var row []string
	var callbackReturn ExcelCallbackReturn
	temp := make(map[string]string)
	for i := model.ValStartRow; i < rowsLen; i++ {
		row = rows[i]
		for keyName, keyIndex := range key {
			temp[keyName] = GetArrIndexString(row, keyIndex)
		}
		callbackReturn = callback(temp)

		if callbackReturn.ErrFlag {
			return callbackReturn.Error
		}

		if callbackReturn.ForType == "continue" {
			continue
		} else if callbackReturn.ForType == "break" {
			break
		}
	}

	return nil
}

// 只处理第一个sheet
func (model *excel) HandleSheetOneRow(title map[string][]string, callback func(info map[string]string) ExcelCallbackReturn) error {
	return model.HandleSheet(0, title, callback)
}


func getExcelTable(rows [][]string, titleKeyValue map[string][]string, titleRow int) (map[string]int,  bool) {
	key := map[string]int{}
	for i := 0; i < titleRow; i++ {
		for colCellk, colCell := range rows[i] {
			for titleKey, titleVal := range titleKeyValue {
				exists, _ := InArray(colCell, titleVal)
				if exists {
					key[titleKey] = colCellk
					if len(key) == len(titleKeyValue) {
						return key, true
					}
				}
			}
		}
	}

	return key, false
}