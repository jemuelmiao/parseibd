package parser

import (
	"fmt"
	"parseibd/proto"
)

func PrintFileHeader(fileHeader *proto.FileHeader) {
	if fileHeader == nil {
		return
	}
	indent := "    "
	fmt.Println("FileHeader....")
	fmt.Println(indent + "checksum:", fileHeader.Checksum)
	fmt.Println(indent + "page no:", fileHeader.Offset)
	if fileHeader.Prev == FilNull {
		fmt.Println(indent + "prev page no: FilNull")
	} else {
		fmt.Println(indent + "prev page no:", fileHeader.Prev)
	}
	if fileHeader.Next == FilNull {
		fmt.Println(indent + "next page no: FilNull")
	} else {
		fmt.Println(indent + "next page no:", fileHeader.Next)
	}
	fmt.Println(indent + "page type:", GetPageTypeName(int(fileHeader.Type)))
	fmt.Println(indent + "space id:", fileHeader.SpaceId)
}

func PrintFileTrailer(fileTrailer *proto.FileTrailer) {
	if fileTrailer == nil {
		return
	}
	indent := "    "
	fmt.Println("FileTrailer....")
	fmt.Println(indent + "checksum:", fileTrailer.Checksum)
}

func PrintPageIndex(pageIndex *proto.PageIndex) {
	if pageIndex == nil {
		return
	}
	fmt.Println("PageIndex....")
	PrintFileHeader(pageIndex.FileHeader)
	PrintPageHeader(pageIndex.PageHeader)
	PrintInfimum(pageIndex.InfimumRecord)
	PrintSupremum(pageIndex.SupremumRecord)
	PrintUserRecordPriNonleafs(pageIndex.UserRecordPriNonleafs)
	PrintUserRecordPriLeafs(pageIndex.UserRecordPriLeafs)
	PrintUserRecordSecNonleafs(pageIndex.UserRecordSecNonleafs)
	PrintUserRecordSecLeafs(pageIndex.UserRecordSecLeafs)
	PrintPageDirectory(pageIndex.PageDirectory)
	PrintFileTrailer(pageIndex.FileTrailer)
}

func PrintPageHeader(pageHeader *proto.PageHeader) {
	if pageHeader == nil {
		return
	}
	indent := "    "
	fmt.Println("PageHeader....")
	fmt.Println(indent + "slot num:", pageHeader.Nslots)
	fmt.Println(indent + "is comp:", pageHeader.IsComp)
	fmt.Println(indent + "heap record num:", pageHeader.Nheap)
	fmt.Println(indent + "user record num:", pageHeader.Nrecs)
	fmt.Println(indent + "page level:", pageHeader.Level)
	fmt.Println(indent + "index id:", pageHeader.IndexId)
}

func PrintRecExtraNew(extraNew *proto.RecExtraNew) {
	indent := "    "
	fmt.Println(indent + "is min record:", extraNew.IsMinRec)
	fmt.Println(indent + "is deleted:", extraNew.IsDeleted)
	fmt.Println(indent + "slot record num:", extraNew.Nowned)
	fmt.Println(indent + "record type:", GetRecordType(int(extraNew.Status)))
	fmt.Println(indent + "next record:", extraNew.Next)
}

func PrintRecExtraOld(extraOld *proto.RecExtraOld) {
	indent := "    "
	fmt.Println(indent + "is min record:", extraOld.IsMinRec)
	fmt.Println(indent + "is deleted:", extraOld.IsDeleted)
	fmt.Println(indent + "slot record num:", extraOld.Nowned)
	fmt.Println(indent + "field num:", extraOld.Nfields)
	if extraOld.OffsetBytes == 1 {
		fmt.Println(indent + "offset byte count: 1")
	} else {
		fmt.Println(indent + "offset byte count: 2")
	}
	fmt.Println(indent + "next record:", extraOld.Next)
}

