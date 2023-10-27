package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"pow/pkg/config"
	"pow/pkg/entity"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

const maxIterations = 1 << 30

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	address := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)

	err = runClient(context.Background(), address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}

func runClient(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// client will send new request every 10 seconds in loop until we will stop it
	for {
		message, err := handlePOWFlow(ctx, conn, conn)
		if err != nil {
			return err
		}
		fmt.Printf("quote: %s \n", message)
		time.Sleep(10 * time.Second)
	}
}

func handlePOWFlow(ctx context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	// First step: requesting challenge
	err := sendMessage(entity.Message{
		Type: entity.RequestChallengeMessageType,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}

	// Second step: got challenge, compute hashcash
	msgRawStr, err := readMessage(reader)
	if err != nil {
		return "", fmt.Errorf("could not read message: %w", err)
	}

	msg, err := entity.ParseMessage(msgRawStr)
	if err != nil {
		return "", fmt.Errorf("could not parse message: %w", err)
	}

	var hc pow.Hashcach
	err = json.Unmarshal([]byte(msg.Payload), &hc)
	if err != nil {
		return "", fmt.Errorf("could not unmarhal hashcash: %w", err)
	}

	hasher, err := hash.NewHasher("sha1")
	if err != nil {
		return "", fmt.Errorf("could not init hasher: %w", err)
	}
	p := pow.New(hasher)

	solution, err := p.Compute(ctx, &hc, maxIterations)
	if err != nil {
		return "", fmt.Errorf("could not coumpute hashcash: %w", err)
	}

	bytesSolution, err := json.Marshal(solution)
	if err != nil {
		return "", fmt.Errorf("could not marshal solution: %w", err)
	}

	// Third step: send challenge solution back to server
	err = sendMessage(entity.Message{
		Type:    entity.RequestResourceMessageType,
		Payload: string(bytesSolution),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// Fourth step: get result quote from server
	msgRawStr, err = readMessage(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	msg, err = entity.ParseMessage(msgRawStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return msg.Payload, nil
}

func readMessage(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func sendMessage(msg entity.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.String())
	_, err := conn.Write([]byte(msgStr))
	return err
}
