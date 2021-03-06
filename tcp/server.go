package tcp

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"unicode/utf8"

	"hipster-cache/common"
	"hipster-cache/hash_table"
	"hipster-cache/hash_table/value_type"
)

const (
	ttlSeconds      = "EX"
	ttlMilliseconds = "PX"
	exitCommand     = "EXIT"
	pingCommand     = "PING"
	endSymbol       = "\n"
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
	var (
		buf           [512]byte
		clientMessage *ClientMessage
	)
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			s.logger.Errorf(`Read message error: "%s"`, err.Error())
			if err = conn.Close(); err != nil {
				s.logger.Errorf(`Close connection error: "%s"`, err.Error())
			}
			return
		}
		command := string(buf[0:n])
		fmt.Printf(`Response "%s"`, command)
		clientMessage, err = s.getClientMessage(command)
		if err != nil {
			conn.Write([]byte(err.Error() + endSymbol))
			return
		}
		if clientMessage.command == exitCommand {
			conn.Close()
			return
		}
		response, err := s.getResponse(clientMessage)
		if err != nil {
			response = err.Error()
		}
		conn.Write([]byte(response + endSymbol))
	}
	return
}

func (s *CacheServer) getClientMessage(command string) (*ClientMessage, error) {
	clientMessage := NewClientMessage()
	if err := clientMessage.Init(command); err != nil {
		return nil, err
	}
	return clientMessage, nil
}

func (s *CacheServer) getResponse(clientMessage *ClientMessage) (string, error) {
	// Check key lenght
	if len(clientMessage.params) >= 1 {
		if int(s.hashTable.MaxKeyLenght) < utf8.RuneCount([]byte(clientMessage.params[0])) {
			return "", fmt.Errorf(`Error: key lenght is more than maximum Lenght "%d"`, s.hashTable.MaxKeyLenght)
		}
	}

	switch clientMessage.command {
	case pingCommand:
		return `"pong"`, nil
	case value_type.GetStringCmdName:
		if len(clientMessage.params) != 1 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "1", was sended "%d"`, len(clientMessage.params))
		}
		getStringOperation := value_type.NewGetStringOperation()
		s.hashTable.GetElement(clientMessage.params[0], getStringOperation)
		return getStringOperation.GetResult()
	case value_type.SetStringCmdName:
		if len(clientMessage.params) != 2 && len(clientMessage.params) != 4 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "2 or 4", was sended "%d"`, len(clientMessage.params))
		}
		var ttl time.Duration
		// This command with TTL
		if len(clientMessage.params) == 4 {
			cmdDuration, _ := strconv.Atoi(clientMessage.params[3])
			if cmdDuration <= 0 {
				return "", fmt.Errorf(`Error: incorrect ttl time, ttl duration must me more  or equal 0, was sended "%s"`, clientMessage.params[3])
			}
			switch clientMessage.params[2] {
			case ttlSeconds:
				ttl = time.Second * time.Duration(cmdDuration)
			case ttlMilliseconds:
				ttl = time.Millisecond * time.Duration(cmdDuration)
			default:
				return "", fmt.Errorf(`Error: incorrect parameter name, must be "%s" or "%s", was sended "%s"r`, ttlSeconds, ttlMilliseconds, clientMessage.params[2])
			}
		}
		setStringOperation := value_type.NewSetStringOperation()
		//time.Now().Unix() + int64(10000)
		s.hashTable.SetElement(clientMessage.params[0], ttl, interface{}(clientMessage.params[1]), setStringOperation)
		return setStringOperation.GetResult()
	case value_type.PushListCmdName:
		if len(clientMessage.params) != 2 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "2", was sended "%d"`, len(clientMessage.params))
		}

		pushListOperation := value_type.NewPushListOperation()
		s.hashTable.SetElement(clientMessage.params[0], 0, interface{}(clientMessage.params[1]), pushListOperation)
		return pushListOperation.GetResult()
	case value_type.RangeListCmdName:
		if len(clientMessage.params) != 3 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "3", was sended "%d"`, len(clientMessage.params))
		}
		indexStart, err := strconv.Atoi(clientMessage.params[1])
		if err != nil {
			return "", fmt.Errorf(`Error: second parameter must be integer, was sended "%d"`, clientMessage.params[1])
		}

		indexEnd, err := strconv.Atoi(clientMessage.params[2])
		if err != nil {
			return "", fmt.Errorf(`Error: third parameter must be integer, was sended "%d"`, clientMessage.params[1])
		}

		rangeListOperation := value_type.NewRangeListOperation(indexStart, indexEnd)
		s.hashTable.GetElement(clientMessage.params[0], rangeListOperation)
		return rangeListOperation.GetResult()
	case value_type.LenghtListCmdName:
		if len(clientMessage.params) != 1 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "1", was sended "%d"`, len(clientMessage.params))
		}

		lenghtListOperation := value_type.NewLenghtListOperation()
		s.hashTable.GetElement(clientMessage.params[0], lenghtListOperation)
		return lenghtListOperation.GetResult()
	case value_type.SetListCmdName:
		if len(clientMessage.params) != 3 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "3", was sended "%d"`, len(clientMessage.params))
		}
		index, err := strconv.Atoi(clientMessage.params[1])
		if err != nil {
			return "", fmt.Errorf(`Error: second parameter must be integer, was sended "%d"`, clientMessage.params[1])
		}

		setListOperation := value_type.NewSetListOperation(index)

		s.hashTable.SetElement(clientMessage.params[0], 0, interface{}(clientMessage.params[2]), setListOperation)
		return setListOperation.GetResult()
	case value_type.SetDictCmdName:
		if len(clientMessage.params) != 3 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "3", was sended "%d"`, len(clientMessage.params))
		}
		setDictOperation := value_type.NewSetDictOperation(clientMessage.params[1])
		s.hashTable.SetElement(clientMessage.params[0], 0, interface{}(clientMessage.params[2]), setDictOperation)
		return setDictOperation.GetResult()
	case value_type.GetDictCmdName:
		if len(clientMessage.params) != 2 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "2", was sended "%d"`, len(clientMessage.params))
		}
		getDictOperation := value_type.NewGetDictOperation(clientMessage.params[1])
		s.hashTable.GetElement(clientMessage.params[0], getDictOperation)
		return getDictOperation.GetResult()
	case value_type.GetAllDictCmdName:
		if len(clientMessage.params) != 1 {
			return "", fmt.Errorf(`Error: incorrect parametes count need "1", was sended "%d"`, len(clientMessage.params))
		}
		getAllDictOperation := value_type.NewGetAllDictOperation()
		s.hashTable.GetElement(clientMessage.params[0], getAllDictOperation)
		return getAllDictOperation.GetResult()

	}
	return "Command not found", nil
}
