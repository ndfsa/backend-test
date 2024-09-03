package main

import (
	"log"

	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/common/repository"
)

type WorkerPool struct {
	jobQueue <-chan model.Transaction
	repo     repository.TransactionsRepository
	workers  int
	queue    int
	logger   *log.Logger
}

func (wp *WorkerPool) worker() {
	for transaction := range wp.jobQueue {
		if err := wp.repo.ExecuteTransaction(transaction); err != nil {
			wp.logger.Printf("transaction %sfailed: \n%s\n", transaction.Id.String(), err)
		}
	}
}

func NewWorkerPool(
	workers int,
	jobQueue <-chan model.Transaction,
	logger *log.Logger,
) WorkerPool {
	if workers < 1 {
		panic("number of workers must be >= 1")
	}

	wp := WorkerPool{
		jobQueue: jobQueue,
		workers:  workers,
	}

	for i := 0; i < workers; i++ {
		go wp.worker()
	}

	return wp
}
