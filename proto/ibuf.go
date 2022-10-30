package proto

type PageIbuf struct {
	FileHeader		*FileHeader
	IbufBitmaps		[]*IbufBitmap
	FileTrailer		*FileTrailer
}

type IbufBitmap struct {
	Free 		uint8 //对应的page空闲的字节数范围，0：[0, 512)，1：[512, 1024)，2：[1024, 2048)，3：[2048, 新compact风格16252或旧风格16247)
	IsBuffered	bool //page是否有ibuf操作缓存
	IsBuf		bool //该page是否是ibuf tree的节点
}