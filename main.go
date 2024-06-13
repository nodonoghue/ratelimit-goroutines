package main

import (
	"fmt"
	"sync"
	"time"
)

var limit = 20

var bucket = make(chan struct{}, limit)
var startTime = time.Now()

func main() {
	//Will be determined from a simple len(records_slice)
	totalRequests := 200

	//Make a WaitGroup that will be used to prevent the application from ending before all of the above records are processed
	var wg sync.WaitGroup

	//Add the total number of records to the WaitGroup, once all processing has completed and marked all as Done() the WaitGroup will unblock
	wg.Add(totalRequests)

	//infinite loop that creates the timer, each tick will refresh the bucket of tokens, setting back to limit
	go func() {
		for {
			// the bucket refill routine
			for i := 0; i < limit; i++ {
				bucket <- struct{}{} //refill the bucket through a channel
			}
			time.Sleep(time.Second) //sleep to enforce a delay before refilling the token bucket
		}
	}()

	// Rate limit to N=limit requests per second
	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			sendRequest(i + 1)
		}(i)
	}

	// Wait until all request are completed
	wg.Wait()
}

func sendRequest(i int) {
	// get "token" from the bucket
	//execution beyond this line will be blocked until a token becomes available
	//because this is inside the processing goroutine there will be no guarantee of the order of operations
	<-bucket // <-channel indicates that this line of code is waiting for a response from the channel
	//Do work here
	fmt.Printf("Completed request %3d at %s\n", i, time.Since(startTime))
}
