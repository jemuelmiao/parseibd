package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"parseibd/proto"
	"strings"
)

type TableInfo struct {
	RowFormat	string
	Charset 	string
}

func GetSession(host, user, password, dbName string) (*sql.DB, error) {
	var dsn string
	if user == "" {
		dsn = fmt.Sprintf("tcp(%v)/", host)
	} else if password == "" {
		dsn = fmt.Sprintf("%v@tcp(%v)/", user, host)
	} else {
		dsn = fmt.Sprintf("%v:%v@tcp(%v)/", user, password, host)
	}
	if dbName != "" {
		dsn += dbName
	}
	dsn += "?readTimeout=10s"
	return sql.Open("mysql", dsn)
}

func GetTableInfo(session *sql.DB, dbName, tbName string) (*TableInfo, error) {
	cmd := fmt.Sprintf(`select t.ROW_FORMAT, c.CHARACTER_SET_NAME
from information_schema.TABLES t join information_schema.COLLATION_CHARACTER_SET_APPLICABILITY c
on t.TABLE_COLLATION=c.COLLATION_NAME 
where t.TABLE_SCHEMA='%v' and t.TABLE_NAME='%v'`, dbName, tbName)
	rows, e := session.Query(cmd)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	for rows.Next() {
		var rowFormat, charset string
		if e := rows.Scan(&rowFormat, &charset); e != nil {
			return nil, e
		}
		tableInfo := &TableInfo{
			RowFormat: strings.ToLower(rowFormat),
			Charset:   strings.ToLower(charset),
		}
		return tableInfo, nil
	}
	return nil, errors.New("not find")
}

func GetFields(session *sql.DB, dbName, tbName string) ([]*proto.Field, error) {
	cmd := fmt.Sprintf(`select sc.NAME, sc.MTYPE, sc.PRTYPE, sc.LEN, c.DATA_TYPE, c.NUMERIC_PRECISION, c.NUMERIC_SCALE, c.DATETIME_PRECISION
from information_schema.INNODB_SYS_COLUMNS sc join information_schema.COLUMNS c
on sc.NAME=c.COLUMN_NAME
where sc.TABLE_ID  in (select TABLE_ID from information_schema.INNODB_SYS_TABLES where NAME='%v/%v')
and c.TABLE_SCHEMA='%v' and c.TABLE_NAME='%v'
order by c.ORDINAL_POSITION asc`, dbName, tbName, dbName, tbName)
	rows, e := session.Query(cmd)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var fields []*proto.Field
	for rows.Next() {
		var name, dataType string
		var mtype, prtype, Len uint32
		var numPrecision, numScale, timePrecision interface{}
		if e := rows.Scan(&name, &mtype, &prtype, &Len, &dataType, &numPrecision, &numScale, &timePrecision); e != nil {
			return nil, e
		}
		field := &proto.Field{
			Name:          name,
			Mtype:         mtype,
			Prtype:        prtype,
			Len:           Len,
			DataType:      dataType,
		}
		if numPrecision != nil {
			var num uint32
			arr := numPrecision.([]uint8)
			for _, c := range arr {
				num = 10*num + uint32(c-48)
			}
			field.NumPrecision = num
		}
		if numScale != nil {
			var num uint32
			arr := numScale.([]uint8)
			for _, c := range arr {
				num = 10*num + uint32(c-48)
			}
			field.NumScale = num
		}
		if timePrecision != nil {
			var num uint32
			arr := timePrecision.([]uint8)
			for _, c := range arr {
				num = 10*num + uint32(c-48)
			}
			field.TimePrecision = num
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func GetPriIndex(session *sql.DB, dbName, tbName string) (*proto.Index, error) {
	cmd := fmt.Sprintf(`select si.TYPE, si.PAGE_NO, sf.NAME
from information_schema.INNODB_SYS_INDEXES si left join information_schema.INNODB_SYS_FIELDS sf
on si.INDEX_ID=sf.INDEX_ID
where si.TABLE_ID in (select TABLE_ID from information_schema.INNODB_SYS_TABLES where NAME='%v/%v')
and si.TYPE in (1, 3)
order by sf.POS ASC`, dbName, tbName)
	rows, e := session.Query(cmd)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var typ, pageNo uint32
	var fieldNames []string
	for rows.Next() {
		var t, p uint32
		var name interface{}
		if e := rows.Scan(&t, &p, &name); e != nil {
			return nil, e
		}
		if name != nil {
			fieldNames = append(fieldNames, string(name.([]uint8)))
		} else {
			//自动生成的row_id主键
			fieldNames = append(fieldNames, "rowid")
		}
		typ = t
		pageNo = p
		//fieldNames = append(fieldNames, name)
	}
	if typ == 0 || pageNo == 0 {
		return nil, errors.New("not find")
	}
	priIndex := &proto.Index{
		Type:       typ,
		PageNo:     pageNo,
		FieldNames: fieldNames,
	}
	return priIndex, nil
}

func GetSecIndexes(session *sql.DB, dbName, tbName string) ([]*proto.Index, error) {
	cmd := fmt.Sprintf(`select si.TYPE, si.PAGE_NO, sf.NAME
from information_schema.INNODB_SYS_INDEXES si join information_schema.INNODB_SYS_FIELDS sf
on si.INDEX_ID=sf.INDEX_ID
where si.TABLE_ID in (select TABLE_ID from information_schema.INNODB_SYS_TABLES where NAME='%v/%v')
and si.TYPE not in (1, 3)
order by sf.POS ASC`, dbName, tbName)

	rows, e := session.Query(cmd)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var secIndexes []*proto.Index
	var typ, pageNo uint32
	var fieldNames []string
	for rows.Next() {
		var t, p uint32
		var name string
		if e := rows.Scan(&t, &p, &name); e != nil {
			return nil, e
		}
		if pageNo != 0 && pageNo != p {
			secIndex := &proto.Index{
				Type:       typ,
				PageNo:     pageNo,
				FieldNames: fieldNames,
			}
			secIndexes = append(secIndexes, secIndex)
			fieldNames = []string{}
		}
		typ = t
		pageNo = p
		fieldNames = append(fieldNames, name)
	}
	if len(fieldNames) > 0 {
		secIndex := &proto.Index{
			Type:       typ,
			PageNo:     pageNo,
			FieldNames: fieldNames,
		}
		secIndexes = append(secIndexes, secIndex)
	}
	return secIndexes, nil
}