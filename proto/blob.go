package proto

type PageBlob struct {
	FileHeader		*FileHeader
	BlobHeader		*BlobHeader
	PayLoad			[]byte
	FileTrailer		*FileTrailer
}

type BlobHeader struct {
	Len 	uint32 //该page上blob长度
	Next 	uint32 //next page no，没有为FilNull
}

type ExternPtr struct {
	SpaceId 	uint32
	PageNo		uint32
	Offset 		uint32 //该page中blob header的偏移量
	ExternLen	uint64 //外部存储blob的长度，最高两位保留
	IsOwner		bool
	IsInherited	bool
}