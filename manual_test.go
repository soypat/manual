package manual

import "testing"

func Test_Allocator(t *testing.T) {
	var ta TestAllocator
	nums := makeIncNumbers(&ta)
	for i := range nums {
		if nums[i] != int64(i) {
			t.Error("failed")
		}
	}
	err := Free(&ta, nums)
	if err != nil {
		t.Error("failed")
	}
	err = Free(&ta, nums) // Double free should return error.
	if err == nil {
		t.Error("failed")
	}
}

func makeIncNumbers(a Allocator) []int64 {
	gz := Malloc[int64](a, 20)
	for i := range gz {
		gz[i] = int64(i)
	}
	return gz
}
