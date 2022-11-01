package parser

import (
	"fmt"
	"parseibd/proto"
	"strings"
	"time"
)

//聚簇索引-非叶子记录
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		dataOffset：实际数据偏移量
//		isComp：是否新格式
//		fields：字段信息列表
//		priIndex：主键索引
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：主键索引非叶子记录
func ReadUserRecordPriNonleaf(buff []byte, path string, dataOffset uint64, isComp bool, fields []*proto.Field,
	priIndex *proto.Index, tableCharset, inputCharset, rowFormat string) *proto.UserRecordPriNonleaf {
	userRecordPriNonleaf := new(proto.UserRecordPriNonleaf)
	headOffset := dataOffset
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[headOffset-RecNewExtraBytes:])
		headOffset -= extraOffset
		userRecordPriNonleaf.ExtraNew = extra
	} else {
		extra, extraOffset := ReadRecExtraOld(buff[headOffset-RecOldExtraBytes:])
		headOffset -= extraOffset
		userRecordPriNonleaf.ExtraOld = extra
	}
	//null标志
	if isComp {
		var nullable uint64
		for _, field := range fields {
			if IsNullable(field.Prtype) {
				nullable++
			}
		}
		byteCount := (nullable+7) >> 3
		headOffset -= byteCount
	} else {
		//旧格式没有单独的null标志
	}
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	//主键值
	if isComp {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			values, hf, df := readFieldValuesNew(buff, path, priIndex.FieldNames, fields, headOffset,
				dataOffset, isComp, tableCharset, inputCharset, rowFormat)
			userRecordPriNonleaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			userRecordPriNonleaf.Primary = append(userRecordPriNonleaf.Primary, rowId)
		}
	} else {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			values, hf, df := readFieldValuesOld(buff, path, priIndex.FieldNames, fields, headOffset, dataOffset,
				userRecordPriNonleaf.ExtraOld.OffsetBytes, true, inputCharset, rowFormat)
			userRecordPriNonleaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			headOffset -= uint64(userRecordPriNonleaf.ExtraOld.OffsetBytes)
			userRecordPriNonleaf.Primary = append(userRecordPriNonleaf.Primary, rowId)
		}
	}
	//子page no
	pageNo := MachReadFrom4(buff[dataOffset:])
	dataOffset += 4
	userRecordPriNonleaf.ChildPageNo = uint32(pageNo)
	return userRecordPriNonleaf
}

