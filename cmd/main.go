package main

import (
	"flag"
	"fmt"
	"os"
	"parseibd"
	"parseibd/parser"
	"parseibd/proto"
	"strings"
)

var host = flag.String("h", "127.0.0.1:3306", "MySQL地址，格式：127.0.0.1:3306")
var user = flag.String("u", "", "MySQL用户名")
var password = flag.String("p", "", "MySQL密码")
var dbName = flag.String("d", "", "MySQL库名")
var tbName = flag.String("t", "", "MySQL表名")
var path = flag.String("f", "", "待解析ibd文件路径")
var inputCharset = flag.String("c", parser.CharsetUtf8, "MySQL数据录入时的字符编码（非表的字符编码）")

func main() {
	flag.Parse()
	if *host == "" || *dbName == "" || *tbName == "" || *path == "" {
		flag.Usage()
		return
	}
	//fmt.Println(*host, *user, *password, *dbName, *tbName, *path, *inputCharset)
	session, e := GetSession(*host, *user, *password, *dbName)
	if e != nil {
		fmt.Println("get session error:", e)
		return
	}
	defer session.Close()
	tableInfo, e := GetTableInfo(session, *dbName, *tbName)
	if e != nil {
		fmt.Println("get table info error:", e)
		return
	}
	fields, e := GetFields(session, *dbName, *tbName)
	if e != nil {
		fmt.Println("get field info error:", e)
		return
	}
	priIndex, e := GetPriIndex(session, *dbName, *tbName)
	if e != nil {
		fmt.Println("get primary index error:", e)
		return
	}
	secIndexes, e := GetSecIndexes(session, *dbName, *tbName)
	if e != nil {
		fmt.Println("get secondary index error:", e)
		return
	}
	if e := os.MkdirAll("./output", 0777); e != nil {
		fmt.Println("create output directory error:", e)
		return
	}
	var indexes []*proto.Index
	indexes = append(indexes, priIndex)
	indexes = append(indexes, secIndexes...)
	for _, currIndex := range indexes {
		parseibd.ShowBtree(*path, uint64(currIndex.PageNo), fields, currIndex, priIndex,
			tableInfo.Charset, strings.ToLower(*inputCharset), strings.ToLower(tableInfo.RowFormat))
		parseibd.ShowRecs(*path, uint64(currIndex.PageNo), fields, currIndex, priIndex,
			tableInfo.Charset, strings.ToLower(*inputCharset), strings.ToLower(tableInfo.RowFormat))
	}
	parseibd.ShowExtents(*path)
	parseibd.ShowInodes(*path)
	parseibd.ShowPages(*path)
	parseibd.ShowSegments(*path)
}