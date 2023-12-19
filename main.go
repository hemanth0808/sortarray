
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)
	http.ListenAndServe(":8000", nil)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		ToSort [][]int `json:"to_sort"`
	}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	sortedArrays := make([][]int, len(inputData.ToSort))
	for i, arr := range inputData.ToSort {
		sortedArrays[i] = make([]int, len(arr))
		copy(sortedArrays[i], arr)
		sort.Ints(sortedArrays[i])
	}

	elapsedTime := time.Since(startTime).Nanoseconds()

	response := createResponse(sortedArrays, elapsedTime)
	encodeJSON(w, response)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		ToSort [][]int `json:"to_sort"`
	}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	sortedArrays := make([][]int, len(inputData.ToSort))
	ch := make(chan int, len(inputData.ToSort))

	for i, arr := range inputData.ToSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sortedArray := make([]int, len(arr))
			copy(sortedArray, arr)
			sort.Ints(sortedArray)
			sortedArrays[i] = sortedArray
			ch <- 1
		}(i, arr)
	}

	wg.Wait()
	close(ch)

	elapsedTime := time.Since(startTime).Nanoseconds()

	response := createResponse(sortedArrays, elapsedTime)
	encodeJSON(w, response)
}

func createResponse(sortedArrays [][]int, elapsedTime int64) map[string]interface{} {
	response := map[string]interface{}{
		"sorted_arrays": sortedArrays,
		"time_ns":       elapsedTime,
	}
	return response
}

func encodeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

