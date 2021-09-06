package grpc

import (
	"time"

	ugrpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

// DialGrpcInsecure will connect to a local non-encrypted gRPC server and will return the connection.
// After this, you will typically do the following with the returned conn:
//
//     grpcClient = service_grpc.NewServiceAPIClient(conn)
//
func DialGrpcInsecure(endpoint, userAgent string) (*ugrpc.ClientConn, error) {
	opts := []ugrpc.DialOption{
		ugrpc.WithInsecure(), // non-encrypted local connection
		ugrpc.WithConnectParams(ugrpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   2 * time.Second,
			},
			MinConnectTimeout: 1 * time.Second,
		}),
		ugrpc.WithNoProxy(),
		ugrpc.WithUserAgent(userAgent),
	}
	return ugrpc.Dial(endpoint, opts...)
}
