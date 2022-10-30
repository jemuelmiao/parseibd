package parseibd

import "testing"

func TestPage1(t *testing.T) {
	ShowPages("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd")
}

func TestPage2(t *testing.T) {
	ShowPages("F:\\SecureFiles\\huge_data.ibd")
}