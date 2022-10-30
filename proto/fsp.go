package proto

type PageFsp struct {
	FileHeader			*FileHeader
	FspHeader			*FspHeader
	ExtentDescriptors	[]*ExtentDescriptor //最多256个
	FileTrailer			*FileTrailer
}

type PageXdes struct {
	FileHeader			*FileHeader
	ExtentDescriptors	[]*ExtentDescriptor //最多256个
	FileTrailer			*FileTrailer
}

type FspHeader struct {
	SpaceId 	uint32
	PageNum		uint32 //当前space中page数量
	FreeLimit	uint32
	SpaceFlags	uint32
	FragUsed	uint32 //free frag列表中使用的page数
	Free 		*FlstBaseNode //空闲extent列表，extent中的page未使用
	FreeFrag	*FlstBaseNode //部分使用的extent列表，extent中的page部分使用
	FullFrag	*FlstBaseNode //完全使用的extent列表，extent中的page全部使用
	NextSegId	uint64 //第一个未使用的segment id
	FullInodes	*FlstBaseNode //完全使用的inode page列表，inode page中segment全部使用
	FreeInodes	*FlstBaseNode //部分可用的inode page列表，inode page中segment部分使用
}

type ExtentDescriptor struct {
	SegId		uint64 //该extent所属segment id
	FlstNode	*FlstNode //双向链表，链接多个extent
	State 		uint32 //该extent状态
	PageFrees	[]bool //管理的page的状态，最多64个page，true：空闲，false：使用
}