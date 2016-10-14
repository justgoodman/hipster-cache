package tcp

import (
	"net"
	"fmt"
//	"time"

        "hipster-cache/common"
)

type CacheServer struct {
	port int
	logger common.ILogger
	listener *net.TCPListener
}

func NewCacheServer(port int, logger common.ILogger) *CacheServer {
	return &CacheServer{port:port, logger:logger}
}

func (this *CacheServer) InitConnection() error {
     tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(`:%d`,this.port))
     if err != nil {
		return err
	}

     this.listener, err = net.ListenTCP("tcp", tcpAddr)
    return err
}

func (this *CacheServer) Run() {
      for {
        conn, err := this.listener.Accept()
        if err != nil {
                this.logger.Errorf(`Connection error: "%s"`, err.Error())
                continue
        }
        go this.handleMessage(conn)
//      conn.Write([]byte("Bratish vse ok"))
//      conn.Close()
      }
}

func (this *CacheServer) handleMessage(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			this.logger.Errorf(`Read message error: "%s"`, err.Error())
		}
		fmt.Println(string(buf[0:n]))
		conn.Write([]byte(string(buf[0:n]) + " Bratish vse ok"))
	//	time.Sleep(time.Second * 10)
//		conn.Close()
//		return
	}
	return
}

