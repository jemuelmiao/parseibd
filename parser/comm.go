package parser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"parseibd/proto"
	"strconv"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func IsExist(list []string, hit string) bool {
	for _, d := range list {
		if d == hit {
			return true
		}
	}
	return false
}

//list1 - list2
func GetListDiff(list1, list2 []string) []string {
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)
	for _, val := range list1 {
		map1[val] = true
	}
	for _, val := range list2 {
		map2[val] = true
	}
	var diff []string
	for val := range map1 {
		if _, ok := map2[val]; !ok {
			diff = append(diff, val)
		}
	}
	return diff
}

func MachReadFrom1(buff []byte) uint64 {
	return MachReadFromN(buff, 1)
}
func MachReadFrom2(buff []byte) uint64 {
	return MachReadFromN(buff, 2)
}
func MachReadFrom3(buff []byte) uint64 {
	return MachReadFromN(buff, 3)
}
func MachReadFrom4(buff []byte) uint64 {
	return MachReadFromN(buff, 4)
}
func MachReadFrom5(buff []byte) uint64 {
	return MachReadFromN(buff, 5)
}
func MachReadFrom6(buff []byte) uint64 {
	return MachReadFromN(buff, 6)
}
func MachReadFrom7(buff []byte) uint64 {
	return MachReadFromN(buff, 7)
}
func MachReadFrom8(buff []byte) uint64 {
	return MachReadFromN(buff, 8)
}
func MachReadFromN(buff []byte, n int) uint64 {
	var r uint64
	m := 0
	for i:=n-1; i>=0; i-- {
		r += uint64(buff[i])<<(m*8)
		m += 1
	}
	return r
}
func MachReadFloat(buff []byte) float32 {
	val := (*float32)(unsafe.Pointer(&buff[0]))
	return *val
}
func MachReadDouble(buff []byte) float64 {
	val := (*float64)(unsafe.Pointer(&buff[0]))
	return *val
}

func IsVariableLength(mtype uint32) bool {
	switch mtype {
	case DataSys, DataChar, DataFixbinary, DataInt, DataFloat, DataDouble, DataPoint:
		return false
	case DataMysql, DataVarchar, DataBinary, DataDecimal, DataVarmysql, DataVarpoint, DataGeometry, DataBlob:
		return true
	default:
		return true
	}
}
func IsNullable(prtype uint32) bool {
	return (prtype & 0x100) == 0
}
func IsUnsigned(prtype uint32) bool {
	return (prtype & 0x200) != 0
}
func IsBinary(prtype uint32) bool {
	return (prtype & 0x400) != 0
}
func IsBigCol(field *proto.Field) bool {
	return field.Len > 255 || field.Mtype == DataBlob || field.Mtype == DataGeometry || field.Mtype == DataVarpoint
}

//offset是实际数据在page中的绝对偏移量
func IsInfimum(offset uint64) bool {
	return offset == RecOldInfimumOffset || offset == RecNewInfimumOffset
}

//offset是实际数据在page中的绝对偏移量
func IsSupremum(offset uint64) bool {
	return offset == RecOldSupremumOffset || offset == RecNewSupremumOffset
}

//读取指定page
func ReadPage(path string, pageNo uint64) ([]byte, error) {
	start := pageNo * PageSize
	fileInfo, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	if start >= uint64(fileInfo.Size()) {
		return nil, ErrPageNotFound
	}
	file, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	if _, e := file.Seek(int64(start), io.SeekStart); e != nil {
		return nil, e
	}
	buff := make([]byte, PageSize)
	n, e := file.Read(buff)
	if e != nil || n != PageSize {
		return nil, e
	}
	return buff, nil
}

