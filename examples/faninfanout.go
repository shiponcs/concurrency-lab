package examples

import (
	"fmt"
	"sync"
	"time"
)

func FanOutFanIn() {
	tasks := []int{10, 20, 30, 40, 50, 60}

	worker1Jobs := make(chan int, 2)
	worker2Jobs := make(chan int, 2)
	worker3Jobs := make(chan int, 2)

	go func() {
		wg := sync.WaitGroup{}
		for i, task := range tasks {
			switch i % 3 {
			case 0:
				fmt.Println("fanning out", task)
				worker1Jobs <- task
			case 1:
				wg.Add(1)
				go func(t int) {
					defer wg.Done()
					time.Sleep(2 * time.Second)
					fmt.Println("fanning out (case 1)", t)
					worker2Jobs <- t
				}(task)
			case 2:
				fmt.Println("fanning out (ignore)", task)
				worker3Jobs <- task
			}
		}
		close(worker1Jobs)
		close(worker3Jobs)
		wg.Wait()
		close(worker2Jobs)
	}()

	worker1Results := make(chan int, 2)
	worker2Results := make(chan int, 2)
	worker3Results := make(chan int, 2)

	// No WaitGroup needed since we're not waiting for workers
	go func() {
		defer close(worker1Results)
		for task := range worker1Jobs {
			result := task * 2
			fmt.Println("processing ", task)
			time.Sleep(1 * time.Second)
			worker1Results <- result
		}
	}()

	go func() {
		defer close(worker2Results)
		for task := range worker2Jobs {
			result := task * 2
			fmt.Println("processing ", task)
			time.Sleep(1 * time.Second)
			worker2Results <- result
		}
	}()

	go func() {
		defer close(worker3Results)
		for task := range worker3Jobs {
			result := task * 2
			fmt.Println("processing ", task)
			time.Sleep(1 * time.Second)
			worker3Results <- result
		}
	}()

	allResults := make(chan int, 6)

	wg1 := sync.WaitGroup{}

	mergeTheResults := func(results <-chan int, output chan<- int) {
		defer wg1.Done()
		for result := range results {
			output <- result
		}
	}

	wg1.Add(3)
	go mergeTheResults(worker1Results, allResults)
	go mergeTheResults(worker2Results, allResults)
	go mergeTheResults(worker3Results, allResults)

	go func() {
		wg1.Wait()
		close(allResults)
	}()

	var finalResults []int
	for result := range allResults {
		finalResults = append(finalResults, result)
	}

	fmt.Println("Final results:", finalResults)

}
