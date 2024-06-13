package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("[ERROR]: %v", err)
		return
	}
}

func run() error {
	amountGorout, err := getAmountGorout()
	if err != nil {
		return err
	}

	data, err := getJSONData()
	if err != nil {
		return err
	}

	sum, err := getSumNumbers(amountGorout, data)
	if err != nil {
		return err
	}

	log.Println("sum of all values in the file:", sum)

	return nil
}

func getAmountGorout() (int, error) {
	if len(os.Args) != 2 {
		return 0, fmt.Errorf("the number of arguments must be to one")
	}

	arg := os.Args[1]

	amountGorout, err := strconv.Atoi(arg)
	if err != nil {
		return 0, err
	}

	if amountGorout < 1 || amountGorout > 1000000 {
		return 0, fmt.Errorf("the number of goroutines can be from 1 to 1000000")
	}

	return amountGorout, nil
}

func getJSONData() ([]Numbers, error) {
	data, err := os.ReadFile("numbers.json")
	if err != nil {
		return nil, err
	}

	var nums []Numbers
	if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&nums); err != nil {
		return nil, err
	}

	return nums, nil
}

func getSumNumbers(amountGorout int, numbers []Numbers) (int, error) {
	var (
		sum int = 0
		wg  sync.WaitGroup
		mut sync.Mutex
	)
	n := len(numbers) / amountGorout

	for i := 0; i < amountGorout; i++ {
		wg.Add(1)

		switch {
		case i == amountGorout-1:
			go func(ind int) {
				defer wg.Done()
				for j := ind * n; j < len(numbers); j++ {
					mut.Lock()
					sum += numbers[j].A + numbers[j].B
					mut.Unlock()
				}
			}(i)
		default:
			go func(ind int) {
				defer wg.Done()
				for j := ind * n; j < (ind+1)*n; j++ {
					mut.Lock()
					sum += numbers[j].A + numbers[j].B
					mut.Unlock()
				}
			}(i)
		}
	}
	wg.Wait()

	return sum, nil
}
