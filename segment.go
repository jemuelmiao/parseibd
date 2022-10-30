package parseibd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"parseibd/parser"
	"parseibd/proto"
	"strconv"
	"strings"
)

func ShowSegments(path string) {
	pageFsp := parser.ReadPageFsp(path, 0)
	var builder bytes.Buffer
	if pageFsp.FspHeader.FreeInodes != nil && pageFsp.FspHeader.FreeInodes.Len > 0 {
		builder.WriteString(GetSegmentText(path, pageFsp.FspHeader.FreeInodes.First))
	}
	if pageFsp.FspHeader.FullInodes != nil && pageFsp.FspHeader.FullInodes.Len > 0 {
		builder.WriteString(GetSegmentText(path, pageFsp.FspHeader.FullInodes.First))
	}
	_ = ioutil.WriteFile("./output/segments", builder.Bytes(), 0666)
}

func GetSegmentText(path string, firstInode *proto.FilAddr) string {
	var builder bytes.Buffer
	pageInodes := GetPageInodes(path, firstInode)
	for _, pageInode := range pageInodes {
		pageNoStr := strconv.FormatInt(int64(pageInode.FileHeader.Offset), 10)
		builder.WriteString("#############################################\n")
		builder.WriteString(fmt.Sprintf("#                   page#%v%v#\n", pageNoStr, getSpaceStr(pageNoStr, 19)))
		builder.WriteString("#############################################\n")
		for i, segment := range pageInode.SegmentInodes {
			freeExtents, fragExtents, fullExtents, fragPages := GetSegmentNos(path, segment)
			builder.WriteString(fmt.Sprintf("============segment#%v============\n", i))
			builder.WriteString("空闲的extent列表：")
			builder.WriteString(strings.Join(freeExtents, ", "))
			builder.WriteString("\n")
			builder.WriteString("部分使用的extent列表：")
			builder.WriteString(strings.Join(fragExtents, ", "))
			builder.WriteString("\n")
			builder.WriteString("全部使用的extent列表：")
			builder.WriteString(strings.Join(fullExtents, ", "))
			builder.WriteString("\n")
			builder.WriteString("碎片page列表：")
			builder.WriteString(strings.Join(fragPages, ", "))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func GetPageInodes(path string, inode *proto.FilAddr) []*proto.PageInode {
	var pageInodes []*proto.PageInode
	for inode != nil && inode.PageNo != parser.FilNull {
		pageInode := parser.ReadPageInode(path, uint64(inode.PageNo))
		pageInodes = append(pageInodes, pageInode)
		inode = pageInode.FlstNode.Next
	}
	return pageInodes
}

func GetSegmentNos(path string, segment *proto.SegmentInode) ([]string, []string, []string, []string) {
	var freeExtents, fragExtents, fullExtents, fragPages []string
	if segment.Free != nil && segment.Free.Len > 0 {
		freeExtents = GetExtentNos(path, segment.Free.First)
	}
	if segment.NotFull != nil && segment.NotFull.Len > 0 {
		fragExtents = GetExtentNos(path, segment.NotFull.First)
	}
	if segment.Full != nil && segment.Full.Len > 0 {
		fullExtents = GetExtentNos(path, segment.Full.First)
	}
	for _, fragPage := range segment.FragPages {
		if fragPage == parser.FilNull {
			fragPages = append(fragPages, "FilNull")
		} else {
			fragPages = append(fragPages, fmt.Sprintf("page#%v", fragPage))
		}
	}
	return freeExtents, fragExtents, fullExtents, fragPages
}