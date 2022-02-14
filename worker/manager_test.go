package worker_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Suicide40/test1/worker"
)

type dummyWebHasher struct {
	CalcFn func(s string) ([]byte, error)
}

func (d dummyWebHasher) Calc(s string) ([]byte, error) { return d.CalcFn(s) }

func Test_workerManager_Run(t *testing.T) {
	t.Run("1 worker, 3 addresses", func(t *testing.T) {
		dummyHash := []byte("blablabla")

		workers := make([]worker.Hasher, 1)
		workers[0] = &dummyWebHasher{CalcFn: func(s string) ([]byte, error) {
			return dummyHash, nil
		}}

		addressCh := make(chan string, 10)
		resultCh := make(chan worker.Result, 10)

		wm := worker.NewManager(workers)
		go wm.Run(addressCh, resultCh)

		addressCh <- "www.test1.com"
		addressCh <- "www.test2.com"
		addressCh <- "www.test3.com"
		close(addressCh)

		for i := 0; i < 3; i++ {
			res := <-resultCh
			if !bytes.Equal(res.Md5hash, dummyHash) {
				t.Logf("results are different %s != %s", res.Md5hash, dummyHash)
				t.FailNow()
			}
		}

		_, ok := <-resultCh
		if ok {
			t.Log("result channel was not closed")
			t.FailNow()
		}
	})

	t.Run("3 worker, 1 address", func(t *testing.T) {
		dummyHash := []byte("blablabla")

		workers := make([]worker.Hasher, 3)
		workers[0] = &dummyWebHasher{CalcFn: func(s string) ([]byte, error) {
			return dummyHash, nil
		}}
		workers[1] = &dummyWebHasher{CalcFn: func(s string) ([]byte, error) {
			return dummyHash, nil
		}}
		workers[2] = &dummyWebHasher{CalcFn: func(s string) ([]byte, error) {
			return dummyHash, nil
		}}

		addressCh := make(chan string, 10)
		resultCh := make(chan worker.Result, 10)

		wm := worker.NewManager(workers)
		go wm.Run(addressCh, resultCh)

		addressCh <- "www.test1.com"
		close(addressCh)

		res := <-resultCh
		if !bytes.Equal(res.Md5hash, dummyHash) {
			t.Logf("results are different %s != %s", res.Md5hash, dummyHash)
			t.FailNow()
		}

		_, ok := <-resultCh
		if ok {
			t.Log("result channel was not closed")
			t.FailNow()
		}
	})

	t.Run("1 worker, 1 address and error", func(t *testing.T) {
		hashError := errors.New("some error text")
		workers := make([]worker.Hasher, 1)
		workers[0] = &dummyWebHasher{CalcFn: func(s string) ([]byte, error) {
			return nil, hashError
		}}

		addressCh := make(chan string, 10)
		resultCh := make(chan worker.Result, 10)

		wm := worker.NewManager(workers)
		go wm.Run(addressCh, resultCh)

		addressCh <- "www.test1.com"
		close(addressCh)

		res := <-resultCh
		if res.Err == nil {
			t.Logf("error is skipped")
			t.FailNow()
		}

		_, ok := <-resultCh
		if ok {
			t.Log("result channel was not closed")
			t.FailNow()
		}
	})
}
