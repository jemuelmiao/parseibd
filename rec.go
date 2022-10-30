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

type BtreeRec struct {
	PageNo		uint64
	UserRecordPriNonleafs 	[]*proto.UserRecordPriNonleaf //聚簇索引-非叶子记录
	UserRecordPriLeafs 		[]*proto.UserRecordPriLeaf //聚簇索引-叶子记录
	UserRecordSecNonleafs 	[]*proto.UserRecordSecNonleaf //二级索引-非叶子记录
	UserRecordSecLeafs 		[]*proto.UserRecordSecLeaf //二级索引-叶子记录
	ChildPages	[]*BtreeRec
}

func ShowRecs(path string, pageNo uint64, fields []*proto.Field, currIndex, priIndex *proto.Index,
	tableCharset, inputCharset, rowFormat string) {
	btreeRec := GetBtreeRec(path, pageNo, fields, currIndex, priIndex, tableCharset, inputCharset, rowFormat)
	//广度遍历
	var builder bytes.Buffer
	var queue []*BtreeRec
	queue = append(queue, btreeRec)
	for len(queue) > 0 {
		first := queue[0]
		queue = queue[1:]
		queue = append(queue, first.ChildPages...)
		if first.UserRecordPriNonleafs != nil {
			//聚簇索引，非叶子节点
			pageNoStr := strconv.FormatInt(int64(first.PageNo), 10)
			builder.WriteString("############################################################\n")
			builder.WriteString(fmt.Sprintf("#               聚簇索引、非叶子page#%v%v#\n", pageNoStr, getSpaceStr(pageNoStr, 22)))
			builder.WriteString("#              格式：格式：主键值，子page编号              #\n")
			builder.WriteString("############################################################\n")
			for _, userRecord := range first.UserRecordPriNonleafs {
				builder.WriteString(fmt.Sprintf("%v, %v\n", userRecord.Primary, userRecord.ChildPageNo))
			}
			builder.WriteString("\n")
		} else if first.UserRecordPriLeafs != nil {
			//聚簇索引，叶子节点
			pageNoStr := strconv.FormatInt(int64(first.PageNo), 10)
			builder.WriteString("############################################################\n")
			builder.WriteString(fmt.Sprintf("#                 聚簇索引、叶子page#%v%v#\n", pageNoStr, getSpaceStr(pageNoStr, 22)))
			builder.WriteString("#         格式：主键值，trx_id，roll_ptr，其他字段         #\n")
			builder.WriteString("############################################################\n")
			for _, userRecord := range first.UserRecordPriLeafs {
				builder.WriteString(fmt.Sprintf("%v, %v, %v, %v\n", userRecord.Primary, userRecord.TrxId,
					parser.FormatHex(userRecord.RollPtr), userRecord.Values))
			}
			builder.WriteString("\n")
		} else if first.UserRecordSecNonleafs != nil {
			//二级索引，非叶子节点
			pageNoStr := strconv.FormatInt(int64(first.PageNo), 10)
			builder.WriteString("############################################################\n")
			builder.WriteString(fmt.Sprintf("#               二级索引、非叶子page#%v%v#\n", pageNoStr, getSpaceStr(pageNoStr, 22)))
			builder.WriteString("#           格式：二级索引值，主键值，子page编号           #\n")
			builder.WriteString("############################################################\n")
			for _, userRecord := range first.UserRecordSecNonleafs {
				builder.WriteString(fmt.Sprintf("%v, %v, %v\n", userRecord.Secondary, userRecord.Primary, userRecord.ChildPageNo))
			}
			builder.WriteString("\n")
		} else if first.UserRecordSecLeafs != nil {
			//二级索引，叶子节点
			pageNoStr := strconv.FormatInt(int64(first.PageNo), 10)
			builder.WriteString("############################################################\n")
			builder.WriteString(fmt.Sprintf("#                 二级索引、叶子page#%v%v#\n", pageNoStr, getSpaceStr(pageNoStr, 22)))
			builder.WriteString("#                 格式：二级索引值，主键值                 #\n")
			builder.WriteString("############################################################\n")
			for _, userRecord := range first.UserRecordSecLeafs {
				builder.WriteString(fmt.Sprintf("%v, %v\n", userRecord.Secondary, userRecord.Primary))
			}
			builder.WriteString("\n")
		}
	}
	filename := fmt.Sprintf("./output/rec_%v", strings.Join(currIndex.FieldNames, "_"))
	_ = ioutil.WriteFile(filename, builder.Bytes(), 0666)
}

func GetBtreeRec(path string, pageNo uint64, fields []*proto.Field, currIndex, priIndex *proto.Index,
	tableCharset, inputCharset, rowFormat string) *BtreeRec {
	pageIndex := parser.ReadPageIndex(path, pageNo, fields, currIndex, priIndex, tableCharset, inputCharset, rowFormat)
	node := new(BtreeRec)
	node.PageNo = pageNo
	node.UserRecordPriNonleafs = pageIndex.UserRecordPriNonleafs
	node.UserRecordPriLeafs = pageIndex.UserRecordPriLeafs
	node.UserRecordSecNonleafs = pageIndex.UserRecordSecNonleafs
	node.UserRecordSecLeafs = pageIndex.UserRecordSecLeafs
	if pageIndex.UserRecordPriNonleafs != nil {
		//聚簇索引，非叶子节点
		for _, userRecord := range pageIndex.UserRecordPriNonleafs {
			childNode := GetBtreeRec(path, uint64(userRecord.ChildPageNo), fields, currIndex,
				priIndex, tableCharset, inputCharset, rowFormat)
			node.ChildPages = append(node.ChildPages, childNode)
		}
	} else if pageIndex.UserRecordPriLeafs != nil {
		//聚簇索引，叶子节点，无需处理
	} else if pageIndex.UserRecordSecNonleafs != nil {
		//二级索引，非叶子节点
		for _, userRecord := range pageIndex.UserRecordSecNonleafs {
			childNode := GetBtreeRec(path, uint64(userRecord.ChildPageNo), fields, currIndex,
				priIndex, tableCharset, inputCharset, rowFormat)
			node.ChildPages = append(node.ChildPages, childNode)
		}
	} else if pageIndex.UserRecordSecLeafs != nil {
		//二级索引，叶子节点，无需处理
	}
	return node
}

func getSpaceStr(pageNoStr string, total int) string {
	spaceNum := total - len(pageNoStr)
	if spaceNum < 0 {
		spaceNum = 0
	}
	var spaces []string
	for i:=0; i<spaceNum; i++ {
		spaces = append(spaces, " ")
	}
	return strings.Join(spaces, "")
}