//聚簇索引-叶子记录
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		dataOffset：实际数据偏移量
//		isComp：是否新格式
//		fields：字段信息列表
//		priIndex：主键索引信息
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：主键索引叶子记录
func ReadUserRecordPriLeaf(buff []byte, path string, dataOffset uint64, isComp bool, fields []*proto.Field,
	priIndex *proto.Index, tableCharset, inputCharset, rowFormat string) *proto.UserRecordPriLeaf {
	userRecordPriLeaf := new(proto.UserRecordPriLeaf)
	headOffset := dataOffset
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[headOffset-RecNewExtraBytes:])
		headOffset -= extraOffset
		userRecordPriLeaf.ExtraNew = extra
	} else {
		extra, extraOffset := ReadRecExtraOld(buff[headOffset-RecOldExtraBytes:])
		headOffset -= extraOffset
		userRecordPriLeaf.ExtraOld = extra
	}
	//null标志
	fieldNull := make(map[string]bool)
	if isComp {
		//只有not null的字段才算到null标志中
		var nullNames []string
		var nullable uint64
		for _, field := range fields {
			if IsNullable(field.Prtype) {
				nullNames = append(nullNames, field.Name)
				nullable++
			}
		}
		byteCount := (nullable+7) >> 3
		if byteCount > 0 {
			for i, name := range nullNames {
				val := buff[headOffset-uint64(1+i/8)] >> (i%8)
				fieldNull[name] = (val & 0x01) != 0
			}
		}
		headOffset -= byteCount
	} else {
		//旧格式没有单独的null标志
	}
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	//字段值
	if isComp {
		//主键值
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			values, hf, df := readFieldValuesNew(buff, path, priIndex.FieldNames, fields, headOffset,
				dataOffset, isComp, tableCharset, inputCharset, rowFormat)
			userRecordPriLeaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			userRecordPriLeaf.Primary = append(userRecordPriLeaf.Primary, rowId)
		}
		//trx_id
		trxId := MachReadFrom6(buff[dataOffset:])
		dataOffset += 6
		userRecordPriLeaf.TrxId = trxId
		//roll_ptr
		rollPtr := MachReadFrom7(buff[dataOffset:])
		dataOffset += 7
		userRecordPriLeaf.RollPtr = rollPtr
		//其他字段，排除主键，排除null
		for _, field := range fields {
			if IsExist(priIndex.FieldNames, field.Name) {
				continue
			}
			if fieldNull[field.Name] {
				//null
				userRecordPriLeaf.Values = append(userRecordPriLeaf.Values, nil)
			} else {
				//not null
				values, hf, df := readFieldValuesNew(buff, path, []string{field.Name}, fields, headOffset,
					dataOffset, isComp, tableCharset, inputCharset, rowFormat)
				userRecordPriLeaf.Values = append(userRecordPriLeaf.Values, values...)
				headOffset = hf
				dataOffset = df
			}
		}
	} else {
		//主键值
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			values, hf, df := readFieldValuesOld(buff, path, priIndex.FieldNames, fields, headOffset, dataOffset,
				userRecordPriLeaf.ExtraOld.OffsetBytes, true, inputCharset, rowFormat)
			userRecordPriLeaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			headOffset -= uint64(userRecordPriLeaf.ExtraOld.OffsetBytes)
			userRecordPriLeaf.Primary = append(userRecordPriLeaf.Primary, rowId)
		}
		//trx_id
		trxId := MachReadFrom6(buff[dataOffset:])
		dataOffset += 6
		headOffset -= uint64(userRecordPriLeaf.ExtraOld.OffsetBytes)
		userRecordPriLeaf.TrxId = trxId
		//roll_ptr
		rollPtr := MachReadFrom7(buff[dataOffset:])
		dataOffset += 7
		headOffset -= uint64(userRecordPriLeaf.ExtraOld.OffsetBytes)
		userRecordPriLeaf.RollPtr = rollPtr
		//其他字段，排除主键，排除null
		for _, field := range fields {
			if IsExist(priIndex.FieldNames, field.Name) {
				continue
			}
			//readFieldValuesOld已处理null
			values, hf, df := readFieldValuesOld(buff, path, []string{field.Name}, fields, headOffset, dataOffset,
				userRecordPriLeaf.ExtraOld.OffsetBytes, false, inputCharset, rowFormat)
			userRecordPriLeaf.Values = append(userRecordPriLeaf.Values, values...)
			headOffset = hf
			dataOffset = df
		}
	}
	return userRecordPriLeaf
}

