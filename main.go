package main

import (
	"fmt"

	"github.com/shipon/concurrency-lab/examples"
)

func main() {
	// fmt.Println("Welcome to Concurrency Lab!")
	// fmt.Println()

	// // Run basic goroutines example
	// examples.BasicGoroutines()
	// fmt.Println()

	// // Run channel examples
	// examples.ChannelBasics()
	// fmt.Println()

	// examples.BufferedChannels()
	// fmt.Println()

	// // Run select statement examples
	// examples.SelectBasics()
	// fmt.Println()

	// examples.SelectWithTimeout()
	// fmt.Println()

	// examples.SelectNonBlocking()
	// fmt.Println()

	// examples.SelectMultiChannelNonBlocking()
	// fmt.Println()

	// examples.CompareWithIfElse()
	// fmt.Println()

	// // Run pipeline examples
	// examples.SimplePipeline()
	// fmt.Println()

	// examples.MultiStagePipeline()
	// fmt.Println()

	// examples.FanOutFanIn()
	// fmt.Println()

	// examples.SimpleFanOutFanIn()
	// fmt.Println()

	// examples.BufferedPipeline()
	// fmt.Println()

	// examples.PipelineWithErrorHandling()
	// fmt.Println()

	fmt.Println("Hello, World!")
	examples.FanOutFanIn()
	fmt.Println("ends")

	fmt.Println("All examples completed!")
}
