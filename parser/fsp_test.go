package parser

import "testing"

func TestFsp1(t *testing.T) {
	pageFsp := ReadPageFsp("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd", 0)
	PrintPageFsp(pageFsp)
}

func TestFsp2(t *testing.T) {
	pageFsp := ReadPageFsp("F:\\SecureFiles\\multi_keys.ibd", 0)
	PrintPageFsp(pageFsp)
}