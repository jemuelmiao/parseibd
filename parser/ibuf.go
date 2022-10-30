package parser

import "parseibd/proto"

func ReadPageIbuf(path string, pageNo uint64) *proto.PageIbuf {
	buff, e := ReadPage(path, pageNo)
	if e != nil {
		return nil
	}
	pageIbuf := new(proto.PageIbuf)
	var offset uint64

	fileHeader, fileHeaderOffset := ReadFileHeader(buff[offset:])
	offset += fileHeaderOffset
	pageIbuf.FileHeader = fileHeader

	//not used，和page header大小相同
	offset += 56

	for i:=0; i<8192; i++ {
		ibufBitmap, ibufBitmapOffset := ReadIbufBitmap(buff[offset:], i%2)
		pageIbuf.IbufBitmaps = append(pageIbuf.IbufBitmaps, ibufBitmap)
		offset += ibufBitmapOffset
	}

	fileTrailer, _ := ReadFileTrailer(buff[PageSize-8:])
	pageIbuf.FileTrailer = fileTrailer

	return pageIbuf
}

func ReadIbufBitmap(buff []byte, bitmapIndex int) (*proto.IbufBitmap, uint64) {
	ibufBitmap := new(proto.IbufBitmap)
	var offset uint64
	val := MachReadFrom1(buff[offset:])
	if bitmapIndex == 0 {
		ibufBitmap.Free = uint8((val & 0xC0) >> 6)
		ibufBitmap.IsBuffered = (val & 0x20) != 0
		ibufBitmap.IsBuf = (val & 0x10) != 0
	} else {
		ibufBitmap.Free = uint8((val & 0x0C) >> 2)
		ibufBitmap.IsBuffered = (val & 0x02) != 0
		ibufBitmap.IsBuf = (val & 0x01) != 0
	}
	return ibufBitmap, uint64(bitmapIndex)
}