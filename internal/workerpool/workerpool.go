package workerpool

import (
	"container/list"
	"sync"
	"time"

	"github.com/CesarDelgadoM/generator-reports/config"
)

// this code is inspired of the github repository: https://github.com/gammazero/workerpool

type Task func()

type WorkerPool struct {
	maxWorkers   int
	idleTimeout  time.Duration
	taskQueue    chan func()
	workerQueue  chan func()
	waitingQueue *list.List
}

func NewWorkerPool(config *config.WorkerConfig) *WorkerPool {
	if config.Pool < 1 {
		config.Pool = 1
	}

	pool := WorkerPool{
		maxWorkers:   config.Pool,
		idleTimeout:  time.Duration(config.Idle) * time.Second,
		taskQueue:    make(chan func()),
		workerQueue:  make(chan func()),
		waitingQueue: list.New(),
	}

	go pool.dispatch()

	return &pool
}

func (w *WorkerPool) Submit(task func()) {
	if task != nil {
		w.taskQueue <- task
	}
}

func (w *WorkerPool) dispatch() {
	var wg sync.WaitGroup
	var workerCount int
	var idle bool

	timeout := time.NewTimer(w.idleTimeout)

	for {
		if w.waitingQueue.Len() != 0 {
			w.processWaitingQueue()
			continue
		}

		select {
		case task := <-w.taskQueue:
			select {
			// if a worker is listenig to the workerqueue, send a task
			case w.workerQueue <- task:

			default:
				if workerCount < w.maxWorkers {
					wg.Add(1)
					go w.worker(task, w.workerQueue, &wg)
					workerCount++
				} else {
					w.waitingQueue.PushBack(task)
				}
			}
			idle = false
		// Time of the idle worker
		case <-timeout.C:
			if idle && workerCount > 0 {
				if w.killedIdleWorker() {
					workerCount--
				}
			}
			idle = true
			timeout.Reset(w.idleTimeout)
		}
	}
}

// Goroutine worker
func (w *WorkerPool) worker(task func(), workerQueue chan func(), wg *sync.WaitGroup) {
	for task != nil {
		task()
		task = <-workerQueue
	}
	wg.Done()
}

func (w *WorkerPool) processWaitingQueue() {
	select {
	case task := <-w.taskQueue:
		w.waitingQueue.PushBack(task)

	case w.workerQueue <- w.waitingQueue.Front().Value.(func()):
		w.waitingQueue.Remove(w.waitingQueue.Front())
	}
}

// Send a nil value to kill an idle worker
func (w *WorkerPool) killedIdleWorker() bool {
	select {
	case w.workerQueue <- nil:
		return true

	default:
		return false
	}
}
