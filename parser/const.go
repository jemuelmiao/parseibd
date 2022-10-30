package parser

import "errors"

var (
	ErrPageNotFound = errors.New("page not found")
)

const (
	//Page大小
	PageSize = 16*1024

	FilNull = 0xFFFFFFFF

	//Page类型
	PageTypeIndex = 17855
	PageTypeRtree = 17854
	PageTypeUndoLog = 2
	PageTypeInode = 3
	PageTypeIbufFreeList = 4
	PageTypeAllocated = 0
	PageTypeIbufBitmap = 5
	PageTypeSys = 6
	PageTypeTrxSys = 7
	PageTypeFspHdr = 8
	PageTypeXdes = 9
	PageTypeBlob = 10
	PageTypeZblob = 11
	PageTypeZblob2 = 12
	PageTypeUnknown = 13
	PageTypeCompressed = 14
	PageTypeEncrypted = 15
	PageTypeCompressedAndEncrypted = 16
	PageTypeEncryptedRtree = 17

	//游标移动方向
	PageLeft = 1
	PageRight = 2
	PageSameRec = 3
	PageSamePage = 4
	PageNoDirection = 5

	//索引类型
	DictClustered = 1
	DictUnique = 2
	DictIbuf = 8
	DictCorrupt = 16

	////隐藏系统列长度
	//RowIdLen = 6
	//TrxIdLen = 6
	//RollPtrLen = 7

	//行格式
	RowFormatDynamic = "dynamic"
	RowFormatCompressed = "compressed"
	RowFormatCompact = "compact"
	RowFormatRedundant = "redundant"
)

const (
	//新、旧风格行格式的固定字节数
	RecOldExtraBytes = 6
	RecNewExtraBytes = 5

	//记录类型
	RecStatusOrdinary = 0 //叶子节点记录
	RecStatusNodePtr = 1 //非叶子节点记录
	RecStatusInfimum = 2 //infimum记录
	RecStatusSupremum = 3 //supremum记录

	//infimum、supremum实际数据在page中偏移量
	RecOldInfimumOffset = 101
	RecOldSupremumOffset = 116
	RecNewInfimumOffset = 99
	RecNewSupremumOffset = 112

	//compact、redundant格式索引页存储的字段最长字节数
	RecAntelopeMaxIndexColLen = 768
)

const (
	//innodb层数据类型
	DataVarchar = 1
	DataChar = 2
	DataFixbinary = 3
	DataBinary = 4
	DataBlob = 5
	DataInt = 6
	DataSys = 8
	DataFloat = 9
	DataDouble = 10
	DataDecimal = 11
	DataVarmysql = 12
	DataMysql = 13
	DataGeometry = 14
	DataPoint = 15
	DataVarpoint = 16
)

const (
	//字符集
	CharsetBig5 = "big5"
	CharsetDec8 = "dec8"
	CharsetCp850 = "cp850"
	CharsetHp8 = "hp8"
	CharsetKoi8r = "koi8r"
	CharsetLatin1 = "latin1"
	CharsetLatin2 = "latin2"
	CharsetSwe7 = "swe7"
	CharsetAscii = "ascii"
	CharsetUjis = "ujis"
	CharsetSjis = "sjis"
	CharsetHebrew = "hebrew"
	CharsetTis620 = "tis620"
	CharsetEuckr = "euckr"
	CharsetKoi8u = "koi8u"
	CharsetGb2312 = "gb2312"
	CharsetGreek = "greek"
	CharsetCp1250 = "cp1250"
	CharsetGbk = "gbk"
	CharsetLatin5 = "latin5"
	CharsetArmscii8 = "armscii8"
	CharsetUtf8 = "utf8"
	CharsetUcs2 = "ucs2"
	CharsetCp866 = "cp866"
	CharsetKeybcs2 = "keybcs2"
	CharsetMacce = "macce"
	CharsetMacroman = "macroman"
	CharsetCp852 = "cp852"
	CharsetLatin7 = "latin7"
	CharsetUtf8mb4 = "utf8mb4"
	CharsetCp1251 = "cp1251"
	CharsetUtf16 = "utf16"
	CharsetUtf16le = "utf16le"
	CharsetCp1256 = "cp1256"
	CharsetCp1257 = "cp1257"
	CharsetUtf32 = "utf32"
	CharsetBinary = "binary"
	CharsetGeostd8 = "geostd8"
	CharsetCp932 = "cp932"
	CharsetEucjpms = "eucjpms"
	CharsetGb18030 = "gb18030"
)

const (
	//extent状态
	ExtentStateFree	= 1 //extent在空闲extent列表中
	ExtentStateFreeFrag	= 2 //extent在部分使用extent列表中
	ExtentStateFullFrag = 3 //extent在完全使用extent列表中
	ExtentStateFseg = 4 //extent已分配给segment
)