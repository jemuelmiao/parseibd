package parser

import (
	"parseibd/proto"
)

//读取PageIndex
//参数：
//		path：ibd文件路径
//		pageNo：page no
//		fields：字段信息列表
//		currIndex：当前索引信息，可能为主键索引或二级索引
//		priIndex：主键索引信息，currIndex为主键索引时不需要
//		tableCharset：表存储字符集
//		inputCharset：输入数据终端字符集
//		rowFormat：行格式
func ReadPageIndex(path string, pageNo uint64, fields []*proto.Field, currIndex, priIndex *proto.Index,
	tableCharset, inputCharset, rowFormat string) *proto.PageIndex {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageIndex := new(proto.PageIndex)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageIndex.FileHeader = fileHeader

	pageHeader, pageHeaderOffset := ReadPageHeader(buff[offset:])
	offset += pageHeaderOffset
	pageIndex.PageHeader = pageHeader
	isComp := pageIndex.PageHeader.IsComp

	infimumRecord, infimumRecordOffset := ReadInfimumRecord(buff[offset:], isComp)
	offset += infimumRecordOffset
	pageIndex.InfimumRecord = infimumRecord

	supremumRecord, supremumRecordOffset := ReadSupremumRecord(buff[offset:], isComp)
	offset += supremumRecordOffset
	pageIndex.SupremumRecord = supremumRecord

	var dataOffset uint64
	if isComp {
		//新格式存储的相对值
		nextOffset := uint64(infimumRecord.ExtraNew.Next)
		dataOffset = (RecNewInfimumOffset + nextOffset) % PageSize
	} else {
		//旧格式存储的绝对值
		dataOffset = uint64(infimumRecord.ExtraOld.Next)
	}
	for !IsSupremum(dataOffset) {
		var extraNew *proto.RecExtraNew
		var extraOld *proto.RecExtraOld
		if (currIndex.Type & DictClustered) != 0 {
			//主键索引
			if pageIndex.PageHeader.Level != 0 {
				//非叶子
				userRecord := ReadUserRecordPriNonleaf(buff, path, dataOffset, isComp, fields, currIndex,
					tableCharset, inputCharset, rowFormat)
				pageIndex.UserRecordPriNonleafs = append(pageIndex.UserRecordPriNonleafs, userRecord)
				extraNew = userRecord.ExtraNew
				extraOld = userRecord.ExtraOld
			} else {
				//叶子
				userRecord := ReadUserRecordPriLeaf(buff, path, dataOffset, isComp, fields, currIndex,
					tableCharset, inputCharset, rowFormat)
				pageIndex.UserRecordPriLeafs = append(pageIndex.UserRecordPriLeafs, userRecord)
				extraNew = userRecord.ExtraNew
				extraOld = userRecord.ExtraOld
			}
		} else {
			//二级索引
			if pageIndex.PageHeader.Level != 0 {
				//非叶子
				userRecord := ReadUserRecordSecNonleaf(buff, path, dataOffset, isComp, fields, currIndex,
					priIndex, tableCharset, inputCharset, rowFormat)
				pageIndex.UserRecordSecNonleafs = append(pageIndex.UserRecordSecNonleafs, userRecord)
				extraNew = userRecord.ExtraNew
				extraOld = userRecord.ExtraOld
			} else {
				//叶子
				userRecord := ReadUserRecordSecLeaf(buff, path, dataOffset, isComp, fields, currIndex,
					priIndex, tableCharset, inputCharset, rowFormat)
				pageIndex.UserRecordSecLeafs = append(pageIndex.UserRecordSecLeafs, userRecord)
				extraNew = userRecord.ExtraNew
				extraOld = userRecord.ExtraOld
			}
		}
		if isComp {
			//新格式存储的相对值
			dataOffset = (dataOffset + uint64(extraNew.Next)) % PageSize
		} else {
			//旧格式存储的绝对值
			dataOffset = uint64(extraOld.Next)
		}
	}

	slotNum := pageIndex.PageHeader.Nslots
	pageDirectory, _ := ReadPageDirectory(buff[PageSize-8-2*slotNum:], slotNum)
	pageIndex.PageDirectory = pageDirectory

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageIndex.FileTrailer = fileTrailer

	return pageIndex
}

