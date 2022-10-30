package proto

type InfimumRecord struct {
	ExtraNew 	*RecExtraNew
	ExtraOld 	*RecExtraOld
	Value 		string //"infimum" + '\0'
}
type SupremumRecord struct {
	ExtraNew 	*RecExtraNew
	ExtraOld 	*RecExtraOld
	Value 		string //新格式："supremum"；旧格式："supremum" + '\0'
}

//旧风格记录的固定信息
type RecExtraOld struct {
	IsMinRec		bool //记录是否同level非叶子节点page中第一个用户记录
	IsDeleted		bool //记录是否标记删除
	Nowned			uint8 //当前slot包含的记录数，该记录是当前slot的最后一个记录
	HeapNo			uint16 //记录的heap no
	Nfields			uint16 //记录中字段数量
	OffsetBytes		uint8 //字段长度偏移量字节数，1或2
	Next 			uint16 //指向page中下一个记录，存储绝对指针
}
//新风格记录的固定信息
type RecExtraNew struct {
	IsMinRec		bool //记录是否同level非叶子节点page中第一个用户记录
	IsDeleted		bool //记录是否标记删除
	Nowned			uint8 //当前slot包含的记录数，该记录是当前slot的最后一个记录
	HeapNo			uint16 //记录的heap no
	Status			uint8 //记录类型
	Next 			uint16 //指向page中下一个记录，存储相对偏移量
}

//聚簇索引-非叶子记录
type UserRecordPriNonleaf struct {
	ExtraNew 	*RecExtraNew //新格式
	ExtraOld 	*RecExtraOld //旧格式
	Primary		[]interface{} //主键值
	ChildPageNo uint32 //子page no
}
//聚簇索引-叶子记录
type UserRecordPriLeaf struct {
	ExtraNew	*RecExtraNew
	ExtraOld 	*RecExtraOld
	Primary		[]interface{} //主键值
	TrxId		uint64
	RollPtr		uint64
	Values 		[]interface{} //非主键字段值
}
//二级索引-非叶子记录
type UserRecordSecNonleaf struct {
	ExtraNew	*RecExtraNew
	ExtraOld 	*RecExtraOld
	Secondary	[]interface{} //二级索引值
	Primary 	[]interface{} //主键值
	ChildPageNo	uint32 //子page no
}
//二级索引-叶子记录
type UserRecordSecLeaf struct {
	ExtraNew	*RecExtraNew
	ExtraOld 	*RecExtraOld
	Secondary	[]interface{} //二级索引值
	Primary		[]interface{} //主键值
}