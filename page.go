package parseibd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"parseibd/parser"
)

func ShowPages(path string) {
	type Page struct {
		PageNo		uint64
		PageType	string
	}
	data, e := ioutil.ReadFile(path)
	if e != nil {
		return
	}
	size := len(data)
	var pages []*Page
	for i:=0; i<size/parser.PageSize; i++ {
		start := i * parser.PageSize
		end := (i+1) * parser.PageSize
		buff := data[start:end]
		fileHeader, _ := parser.ReadFileHeader(buff)
		page := &Page{
			PageNo:   uint64(i),
			PageType: parser.GetPageTypeName(int(fileHeader.Type)),
		}
		pages = append(pages, page)
	}
	var builder bytes.Buffer
	for _, page := range pages {
		builder.WriteString(fmt.Sprintf("page#%v    %v\n", page.PageNo, page.PageType))
	}
	_ = ioutil.WriteFile("./output/pages", builder.Bytes(), 0666)
}