package goroutine

import (
	"sync"
)

func Sum(arr []int) int {
	arr1 := arr[0 : len(arr)/2]
	arr2 := arr[len(arr)/2:]

	var sum1 int
	var sum2 int

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum1 = sum(arr1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum2 = sum(arr2)
	}()
	wg.Wait()
	return sum1 + sum2
}

func Sum2(arr []int) int {
	arr1 := arr[0 : len(arr)/2]
	arr2 := arr[len(arr)/2:]

	var sum1 int
	var sum2 int

	chanel := make(chan int, 2)
	go func() {
		sum1 = sum(arr1)
		chanel <- 1
	}()
	go func() {
		sum2 = sum(arr2)
		chanel <- 1
	}()
	<-chanel
	<-chanel
	return sum1 + sum2
}

func sum(arr []int) int {
	var ret int
	for _, v := range arr {
		ret += v
	}
	return ret
}
