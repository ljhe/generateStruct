package main

import (
	"flag"
	"fmt"
	"generateStruct/tool"
)

var (
	savePath = flag.String("savePath", "", "Path to save the makefile")
	readPath = flag.String("readPath", "", "The path of reading Excel")
	allType  = flag.String("allType", "", "Specified field type")
)

func main() {
	flag.Parse()
	if *savePath == "" || *readPath == "" || *allType == "" {
		fmt.Println("savePath, readPath or allType is nil")
		return
	}
	gt := tool.Generate{}
	err := gt.ReadExcel(*readPath, *savePath, *allType)
	if err != nil {
		fmt.Printf("something err:%v\n", err)
		return
	}
}
