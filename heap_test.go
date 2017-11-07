package easypool

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	InitialCap = 5
	MaxCap     = 30
	MaxIdle    = 5
	network    = "tcp"
	address    = "127.0.0.1:7777"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
)

func init() {
	// used for factory function
	go simpleTCPServer()
	time.Sleep(time.Second) // wait until tcp server has been settled
}

func TestNew(t *testing.T) {
	_, err := newHeapPool()
	if err != nil {
		t.Errorf("New error: %v", err)
	}
}

func TestPool(t *testing.T) {
	p, _ := newHeapPool()
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	_, ok := conn.(*PoolConn)
	if !ok {
		t.Errorf("Conn is not of type poolConn")
	}

	if err := conn.Close(); err != nil {
		t.Errorf("Pool Conn close error:%v", err)
	}

	if p.Len() != InitialCap {
		t.Errorf("Pool size is invalid, size:%v", p.Len())
	}
}

func newHeapPool() (Pool, error) {
	return NewHeapPool(&PoolConfig{
		InitialCap:  InitialCap,
		MaxCap:      MaxCap,
		MaxIdle:     MaxIdle,
		Idletime:    10 * time.Second,
		MaxLifetime: time.Minute,
		Factory:     factory,
	})
}

func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}()
	}
}
