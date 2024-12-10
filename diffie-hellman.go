package main

import (
	"crypto/des"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"sync"
)

const q = 11027
const a = 53

func DiffieHellman(send chan<- Message, recv <-chan Message, wg *sync.WaitGroup, name string) {
	defer wg.Done()

	fmt.Printf("[%s] q = %d, a = %d\n", name, q, a)

	privateKey, err := rand.Int(rand.Reader, big.NewInt(q))
	if err != nil {
		panic("DH rand error")
	}
	publicKey := int64(math.Pow(a, float64(privateKey.Int64()))) % q

	send <- Message{
		messageType: "PK",
		value:       publicKey,
	}

	otherPublicKey, ok := <-recv
	if !ok {
		panic("DH failed to recieve public key")
	}

	K := int64(math.Pow(float64(otherPublicKey.value), float64(privateKey.Int64()))) % q

	fmt.Printf("[%s] Found K: %d\n", name, K)

	value := make([]byte, 8)
	rand.Read(value)

	valueInt := int64(binary.LittleEndian.Uint64(value))

	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(K))
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	toSend := make([]byte, 8)
	block.Encrypt(toSend, value)

	toSendInt := int64(binary.LittleEndian.Uint64(toSend))

	fmt.Printf("[%s] Sending %d ciphered as %d\n", name, valueInt, toSendInt)

	send <- Message{
		messageType: "Value",
		value:       toSendInt,
	}

	msg := <-recv

	binary.LittleEndian.PutUint64(value, uint64(msg.value))

	result := make([]byte, 8)
	block.Decrypt(result, value)

	resultInt := int64(binary.LittleEndian.Uint64(result))

	fmt.Printf("[%s] Received %d which value is %d\n", name, msg.value, resultInt)
}
