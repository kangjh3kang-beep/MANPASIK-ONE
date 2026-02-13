package proxy

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCPool manages reusable gRPC connections to backend microservices.
type GRPCPool struct {
	conns map[string]*grpc.ClientConn
}

// NewGRPCPool creates a connection pool to the given services.
// services is a map of serviceName → "host:port".
func NewGRPCPool(services map[string]string) (*GRPCPool, error) {
	pool := &GRPCPool{conns: make(map[string]*grpc.ClientConn)}

	for name, addr := range services {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err != nil {
			// Non-fatal: service might not be running yet.
			log.Printf("[grpc-pool] warning: could not connect to %s at %s: %v", name, addr, err)
			continue
		}
		log.Printf("[grpc-pool] connected to %s at %s", name, addr)
		pool.conns[name] = conn
	}

	return pool, nil
}

// GetConn returns a connection for the named service.
func (p *GRPCPool) GetConn(name string) (*grpc.ClientConn, bool) {
	conn, ok := p.conns[name]
	return conn, ok
}

// Close closes all connections in the pool.
func (p *GRPCPool) Close() {
	for name, conn := range p.conns {
		if err := conn.Close(); err != nil {
			log.Printf("[grpc-pool] error closing %s: %v", name, err)
		}
	}
}