//二级索引-非叶子记录
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		dataOffset：实际数据偏移量
//		isComp：是否新格式
//		fields：字段信息列表
//		secIndex：二级索引
//		priIndex：主键索引
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：二级索引非叶子记录
func ReadUserRecordSecNonleaf(buff []byte, path string, dataOffset uint64, isComp bool, fields []*proto.Field,
	secIndex, priIndex *proto.Index, tableCharset, inputCharset, rowFormat string) *proto.UserRecordSecNonleaf {
	userRecordSecNonleaf := new(proto.UserRecordSecNonleaf)
	headOffset := dataOffset
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[headOffset-RecNewExtraBytes:])
		headOffset -= extraOffset
		userRecordSecNonleaf.ExtraNew = extra
	} else {
		extra, extraOffset := ReadRecExtraOld(buff[headOffset-RecOldExtraBytes:])
		headOffset -= extraOffset
		userRecordSecNonleaf.ExtraOld = extra
	}
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	//null标志
	fieldNull := make(map[string]bool)
	if isComp {
		//只有not null的字段才算到null标志中
		var nullNames []string
		var nullable uint64
		for _, name := range secIndex.FieldNames {
			field := fieldMap[name]
			if IsNullable(field.Prtype) {
				nullNames = append(nullNames, name)
				nullable++
			}
		}
		byteCount := (nullable+7) >> 3
		if byteCount > 0 {
			for i, name := range nullNames {
				val := buff[headOffset-uint64(1+i/8)] >> (i%8)
				fieldNull[name] = (val & 0x01) != 0
			}
		}
		headOffset -= byteCount
	} else {
		//旧格式没有单独的null标志
	}
	//二级索引值
	if isComp {
		for _, name := range secIndex.FieldNames {
			//排除null
			if fieldNull[name] {
				//null
				userRecordSecNonleaf.Secondary = append(userRecordSecNonleaf.Secondary, nil)
			} else {
				//not null
				values, hf, df := readFieldValuesNew(buff, path, []string{name}, fields, headOffset,
					dataOffset, isComp, tableCharset, inputCharset, rowFormat)
				userRecordSecNonleaf.Secondary = append(userRecordSecNonleaf.Secondary, values...)
				headOffset = hf
				dataOffset = df
			}
		}
	} else {
		//readFieldValuesOld已处理null
		values, hf, df := readFieldValuesOld(buff, path, secIndex.FieldNames, fields, headOffset, dataOffset,
			userRecordSecNonleaf.ExtraOld.OffsetBytes, true, inputCharset, rowFormat)
		userRecordSecNonleaf.Secondary = values
		headOffset = hf
		dataOffset = df
	}
	//主键值
	if isComp {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			diffNames := GetListDiff(priIndex.FieldNames, secIndex.FieldNames)
			values, hf, df := readFieldValuesNew(buff, path, diffNames, fields, headOffset,
				dataOffset, isComp, tableCharset, inputCharset, rowFormat)
			userRecordSecNonleaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			userRecordSecNonleaf.Primary = append(userRecordSecNonleaf.Primary, rowId)
		}
	} else {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			diffNames := GetListDiff(priIndex.FieldNames, secIndex.FieldNames)
			values, hf, df := readFieldValuesOld(buff, path, diffNames, fields, headOffset, dataOffset,
				userRecordSecNonleaf.ExtraOld.OffsetBytes, false, inputCharset, rowFormat)
			userRecordSecNonleaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			headOffset -= uint64(userRecordSecNonleaf.ExtraOld.OffsetBytes)
			userRecordSecNonleaf.Primary = append(userRecordSecNonleaf.Primary, rowId)
		}
	}
	//子page no
	pageNo := MachReadFrom4(buff[dataOffset:])
	dataOffset += 4
	userRecordSecNonleaf.ChildPageNo = uint32(pageNo)
	return userRecordSecNonleaf
}

