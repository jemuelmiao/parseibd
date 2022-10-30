package proto

type PageInode struct {
	FileHeader		*FileHeader
	FlstNode		*FlstNode //双向链表，链接多个page inode
	SegmentInodes	[]*SegmentInode //page inode管理的segment数组，一个page inode最多管理85个segment
	FileTrailer		*FileTrailer
}

type SegmentInode struct {
	SegId		uint64 //该segment id
	NotFullUsed	uint32 //not full列表中使用的page数
	Free 		*FlstBaseNode //该segment中空闲的extent列表
	NotFull		*FlstBaseNode //该segment中部分使用的extent列表
	Full		*FlstBaseNode //该segment中完全使用的extent列表
	FragPages	[]uint32 //32个碎片page no，没有则为FIL_NULL
}