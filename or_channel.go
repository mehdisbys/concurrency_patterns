package main

import (
	"fmt"
	"time"
)

func Or() func(channels ...<-chan interface{}) <-chan interface{} {

	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {

			defer close(orDone)

			switch len(channels) {
			case 2:

				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:

				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-or(append(channels[2:], orDone)...):

				}
			}
		}()
		return orDone
	}
	return or
}

func OrDemo() {

	start := time.Now()
	or := Or()
	<-or(
		timer(2*time.Hour),
		timer(5*time.Minute),
		timer(1*time.Second),
		timer(1*time.Hour),
		timer(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func timer(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}
