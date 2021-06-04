package main

import (
	"database/sql"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
)

func exportToExcel(db *sql.DB) (*excelize.File, error) {
	file := excelize.NewFile()
	for _type, name := range _Map {
		list, err := findAllByGachaType(db, _type)
		if err != nil {
			return nil, err
		}
		err = toExcel(file, name, list)
		if err != nil {
			return nil, err
		}
	}
	file.DeleteSheet("Sheet1")
	return file, nil
}

func toExcel(excel *excelize.File, sheetName string, list RealDataList) error {
	excel.SetActiveSheet(excel.NewSheet(sheetName))
	slice := &[]interface{}{"时间", "名称", "类别", "星级", "总次数", "保底内"}
	err := excel.SetSheetRow(sheetName, "A1", slice)
	if err != nil {
		return err
	}
	count := 1
	for i, data := range list {
		slice := &[]interface{}{data.Time, data.Name, data.ItemType, data.RankType, i + 1, count}
		err := excel.SetSheetRow(sheetName, "A"+strconv.Itoa(i+2), slice)
		if err != nil {
			return err
		}
		if data.RankType == "5" {
			count = 1
		} else {
			count++
		}
	}
	return nil
}