func PrintInfimum(infimumRecord *proto.InfimumRecord) {
	if infimumRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("InfimumRecord....")
	if infimumRecord.ExtraNew != nil {
		PrintRecExtraNew(infimumRecord.ExtraNew)
	} else if infimumRecord.ExtraOld != nil {
		PrintRecExtraOld(infimumRecord.ExtraOld)
	}
	fmt.Println(indent + "value:", infimumRecord.Value)
}

func PrintSupremum(supremumRecord *proto.SupremumRecord) {
	if supremumRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("SupremumRecord....")
	if supremumRecord.ExtraNew != nil {
		PrintRecExtraNew(supremumRecord.ExtraNew)
	} else if supremumRecord.ExtraOld != nil {
		PrintRecExtraOld(supremumRecord.ExtraOld)
	}
	fmt.Println(indent + "value:", supremumRecord.Value)
}

func PrintPageDirectory(pageDirectory *proto.PageDirectory) {
	if pageDirectory == nil {
		return
	}
	indent := "    "
	fmt.Println("PageDirectory....")
	fmt.Println(indent + "slots:", pageDirectory.Slots)
}

func PrintUserRecordPriNonleaf(userRecord *proto.UserRecordPriNonleaf) {
	if userRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("UserRecordPriNonleaf....")
	if userRecord.ExtraNew != nil {
		PrintRecExtraNew(userRecord.ExtraNew)
	} else if userRecord.ExtraOld != nil {
		PrintRecExtraOld(userRecord.ExtraOld)
	}
	fmt.Println(indent + "primary value:", userRecord.Primary)
	fmt.Println(indent + "child page no:", userRecord.ChildPageNo)
}

func PrintUserRecordPriNonleafs(userRecords []*proto.UserRecordPriNonleaf) {
	for _, userRecord := range userRecords {
		PrintUserRecordPriNonleaf(userRecord)
	}
}

func PrintUserRecordPriLeaf(userRecord *proto.UserRecordPriLeaf) {
	if userRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("UserRecordPriLeaf....")
	if userRecord.ExtraNew != nil {
		PrintRecExtraNew(userRecord.ExtraNew)
	} else if userRecord.ExtraOld != nil {
		PrintRecExtraOld(userRecord.ExtraOld)
	}
	fmt.Println(indent + "primary value:", userRecord.Primary)
	fmt.Println(indent + "trx id:", userRecord.TrxId)
	fmt.Println(indent + "roll ptr:", FormatHex(userRecord.RollPtr))
	fmt.Println(indent + "field values:", userRecord.Values)
}

func PrintUserRecordPriLeafs(userRecords []*proto.UserRecordPriLeaf) {
	for _, userRecord := range userRecords {
		PrintUserRecordPriLeaf(userRecord)
	}
}

func PrintUserRecordSecNonleaf(userRecord *proto.UserRecordSecNonleaf) {
	if userRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("UserRecordSecNonleaf....")
	if userRecord.ExtraNew != nil {
		PrintRecExtraNew(userRecord.ExtraNew)
	} else if userRecord.ExtraOld != nil {
		PrintRecExtraOld(userRecord.ExtraOld)
	}
	fmt.Println(indent + "secondary value:", userRecord.Secondary)
	fmt.Println(indent + "primary value:", userRecord.Primary)
	fmt.Println(indent + "child page no:", userRecord.ChildPageNo)
}

func PrintUserRecordSecNonleafs(userRecords []*proto.UserRecordSecNonleaf)  {
	for _, userRecord := range userRecords {
		PrintUserRecordSecNonleaf(userRecord)
	}
}

func PrintUserRecordSecLeaf(userRecord *proto.UserRecordSecLeaf) {
	if userRecord == nil {
		return
	}
	indent := "    "
	fmt.Println("UserRecordSecLeaf....")
	if userRecord.ExtraNew != nil {
		PrintRecExtraNew(userRecord.ExtraNew)
	} else  if userRecord.ExtraOld != nil {
		PrintRecExtraOld(userRecord.ExtraOld)
	}
	fmt.Println(indent + "secondary value:", userRecord.Secondary)
	fmt.Println(indent + "primary value:", userRecord.Primary)
}

