/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import (
	"sync/atomic"
	"time"
)

/***********************
 *      WorkerPool     *
 ***********************/

type WorkerTask func()

type WorkerPool interface {
	// returns false if task dropped
	Submit(task WorkerTask) bool
	Shutdown()
}

/***********************
 *  Default WorkerPool *
 ***********************/

type defaultWorkerPool struct {
	logger Logger

	minWorkers int
	maxWorkers int
	queueCap   int

	workerCount        int32
	queue              chan WorkerTask
	queueGrowThreshold float64

	stopSignal   chan struct{}
	shutdownOnce atomic.Bool
	idleTimeout  time.Duration
}

type workerOption func(*defaultWorkerPool)

// WithMinWorkers sets min workers
func WithMinWorkers(_min int) workerOption {
	return func(p *defaultWorkerPool) {
		p.minWorkers = _min
	}
}

// WithMaxWorkers sets max workers
func WithMaxWorkers(_max int) workerOption {
	return func(p *defaultWorkerPool) {
		p.maxWorkers = _max
	}
}

// WithQueueCap sets queue capacity
func WithQueueCap(_cap int) workerOption {
	return func(p *defaultWorkerPool) {
		p.queueCap = _cap
	}
}

// WithIdleTimeout sets idle timeout for workers
func WithIdleTimeout(d time.Duration) workerOption {
	return func(p *defaultWorkerPool) {
		p.idleTimeout = d
	}
}

// WithQueueGrowThreshold sets the queue usage threshold at which
// the pool attempts to dynamically spawn a new worker.
// A value of 0.75 means new workers are added when the queue is 75% full.
func WithQueueGrowThreshold(threshold float64) workerOption {
	return func(p *defaultWorkerPool) {
		p.queueGrowThreshold = threshold
	}
}

// NewDefaultWorkerPool creates a new worker pool with options.
func NewDefaultWorkerPool(logger Logger, opts ...workerOption) WorkerPool {
	p := &defaultWorkerPool{
		logger:             logger,
		minWorkers:         10,
		maxWorkers:         300,
		queueCap:           200,
		idleTimeout:        10 * time.Second,
		stopSignal:         make(chan struct{}),
		queueGrowThreshold: 0.75,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.queue = make(chan WorkerTask, p.queueCap)

	for i := 0; i < p.minWorkers; i++ {
		p.addWorker()
	}

	return p
}

func (p *defaultWorkerPool) addWorker() {
	atomic.AddInt32(&p.workerCount, 1)

	go func() {
		idleTimer := time.NewTimer(p.idleTimeout)
		defer idleTimer.Stop()

		for {
			select {
			case task := <-p.queue:
				task()

				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(p.idleTimeout)

			case <-idleTimer.C:
				if atomic.LoadInt32(&p.workerCount) > int32(p.minWorkers) {
					atomic.AddInt32(&p.workerCount, -1)
					p.logger.Debug("WorkerPool: worker exited due to idle timeout")
					return
				}
				idleTimer.Reset(p.idleTimeout)

			case <-p.stopSignal:
				return
			}
		}
	}()
}

// Submit submits a task to the pool.
// Returns false if the queue is full and task dropped.
func (p *defaultWorkerPool) Submit(task WorkerTask) bool {
	if p.shutdownOnce.Load() {
		return false
	}

	if float64(len(p.queue)) >= float64(p.queueCap)*p.queueGrowThreshold {
		if atomic.LoadInt32(&p.workerCount) < int32(p.maxWorkers) {
			p.addWorker()
			p.logger.Debug("WorkerPool: spawned new worker due to high queue usage")
		}
	}

	select {
	case p.queue <- task:
		return true
	default:
		p.logger.Debug("WorkerPool: dropping task due to full queue")
		return false
	}
}

// Shutdown stops the pool immediately; no waiting for workers.
func (p *defaultWorkerPool) Shutdown() {
	if p.shutdownOnce.CompareAndSwap(false, true) {
		close(p.stopSignal)
	}
}
