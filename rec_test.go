package parseibd

import (
	"parseibd/parser"
	"parseibd/proto"
	"testing"
)

func TestRec1(t *testing.T) {
	/*
		CREATE TABLE `pri_sec_table_redundant` (
		  `id` int(11) NOT NULL,
		  `name` varchar(20) NOT NULL,
		  `age` int(11) DEFAULT NULL,
		  `city` varchar(10) DEFAULT NULL,
		  `addr` varchar(50) DEFAULT NULL,
		  PRIMARY KEY (`id`,`name`),
		  KEY `name` (`name`,`city`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=REDUNDANT
	*/

	fields := []*proto.Field{
		{
			Name:          "id",
			Mtype:         6,
			Prtype:        1283,
			Len:           4,
			DataType:      "int",
			NumPrecision:  10,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "name",
			Mtype:         12,
			Prtype:        2949391,
			Len:           80,
			DataType:      "varchar",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "age",
			Mtype:         6,
			Prtype:        1027,
			Len:           4,
			DataType:      "int",
			NumPrecision:  10,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "city",
			Mtype:         12,
			Prtype:        2949135,
			Len:           40,
			DataType:      "varchar",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "addr",
			Mtype:         12,
			Prtype:        2949135,
			Len:           200,
			DataType:      "varchar",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       3,
		PageNo:     3,
		FieldNames: []string{"id", "name"},
	}

	tableCharset := parser.CharsetUtf8mb4
	inputCharset := parser.CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\pri_sec_table_redundant.ibd"
	rowFormat := parser.RowFormatRedundant

	////聚簇索引，非叶子节点
	//ShowRecs(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)

	secIndex1 := &proto.Index{
		Type:       0,
		PageNo:     4,
		FieldNames: []string{"name", "city"},
	}

	//二级索引，非叶子节点
	ShowRecs(path, 4, fields, secIndex1, priIndex, tableCharset, inputCharset, rowFormat)
}