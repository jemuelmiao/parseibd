package parser

import "testing"

func TestInode1(t *testing.T) {
	pageInode := ReadPageInode("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd", 2)
	PrintPageInode(pageInode)
}