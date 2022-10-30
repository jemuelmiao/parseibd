package parser

import (
	"parseibd/proto"
)

func ReadFileHeader(buff []byte) (*proto.FileHeader, uint64) {
	fileHeader := new(proto.FileHeader)
	var offset uint64
	fileHeader.Checksum = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileHeader.Offset = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileHeader.Prev = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileHeader.Next = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileHeader.Lsn = MachReadFrom8(buff[offset:])
	offset += 8
	fileHeader.Type = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	if fileHeader.Type == PageTypeCompressed {
		fileHeader.CompressCtrlInfo = &proto.FileCompressCtrlInfo{
			Version:        uint8(MachReadFrom1(buff[offset:])),
			AlgorithmV1:    uint8(MachReadFrom1(buff[offset+1:])),
			OriginalTypeV1: uint16(MachReadFrom2(buff[offset+2:])),
			OriginalSizeV1: uint16(MachReadFrom2(buff[offset+4:])),
			CompressSizeV1: uint16(MachReadFrom2(buff[offset+6:])),
		}
	} else {
		fileHeader.FlushLsn = MachReadFrom8(buff[offset:])
	}
	offset += 8
	fileHeader.SpaceId = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	return fileHeader, offset
}

func ReadFileTrailer(buff []byte) (*proto.FileTrailer, uint64) {
	fileTrailer := new(proto.FileTrailer)
	var offset uint64
	fileTrailer.Checksum = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	fileTrailer.Lsn = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	return fileTrailer, offset
}