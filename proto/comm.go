package proto

type FlstBaseNode struct {
	Len 	uint32 //链表长度
	First 	*FilAddr //链表头节点
	Last 	*FilAddr //链表尾节点
}

type FlstNode struct {
	Prev	*FilAddr //链表前一节点
	Next	*FilAddr //链表后一节点
}

type FilAddr struct {
	PageNo		uint32 //page no
	Offset 		uint16 //page中字节偏移量
}