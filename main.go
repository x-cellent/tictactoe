package main

import (
	"flag"
	"gitlab.com/kolls/networking/grpc/client"
	"gitlab.com/kolls/networking/grpc/server"
	"sync"
)

var (
	numClients int
)

func init() {
	flag.IntVar(&numClients, "c", 1, "number of clients")
}

func main() {
	flag.Parse()

	listener := server.Run()
	wg := new(sync.WaitGroup)
	wg.Add(numClients)
	for i := 1; i <= numClients; i++ {
		i := i // rebidding i so the goroutines will not share it
		go func() {
			client.Run(i, listener)
			wg.Done()
		}()
	}
	wg.Wait()
}
