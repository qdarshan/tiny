package main

import (
	"fmt"
	"sync"
)

func Filter(in chan int, prime int, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Printf("Prime: %d\n", prime)
	var out chan int

	for i := range in {
		if i%prime != 0 {

			if out == nil {
				out = make(chan int)
				wg.Add(1)
				go Filter(out, i, wg)
			} else {
				out <- i
			}
		}
	}
	if out != nil {
		close(out)
	}
}

func main() {
	var wg sync.WaitGroup

	origin := make(chan int)

	wg.Add(1)
	go Filter(origin, 2, &wg)

	for i := 3; i < 100; i++ {
		origin <- i
	}
	close(origin)
	wg.Wait()
}
