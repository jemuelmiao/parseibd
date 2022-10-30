package parser

import "parseibd/proto"

func ReadPageBlob(path string, pageNo uint64) *proto.PageBlob {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageBlob := new(proto.PageBlob)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageBlob.FileHeader = fileHeader

	blobHeader, blobHeaderOffset := ReadBlobHeader(buff[offset:])
	offset += blobHeaderOffset
	pageBlob.BlobHeader = blobHeader

	pageBlob.PayLoad = buff[offset:offset+uint64(pageBlob.BlobHeader.Len)]

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageBlob.FileTrailer = fileTrailer

	return pageBlob
}

func ReadBlobHeader(buff []byte) (*proto.BlobHeader, uint64) {
	blobHeader := new(proto.BlobHeader)
	var offset uint64
	blobHeader.Len = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	blobHeader.Next = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	return blobHeader, offset
}

func ReadExternPtr(buff []byte) (*proto.ExternPtr, uint64) {
	externPtr := new(proto.ExternPtr)
	var offset uint64
	externPtr.SpaceId = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	externPtr.PageNo = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	externPtr.Offset = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	val := MachReadFrom8(buff[offset:])
	externPtr.ExternLen = val & (^(uint64(0xC0) << 56))
	externPtr.IsOwner = ((val >> 63) & 0x01) != 0
	externPtr.IsInherited = ((val >> 62) & 0x01) != 0
	offset += 8
	return externPtr, offset
}