func PrintUserRecordSecLeafs(userRecords []*proto.UserRecordSecLeaf) {
	for _, userRecord := range userRecords {
		PrintUserRecordSecLeaf(userRecord)
	}
}

func PrintPageBlob(pageBlob *proto.PageBlob) {
	if pageBlob == nil {
		return
	}
	fmt.Println("PageBlob....")
	PrintFileHeader(pageBlob.FileHeader)
	PrintBlobHeader(pageBlob.BlobHeader)
	PrintFileTrailer(pageBlob.FileTrailer)
}

func PrintBlobHeader(blobHeader *proto.BlobHeader) {
	indent := "    "
	fmt.Println(indent + "len:", blobHeader.Len)
	if blobHeader.Next == FilNull {
		fmt.Println(indent + "next: FilNull")
	} else {
		fmt.Println(indent + "next:", blobHeader.Next)
	}
}

func PrintPageFsp(pageFsp *proto.PageFsp) {
	if pageFsp == nil {
		return
	}
	fmt.Println("PageFsp....")
	PrintFileHeader(pageFsp.FileHeader)
	PrintFspHeader(pageFsp.FspHeader)
	PrintExtentDescriptors(pageFsp.ExtentDescriptors)
	PrintFileTrailer(pageFsp.FileTrailer)
}

func PrintPageXdes(pageXdes *proto.PageXdes) {
	if pageXdes == nil {
		return
	}
	fmt.Println("PageXdes...")
	PrintFileHeader(pageXdes.FileHeader)
	PrintExtentDescriptors(pageXdes.ExtentDescriptors)
	PrintFileTrailer(pageXdes.FileTrailer)
}

func PrintFspHeader(fspHeader *proto.FspHeader) {
	if fspHeader == nil {
		return
	}
	indent := "    "
	fmt.Println("FspHeader....")
	fmt.Println(indent + "space id:", fspHeader.SpaceId)
	fmt.Println(indent + "page num:", fspHeader.PageNum)
	fmt.Println(indent + "free limit:", fspHeader.FreeLimit)
	fmt.Println(indent + "space flags:", fspHeader.SpaceFlags)
	fmt.Println(indent + "frag used:", fspHeader.FragUsed)
	fmt.Println(indent + "free:", StringfyFlstBaseNode(fspHeader.Free))
	fmt.Println(indent + "free frag:", StringfyFlstBaseNode(fspHeader.FreeFrag))
	fmt.Println(indent + "full frag:", StringfyFlstBaseNode(fspHeader.FullFrag))
	fmt.Println(indent + "next segment id:", fspHeader.NextSegId)
	fmt.Println(indent + "full inodes:", StringfyFlstBaseNode(fspHeader.FullInodes))
	fmt.Println(indent + "free inodes:", StringfyFlstBaseNode(fspHeader.FreeInodes))
}

func PrintExtentDescriptor(extentDescriptor *proto.ExtentDescriptor) {
	if extentDescriptor == nil {
		return
	}
	indent := "    "
	fmt.Println("ExtentDescriptor...")
	fmt.Println(indent + "segment id:", extentDescriptor.SegId)
	fmt.Println(indent + "node:", StringfyFlstNode(extentDescriptor.FlstNode))
	switch extentDescriptor.State {
	case ExtentStateFree:
		fmt.Println(indent + "state: free")
	case ExtentStateFreeFrag:
		fmt.Println(indent + "state: free frag")
	case ExtentStateFullFrag:
		fmt.Println(indent + "state: full frag")
	case ExtentStateFseg:
		fmt.Println(indent + "state: fseg")
	}
	fmt.Println(indent + "page frees:", extentDescriptor.PageFrees)
}

