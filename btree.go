package parseibd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"parseibd/parser"
	"parseibd/proto"
	"strings"
)

type BtreeNode struct {
	PageNo		uint64
	MinIndex	interface{}
	MaxIndex	interface{}
	ChildPages	[]*BtreeNode
	Level		int
	SegLeaf		*proto.FileSegmentHeader
	SegNonLeaf	*proto.FileSegmentHeader
}

func ShowBtree(path string, pageNo uint64, fields []*proto.Field, currIndex, priIndex *proto.Index,
	tableCharset, inputCharset, rowFormat string) {
	btreeNode := GetBtreeNode(path, pageNo, fields, currIndex, priIndex, tableCharset, inputCharset, rowFormat, 0)
	//深度遍历
	var builder bytes.Buffer
	builder.WriteString("#############################################\n")
	builder.WriteString("#                索引segment                #\n")
	builder.WriteString("#############################################\n")
	builder.WriteString(fmt.Sprintf("非叶子节点segment：page#%v@segment#%v\n", btreeNode.SegNonLeaf.PageNo,
		parser.GetSegmentNo(int(btreeNode.SegNonLeaf.Offset))))
	builder.WriteString(fmt.Sprintf("叶子节点segment：page#%v@segment#%v\n", btreeNode.SegLeaf.PageNo,
		parser.GetSegmentNo(int(btreeNode.SegLeaf.Offset))))
	builder.WriteString("\n#############################################\n")
	builder.WriteString("#                 索引btree                 #\n")
	builder.WriteString("#############################################\n")
	var queue []*BtreeNode
	queue = append(queue, btreeNode)
	for len(queue) > 0 {
		first := queue[0]
		queue = queue[1:]
		for i:=0; i<first.Level-1; i++ {
			builder.WriteString("|        ")
		}
		if first.Level > 0 {
			builder.WriteString("|--------")
		}
		builder.WriteString(fmt.Sprintf("page#%v 索引范围:(%v, %v)\n", first.PageNo, first.MinIndex, first.MaxIndex))
		if len(first.ChildPages) > 0 {
			queue = append(first.ChildPages, queue...)
		}
	}
	filename := fmt.Sprintf("./output/btree_%v", strings.Join(currIndex.FieldNames, "_"))
	_ = ioutil.WriteFile(filename, builder.Bytes(), 0666)
}

func GetBtreeNode(path string, pageNo uint64, fields []*proto.Field, currIndex, priIndex *proto.Index,
	tableCharset, inputCharset, rowFormat string, level int) *BtreeNode {
	pageIndex := parser.ReadPageIndex(path, pageNo, fields, currIndex, priIndex, tableCharset, inputCharset, rowFormat)
	node := new(BtreeNode)
	node.PageNo = pageNo
	node.Level = level
	node.SegLeaf = pageIndex.PageHeader.SegLeaf
	node.SegNonLeaf = pageIndex.PageHeader.SegNonLeaf
	if pageIndex.UserRecordPriNonleafs != nil {
		//聚簇索引，非叶子节点
		n := len(pageIndex.UserRecordPriNonleafs)
		node.MinIndex = pageIndex.UserRecordPriNonleafs[0].Primary
		node.MaxIndex = pageIndex.UserRecordPriNonleafs[n-1].Primary
		for _, userRecord := range pageIndex.UserRecordPriNonleafs {
			childNode := GetBtreeNode(path, uint64(userRecord.ChildPageNo), fields, currIndex,
				priIndex, tableCharset, inputCharset, rowFormat, level+1)
			node.ChildPages = append(node.ChildPages, childNode)
		}
	} else if pageIndex.UserRecordPriLeafs != nil {
		//聚簇索引，叶子节点
		n := len(pageIndex.UserRecordPriLeafs)
		node.MinIndex = pageIndex.UserRecordPriLeafs[0].Primary
		node.MaxIndex = pageIndex.UserRecordPriLeafs[n-1].Primary
	} else if pageIndex.UserRecordSecNonleafs != nil {
		//二级索引，非叶子节点
		n := len(pageIndex.UserRecordSecNonleafs)
		node.MinIndex = pageIndex.UserRecordSecNonleafs[0].Secondary
		node.MaxIndex = pageIndex.UserRecordSecNonleafs[n-1].Secondary
		for _, userRecord := range pageIndex.UserRecordSecNonleafs {
			childNode := GetBtreeNode(path, uint64(userRecord.ChildPageNo), fields, currIndex,
				priIndex, tableCharset, inputCharset, rowFormat, level+1)
			node.ChildPages = append(node.ChildPages, childNode)
		}
	} else if pageIndex.UserRecordSecLeafs != nil {
		//二级索引，叶子节点
		n := len(pageIndex.UserRecordSecLeafs)
		node.MinIndex = pageIndex.UserRecordSecLeafs[0].Secondary
		node.MaxIndex = pageIndex.UserRecordSecLeafs[n-1].Secondary
	}
	return node
}