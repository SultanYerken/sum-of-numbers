package main

import (
	"github.com/SultanYerken/sum-of-numbers/internal/service"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("[ERROR]: %v", err)
		return
	}
}

func run() error {
	service := service.NewService()

	if err := service.GetNumbersSum(); err != nil {
		return err
	}

	return nil
}
