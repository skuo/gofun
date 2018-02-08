package util

import (
	"fmt"
	"strings"
	"sync"
)

var data = []string{
	"The yellow fish swims slowly in the water",
	"The brown dog barks loudly after a drink from its water bowl",
	"The dark bird of prey lands on a small tree after hunting for fish",
}

// =========================
// Two chans
// generator function that produces data
func words(stopCh chan struct{}, data []string) <-chan string {
	out := make(chan string)

	// splits line and emit words
	go func() {
		defer close(out) // closes channel upon fn return
		//count := 0
		for _, line := range data {
			words := strings.Split(line, " ")
			for _, word := range words {
				word = strings.ToLower(word)
				select {
				case out <- word:
					// no statement is needed
					//count++
				case <-stopCh:
					//fmt.Println("total words processed=", count)
					return
				}
			}
		}
	}()

	return out
}

func tryTwoChans() {

	histogram := make(map[string]int)
	stopCh := make(chan struct{})

	words := words(stopCh, data) // returns handle to channel
	for word := range words {
		if histogram["the"] == 3 {
			close(stopCh)
		}
		histogram[word]++
	}

	for k, v := range histogram {
		fmt.Printf("%s\t(%d)\n", k, v)
	}
}

// =========================
// chained chans
type histogram struct {
	total int
	freq  map[string]int
}

func (h *histogram) ingest() <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, line := range data {
			out <- line
		}
	}()
	return out
}

func (h *histogram) split(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for line := range in {
			for _, word := range strings.Split(line, " ") {
				out <- strings.ToLower(word)
			}
		}
	}()
	return out
}

func (h *histogram) count(in <-chan string) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for word := range in {
			h.freq[word]++
			h.total++
		}
		for k, v := range h.freq {
			fmt.Printf("chainedChan %s\t(%d)\n", k, v)
		}
	}()
	return done
}

func tryChainedChans() {
	h := &histogram{freq: make(map[string]int)}
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-h.count(h.split(h.ingest()))
	}()
	<-done
	fmt.Printf("Counted %d words!\n", h.total)
}

// =========================
// mutex

type Service struct {
	started bool
	stpCh   chan struct{}
	mutex   sync.RWMutex
	cache   map[int]string
}

func (s *Service) Start() {
	fmt.Println("In Start()")
	s.stpCh = make(chan struct{})
	s.cache = make(map[int]string)
	go func() {
		s.mutex.Lock()
		s.started = true
		s.cache[1] = "Hello World"
		s.cache[2] = "Hello Universe"
		s.cache[3] = "Hello Galaxy!"
		fmt.Println("Start() finished initialization")
		s.mutex.Unlock()
		<-s.stpCh // wait to be closed.
	}()
}

func (s *Service) Stop() {
	fmt.Println("In Stop()")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.started {
		s.started = false
		close(s.stpCh)
		fmt.Println("Stop() closed stopCh")
	}
}

func (s *Service) Serve(id int) {
	s.mutex.Lock()
	msg := s.cache[id]
	s.mutex.Unlock()
	if msg != "" {
		fmt.Println(msg)
	} else {
		fmt.Println("Hello, goodbye!")
	}
}

// tryMutex() does not really work as the goroutine in Start() does not
// get run before s.Serve(3) and s.Stop() is called.
func tryMutex() {
	s := &Service{}
	s.Start()
	s.Serve(3) // do some work
	s.Stop()
}

// =========================
// workgroup
const MAX = 300
const workers = 3

func tryWorkgroup() {
	fmt.Printf("\n\n--- tryWorkgroup() ---\n")
	values := make(chan int)
	result := make(chan int, workers)
	var wg sync.WaitGroup

	go func() { // gen multiple of 3 & 5 values
		for i := 1; i < MAX; i++ {
			if (i%3) == 0 && (i%5) == 0 {
				values <- i // push downstream
			}
		}
		close(values)
	}()

	work := func(index int) { // work unit, calc partial result
		defer wg.Done()
		r := 0
		for i := range values {
		    fmt.Println("index=", index, "i=", i)
			r += i
		}
		result <- r
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go work(i) // execute on its own thread
	}

	wg.Wait() // wait for both groutines
	close(result)
	total := 0
	for pr := range result {
		total += pr
	}
	fmt.Println("Total:", total)
}

// =========================
func TryConcurrent() {
	//tryTwoChans()
	//tryChainedChans()
	//tryMutex()
	tryWorkgroup()
}
