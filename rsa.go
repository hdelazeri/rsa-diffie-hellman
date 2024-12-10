package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"github.com/hdelazeri/rsa-diffie-hellman/math"
)

func RSA(send chan<- Message, recv <-chan Message, wg *sync.WaitGroup, name string) {
	defer wg.Done()

	p, err := rand.Prime(rand.Reader, 8)
	if err != nil {
		panic("RSA: failed to generate prime p")
	}

	q, err := rand.Prime(rand.Reader, 8)
	if err != nil {
		panic("RSA: failed to generate prime q")
	}

	n := p.Int64() * q.Int64()

	lambda_n := (p.Int64() - 1) * (q.Int64() - 1)

	e, _ := rand.Int(rand.Reader, big.NewInt(lambda_n))

	for gcd, _, _ := math.GCDExtended(int(e.Int64()), int(lambda_n)); gcd != 1; {
		e, _ = rand.Int(rand.Reader, big.NewInt(lambda_n))
		gcd, _, _ = math.GCDExtended(int(e.Int64()), int(lambda_n))
	}

	d := math.ModularInverse(int(e.Int64()), int(lambda_n))

	fmt.Printf("[%s] Public key is (%d,%d) and private key (%d, %d)\n", name, e.Int64(), n, d, n)

	send <- Message{
		messageType: "PK-E",
		value:       e.Int64(),
	}

	send <- Message{
		messageType: "PK-N",
		value:       n,
	}

	otherE := <-recv
	otherN := <-recv

	fmt.Printf("[%s] Other public key is (%d,%d)\n", name, otherE.value, otherN.value)

	toSend, _ := rand.Int(rand.Reader, big.NewInt(1000))
	ciphered := math.Exponentiation(toSend.Int64(), otherE.value, otherN.value)

	fmt.Printf("[%s] Sending %d ciphered as %d\n", name, toSend.Int64(), ciphered)

	send <- Message{
		messageType: "Value",
		value:       ciphered,
	}

	msg := <-recv

	deciphered := math.Exponentiation(msg.value, int64(d), n)

	fmt.Printf("[%s] Received %d wich value is %d\n", name, msg.value, deciphered)
}
