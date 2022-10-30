package parseibd

import "testing"

func TestSegment1(t *testing.T) {
	ShowSegments("E:\\bdc\\mysql-server\\data\\jemuel\\test_date.ibd")
}

func TestSegment2(t *testing.T) {
	ShowSegments("F:\\SecureFiles\\huge_data.ibd")
}