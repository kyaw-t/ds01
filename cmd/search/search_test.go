package search

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var host = "http://localhost:3000/"
var searchTerms = []string{
	"kubernetes",
	"istio",
	"prometheus",
	"nginx",
	"redis",
	"ls;adkls;kdfsl;adfld;kjsf;lskdj",
}
var registries = []string{
	"registry1",
	"registry2",
	"registry3",
	"registry4",
	"registry5",
	"registry6",
	"registry7",
	"n-registry1",
	"n-registry2",
	"n-registry3",
	"n-registry4",
	"n-registry5",
	"n-registry6",
	"n-registry7",
}

func BenchmarkNaiveSearch(b *testing.B) {

	var times []time.Duration
	var found int
	for _, searchTerm := range searchTerms {
		start := time.Now()
		results, err := naiveSearch(host, searchTerm, registries)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Since(start)
		times = append(times, t)
		found += len(results)
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	totalSeconds := total.Seconds() / float64(len(searchTerms))
	fmt.Printf("Found: %v. Average time: %v\n", found, totalSeconds)
}

func BenchmarkRoutinedSearch(b *testing.B) {

	var times []time.Duration
	var found int
	for _, searchTerm := range searchTerms {
		start := time.Now()
		results, err := routinedSearch(host, searchTerm, registries)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Since(start)
		times = append(times, t)
		found += len(results)
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	totalSeconds := total.Seconds() / float64(len(searchTerms))
	fmt.Printf("Found: %v. Average time: %v\n", found, totalSeconds)
}
