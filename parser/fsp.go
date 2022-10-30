package parser

import (
	"parseibd/proto"
)

func ReadPageFsp(path string, pageNo uint64) *proto.PageFsp {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageFsp := new(proto.PageFsp)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageFsp.FileHeader = fileHeader

	fspHeader, fspHeaderOffset := ReadFspHeader(buff[offset:])
	pageFsp.FspHeader = fspHeader
	offset += fspHeaderOffset

	//不一定有256个，无效的extent descriptor表示还未分配
	for i:=0; i<256; i++ {
		extentDesc, extentDescOffset := ReadExtentDescriptor(buff[offset:], path, pageNo+1, i)
		if extentDesc.State == 0 {
			//无效extent
			break
		}
		offset += extentDescOffset
		pageFsp.ExtentDescriptors = append(pageFsp.ExtentDescriptors, extentDesc)
	}

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageFsp.FileTrailer = fileTrailer

	return pageFsp
}

func ReadPageXdes(path string, pageNo uint64) *proto.PageXdes {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageXdes := new(proto.PageXdes)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageXdes.FileHeader = fileHeader

	//not used，和fsp header大小相同
	offset += 112

	//不一定有256个，无效的extent descriptor表示还未分配
	for i:=0; i<256; i++ {
		extentDesc, extentDescOffset := ReadExtentDescriptor(buff[offset:], path, pageNo+1, i)
		if extentDesc.State == 0 {
			//无效extent
			break
		}
		offset += extentDescOffset
		pageXdes.ExtentDescriptors = append(pageXdes.ExtentDescriptors, extentDesc)
	}

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageXdes.FileTrailer = fileTrailer

	return pageXdes
}

func ReadFspHeader(buff []byte) (*proto.FspHeader, uint64) {
	fspHeader := new(proto.FspHeader)
	var offset uint64
	fspHeader.SpaceId = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	offset += 4 //not used

	fspHeader.PageNum = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	fspHeader.FreeLimit = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	fspHeader.SpaceFlags = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	fspHeader.FragUsed = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	free, freeOffset := ReadFlstBaseNode(buff[offset:])
	fspHeader.Free = free
	offset += freeOffset

	freeFrag, freeFragOffset := ReadFlstBaseNode(buff[offset:])
	fspHeader.FreeFrag = freeFrag
	offset += freeFragOffset

	fullFrag, fullFragOffset := ReadFlstBaseNode(buff[offset:])
	fspHeader.FullFrag = fullFrag
	offset += fullFragOffset

	fspHeader.NextSegId = MachReadFrom8(buff[offset:])
	offset += 8

	fullInodes, fullInodesOffset := ReadFlstBaseNode(buff[offset:])
	fspHeader.FullInodes = fullInodes
	offset += fullInodesOffset

	freeInodes, freeInodesOffset := ReadFlstBaseNode(buff[offset:])
	fspHeader.FreeInodes = freeInodes
	offset += freeInodesOffset

	return fspHeader, offset
}

func ReadExtentDescriptor(buff []byte, path string, basePageNo uint64, extentIndex int) (*proto.ExtentDescriptor, uint64) {
	extentDescriptor := new(proto.ExtentDescriptor)
	var offset uint64
	extentDescriptor.SegId = MachReadFrom8(buff[offset:])
	offset += 8

	node, nodeOffset := ReadFlstNode(buff[offset:])
	extentDescriptor.FlstNode = node
	offset += nodeOffset

	extentDescriptor.State = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	//不一定有64个，不存在的page表示还未分配
	for i:=0; i<16; i++ {
		pageNo := basePageNo + uint64(extentIndex*64 + i*4)
		val := MachReadFrom1(buff[offset:])
		if _, e := ReadPage(path, pageNo); e != nil {
			break
		}
		extentDescriptor.PageFrees = append(extentDescriptor.PageFrees, (val & 0x80) != 0)
		if _, e := ReadPage(path, pageNo+1); e != nil {
			break
		}
		extentDescriptor.PageFrees = append(extentDescriptor.PageFrees, (val & 0x20) != 0)
		if _, e := ReadPage(path, pageNo+2); e != nil {
			break
		}
		extentDescriptor.PageFrees = append(extentDescriptor.PageFrees, (val & 0x08) != 0)
		if _, e := ReadPage(path, pageNo+3); e != nil {
			break
		}
		extentDescriptor.PageFrees = append(extentDescriptor.PageFrees, (val & 0x02) != 0)
		offset += 1
	}

	return extentDescriptor, offset
}