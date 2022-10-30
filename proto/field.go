package proto

type Field struct {
	Name			string
	Mtype			uint32
	Prtype			uint32
	Len				uint32 //字段字节长度
	DataType		string //类型
	NumPrecision	uint32 //decimal
	NumScale		uint32 //decimal
	TimePrecision	uint32 //date、datetime、time、timestamp，取值：[0, 6]
}

type Index struct {
	Type 		uint32 //0：一般索引；1：自动生成的主键索引；2：unique索引；3：用户的主键索引；...
	PageNo		uint32 //索引root page no
	FieldNames	[]string //字段列表，自动生成的主键索引不需要，按照pos排序
}