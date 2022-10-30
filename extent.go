package parseibd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"parseibd/parser"
	"parseibd/proto"
)

func ShowExtents(path string) {
	pageFsp := parser.ReadPageFsp(path, 0)
	var builder bytes.Buffer
	if pageFsp.FspHeader.Free != nil && pageFsp.FspHeader.Free.Len > 0 {
		extents := GetExtentNos(path, pageFsp.FspHeader.Free.First)
		builder.WriteString("#############################################\n")
		builder.WriteString("#           全局空闲的extent列表            #\n")
		builder.WriteString("#############################################\n")
		for _, extent := range extents {
			builder.WriteString(extent)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	if pageFsp.FspHeader.FreeFrag != nil && pageFsp.FspHeader.FreeFrag.Len > 0 {
		extents := GetExtentNos(path, pageFsp.FspHeader.FreeFrag.First)
		builder.WriteString("#############################################\n")
		builder.WriteString("#          全局部分使用的extent列表         #\n")
		builder.WriteString("#############################################\n")
		for _, extent := range extents {
			builder.WriteString(extent)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	if pageFsp.FspHeader.FullFrag != nil && pageFsp.FspHeader.FullFrag.Len > 0 {
		extents := GetExtentNos(path, pageFsp.FspHeader.FullFrag.First)
		builder.WriteString("#############################################\n")
		builder.WriteString("#          全局全部使用的extent列表         #\n")
		builder.WriteString("#############################################\n")
		for _, extent := range extents {
			builder.WriteString(extent)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	_ = ioutil.WriteFile("./output/extents", builder.Bytes(), 0666)
}

func GetExtentNos(path string, extent *proto.FilAddr) []string {
	var extentNos []string
	for extent != nil && extent.PageNo != parser.FilNull {
		extentNo := parser.GetExtentNo(int(extent.Offset))
		extentNos = append(extentNos, fmt.Sprintf("page#%v@extent#%v", extent.PageNo, extentNo))
		buff, e := parser.ReadPage(path, uint64(extent.PageNo))
		if e != nil {
			break
		}
		fileHeader, _ := parser.ReadFileHeader(buff)
		if fileHeader.Type == parser.PageTypeFspHdr {
			page := parser.ReadPageFsp(path, uint64(extent.PageNo))
			extent = page.ExtentDescriptors[extentNo].FlstNode.Next
		} else if fileHeader.Type == parser.PageTypeXdes {
			page := parser.ReadPageXdes(path, uint64(extent.PageNo))
			extent = page.ExtentDescriptors[extentNo].FlstNode.Next
		} else {
			extent = nil
		}
	}
	return extentNos
}