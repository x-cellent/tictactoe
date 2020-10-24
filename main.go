package main

import (
	"flag"
	"gitlab.com/kolls/networking/grpc/client"
	"gitlab.com/kolls/networking/grpc/server"
	"sync"
)

var (
	numClients int
	inMemory   bool
)

func init() {
	flag.IntVar(&numClients, "c", 1, "number of clients")
	flag.BoolVar(&inMemory, "m", false, "whether to run gRPC connections in-memory")
}

func main() {
	flag.Parse()

	listener := server.Run(inMemory)
	wg := new(sync.WaitGroup)
	wg.Add(numClients)
	for i := 1; i <= numClients; i++ {
		i := i // rebidding i so the goroutines will not share it
		go func() {
			client.Connect(i, listener, inMemory)
			wg.Done()
		}()
	}
	wg.Wait()
}
