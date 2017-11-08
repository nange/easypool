package easypool

import (
	"log"
	"net"
	"time"
)

type PoolConn struct {
	net.Conn
	hp          *heapPool
	updatedtime time.Time
}

func (pc *PoolConn) Close() error {
	pc.updatedtime = time.Now()

	if err := pc.hp.put(pc); err != nil {
		log.Printf("put conn failed:%v\n", err)
		pc.hp = nil
		return pc.close()
	}
	return nil
}

func (pc *PoolConn) close() error {
	return pc.Conn.Close()
}
