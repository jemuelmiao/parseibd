package parser

import (
	"testing"
	"parseibd/proto"
)

func TestIndex1(t *testing.T) {
	/*
	CREATE TABLE `test_redundant` (
		`id` int(11) NOT NULL,
		`name` varchar(100) DEFAULT NULL,
		`addr` char(100) DEFAULT NULL,
		PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1 ROW_FORMAT=REDUNDANT
	 */

	fields := []*proto.Field{
		{
			Name:   "id",
			Mtype:  6,
			Prtype: 1283,
			Len:    4,
		},
		{
			Name:   "name",
			Mtype:  1,
			Prtype: 524303,
			Len:    100,
		},
		{
			Name:   "addr",
			Mtype:  2,
			Prtype: 524542,
			Len:    100,
		},
	}
	
	priIndex := &proto.Index{
		Type:       3,
		PageNo:     3,
		FieldNames: []string{"id"},
	}

	tableCharset := CharsetLatin1
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_redundant.ibd"
	rowFormat := RowFormatRedundant

	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex2(t *testing.T) {
	/*
		CREATE TABLE `test_dynamic` (
			`id` int(11) NOT NULL,
			`name` varchar(100) DEFAULT NULL,
			`addr` char(100) DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1
	*/

	fields := []*proto.Field{
		{
			Name:   "id",
			Mtype:  6,
			Prtype: 1283,
			Len:    4,
		},
		{
			Name:   "name",
			Mtype:  1,
			Prtype: 524303,
			Len:    100,
		},
		{
			Name:   "addr",
			Mtype:  2,
			Prtype: 524542,
			Len:    100,
		},
	}

	priIndex := &proto.Index{
		Type:       3,
		PageNo:     3,
		FieldNames: []string{"id"},
	}
	
	tableCharset := CharsetLatin1
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_dynamic.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex3(t *testing.T) {
	/*
		CREATE TABLE `test_dynamic1` (
			`id` int(11) NOT NULL,
			`name` varchar(100) DEFAULT NULL,
			`addr` char(100) DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8
	*/

	fields := []*proto.Field{
		{
			Name:   "id",
			Mtype:  6,
			Prtype: 1283,
			Len:    4,
		},
		{
			Name:   "name",
			Mtype:  12,
			Prtype: 2166799,
			Len:    300,
		},
		{
			Name:   "addr",
			Mtype:  13,
			Prtype: 2162942,
			Len:    300,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_dynamic1.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex4(t *testing.T) {
	/*
	CREATE TABLE `test_type_dynamic` (
		`f1` tinyint(4) DEFAULT NULL,
		`f2` smallint(6) DEFAULT NULL,
		`f3` mediumint(9) DEFAULT NULL,
		`f4` int(11) DEFAULT NULL,
		`f5` bigint(20) DEFAULT NULL,
		`f6` float DEFAULT NULL,
		`f7` double DEFAULT NULL,
		`f8` decimal(12,3) DEFAULT NULL,
		`f9` date DEFAULT NULL,
		`f10` time DEFAULT NULL,
		`f11` year(4) DEFAULT NULL,
		`f12` datetime DEFAULT NULL,
		`f13` timestamp NULL DEFAULT NULL,
		`f14` char(100) DEFAULT NULL,
		`f15` varchar(100) DEFAULT NULL,
		`f16` tinyblob,
		`f17` tinytext,
		`f18` blob,
		`f19` text,
		`f20` mediumblob,
		`f21` mediumtext,
		`f22` longblob,
		`f23` longtext
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	 */

	fields := []*proto.Field{
		{
			Name:   "f1",
			Mtype:  6,
			Prtype: 1025,
			Len:    1,
			DataType: "tinyint",
			NumPrecision: 3,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f2",
			Mtype:  6,
			Prtype: 1026,
			Len:    2,
			DataType: "smallint",
			NumPrecision: 5,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f3",
			Mtype:  6,
			Prtype: 1033,
			Len:    3,
			DataType: "mediumint",
			NumPrecision: 7,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f4",
			Mtype:  6,
			Prtype: 1027,
			Len:    4,
			DataType: "int",
			NumPrecision: 10,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f5",
			Mtype:  6,
			Prtype: 1032,
			Len:    8,
			DataType: "bigint",
			NumPrecision: 19,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f6",
			Mtype:  9,
			Prtype: 1028,
			Len:    4,
			DataType: "float",
			NumPrecision: 12,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f7",
			Mtype:  10,
			Prtype: 1029,
			Len:    8,
			DataType: "double",
			NumPrecision: 22,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f8",
			Mtype:  3,
			Prtype: 525558,
			Len:    6,
			DataType: "decimal",
			NumPrecision: 12,
			NumScale: 3,
			TimePrecision: 0,
		},
		{
			Name:   "f9",
			Mtype:  6,
			Prtype: 1034,
			Len:    3,
			DataType: "date",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f10",
			Mtype:  3,
			Prtype: 525323,
			Len:    3,
			DataType: "time",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f11",
			Mtype:  6,
			Prtype: 1549,
			Len:    1,
			DataType: "year",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f12",
			Mtype:  3,
			Prtype: 525324,
			Len:    5,
			DataType: "datetime",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f13",
			Mtype:  3,
			Prtype: 525319,
			Len:    4,
			DataType: "timestamp",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f14",
			Mtype:  13,
			Prtype: 2162942,
			Len:    300,
			DataType: "char",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f15",
			Mtype:  12,
			Prtype: 2166799,
			Len:    300,
			DataType: "varchar",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f16",
			Mtype:  5,
			Prtype: 4130044,
			Len:    9,
			DataType: "tinyblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f17",
			Mtype:  5,
			Prtype: 2162940,
			Len:    9,
			DataType: "tinytext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f18",
			Mtype:  5,
			Prtype: 4130044,
			Len:    10,
			DataType: "blob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f19",
			Mtype:  5,
			Prtype: 2162940,
			Len:    10,
			DataType: "text",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f20",
			Mtype:  5,
			Prtype: 4130044,
			Len:    11,
			DataType: "mediumblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f21",
			Mtype:  5,
			Prtype: 2162940,
			Len:    11,
			DataType: "mediumtext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f22",
			Mtype:  5,
			Prtype: 4130044,
			Len:    12,
			DataType: "longblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f23",
			Mtype:  5,
			Prtype: 2162940,
			Len:    12,
			DataType: "longtext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	rowFormat := RowFormatDynamic
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_type_dynamic.ibd"
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex5(t *testing.T) {
	/*
	CREATE TABLE `test_type_redundant` (
		`f1` tinyint(4) DEFAULT NULL,
		`f2` smallint(6) DEFAULT NULL,
		`f3` mediumint(9) DEFAULT NULL,
		`f4` int(11) DEFAULT NULL,
		`f5` bigint(20) DEFAULT NULL,
		`f6` float DEFAULT NULL,
		`f7` double DEFAULT NULL,
		`f8` decimal(12,3) DEFAULT NULL,
		`f9` date DEFAULT NULL,
		`f10` time DEFAULT NULL,
		`f11` year(4) DEFAULT NULL,
		`f12` datetime DEFAULT NULL,
		`f13` timestamp NULL DEFAULT NULL,
		`f14` char(100) DEFAULT NULL,
		`f15` varchar(100) DEFAULT NULL,
		`f16` tinyblob,
		`f17` tinytext,
		`f18` blob,
		`f19` text,
		`f20` mediumblob,
		`f21` mediumtext,
		`f22` longblob,
		`f23` longtext
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=REDUNDANT
	 */

	fields := []*proto.Field{
		{
			Name:   "f1",
			Mtype:  6,
			Prtype: 1025,
			Len:    1,
			DataType: "tinyint",
			NumPrecision: 3,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f2",
			Mtype:  6,
			Prtype: 1026,
			Len:    2,
			DataType: "smallint",
			NumPrecision: 5,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f3",
			Mtype:  6,
			Prtype: 1033,
			Len:    3,
			DataType: "mediumint",
			NumPrecision: 7,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f4",
			Mtype:  6,
			Prtype: 1027,
			Len:    4,
			DataType: "int",
			NumPrecision: 10,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f5",
			Mtype:  6,
			Prtype: 1032,
			Len:    8,
			DataType: "bigint",
			NumPrecision: 19,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f6",
			Mtype:  9,
			Prtype: 1028,
			Len:    4,
			DataType: "float",
			NumPrecision: 12,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f7",
			Mtype:  10,
			Prtype: 1029,
			Len:    8,
			DataType: "double",
			NumPrecision: 22,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f8",
			Mtype:  3,
			Prtype: 525558,
			Len:    6,
			DataType: "decimal",
			NumPrecision: 12,
			NumScale: 3,
			TimePrecision: 0,
		},
		{
			Name:   "f9",
			Mtype:  6,
			Prtype: 1034,
			Len:    3,
			DataType: "date",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f10",
			Mtype:  3,
			Prtype: 525323,
			Len:    3,
			DataType: "time",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f11",
			Mtype:  6,
			Prtype: 1549,
			Len:    1,
			DataType: "year",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f12",
			Mtype:  3,
			Prtype: 525324,
			Len:    5,
			DataType: "datetime",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f13",
			Mtype:  3,
			Prtype: 525319,
			Len:    4,
			DataType: "timestamp",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f14",
			Mtype:  13,
			Prtype: 2162942,
			Len:    300,
			DataType: "char",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f15",
			Mtype:  12,
			Prtype: 2166799,
			Len:    300,
			DataType: "varchar",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f16",
			Mtype:  5,
			Prtype: 4130044,
			Len:    9,
			DataType: "tinyblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f17",
			Mtype:  5,
			Prtype: 2162940,
			Len:    9,
			DataType: "tinytext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f18",
			Mtype:  5,
			Prtype: 4130044,
			Len:    10,
			DataType: "blob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f19",
			Mtype:  5,
			Prtype: 2162940,
			Len:    10,
			DataType: "text",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f20",
			Mtype:  5,
			Prtype: 4130044,
			Len:    11,
			DataType: "mediumblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f21",
			Mtype:  5,
			Prtype: 2162940,
			Len:    11,
			DataType: "mediumtext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f22",
			Mtype:  5,
			Prtype: 4130044,
			Len:    12,
			DataType: "longblob",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
		{
			Name:   "f23",
			Mtype:  5,
			Prtype: 2162940,
			Len:    12,
			DataType: "longtext",
			NumPrecision: 0,
			NumScale: 0,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_type_redundant.ibd"
	rowFormat := RowFormatRedundant
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex6(t *testing.T) {
	/*
	CREATE TABLE `test_float` (
		`id` int(11) DEFAULT NULL,
		`name` varchar(100) DEFAULT NULL,
		`price1` float DEFAULT NULL,
		`price2` double DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	 */

	fields := []*proto.Field{
		{
			Name:   "id",
			Mtype:  6,
			Prtype: 1027,
			Len:    4,
		},
		{
			Name:   "name",
			Mtype:  12,
			Prtype: 2166799,
			Len:    300,
		},
		{
			Name:   "price1",
			Mtype:  9,
			Prtype: 1028,
			Len:    4,
		},
		{
			Name:   "price2",
			Mtype:  10,
			Prtype: 1029,
			Len:    8,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_float.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex7(t *testing.T) {
	/*
	CREATE TABLE `test_decimal` (
	  `f1` decimal(10,0) DEFAULT NULL,
	  `f2` decimal(13,3) DEFAULT NULL,
	  `f3` decimal(12,0) DEFAULT NULL,
	  `f4` decimal(6,6) DEFAULT NULL,
	  `f5` decimal(12,12) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	 */

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525558,
			Len:           5,
			DataType:      "decimal",
			NumPrecision:  10,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         3,
			Prtype:        525558,
			Len:           7,
			DataType:      "decimal",
			NumPrecision:  13,
			NumScale:      3,
			TimePrecision: 0,
		},
		{
			Name:          "f3",
			Mtype:         3,
			Prtype:        525558,
			Len:           6,
			DataType:      "decimal",
			NumPrecision:  12,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f4",
			Mtype:         3,
			Prtype:        525558,
			Len:           3,
			DataType:      "decimal",
			NumPrecision:  6,
			NumScale:      6,
			TimePrecision: 0,
		},
		{
			Name:          "f5",
			Mtype:         3,
			Prtype:        525558,
			Len:           6,
			DataType:      "decimal",
			NumPrecision:  12,
			NumScale:      12,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_decimal.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex8(t *testing.T) {
	/*
	CREATE TABLE `test_decimal_redundant` (
	  `f1` decimal(10,0) DEFAULT NULL,
	  `f2` decimal(13,3) DEFAULT NULL,
	  `f3` decimal(12,0) DEFAULT NULL,
	  `f4` decimal(6,6) DEFAULT NULL,
	  `f5` decimal(12,12) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=REDUNDANT
	 */

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525558,
			Len:           5,
			DataType:      "decimal",
			NumPrecision:  10,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         3,
			Prtype:        525558,
			Len:           7,
			DataType:      "decimal",
			NumPrecision:  13,
			NumScale:      3,
			TimePrecision: 0,
		},
		{
			Name:          "f3",
			Mtype:         3,
			Prtype:        525558,
			Len:           6,
			DataType:      "decimal",
			NumPrecision:  12,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f4",
			Mtype:         3,
			Prtype:        525558,
			Len:           3,
			DataType:      "decimal",
			NumPrecision:  6,
			NumScale:      6,
			TimePrecision: 0,
		},
		{
			Name:          "f5",
			Mtype:         3,
			Prtype:        525558,
			Len:           6,
			DataType:      "decimal",
			NumPrecision:  12,
			NumScale:      12,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_decimal_redundant.ibd"
	rowFormat := RowFormatRedundant
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex9(t *testing.T) {
	/*
	CREATE TABLE `test_decimal2` (
	  `f1` decimal(14,4) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525558,
			Len:           7,
			DataType:      "decimal",
			NumPrecision:  14,
			NumScale:      4,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_decimal2.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex10(t *testing.T) {
	/*
	CREATE TABLE `test_decimal3` (
	  `f1` decimal(25,12) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525558,
			Len:           12,
			DataType:      "decimal",
			NumPrecision:  25,
			NumScale:      12,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_decimal3.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex11(t *testing.T) {
	/*
	CREATE TABLE `test_year` (
	  `f1` year(4) DEFAULT NULL,
	  `f2` year(4) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         6,
			Prtype:        1549,
			Len:           1,
			DataType:      "year",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         6,
			Prtype:        1549,
			Len:           1,
			DataType:      "year",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_year.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex12(t *testing.T) {
	/*
	CREATE TABLE `test_time` (
	  `f1` time DEFAULT NULL,
	  `f2` time(1) DEFAULT NULL,
	  `f3` time(2) DEFAULT NULL,
	  `f4` time(3) DEFAULT NULL,
	  `f5` time(4) DEFAULT NULL,
	  `f6` time(5) DEFAULT NULL,
	  `f7` time(6) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525323,
			Len:           3,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         3,
			Prtype:        525323,
			Len:           4,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 1,
		},
		{
			Name:          "f3",
			Mtype:         3,
			Prtype:        525323,
			Len:           4,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 2,
		},
		{
			Name:          "f4",
			Mtype:         3,
			Prtype:        525323,
			Len:           5,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 3,
		},
		{
			Name:          "f5",
			Mtype:         3,
			Prtype:        525323,
			Len:           5,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 4,
		},
		{
			Name:          "f6",
			Mtype:         3,
			Prtype:        525323,
			Len:           6,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 5,
		},
		{
			Name:          "f7",
			Mtype:         3,
			Prtype:        525323,
			Len:           6,
			DataType:      "time",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 6,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_time.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex13(t *testing.T) {
	/*
	CREATE TABLE `test_datetime` (
	  `f1` datetime DEFAULT NULL,
	  `f2` datetime(1) DEFAULT NULL,
	  `f3` datetime(2) DEFAULT NULL,
	  `f4` datetime(3) DEFAULT NULL,
	  `f5` datetime(4) DEFAULT NULL,
	  `f6` datetime(5) DEFAULT NULL,
	  `f7` datetime(6) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525324,
			Len:           5,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         3,
			Prtype:        525324,
			Len:           6,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 1,
		},
		{
			Name:          "f3",
			Mtype:         3,
			Prtype:        525324,
			Len:           6,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 2,
		},
		{
			Name:          "f4",
			Mtype:         3,
			Prtype:        525324,
			Len:           7,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 3,
		},
		{
			Name:          "f5",
			Mtype:         3,
			Prtype:        525324,
			Len:           7,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 4,
		},
		{
			Name:          "f6",
			Mtype:         3,
			Prtype:        525324,
			Len:           8,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 5,
		},
		{
			Name:          "f7",
			Mtype:         3,
			Prtype:        525324,
			Len:           8,
			DataType:      "datetime",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 6,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_datetime.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex14(t *testing.T) {
	/*
	CREATE TABLE `test_timestamp` (
	  `f1` timestamp NULL DEFAULT NULL,
	  `f2` timestamp(1) NULL DEFAULT NULL,
	  `f3` timestamp(2) NULL DEFAULT NULL,
	  `f4` timestamp(3) NULL DEFAULT NULL,
	  `f5` timestamp(4) NULL DEFAULT NULL,
	  `f6` timestamp(5) NULL DEFAULT NULL,
	  `f7` timestamp(6) NULL DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         3,
			Prtype:        525319,
			Len:           4,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
		{
			Name:          "f2",
			Mtype:         3,
			Prtype:        525319,
			Len:           5,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 1,
		},
		{
			Name:          "f3",
			Mtype:         3,
			Prtype:        525319,
			Len:           5,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 2,
		},
		{
			Name:          "f4",
			Mtype:         3,
			Prtype:        525319,
			Len:           6,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 3,
		},
		{
			Name:          "f5",
			Mtype:         3,
			Prtype:        525319,
			Len:           6,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 4,
		},
		{
			Name:          "f6",
			Mtype:         3,
			Prtype:        525319,
			Len:           7,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 5,
		},
		{
			Name:          "f7",
			Mtype:         3,
			Prtype:        525319,
			Len:           7,
			DataType:      "timestamp",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 6,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_timestamp.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex15(t *testing.T) {
	/*
	CREATE TABLE `test_date` (
	  `f1` date DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC
	*/

	fields := []*proto.Field{
		{
			Name:          "f1",
			Mtype:         6,
			Prtype:        1034,
			Len:           3,
			DataType:      "date",
			NumPrecision:  0,
			NumScale:      0,
			TimePrecision: 0,
		},
	}

	priIndex := &proto.Index{
		Type:       1,
		PageNo:     3,
		FieldNames: nil,
	}

	tableCharset := CharsetUtf8
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd"
	rowFormat := RowFormatDynamic
	pageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(pageIndex)
}

func TestIndex16(t *testing.T) {
	/*
	CREATE TABLE `pri_sec_table` (
	  `id` int(11) NOT NULL,
	  `name` varchar(20) NOT NULL,
	  `age` int(11) DEFAULT NULL,
	  `city` varchar(10) DEFAULT NULL,
	  `addr` varchar(50) DEFAULT NULL,
	  PRIMARY KEY (`id`,`name`),
	  KEY `name` (`name`,`city`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
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

	tableCharset := CharsetUtf8mb4
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\pri_sec_table.ibd"
	rowFormat := RowFormatDynamic

	////聚簇索引，非叶子节点
	//priPageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(priPageIndex)

	////聚簇索引，叶子节点
	//priPageIndex := ReadPageIndex(path, 6, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(priPageIndex)

	secIndex1 := &proto.Index{
		Type:       0,
		PageNo:     4,
		FieldNames: []string{"name", "city"},
	}

	////二级索引，非叶子节点
	//secPageIndex := ReadPageIndex(path, 4, fields, secIndex1, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(secPageIndex)

	//二级索引，叶子节点
	secPageIndex := ReadPageIndex(path, 9, fields, secIndex1, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(secPageIndex)
}

func TestIndex17(t *testing.T) {
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

	tableCharset := CharsetUtf8mb4
	inputCharset := CharsetGbk
	path := "E:\\bdc\\mysql-server\\data\\jemuel\\pri_sec_table_redundant.ibd"
	rowFormat := RowFormatRedundant

	////聚簇索引，非叶子节点
	//priPageIndex := ReadPageIndex(path, 3, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(priPageIndex)

	////聚簇索引，叶子节点
	//priPageIndex := ReadPageIndex(path, 6, fields, priIndex, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(priPageIndex)

	secIndex1 := &proto.Index{
		Type:       0,
		PageNo:     4,
		FieldNames: []string{"name", "city"},
	}

	////二级索引，非叶子节点
	//secPageIndex := ReadPageIndex(path, 4, fields, secIndex1, priIndex, tableCharset, inputCharset, rowFormat)
	//PrintPageIndex(secPageIndex)

	//二级索引，叶子节点
	secPageIndex := ReadPageIndex(path, 9, fields, secIndex1, priIndex, tableCharset, inputCharset, rowFormat)
	PrintPageIndex(secPageIndex)
}