//二级索引-叶子记录
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		dataOffset：实际数据偏移量
//		isComp：是否新格式
//		fields：字段信息列表
//		secIndex：二级索引
//		priIndex：主键索引
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：二级索引叶子记录
func ReadUserRecordSecLeaf(buff []byte, path string, dataOffset uint64, isComp bool, fields []*proto.Field,
	secIndex, priIndex *proto.Index, tableCharset, inputCharset, rowFormat string) *proto.UserRecordSecLeaf {
	userRecordSecLeaf := new(proto.UserRecordSecLeaf)
	headOffset := dataOffset
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[headOffset-RecNewExtraBytes:])
		headOffset -= extraOffset
		userRecordSecLeaf.ExtraNew = extra
	} else {
		extra, extraOffset := ReadRecExtraOld(buff[headOffset-RecOldExtraBytes:])
		headOffset -= extraOffset
		userRecordSecLeaf.ExtraOld = extra
	}
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	//null标志
	fieldNull := make(map[string]bool)
	if isComp {
		//只有not null的字段才算到null标志中
		var nullNames []string
		var nullable uint64
		for _, name := range secIndex.FieldNames {
			field := fieldMap[name]
			if IsNullable(field.Prtype) {
				nullNames = append(nullNames, name)
				nullable++
			}
		}
		byteCount := (nullable+7) >> 3
		if byteCount > 0 {
			for i, name := range nullNames {
				val := buff[headOffset-uint64(1+i/8)] >> (i%8)
				fieldNull[name] = (val & 0x01) != 0
			}
		}
		headOffset -= byteCount
	} else {
		//旧格式没有单独的null标志
	}
	//二级索引值
	if isComp {
		for _, name := range secIndex.FieldNames {
			//排除null
			if fieldNull[name] {
				//null
				userRecordSecLeaf.Secondary = append(userRecordSecLeaf.Secondary, nil)
			} else {
				//not null
				values, hf, df := readFieldValuesNew(buff, path, []string{name}, fields, headOffset,
					dataOffset, isComp, tableCharset, inputCharset, rowFormat)
				userRecordSecLeaf.Secondary = append(userRecordSecLeaf.Secondary, values...)
				headOffset = hf
				dataOffset = df
			}
		}
	} else {
		values, hf, df := readFieldValuesOld(buff, path, secIndex.FieldNames, fields, headOffset, dataOffset,
			userRecordSecLeaf.ExtraOld.OffsetBytes, true, inputCharset, rowFormat)
		userRecordSecLeaf.Secondary = values
		headOffset = hf
		dataOffset = df
	}
	//主键值
	if isComp {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			diffNames := GetListDiff(priIndex.FieldNames, secIndex.FieldNames)
			values, hf, df := readFieldValuesNew(buff, path, diffNames, fields, headOffset,
				dataOffset, isComp, tableCharset, inputCharset, rowFormat)
			userRecordSecLeaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			userRecordSecLeaf.Primary = append(userRecordSecLeaf.Primary, rowId)
		}
	} else {
		if (priIndex.Type & DictUnique) != 0 {
			//用户设置的主键
			diffNames := GetListDiff(priIndex.FieldNames, secIndex.FieldNames)
			values, hf, df := readFieldValuesOld(buff, path, diffNames, fields, headOffset, dataOffset,
				userRecordSecLeaf.ExtraOld.OffsetBytes, false, inputCharset, rowFormat)
			userRecordSecLeaf.Primary = values
			headOffset = hf
			dataOffset = df
		} else {
			//系统自动生成的主键row_id
			rowId := MachReadFrom6(buff[dataOffset:])
			dataOffset += 6
			headOffset -= uint64(userRecordSecLeaf.ExtraOld.OffsetBytes)
			userRecordSecLeaf.Primary = append(userRecordSecLeaf.Primary, rowId)
		}
	}
	return userRecordSecLeaf
}

