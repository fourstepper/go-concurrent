package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	UserID    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func getTodo(id int, ch chan Todo, wg *sync.WaitGroup) {
	defer wg.Done()
	var s Todo
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%v", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(data, &s)
	ch <- s
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan Todo)

	for i := 1; i < 6; i++ {
		wg.Add(1)
		go getTodo(i, ch, &wg) // produce results into channel "ch"
	}

	go func() {
		wg.Wait() // wait until all job results are "produced" and close the channel
		close(ch)
	}()

	for elem := range ch { // consume all results until channel is closed
		fmt.Println(elem)
	}
}
