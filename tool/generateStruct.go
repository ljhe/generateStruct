package tool

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

/**
 * 将excel中的前四列转化为struct
 * 第一列字段类型		如 int
 * 第二列字段名称		如 显示顺序
 * 第三列字段名		如 id
 * 第四列s,c,all 	s表示服务端使用 c表示客户端使用 all表示都使用
 */

var (
	lineNumber           = 4                                     // 每个工作表需要读取的行数
	structBegin          = "type %s struct {\n"                  // 结构体开始
	structValue          = "    %s %s	`col:\"%s\" client:\"%s\"`" // 结构体的内容
	structValueForServer = "    %s %s	`col:\"%s\"`"              // 服务端使用的结构体内容
	structRemarks        = "	 // %s"                             // 结构体备注
	structValueEnd       = "\n"                                  // 结构体内容结束
	structEnd            = "}\n"                                 // 结构体结束
	header               = "package %s\n\r"                      // 文件头
)

type Generate struct {
	savePath string // 生成文件的保存路径
	data     string // 生成文件的内容
	allType  string // 文件当中的数据类型
}

// 按行读取文件
func (this *Generate) ReadExcelOld(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	read := bufio.NewReader(file)
	for {
		line, _, err := read.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Printf("%v\n", string(line))
	}
	return nil
}

// 读取excel
func (this *Generate) ReadExcel(readPath, savePath, allType string) error {
	if savePath == "" || allType == "" {
		return fmt.Errorf("ReadExcel|savePath or allType is nil")
	}
	this.savePath = savePath
	this.allType = allType
	files, err := ioutil.ReadDir(readPath)
	if err != nil {
		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
	}
	for _, file := range files {
		if path.Ext(file.Name()) != ".xlsx" || hasChineseOrDefault(file.Name()) {
			continue
		}
		wb, err := xlsx.OpenFile(readPath + "\\" + file.Name())
		if err != nil {
			return fmt.Errorf("ReadExcel|xlsx.OpenFile is err :%v", err)
		}
		// 遍历工作表
		for _, sheet := range wb.Sheets {
			if hasChineseOrDefault(sheet.Name) {
				continue
			}
			sheetData := make([][]string, 0)
			// 判断表格中内容的行数是否小于需要读取的行数
			if sheet.MaxRow < lineNumber {
				return fmt.Errorf("ReadExcel|sheet.MaxRow:%d < lineNumber:%d", sheet.MaxRow, lineNumber)
			}
			// 遍历列
			for i := 0; i < sheet.MaxCol; i++ {
				// 判断某一列的第一行是否为空
				if sheet.Cell(0, i).Value == "" {
					continue
				}
				cellData := make([]string, 0)
				// 遍历行
				for j := 0; j < lineNumber; j++ {
					cellData = append(cellData, sheet.Cell(j, i).Value)
				}
				sheetData = append(sheetData, cellData)
			}
			err := this.SplicingData(sheetData, sheet.Name)
			if err != nil {
				return fmt.Errorf("fileName:\"%v\" is err:%v", file.Name(), err)
			}
		}
	}
	if this.data == "" {
		return fmt.Errorf("ReadExcel|this.data is nil")
	}
	err = this.WriteNewFile(this.data)
	if err != nil {
		return err
	}
	return nil
}

// 拼装struct
func (this *Generate) SplicingData(data [][]string, structName string) error {
	structData := fmt.Sprintf(structBegin, firstRuneToUpper(structName))
	for _, value := range data {
		if len(value) != lineNumber {
			return fmt.Errorf("SplicingData|sheetName:%v col's len:%d is err", value, len(value))
		}
		err := this.CheckType(value[0], structName)
		if err != nil {
			return err
		}
		switch value[3] {
		case "all":
			structData += fmt.Sprintf(structValue, firstRuneToUpper(value[2]), value[0], value[2], value[2])
			if value[1] != "" {
				structData += fmt.Sprintf(structRemarks, value[1])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "s":
			structData += fmt.Sprintf(structValueForServer, firstRuneToUpper(value[2]), value[0], value[2])
			if value[1] != "" {
				structData += fmt.Sprintf(structRemarks, value[1])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "c":
			continue
		default:
			return fmt.Errorf("SplicingData|value[3]:\"%v\" is not in s,c,all", value[3])
		}
	}
	structData += structEnd
	this.data += structData
	return nil
}

// 拼装好的struct写入新的文件
func (this *Generate) WriteNewFile(data string) error {
	str := strings.Split(this.savePath, "\\")
	if len(str) == 0 {
		return fmt.Errorf("WriteNewFile|len(str) is 0")
	}
	header = fmt.Sprintf(header, str[len(str)-1])
	data = header + data
	fw, err := os.OpenFile(this.savePath+"\\objs.go", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
	}
	defer fw.Close()
	_, err = fw.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("WriteNewFile|Write is err:%v", err)
	}
	return nil
}

// 检测解析出来的字段类型是否符合要求
func (this *Generate) CheckType(dataType, structName string) error {
	res := strings.Index(this.allType, dataType)
	if res == -1 {
		return fmt.Errorf("CheckType|struct:\"%v\" dataType:\"%v\" is not in provide dataType", structName, dataType)
	}
	return nil
}

// 字符串首字母转换成大写
func firstRuneToUpper(str string) string {
	data := []byte(str)
	for k, v := range data {
		if k == 0 {
			first := []byte(strings.ToUpper(string(v)))
			newData := data[1:]
			data = append(first, newData...)
			break
		}
	}
	return string(data[:])
}

// 判断是否存在汉字或者是否为默认的工作表
func hasChineseOrDefault(r string) bool {
	if strings.Index(r, "Sheet") != -1 {
		return true
	}
	for _, v := range []rune(r) {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}
