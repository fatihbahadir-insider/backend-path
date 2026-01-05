package workers

import (
	"backend-path/app/models"
	"backend-path/utils"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type TransactionJob struct {
	ID uuid.UUID
	Type models.TransactionType
	FromUserID *uuid.UUID
	ToUserID *uuid.UUID
	Amount float64
	ResultChan chan TransactionResult
}

type TransactionResult struct {
	Transaction *models.Transaction
	Error error
}

type TransactionStats struct {
	TotalProcessed int64
	TotalSuccessful int64
	TotalFailed int64
	TotalCredited int64
	TotalDebited int64
	TotalTransferred int64
}

type TransactionWorkerPool struct {
	jobQueue chan TransactionJob
	workerCount int
	wg sync.WaitGroup
	stats *TransactionStats
	processor func(job TransactionJob) TransactionResult
	running bool
	mu sync.RWMutex
}

func NewTransactionWorkerPool(workerCount, queueSize int, processor func(job TransactionJob) TransactionResult) *TransactionWorkerPool {
	return &TransactionWorkerPool{
		jobQueue: make(chan TransactionJob, queueSize),
		workerCount: workerCount,
		stats: &TransactionStats{},
		processor: processor,
	}
}

func (p *TransactionWorkerPool) Start() {
	p.mu.Lock()

	if p.running {
		p.mu.Unlock()
		return
	}

	p.running = true
	p.mu.Unlock()

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	utils.Logger.Info("Transaction worker pool started with " + string(rune(p.workerCount+'0')) + " workers")
}

func (p *TransactionWorkerPool) worker(id int) {
	defer p.wg.Done()

	for job := range p.jobQueue {
		startTime := time.Now()
		result := p.processor(job)

		atomic.AddInt64(&p.stats.TotalProcessed, 1)
		if result.Error == nil {
			atomic.AddInt64(&p.stats.TotalSuccessful, 1)
			p.updateAmountStats(job)
		} else {
			atomic.AddInt64(&p.stats.TotalFailed, 1)
		}

		if job.ResultChan != nil {
			job.ResultChan <- result
		}

		utils.Logger.Info("Worker " + string(rune(id+'0')) + " processed job in " + time.Since(startTime).String())
	}
}

func (p *TransactionWorkerPool) updateAmountStats(job TransactionJob) {
	amountCents := int64(job.Amount * 100)
	switch job.Type {
	case models.TxTypeDeposit:
		atomic.AddInt64(&p.stats.TotalCredited, amountCents)
	case models.TxTypeWithdraw:
		atomic.AddInt64(&p.stats.TotalDebited, amountCents)
	case models.TxTypeTransfer:
		atomic.AddInt64(&p.stats.TotalTransferred, amountCents)
	}
}

func (p *TransactionWorkerPool) Submit(job TransactionJob) {
	p.jobQueue <- job
}

func (p *TransactionWorkerPool) SubmitAndWait(job TransactionJob) TransactionResult {
	job.ResultChan = make(chan TransactionResult, 1)
	p.jobQueue <- job
	return <-job.ResultChan
}

func (p *TransactionWorkerPool) GetStats() TransactionStats {
	return TransactionStats{
		TotalProcessed:   atomic.LoadInt64(&p.stats.TotalProcessed),
		TotalSuccessful:  atomic.LoadInt64(&p.stats.TotalSuccessful),
		TotalFailed:      atomic.LoadInt64(&p.stats.TotalFailed),
		TotalCredited:    atomic.LoadInt64(&p.stats.TotalCredited),
		TotalDebited:     atomic.LoadInt64(&p.stats.TotalDebited),
		TotalTransferred: atomic.LoadInt64(&p.stats.TotalTransferred),
	}
}

func (p *TransactionWorkerPool) QueueLength() int {
	return len(p.jobQueue)
}

func (p *TransactionWorkerPool) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	close(p.jobQueue)
	p.wg.Wait()
	utils.Logger.Info("Transaction worker pool stopped")
}