package examples

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

const (
	nWorker = 5
	Second  = int64(time.Second)
)

type Request struct {
	fn func() int // The operation to perform
	c  chan int   // The channel to return the result
}

type Worker struct {
	requests chan Request // work to do (buffered channel)
	pending  int          // count of pending tasks
	index    int          // index in the heap (required by heap.Interface)
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests // get Request from balancer
		req.c <- req.fn()   // call fn and send result
		done <- w           // we've finished this request
	}
}

type Pool []*Worker

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *Pool) Push(x any) {
	n := len(*p)
	worker := x.(*Worker)
	worker.index = n
	*p = append(*p, worker)
}

func (p *Pool) Pop() any {
	old := *p
	n := len(old)
	worker := old[n-1]
	worker.index = -1
	*p = old[0 : n-1]
	return worker
}

type Balancer struct {
	pool Pool
	done chan *Worker
}

func NewBalancer(numWorkers int) *Balancer {
	done := make(chan *Worker, numWorkers)
	balancer := &Balancer{
		pool: make(Pool, 0, numWorkers),
		done: done,
	}

	for i := 0; i < numWorkers; i++ {
		w := &Worker{
			requests: make(chan Request, 10),
			pending:  0,
			index:    i,
		}

		go w.work(done)

		heap.Push(&balancer.pool, w)
	}

	return balancer
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	if len(b.pool) == 0 {
		return
	}

	w := heap.Pop(&b.pool).(*Worker)
	w.requests <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
}

func (b *Balancer) PrintStatus() {
	fmt.Printf("Load Balancer Status - Workers: %d\n", len(b.pool))
	for i, w := range b.pool {
		fmt.Printf("  Worker %d: pending=%d, index=%d\n", i, w.pending, w.index)
	}
	fmt.Println()
}

func requester(work chan<- Request, clientID int, numRequests int) {
	for i := 0; i < numRequests; i++ {
		c := make(chan int)

		time.Sleep(time.Duration(rand.Int63n(nWorker*2)) * time.Millisecond * 100)

		work <- Request{workFn, c}

		result := <-c

		furtherProcess(result, clientID, i)
	}
}

func workFn() int {
	workTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
	time.Sleep(workTime)

	return rand.Intn(1000)
}

func furtherProcess(result int, clientID int, requestID int) {
	fmt.Printf("Client %d: Request %d completed with result: %d\n",
		clientID, requestID, result)
}

func LoadBalancerDemo() {
	work := make(chan Request, 100)

	balancer := NewBalancer(nWorker)

	go balancer.balance(work)

	fmt.Println("Initial status:")
	balancer.PrintStatus()

	numClients := 3
	requestsPerClient := 5

	fmt.Printf("Starting %d clients, each making %d requests...\n\n",
		numClients, requestsPerClient)

	for i := 0; i < numClients; i++ {
		go requester(work, i+1, requestsPerClient)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("\nFinal status:")
	balancer.PrintStatus()
}
