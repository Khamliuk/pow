package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"pow/pkg/config"
	"pow/pkg/entity"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	address := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)

	err = run(context.Background(), address)
	if err != nil {
		fmt.Println("server error:", err)
	}
}

func run(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %w", err)
		}
		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := processMessage(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMessage(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

func processMessage(_ context.Context, msgStr string, resource string) (*entity.Message, error) {
	msg, err := entity.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	hasher, err := hash.NewHasher("sha1")
	if err != nil {
		return nil, fmt.Errorf("could not init hasher: %w", err)
	}
	p := pow.New(hasher)

	switch msg.Type {
	case entity.Quit:
		return nil, errors.New("connection was closed by client")
	case entity.RequestChallengeMessageType:
		secret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(100000))))
		hc, err := pow.InitHashcash(5, resource, pow.SignExt(secret, hasher))
		if err != nil {
			return nil, fmt.Errorf("could not init hashcahe: %d", err)
		}

		hcMarshaled, err := json.Marshal(hc)
		if err != nil {
			return nil, fmt.Errorf("could not marshal hashcash: %w", err)
		}

		return &entity.Message{
			Type:    entity.ResponseChallengeMessageType,
			Payload: string(hcMarshaled),
		}, nil

	case entity.RequestResourceMessageType:
		var hc pow.Hashcach
		err = json.Unmarshal([]byte(msg.Payload), &hc)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}

		err = p.Verify(&hc, resource)
		if err != nil {
			return nil, fmt.Errorf("invalid hashcash")
		}

		return &entity.Message{
			Type:    entity.ResponseResourceMessageType,
			Payload: entity.Quotes[rand.Intn(10)],
		}, nil
	default:
		return nil, fmt.Errorf("unknown message type")
	}
}

func sendMessage(msg entity.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.String())
	_, err := conn.Write([]byte(msgStr))
	return err
}
