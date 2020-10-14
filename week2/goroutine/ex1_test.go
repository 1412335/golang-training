package goroutine

import "testing"

func TestSum(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := Sum(arr); got != 45 {
		t.Errorf("Sum() = %v, want %v", got, 45)
	}
}