func GetPageTypeName(pageType int) string {
	switch pageType {
	case PageTypeIndex:
		return "FIL_PAGE_INDEX"
	case PageTypeRtree:
		return "FIL_PAGE_RTREE"
	case PageTypeUndoLog:
		return "FIL_PAGE_UNDO_LOG"
	case PageTypeInode:
		return "FIL_PAGE_INODE"
	case PageTypeIbufFreeList:
		return "FIL_PAGE_IBUF_FREE_LIST"
	case PageTypeAllocated:
		return "FIL_PAGE_TYPE_ALLOCATED"
	case PageTypeIbufBitmap:
		return "FIL_PAGE_IBUF_BITMAP"
	case PageTypeSys:
		return "FIL_PAGE_TYPE_SYS"
	case PageTypeTrxSys:
		return "FIL_PAGE_TYPE_TRX_SYS"
	case PageTypeFspHdr:
		return "FIL_PAGE_TYPE_FSP_HDR"
	case PageTypeXdes:
		return "FIL_PAGE_TYPE_XDES"
	case PageTypeBlob:
		return "FIL_PAGE_TYPE_BLOB"
	case PageTypeZblob:
		return "FIL_PAGE_TYPE_ZBLOB"
	case PageTypeZblob2:
		return "FIL_PAGE_TYPE_ZBLOB2"
	case PageTypeUnknown:
		return "FIL_PAGE_TYPE_UNKNOWN"
	case PageTypeCompressed:
		return "FIL_PAGE_COMPRESSED"
	case PageTypeEncrypted:
		return "FIL_PAGE_ENCRYPTED"
	case PageTypeCompressedAndEncrypted:
		return "FIL_PAGE_COMPRESSED_AND_ENCRYPTED"
	case PageTypeEncryptedRtree:
		return "FIL_PAGE_ENCRYPTED_RTREE"
	default:
		return ""
	}
}

func GetRecordType(status int) string {
	switch status {
	case RecStatusOrdinary:
		return "leaf"
	case RecStatusNodePtr:
		return "non leaf"
	case RecStatusInfimum:
		return "infimum"
	case RecStatusSupremum:
		return "supremum"
	default:
		return ""
	}
}

func GetMaxBytesPerChar(charset string) int {
	switch charset {
	case CharsetBig5:
		return 2
	case CharsetDec8:
		return 1
	case CharsetCp850:
		return 1
	case CharsetHp8:
		return 1
	case CharsetKoi8r:
		return 1
	case CharsetLatin1:
		return 1
	case CharsetLatin2:
		return 1
	case CharsetSwe7:
		return 1
	case CharsetAscii:
		return 1
	case CharsetUjis:
		return 3
	case CharsetSjis:
		return 2
	case CharsetHebrew:
		return 1
	case CharsetTis620:
		return 1
	case CharsetEuckr:
		return 2
	case CharsetKoi8u:
		return 1
	case CharsetGb2312:
		return 2
	case CharsetGreek:
		return 1
	case CharsetCp1250:
		return 1
	case CharsetGbk:
		return 2
	case CharsetLatin5:
		return 1
	case CharsetArmscii8:
		return 1
	case CharsetUtf8:
		return 3
	case CharsetUcs2:
		return 2
	case CharsetCp866:
		return 1
	case CharsetKeybcs2:
		return 1
	case CharsetMacce:
		return 1
	case CharsetMacroman:
		return 1
	case CharsetCp852:
		return 1
	case CharsetLatin7:
		return 1
	case CharsetUtf8mb4:
		return 4
	case CharsetCp1251:
		return 1
	case CharsetUtf16:
		return 4
	case CharsetUtf16le:
		return 4
	case CharsetCp1256:
		return 1
	case CharsetCp1257:
		return 1
	case CharsetUtf32:
		return 4
	case CharsetBinary:
		return 1
	case CharsetGeostd8:
		return 1
	case CharsetCp932:
		return 2
	case CharsetEucjpms:
		return 3
	case CharsetGb18030:
		return 4
	default:
		return 1
	}
}

func GetDecimalBytes(precision, scale uint32) (uint32, uint32) {
	intg0 := (precision-scale)/9
	frac0 := scale/9
	dig2bytes := []uint32{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}
	intg0x := (precision-scale) - intg0*9
	frac0x := scale - frac0*9
	return intg0*4 + dig2bytes[intg0x], frac0*4 + dig2bytes[frac0x]
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Assert(expr bool) {
	if !expr {
		fmt.Println("assert fail")
		os.Exit(1)
	}
}

func FormatHex(val uint64) string {
	return "0x" + strconv.FormatUint(val, 16)
}

func GetSegmentNo(offset int) int {
	return (offset - 38 - 12) / 192
}

func GetExtentNo(offset int) int {
	return (offset - 38 - 112) / 40
}