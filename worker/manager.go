package worker

import (
	"sync"
)

type Hasher interface {
	Calc(string) ([]byte, error)
}

type Result struct {
	Address string
	Md5hash []byte
	Err     error
}

type Manager struct {
	workers []Hasher
}

func NewManager(workers []Hasher) *Manager {
	return &Manager{workers: workers}
}

// Run runs workers which one by one get items from input channel, calculate md5, and put results to output channel
func (m *Manager) Run(input <-chan string, output chan<- Result) {
	wg := sync.WaitGroup{}
	for _, w := range m.workers {
		wg.Add(1)
		go func(w Hasher) {
			for addr := range input {
				res := Result{
					Address: addr,
				}

				var err error
				res.Md5hash, err = w.Calc(addr)
				if err != nil {
					res.Err = err
				}
				output <- res
			}

			wg.Done()
		}(w)
	}

	wg.Wait()
	close(output)
}
