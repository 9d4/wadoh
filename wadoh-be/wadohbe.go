package wadohbe

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/9d4/wadoh/wadoh-be/pb"
)

// Client is wrapper around wadoh-be grpc client.
type Client struct {
	conn    *grpc.ClientConn
	Service pb.ControllerServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}
	client := &Client{}
	client.conn = conn
	client.Service = pb.NewControllerServiceClient(conn)
	return client, nil
}

type EventMessage interface {
	GetFrom() string
	GetJid() string
	GetMessage() string
}

// ReceiveMessage wraps pb.ControllerServiceClient.ReceiveMessage call with
// additional retry without closing channel when got retryable error code
// such Unavailable. But after maxretries limit hit, channel will be closed.
func (c *Client) ReceiveMessage() <-chan EventMessage {
	const backoff = 5 * time.Second
	const maxRetries = 10

	var recv pb.ControllerService_ReceiveMessageClient
	var connectErr error

	eventc := make(chan EventMessage, 10)

	retryableStatusCodes := map[codes.Code]bool{
		codes.Unavailable: true, // etc
	}

	connect := func() bool {
		for i := 0; i < maxRetries; i++ {
			recv, connectErr = c.Service.ReceiveMessage(context.Background(), nil)
			if retryableStatusCodes[status.Code(connectErr)] {
				log.Err(connectErr).Send()
				time.Sleep(backoff)
				continue
			}
			log.Info().Msg("wadohbe rpc: ReceiveMessage connection established")
			return true
		}
		close(eventc)
		log.Info().Msg("wadohbe rpc: closing event channel, max retries exceeded")
		return false
	}

	go func() {
		if !connect() {
			return
		}
		for {
			msg, err := recv.Recv()
			if err != nil {
				log.Error().Caller().Err(err).Send()
				if retryableStatusCodes[status.Code(err)] {
					if !connect() {
						return
					}
				}
				continue
			}
			log.Debug().Caller().Any("recv:msg", msg).Send()
			eventc <- msg
		}
	}()

	return eventc
}
