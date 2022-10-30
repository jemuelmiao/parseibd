package proto

type PageIndex struct {
	FileHeader		*FileHeader
	PageHeader		*PageHeader
	InfimumRecord	*InfimumRecord
	SupremumRecord	*SupremumRecord
	//用户记录
	UserRecordPriNonleafs 	[]*UserRecordPriNonleaf //聚簇索引-非叶子记录
	UserRecordPriLeafs 		[]*UserRecordPriLeaf //聚簇索引-叶子记录
	UserRecordSecNonleafs 	[]*UserRecordSecNonleaf //二级索引-非叶子记录
	UserRecordSecLeafs 		[]*UserRecordSecLeaf //二级索引-叶子记录
	PageDirectory	*PageDirectory
	FileTrailer		*FileTrailer
}

//PageTypeIndex头
type PageHeader struct {
	Nslots 		uint16 //slot数量
	HeapTop		uint16
	IsComp		bool //是否为新格式
	Nheap		uint16 //堆中记录数，包含系统记录、标记删除记录、用户记录
	Free 		uint16 //空闲记录列，删除的记录会加入其中
	Garbage		uint16
	LastInsert	uint16
	Direction	uint16
	Ndirection	uint16
	Nrecs		uint16 //有效用户记录数
	MaxTrxId	uint64
	Level		uint16
	IndexId		uint64
	SegLeaf		*FileSegmentHeader
	SegNonLeaf	*FileSegmentHeader
}

type FileSegmentHeader struct {
	SpaceId		uint32
	PageNo		uint32
	Offset 		uint16
}

type PageDirectory struct {
	Slots 	[]uint16 //N个slot，每个slot 2个字节，每个slot中存储最后一个记录的页内偏移量，正序
}