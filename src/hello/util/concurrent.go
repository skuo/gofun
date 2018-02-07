package util

import (
	"fmt"
	"strings"
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


// =========================
func TryConcurrent() {
	tryTwoChans()
	tryChainedChans()
}
