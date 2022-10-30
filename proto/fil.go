package proto

//文件头
type FileHeader struct {
	Checksum	uint32
	Offset 		uint32 //page no
	Prev 		uint32 //prev page no，没有为FilNull
	Next 		uint32 //next page no，没有为FilNull
	Lsn 		uint64
	Type 		uint16 //page类型
	FlushLsn	uint64
	CompressCtrlInfo 	*FileCompressCtrlInfo
	SpaceId		uint32
}

type FileCompressCtrlInfo struct {
	Version 		uint8
	AlgorithmV1		uint8
	OriginalTypeV1	uint16
	OriginalSizeV1	uint16
	CompressSizeV1	uint16
}
//文件尾
type FileTrailer struct {
	Checksum	uint32
	Lsn 		uint32
}