func PrintExtentDescriptors(extentDescriptors []*proto.ExtentDescriptor) {
	for _, extentDescriptor := range extentDescriptors {
		PrintExtentDescriptor(extentDescriptor)
	}
}

func PrintPageInode(pageInode *proto.PageInode) {
	if pageInode == nil {
		return
	}
	indent := "    "
	fmt.Println("PageInode...")
	PrintFileHeader(pageInode.FileHeader)
	fmt.Println(indent + "node:", StringfyFlstNode(pageInode.FlstNode))
	PrintSegmentInodes(pageInode.SegmentInodes)
	PrintFileTrailer(pageInode.FileTrailer)
}

func PrintSegmentInode(segInode *proto.SegmentInode) {
	if segInode == nil {
		return
	}
	indent := "    "
	fmt.Println("SegmentInode...")
	fmt.Println(indent + "segment id:", segInode.SegId)
	fmt.Println(indent + "not full used:", segInode.NotFullUsed)
	fmt.Println(indent + "free:", StringfyFlstBaseNode(segInode.Free))
	fmt.Println(indent + "not full:", StringfyFlstBaseNode(segInode.NotFull))
	fmt.Println(indent + "full:", StringfyFlstBaseNode(segInode.Full))
	var fragPages []interface{}
	for _, pageNo := range segInode.FragPages {
		if pageNo == FilNull {
			fragPages = append(fragPages, "FilNull")
		} else {
			fragPages = append(fragPages, pageNo)
		}
	}
	fmt.Println(indent + "frag pages:", fragPages)
}

func PrintSegmentInodes(segInodes []*proto.SegmentInode) {
	for _, segInode := range segInodes {
		PrintSegmentInode(segInode)
	}
}

func PrintPageIbuf(pageIbuf *proto.PageIbuf) {
	if pageIbuf == nil {
		return
	}
	fmt.Println("PageIbuf...")
	PrintFileHeader(pageIbuf.FileHeader)
	PrintIbufBitmaps(pageIbuf.IbufBitmaps)
	PrintFileTrailer(pageIbuf.FileTrailer)
}

func PrintIbufBitmap(ibufBitmap *proto.IbufBitmap) {
	if ibufBitmap == nil {
		return
	}
	indent := "    "
	fmt.Println("IbufBitmap...")
	switch ibufBitmap.Free {
	case 0:
		fmt.Println(indent + "free:[0, 512)")
	case 1:
		fmt.Println(indent + "free:[512, 1024)")
	case 2:
		fmt.Println(indent + "free:[1024, 2048)")
	case 3:
		fmt.Println(indent + "free:[2048, 16252或16247)")
	}
	fmt.Println(indent + "is buffered:", ibufBitmap.IsBuffered)
	fmt.Println(indent + "is buf tree:", ibufBitmap.IsBuf)
}

func PrintIbufBitmaps(ibufBitmaps []*proto.IbufBitmap) {
	for _, ibufBitmap := range ibufBitmaps {
		PrintIbufBitmap(ibufBitmap)
	}
}

func StringfyFlstBaseNode(baseNode *proto.FlstBaseNode) string {
	return fmt.Sprintf("len:%v, first:(%v), last:(%v)", baseNode.Len,
		StringfyFilAddr(baseNode.First), StringfyFilAddr(baseNode.Last))
}

func StringfyFlstNode(node *proto.FlstNode) string {
	return fmt.Sprintf("prev:(%v), next:(%v)", StringfyFilAddr(node.Prev), StringfyFilAddr(node.Next))
}

func StringfyFilAddr(filAddr *proto.FilAddr) string {
	if filAddr.PageNo == FilNull {
		return fmt.Sprintf("page no:FilNull, offset:%v", filAddr.Offset)
	} else {
		return fmt.Sprintf("page no:%v, offset:%v", filAddr.PageNo, filAddr.Offset)
	}
}