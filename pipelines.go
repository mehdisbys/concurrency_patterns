package main

import (
	"fmt"
)

func Pipeline() {
	done := make(chan bool)
	defer close(done)

	multiplier := func(m int) func(i int) int {
		return func(i int) int {
			return i * m
		}
	}

	adder := func(m int) func(i int) int {
		return func(i int) int {
			return i + m
		}
	}

	multiply2 := stageFactory(multiplier(2))

	add1 := stageFactory(adder(1))

	intStream := generator(done, 1, 2, 3, 4)

	pipeline := multiply2(done, add1(done, multiply2(done, intStream)))

	for v := range pipeline {
		fmt.Println(v)
	}
}

func stageFactory(f func(int) int) func(done <-chan bool, intStream <-chan int) <-chan int {
	return func(
		done <-chan bool,
		intStream <-chan int,
	) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for i := range intStream {
				select {
				case <-done:
					return
				case stream <- f(i):
				}
			}
		}()
		return stream
	}
}

func generator(done <-chan bool, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}
