package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Suicide40/test1/web"
	"github.com/Suicide40/test1/worker"
)

type inputData struct {
	workersLimit int
	addresses    []string
}

func main() {
	params := parseInputData()

	// build workers and prepare dependencies
	workers := make([]worker.Hasher, 0, params.workersLimit)
	for i := 0; i < params.workersLimit; i++ {
		newWorker := web.NewHasher(http.DefaultClient, md5.New())
		workers = append(workers, newWorker)
	}

	// prepare in and out channels
	addrCh := make(chan string, len(params.addresses))
	resCh := make(chan worker.Result, len(addrCh))

	// run worker manager
	wm := worker.NewManager(workers)
	go wm.Run(addrCh, resCh)

	// fill "in" channel
	for _, address := range params.addresses {
		addrCh <- address
	}
	close(addrCh)

	// receive results from the "out" channel
	for res := range resCh {
		if res.Err != nil {
			fmt.Printf("%s %v\n", res.Address, res.Err)
			continue
		}

		fmt.Printf("%s %s\n", res.Address, hex.EncodeToString(res.Md5hash))
	}
}
