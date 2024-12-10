package main

import (
	"fmt"
	"sync"
)

type Message struct {
	messageType string
	value       int64
}

func main() {
	var wg sync.WaitGroup
	commChan1 := make(chan Message, 10)
	commChan2 := make(chan Message, 10)

	fmt.Println("Diffie Hellman")

	wg.Add(2)

	go DiffieHellman(commChan1, commChan2, &wg, "Alice")
	go DiffieHellman(commChan2, commChan1, &wg, "Bob")

	wg.Wait()

	fmt.Println("\nRSA")

	wg.Add(2)

	go RSA(commChan1, commChan2, &wg, "Alice")
	go RSA(commChan2, commChan1, &wg, "Bob")

	wg.Wait()

	close(commChan1)
	close(commChan2)
}
