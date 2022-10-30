package parser

import "parseibd/proto"

func ReadFlstBaseNode(buff []byte) (*proto.FlstBaseNode, uint64) {
	flstBaseNode := new(proto.FlstBaseNode)
	var offset uint64
	flstBaseNode.Len = uint32(MachReadFrom4(buff[offset:]))
	offset += 4

	first, firstOffset := ReadFilAddr(buff[offset:])
	flstBaseNode.First = first
	offset += firstOffset

	last, lastOffset := ReadFilAddr(buff[offset:])
	flstBaseNode.Last = last
	offset += lastOffset

	return flstBaseNode, offset
}

func ReadFlstNode(buff []byte) (*proto.FlstNode, uint64) {
	flstNode := new(proto.FlstNode)
	var offset uint64
	prev, prevOffset := ReadFilAddr(buff[offset:])
	flstNode.Prev = prev
	offset += prevOffset

	next, nextOffset := ReadFilAddr(buff[offset:])
	flstNode.Next = next
	offset += nextOffset

	return flstNode, offset
}

func ReadFilAddr(buff []byte) (*proto.FilAddr, uint64) {
	filAddr := new(proto.FilAddr)
	var offset uint64
	filAddr.PageNo = uint32(MachReadFrom4(buff[offset:]))
	offset += 4
	filAddr.Offset = uint16(MachReadFrom2(buff[offset:]))
	offset += 2
	return filAddr, offset
}