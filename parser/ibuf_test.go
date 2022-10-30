package parser

import "testing"

func TestIbuf1(t *testing.T) {
	pageIbuf := ReadPageIbuf("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd", 1)
	PrintPageIbuf(pageIbuf)
}