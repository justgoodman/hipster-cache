package tcp

import (
	"fmt"
	"net"
	//	"time"

	"hipster-cache/common"
	"hipster-cache/hash_table"
)

type CacheServer struct {
	port      int
	logger    common.ILogger
	listener  *net.TCPListener
	hashTable *hash_table.HashTable
}

func NewCacheServer(hashTable *hash_table.HashTable, logger common.ILogger, port int) *CacheServer {
	return &CacheServer{port: port, logger: logger, hashTable: hashTable}
}

func (s *CacheServer) InitConnection() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(`:%d`, s.port))
	if err != nil {
		return err
	}

	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	return err
}

func (s *CacheServer) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Errorf(`Connection error: "%s"`, err.Error())
			continue
		}
		go s.handleMessage(conn)
		//      conn.Write([]byte("Bratish vse ok"))
		//      conn.Close()
	}
}

func (s *CacheServer) handleMessage(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			s.logger.Errorf(`Read message error: "%s"`, err.Error())
			if err = conn.Close(); err != nil {
				s.logger.Errorf(`Close connection error: "%s"`, err.Error())
			}
			return
		}
		fmt.Println(string(buf[0:n]))
		conn.Write([]byte(string(buf[0:n]) + " Bratish vse ok"))
		//	time.Sleep(time.Second * 10)
		//		conn.Close()
		//		return
	}
	return
}
