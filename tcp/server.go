package tcp

import (
	"fmt"
	"net"
	"strings"

	"hipster-cache/common"
	"hipster-cache/hash_table"
	"hipster-cache/hash_table/value_type"
)

type CacheServer struct {
	port      int
	logger    common.ILogger
	listener  *net.TCPListener
	hashTable *hash_table.HashTable
}

type ClientMessage struct {
	command string
	params  []string
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
		// Remove last symbol => "\n"
		command := string(buf[0:n])
		// SET mykey "Hello sdf sdf d df df"
		// SET
		// get and replace
		// split раз и split
		// GET nonexisting
		fmt.Printf(`Response "%s"`, command)
		response,err := s.getResponse(command)
		if err !=  nil {
			response = err.Error()
		}
		conn.Write([]byte(response + "\n"))
		//	time.Sleep(time.Second * 10)
		//		conn.Close()
		//		return
	}
	return
}

func (s *CacheServer) getResponse(command string) (string,error) {
	clientMessage := NewClientMessage()
	if err := clientMessage.Init(command); err != nil {
		return "",err
	}

	switch clientMessage.command {
		case "GET":
			if len(clientMessage.params) != 2 {
				return "",fmt.Errorf("Error: incorrect parametes count")
			}
			GetStringOperation := value_type.NewGetStringOperation()
			s.hashTable.GetElement(clientMessage.params[0],clientMessage.params[1], GetStringOperation)
			return GetStringOperation.GetResult()
	}
	return "No error",nil
}

func NewClientMessage() *ClientMessage {
	return &ClientMessage{}
}

func (m *ClientMessage) Init(value string) error {
	words := m.splitMessageBySpaceAndQuates(value)
	fmt.Printf(`\n Words :"%#v"`, words)
	if len(words) == 0 {
		return fmt.Errorf(`Error: you don't set the command`)
	}
	if len(words) < 2 {
		return fmt.Errorf(`You don't set any parameters`)
	}
	m.command = strings.ToUpper(words[0])
	m.params = words[1:]
	return nil
}

func (m *ClientMessage) splitMessageBySpaceAndQuates(message string) []string {
	words := []string{}
	var word string
	var character string
	delimeter := ""
	for _, characterCode := range message {
		character = string(characterCode)
		switch character {
		case ` `:
			if delimeter == "" {
				words = append(words, word)
				word = ""
				break
			}
			word += character
		case `"`:
			if delimeter == character {
				delimeter = ""
				break
			}
			if delimeter == "" {
				delimeter = character
				break
			}
		case "\n", "\r":
		default:
			word += character

		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}