//读取新格式指定字段值列表
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		fnames：待读取字段名列表
//		fields：所有字段信息列表
//		headOffset：头信息偏移量
//		dataOffset：数据信息偏移量
//		isComp：是否新格式
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：字段值列表、新的头信息偏移量、新的数据信息偏移量
func readFieldValuesNew(buff []byte, path string, fnames []string, fields []*proto.Field, headOffset, dataOffset uint64,
	isComp bool, tableCharset, inputCharset, rowFormat string) ([]interface{}, uint64, uint64) {
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	var values []interface{}
	for _, name := range fnames {
		field := fieldMap[name]
		switch field.Mtype {
		case DataSys:
			//TODO
			//定长
		case DataFixbinary:
			//TODO
			//定长
			if field.DataType == "decimal" {
				val := readDecimal(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "time" {
				val := readTime(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "datetime" {
				val := readDatetime(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "timestamp" {
				val := readTimestamp(buff, field, dataOffset)
				values = append(values, val)
			}
			dataOffset += uint64(field.Len)

		case DataInt:
			//定长
			if field.DataType == "year" {
				val := readYear(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "date" {
				val := readDate(buff, field, dataOffset)
				values = append(values, val)
			} else {
				val := readInt(buff, field, dataOffset)
				values = append(values, val)
			}
			dataOffset += uint64(field.Len)
		case DataFloat:
			//定长
			val := MachReadFloat(buff[dataOffset:])
			values = append(values, val)
			dataOffset += uint64(field.Len)
		case DataDouble:
			//定长
			val := MachReadDouble(buff[dataOffset:])
			values = append(values, val)
			dataOffset += uint64(field.Len)
		case DataPoint:
			//TODO
			//定长
		case DataChar:
			//定长
			dataLen := uint64(field.Len)
			values = append(values, string(buff[dataOffset:dataOffset+dataLen]))
			dataOffset += dataLen
		case DataMysql:
			//定长或变长
			if IsBinary(field.Prtype) || !isComp || GetMaxBytesPerChar(tableCharset)==1 {
				dataLen := uint64(field.Len)
				values = append(values, string(buff[dataOffset:dataOffset+dataLen]))
				dataOffset += dataLen
				break
			}
			fallthrough
		case DataVarchar, DataVarmysql, DataBlob:
			//变长类型
			var useExternal bool
			var dataLen uint64
			if IsBigCol(field) {
				//最大长度>255
				if (buff[headOffset-1] & 0x80) != 0 {
					//实际长度>127，使用两个字节
					//处理外部存储标志
					//高位在高地址，低位在低地址，小端序
					useExternal = (buff[headOffset-1] & 0x40) != 0
					//使用两个字节
					var bt []byte
					bt = append(bt, buff[headOffset-1] & 0x3F)
					bt = append(bt, buff[headOffset-2])
					dataLen = MachReadFrom2(bt)
					headOffset -= 2
				} else {
					//实际长度<=127，使用一个字节
					dataLen = MachReadFrom1(buff[headOffset-1:])
					headOffset -= 1
				}
			} else {
				//最大长度<=255，使用一个字节
				dataLen = MachReadFrom1(buff[headOffset-1:])
				headOffset -= 1
			}
			val := readString(buff, path, field, dataOffset, useExternal, dataLen, inputCharset, rowFormat)
			values = append(values, val)
			dataOffset += dataLen
		case DataBinary:
			//TODO
			//变长
		case DataDecimal:
			//TODO
			//变长
		case DataVarpoint:
			//TODO
			//变长
		case DataGeometry:
			//TODO
			//变长
		}
	}
	return values, headOffset, dataOffset
}

//读取旧格式指定字段值列表
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		fnames：待读取字段名列表
//		fields：所有字段信息列表
//		headOffset：头信息偏移量
//		dataOffset：数据信息偏移量
//		offsetBytes：字段长度偏移量字节数
//		firstOffset：fnames中是否包含第一个字段，第一个字段的前一偏移量为0
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：字段值列表、新的头信息偏移量、新的数据信息偏移量
func readFieldValuesOld(buff []byte, path string, fnames []string, fields []*proto.Field, headOffset, dataOffset uint64,
	offsetBytes uint8, firstOffset bool, inputCharset, rowFormat string) ([]interface{}, uint64, uint64) {
	fieldMap := make(map[string]*proto.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	var prevLen, currLen, diffLen uint64
	var values []interface{}
	for i, name := range fnames {
		//前一偏移量
		if firstOffset && i == 0 {
			prevLen = 0
		} else {
			if offsetBytes == 1 {
				//1字节长度
				prevLen = MachReadFrom1([]byte{buff[headOffset] & 0x7F})
			} else {
				//2字节长度
				prevLen = MachReadFrom2([]byte{buff[headOffset] & 0x3F, buff[headOffset+1]})
			}
		}
		//数据长度
		//当前偏移量长度，高位在低地址，低位在高地址，大端序
		var isNull bool
		var useExternal bool
		if offsetBytes == 1 {
			//1字节长度
			isNull = (buff[headOffset-1] & 0x80) != 0
			currLen = MachReadFrom1([]byte{buff[headOffset-1] & 0x7F})
			headOffset -= 1
		} else {
			//2字节长度
			useExternal = (buff[headOffset-2] & 0x40) != 0
			isNull = (buff[headOffset-2] & 0x80) != 0
			currLen = MachReadFrom2([]byte{buff[headOffset-2] & 0x3F, buff[headOffset-1]})
			headOffset -= 2
		}
		//数据长度
		diffLen = currLen - prevLen
		prevLen = currLen
		if isNull {
			dataOffset += diffLen
			values = append(values, nil)
			continue
		}
		field := fieldMap[name]
		switch field.Mtype {
		case DataSys:
			//TODO
			//定长
		case DataFixbinary:
			//TODO
			//定长
			if field.DataType == "decimal" {
				val := readDecimal(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "time" {
				val := readTime(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "datetime" {
				val := readDatetime(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "timestamp" {
				val := readTimestamp(buff, field, dataOffset)
				values = append(values, val)
			}

		case DataInt:
			//定长
			if field.DataType == "year" {
				val := readYear(buff, field, dataOffset)
				values = append(values, val)
			} else if field.DataType == "date" {
				val := readDate(buff, field, dataOffset)
				values = append(values, val)
			} else {
				val := readInt(buff, field, dataOffset)
				values = append(values, val)
			}
		case DataFloat:
			//定长
			val := MachReadFloat(buff[dataOffset:])
			values = append(values, val)
		case DataDouble:
			//定长
			val := MachReadDouble(buff[dataOffset:])
			values = append(values, val)
		case DataPoint:
			//TODO
			//定长
		case DataMysql, DataVarmysql, DataChar, DataVarchar, DataBlob:
			val := readString(buff, path, field, dataOffset, useExternal, diffLen, inputCharset, rowFormat)
			values = append(values, val)
		case DataBinary:
			//TODO
			//变长
		case DataDecimal:
			//TODO
			//变长
		case DataVarpoint:
			//TODO
			//变长
		case DataGeometry:
			//TODO
			//变长
		}
		dataOffset += diffLen
	}
	return values, headOffset, dataOffset
}

//读取指定decimal字段的数据值
//参数：
//		buff：数据缓存
//		field：decimal字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readDecimal(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	intBytes, fracBytes := GetDecimalBytes(field.NumPrecision, field.NumScale)
	//负数，全部取反
	isNeg := (buff[dataOffset] & 0x80) == 0
	var bt []byte
	for i:=uint32(0); i<intBytes+fracBytes; i++ {
		b := buff[dataOffset+uint64(i)]
		if isNeg {
			b = b ^ 0xFF
		}
		if i == 0 {
			bt = append(bt, b ^ 0x80)
		} else {
			bt = append(bt, b)
		}
	}
	//整数部分
	var ints []string
	for i:=int64(intBytes); i>=0; i-=4 {
		end := i
		var start int64
		if i-4 < 0 {
			start = 0
		} else {
			start = i - 4
		}
		v := MachReadFromN(bt[start:end], int(end-start))
		//整数部分最前面不需要补0
		if start == 0 {
			ints = append([]string{fmt.Sprintf("%d", v)}, ints...)
		} else {
			ints = append([]string{fmt.Sprintf("%09d", v)}, ints...)
		}
	}
	//小数部分
	var fracs []string
	for i:=intBytes; i<intBytes+fracBytes; i+=4 {
		start := uint64(i)
		var end uint64
		if i+4 > intBytes+fracBytes {
			end = uint64(intBytes+fracBytes)
		} else {
			end = uint64(i) + 4
		}
		v := MachReadFromN(bt[start:end], int(end-start))
		if end-start == 4 {
			fracs = append(fracs, fmt.Sprintf("%09d", v))
		} else {
			fracs = append(fracs, fmt.Sprintf(fmt.Sprintf("%%0%dd", field.NumScale%9), v))
		}
	}
	var val string
	if len(ints) > 0 {
		val = strings.Join(ints, "")
		val = strings.TrimLeft(val, "0")
	}
	if val == "" {
		val = "0"
	}
	if len(fracs) > 0 {
		val += "."
		val += strings.Join(fracs, "")
	}
	if isNeg {
		val = "-" + val
	}
	return val
}

//读取指定整数字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readInt(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	//固定长度类型
	if !IsUnsigned(field.Prtype) {
		//有符号，翻转符号位
		var bt []byte
		bt = append(bt, buff[dataOffset] ^ 0x80)
		bt = append(bt, buff[dataOffset+1:]...)
		val := MachReadFromN(bt, int(field.Len))
		return int64(val)
	} else {
		//无符号
		return MachReadFromN(buff[dataOffset:], int(field.Len))
	}
}

//读取指定year字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readYear(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	val := readInt(buff, field, dataOffset)
	if val.(uint64) != 0 {
		val = val.(uint64) + 1900
	}
	return val
}

//读取指定time字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readTime(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	format := "%02d:%02d:%02d"
	var ret string
	switch field.TimePrecision {
	case 1, 2:
		intpart := int64(MachReadFrom3(buff[dataOffset:])) - 0x800000
		frac := int64(MachReadFrom1(buff[dataOffset+3:]))
		if intpart < 0 && frac != 0 {
			intpart++
			frac -= 0x100
		}
		val := (intpart << 24) + frac*10000
		isNeg := val < 0
		if isNeg {
			if intpart < 0 {
				intpart = -intpart
			}
			if frac < 0 {
				frac = -frac
			}
		}
		hour := (intpart >> 12) & 0x1F
		minute := (intpart >> 6) & 0x3F
		second := intpart & 0x3F
		format += ".%06d"
		if isNeg {
			format = "-" + format
		}
		ret = fmt.Sprintf(format, hour, minute, second, frac)
	case 3, 4:
		intpart := int64(MachReadFrom3(buff[dataOffset:])) - 0x800000
		frac := int64(MachReadFrom2(buff[dataOffset+3:]))
		if intpart < 0 && frac != 0 {
			intpart++
			frac -= 0x10000
		}
		val := (intpart << 24) + frac*100
		isNeg := val < 0
		if isNeg {
			if intpart < 0 {
				intpart = -intpart
			}
			if frac < 0 {
				frac = -frac
			}
		}
		hour := (intpart >> 12) & 0x1F
		minute := (intpart >> 6) & 0x3F
		second := intpart & 0x3F
		format += ".%06d"
		if isNeg {
			format = "-" + format
		}
		ret = fmt.Sprintf(format, hour, minute, second, frac)
	case 5, 6:
		val := int64(MachReadFrom6(buff[dataOffset:])) - (1 << 47)
		isNeg := val < 0
		if isNeg {
			val = -val
		}
		hour := (val >> 36) & 0x1F
		minute := (val >> 30) & 0x3F
		second := (val >> 24) & 0x3F
		frac := val & 0xFFFFFF
		format += ".%06d"
		if isNeg {
			format = "-" + format
		}
		ret = fmt.Sprintf(format, hour, minute, second, frac)
	default:
		val := int64(MachReadFrom3(buff[dataOffset:])) - 0x800000
		isNeg := val < 0
		if isNeg {
			val = -val
		}
		hour := (val >> 12) & 0x1F
		minute := (val >> 6) & 0x3F
		second := val & 0x3F
		if isNeg {
			format = "-" + format
		}
		ret = fmt.Sprintf(format, hour, minute, second)
	}
	index := strings.LastIndex(ret, ".")
	if index != -1 && uint32(index)+1+field.TimePrecision < uint32(len(ret)) {
		return ret[0:uint32(index)+1+field.TimePrecision]
	} else {
		return ret
	}
}

//读取指定datetime字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readDatetime(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	intpart := int64(MachReadFrom5(buff[dataOffset:])) - (1 << 39)
	ym := (intpart >> 22) & 0x1FFFF
	year := ym/13
	month := ym%13
	day := (intpart >> 17) & 0x1F
	hour := (intpart >> 12) & 0x1F
	minute := (intpart >> 6) & 0x3F
	second := intpart & 0x3F
	format := "%04d-%02d-%02d %02d:%02d:%02d"
	var ret string
	switch field.TimePrecision {
	case 1, 2:
		frac := int64(MachReadFrom1(buff[dataOffset+5:])) * 10000
		format += ".%06d"
		ret = fmt.Sprintf(format, year, month, day, hour, minute, second, frac)
	case 3, 4:
		frac := int64(MachReadFrom2(buff[dataOffset+5:])) * 100
		format += ".%06d"
		ret = fmt.Sprintf(format, year, month, day, hour, minute, second, frac)
	case 5, 6:
		frac := int64(MachReadFrom3(buff[dataOffset+5:]))
		format += ".%06d"
		ret = fmt.Sprintf(format, year, month, day, hour, minute, second, frac)
	default:
		ret = fmt.Sprintf(format, year, month, day, hour, minute, second)
	}
	index := strings.LastIndex(ret, ".")
	if index != -1 && uint32(index)+1+field.TimePrecision < uint32(len(ret)) {
		return ret[0:uint32(index)+1+field.TimePrecision]
	} else {
		return ret
	}
}

//读取指定timestamp字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readTimestamp(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	sec := MachReadFrom4(buff[dataOffset:])
	ret := time.Unix(int64(sec), 0).Format("2006-01-02 15:04:05")
	switch field.TimePrecision {
	case 1, 2:
		usec := int64(MachReadFrom1(buff[dataOffset+4:])) * 10000
		ret += fmt.Sprintf(".%06d", usec)
	case 3, 4:
		usec := int64(MachReadFrom2(buff[dataOffset+4:])) * 100
		ret += fmt.Sprintf(".%06d", usec)
	case 5, 6:
		usec := int64(MachReadFrom3(buff[dataOffset+4:]))
		ret += fmt.Sprintf(".%06d", usec)
	}
	index := strings.LastIndex(ret, ".")
	if index != -1 && uint32(index)+1+field.TimePrecision < uint32(len(ret)) {
		return ret[0:uint32(index)+1+field.TimePrecision]
	} else {
		return ret
	}
}

//读取指定date字段的数据值
//参数：
//		buff：数据缓存
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//返回：字段值
func readDate(buff []byte, field *proto.Field, dataOffset uint64) interface{} {
	val := MachReadFrom3(buff[dataOffset:])
	year := (val >> 9) & 0x3FFF
	month := (val >> 5) & 0x0F
	day := val & 0x1F
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

//读取指定字符串字段的数据值
//参数：
//		buff：数据缓存
//		path：ibd文件路径
//		field：整数字段信息
//		dataOffset：数据信息偏移量
//		useExternal：是否使用外部页存储
//		dataLen：该字段值在当前索引页中长度
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
//返回：字段值
func readString(buff []byte, path string, field *proto.Field, dataOffset uint64, useExternal bool,
	dataLen uint64, inputCharset, rowFormat string) interface{} {
	var bt []byte
	if useExternal {
		//使用外部页存储
		var externPtr *proto.ExternPtr
		switch strings.ToLower(rowFormat) {
		case RowFormatDynamic, RowFormatCompressed:
			Assert(dataLen == 20)
			externPtr, _ = ReadExternPtr(buff[dataOffset:])
		case RowFormatCompact, RowFormatRedundant:
			Assert(dataLen == RecAntelopeMaxIndexColLen + 20)
			bt = append(bt, buff[dataOffset:dataOffset+RecAntelopeMaxIndexColLen]...)
			externPtr, _ = ReadExternPtr(buff[dataOffset+RecAntelopeMaxIndexColLen:])
		}
		pageNo := externPtr.PageNo
		for pageNo != FilNull {
			pageBlob := ReadPageBlob(path, uint64(pageNo))
			bt = append(bt, pageBlob.PayLoad...)
			pageNo = pageBlob.BlobHeader.Next
		}
	} else {
		//未使用外部页存储
		bt = append(bt, buff[dataOffset:dataOffset+dataLen]...)
	}
	var val string
	if strings.Contains(field.DataType, "blob") {
		switch inputCharset {
		case CharsetGbk:
			if tmp, e := GbkToUtf8(bt); e == nil {
				val = string(tmp)
			}
		default:
			val = string(bt)
		}
	} else {
		val = string(bt)
	}
	return val
}