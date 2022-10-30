package parser

import "parseibd/proto"

func ReadPageInode(path string, pageNo uint64) *proto.PageInode {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageInode := new(proto.PageInode)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageInode.FileHeader = fileHeader

	node, nodeOffset := ReadFlstNode(buff[offset:])
	pageInode.FlstNode = node
	offset += nodeOffset

	for i:=0; i<85; i++ {
		segInode, segInodeOffset := ReadSegmentInode(buff[offset:])
		pageInode.SegmentInodes = append(pageInode.SegmentInodes, segInode)
		offset += segInodeOffset
	}

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageInode.FileTrailer = fileTrailer

	return pageInode
}

func ReadSegmentInode(buff []byte) (*proto.SegmentInode, uint64) {
	segmentInode := new(proto.SegmentInode)
	var offset uint64
	segmentInode.SegId = MachReadFrom8(buff[offset:])
	offset += 8
	segmentInode.NotFullUsed = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	free, freeOffset := ReadFlstBaseNode(buff[offset:])
	segmentInode.Free = free
	offset += freeOffset

	notFull, notFullOffset := ReadFlstBaseNode(buff[offset:])
	segmentInode.NotFull = notFull
	offset += notFullOffset

	full, fullOffset := ReadFlstBaseNode(buff[offset:])
	segmentInode.Full = full
	offset += fullOffset

	offset += 4 //magic number

	for i:=0; i<32; i++ {
		segmentInode.FragPages = append(segmentInode.FragPages, uint32(MachReadFrom4(buff[offset:])))
		offset += 4
	}

	return segmentInode, offset
}