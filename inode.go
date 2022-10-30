package parseibd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"parseibd/parser"
	"parseibd/proto"
)

func ShowInodes(path string) {
	pageFsp := parser.ReadPageFsp(path, 0)
	var builder bytes.Buffer
	if pageFsp.FspHeader.FreeInodes != nil && pageFsp.FspHeader.FreeInodes.Len > 0 {
		inodes := GetInodeNos(path, pageFsp.FspHeader.FreeInodes.First)
		builder.WriteString("#############################################\n")
		builder.WriteString("#          全局部分使用的inode列表          #\n")
		builder.WriteString("#############################################\n")
		for _, inode := range inodes {
			builder.WriteString(inode)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	if pageFsp.FspHeader.FullInodes != nil && pageFsp.FspHeader.FullInodes.Len > 0 {
		inodes := GetInodeNos(path, pageFsp.FspHeader.FullInodes.First)
		builder.WriteString("#############################################\n")
		builder.WriteString("#          全局全部使用的inode列表          #\n")
		builder.WriteString("#############################################\n")
		for _, inode := range inodes {
			builder.WriteString(inode)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	_ = ioutil.WriteFile("./output/inodes", builder.Bytes(), 0666)
}

func GetInodeNos(path string, inode *proto.FilAddr) []string {
	var inodeNos []string
	for inode != nil && inode.PageNo != parser.FilNull {
		inodeNos = append(inodeNos, fmt.Sprintf("page#%v", inode.PageNo))
		buff, e := parser.ReadPage(path, uint64(inode.PageNo))
		if e != nil {
			break
		}
		node, _ := parser.ReadFlstNode(buff[inode.Offset:])
		inode = node.Next
	}
	return inodeNos
}