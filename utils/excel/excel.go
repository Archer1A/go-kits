package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

type Client[T any] struct {
	currentSheet string
	excel        *excelize.File
	headerMap    map[string]int
	error        error
}

func (g *Client[T]) Header(i interface{}) Generator[T] {
	header := GetHeader(i)
	g.headerMap = make(map[string]int)
	g.excel.SetSheetRow(g.currentSheet, "A1", &header)
	for index, h := range header {
		g.headerMap[fmt.Sprintf("%v", h)] = index
	}
	return g
}

func (g *Client[T]) Save(name string) error {
	err := g.excel.SaveAs(name)
	if err != nil {
		return err
	}
	err = g.excel.Close()
	if err != nil {
		return err
	}
	return nil
}

type Generator[T any] interface {
	Sheet(sheet string) Generator[T]
	Header(i interface{}) Generator[T]
	Save(name string) error
	Rows(rows []T) Generator[T]
}

func (g *Client[T]) Rows(rows []T) Generator[T] {
	for index, row := range rows {

		headerMap := MapValueToHeader(row)
		for k, v := range headerMap {
			colIndex := g.headerMap[k]
			name, _ := excelize.CoordinatesToCellName(colIndex+1, index+2)
			g.excel.SetCellValue(g.currentSheet, name, v)
		}
	}
	return g
}

func (g *Client[T]) Sheet(sheet string) Generator[T] {
	if g.excel == nil {
		g.excel = excelize.NewFile()
	}
	index, err := g.excel.NewSheet(sheet)
	if err != nil {
		g.error = err
		return g
	}
	g.currentSheet = sheet
	g.excel.SetActiveSheet(index)
	return g
}

func Create() {
	workbook := excelize.NewFile()
	index, err := workbook.NewSheet("Sheet1")
	if err != nil {
		return
	}
	workbook.SetActiveSheet(index)

	workbook.SetCellValue("Sheet1", "A1", "Hello World")
	workbook.SaveAs("Book1.xlsx")
	workbook.Close()

}

type User struct {
	Username string `json:"username" excel:"header:用户名"`
	Sex      string `json:"sex" excel:"header:性别"`
	Age      int    `json:"age" excel:"header:年龄"`
	Phone    string `json:"phone" excel:"header:手机号"`
	Email    string `json:"email" excel:"header:邮箱"`
}

func (u *User) Sheet() string {
	return "Sheet1"
}

func MapValueToHeader(object interface{}) map[string]interface{} {
	headerMap := make(map[string]interface{})
	v := reflect.ValueOf(object)
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := getTagOptionFromKey(field.Tag.Get("excel"), "header")
		if key == "-" {
			continue
		}
		if key == "" {
			key = field.Name
		}
		headerMap[key] = v.Field(i).Interface()
	}
	return headerMap
}

func isFieldPublic(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}
	// 检查首字母是否大写
	firstChar := fieldName[0]
	return firstChar >= 'A' && firstChar <= 'Z'
}

func getTag(object interface{}) {
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		get := field.Tag.Get("excel")
		if get == "" {
			continue
		}
		fmt.Println(get)
	}
}

func GetHeader(object interface{}) []string {
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	headers := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !isFieldPublic(field.Name) {
			continue
		}
		excelTag := field.Tag.Get("excel")
		if excelTag == "" {
			continue
		}
		allOptions := strings.Split(excelTag, ";")
		for _, tag := range allOptions {
			header := getTagOptionFromKey(tag, "header")
			if header == "-" {
				continue
			}
			if header == "" {
				header = field.Name
			}
			headers = append(headers, header)
		}
	}
	return headers
}

func getTagOptionFromKey(tag string, key string) string {
	allOptions := strings.Split(tag, ";")
	for _, option := range allOptions {
		if strings.HasPrefix(option, key+":") {
			return option[len(key)+1:]
		}
	}
	return ""
}
