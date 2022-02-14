package main

import (
	"flag"
)

const defaultWorkers = 10

func parseInputData() inputData {
	params := inputData{
		workersLimit: defaultWorkers,
	}

	parallel := flag.Int("parallel", defaultWorkers, "limit the number of parallel requests")
	flag.Parse()
	if parallel != nil {
		params.workersLimit = *parallel
	}

	params.addresses = flag.Args()

	return params
}
