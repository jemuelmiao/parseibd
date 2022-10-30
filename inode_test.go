package parseibd

import "testing"

func TestInode1(t *testing.T) {
	ShowInodes("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd")
}

func TestInode2(t *testing.T) {
	ShowInodes("F:\\SecureFiles\\huge_data.ibd")
}