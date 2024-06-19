package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/SultanYerken/sum-of-numbers/internal/models"
	"log"
	"os"
	"strconv"
	"sync"
)

type Service interface {
	GetNumbersSum() error
}

type ServiceImpl struct{}

func NewService() Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) GetNumbersSum() error {
	amountGorout, err := s.getAmountGorout()
	if err != nil {
		return err
	}

	sum, err := s.sumNumbers(amountGorout)
	if err != nil {
		return err
	}

	log.Println("sum of all values in the file:", sum)

	return nil
}

func (s *ServiceImpl) getAmountGorout() (int, error) {
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

func (s *ServiceImpl) sumNumbers(amountGorout int) (int, error) {
	file, err := os.Open("numbers.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var (
		allSum int
		wg     sync.WaitGroup
		mut    sync.Mutex
		ch     = make(chan models.Numbers)
	)

	for i := 0; i < amountGorout; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			localSum := 0
			for num := range ch {
				localSum += num.A + num.B
			}

			mut.Lock()
			allSum += localSum
			mut.Unlock()
		}()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case len(line) == 0 || line == "[" || line == "]":
			continue
		case line[len(line)-1] == ',':
			line = line[:len(line)-1]
		}

		var num models.Numbers
		if err := json.Unmarshal([]byte(line), &num); err != nil {
			return 0, fmt.Errorf("error decoding JSON: %w, line: %s", err, line)
		}
		ch <- num
	}

	close(ch)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return allSum, nil
}
