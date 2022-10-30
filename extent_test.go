package parseibd

import "testing"

func TestExtent1(t *testing.T) {
	ShowExtents("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd")
}

func TestExtent2(t *testing.T) {
	ShowExtents("F:\\SecureFiles\\huge_data.ibd")
}