func ReadPageHeader(buff []byte) (*proto.PageHeader, uint64) {
	pageHeader := new(proto.PageHeader)
	var offset uint64
	pageHeader.Nslots = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.HeapTop = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	Nheap := uint16(MachReadFrom2(buff[offset:]))
	pageHeader.IsComp = (Nheap & 0x8000) != 0
	pageHeader.Nheap = Nheap & 0x7FFF
	offset += 2
	pageHeader.Free = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.Garbage = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.LastInsert = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.Direction = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.Ndirection = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.Nrecs = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.MaxTrxId = MachReadFrom8(buff[offset:])
	offset += 8
	pageHeader.Level = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	pageHeader.IndexId = MachReadFrom8(buff[offset:])
	offset += 8
	var segHeader *proto.FileSegmentHeader
	var segOffset uint64
	segHeader, segOffset = ReadFileSegmentHeader(buff[offset:])
	if segHeader.SpaceId != 0 || segHeader.PageNo != 0 || segHeader.Offset != 0 {
		pageHeader.SegLeaf = segHeader
	}
	offset += segOffset
	segHeader, segOffset = ReadFileSegmentHeader(buff[offset:])
	if segHeader.SpaceId != 0 || segHeader.PageNo != 0 || segHeader.Offset != 0 {
		pageHeader.SegNonLeaf = segHeader
	}
	offset += segOffset
	return pageHeader, offset
}

func ReadFileSegmentHeader(buff []byte) (*proto.FileSegmentHeader, uint64) {
	fileSegmentHeader := new(proto.FileSegmentHeader)
	var offset uint64
	fileSegmentHeader.SpaceId = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileSegmentHeader.PageNo = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileSegmentHeader.Offset = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	return fileSegmentHeader, offset
}

func ReadRecExtraOld(buff []byte) (*proto.RecExtraOld, uint64) {
	recOldExtra := new(proto.RecExtraOld)
	var offset uint64
	info := uint32(MachReadFrom4(buff[offset:]))
	recOldExtra.IsMinRec = (info & 0x10000000) != 0
	recOldExtra.IsDeleted = (info & 0x20000000) != 0
	recOldExtra.Nowned = uint8((info & 0x0F000000) >> 24)
	recOldExtra.HeapNo = uint16((info & 0x00FFF800) >> 11)
	recOldExtra.Nfields = uint16((info & 0x000007FE) >> 1)
	if (info & 0x01) != 0 {
		recOldExtra.OffsetBytes = 1
	} else {
		recOldExtra.OffsetBytes = 2
	}
	offset += 4
	recOldExtra.Next = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	return recOldExtra, offset
}

func ReadRecExtraNew(buff []byte) (*proto.RecExtraNew, uint64) {
	recNewExtra := new(proto.RecExtraNew)
	var offset uint64
	info := uint32(MachReadFrom3(buff[offset:]))
	recNewExtra.IsMinRec = (info & 0x100000) != 0
	recNewExtra.IsDeleted = (info & 0x200000) != 0
	recNewExtra.Nowned = uint8((info & 0x0F0000) >> 16)
	recNewExtra.HeapNo = uint16((info & 0x00FFF8) >> 3)
	recNewExtra.Status = uint8(info & 0x000007)
	offset += 3
	recNewExtra.Next = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	return recNewExtra, offset
}

func ReadInfimumRecord(buff []byte, isComp bool) (*proto.InfimumRecord, uint64) {
	infimumRecord := new(proto.InfimumRecord)
	var offset uint64
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[offset:])
		offset += extraOffset
		infimumRecord.ExtraNew = extra
	} else {
		offset += 1 //旧格式：1字节end offset
		extra, extraOffset := ReadRecExtraOld(buff[offset:])
		offset += extraOffset
		infimumRecord.ExtraOld = extra
	}
	infimumRecord.Value = "infimum"
	offset += 8
	return infimumRecord, offset
}

func ReadSupremumRecord(buff []byte, isComp bool) (*proto.SupremumRecord, uint64) {
	supremumRecord := new(proto.SupremumRecord)
	var offset uint64
	if isComp {
		extra, extraOffset := ReadRecExtraNew(buff[offset:])
		offset += extraOffset
		supremumRecord.ExtraNew = extra
		offset += 8
	} else {
		offset += 1 //旧格式：1字节end offset
		extra, extraOffset := ReadRecExtraOld(buff[offset:])
		offset += extraOffset
		supremumRecord.ExtraOld = extra
		offset += 9
	}
	supremumRecord.Value = "supremum"
	return supremumRecord, offset
}

func ReadPageDirectory(buff []byte, slotNum uint16) (*proto.PageDirectory, uint64) {
	pageDirectory := new(proto.PageDirectory)
	var offset uint64
	//顺序翻转
	for i:=slotNum; i>0; i-- {
		pageDirectory.Slots = append(pageDirectory.Slots, uint16(MachReadFrom2(buff[2*(i-1):])))
		offset += 2
	}
	return pageDirectory, offset
}
