package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Server struct {
	logger *log.Logger
	config *Config
	pool   *SSDBPool
}

func NewServer(logger *log.Logger, config *Config, pool *SSDBPool) *Server {
	return &Server{
		logger: logger,
		config: config,
		pool:   pool,
	}
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		s.logger.Fatal(err)
	}
	defer l.Close()
	ln := l.(*net.TCPListener)
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			s.logger.Println("AcceptTCP error: ", err)
			continue
		}
		go s.pipe(conn)
	}
}

func (s *Server) pipe(conn *net.TCPConn) {
	addr := conn.RemoteAddr().String()
	defer func() {
		s.logger.Println("disconnected: ", addr)
		conn.Close()
	}()

	for {
		data, err := recv(conn)
		if err != nil {
			s.logger.Println("recv error: ", err)
			break
		}

		s.logger.Println("request  ", len(data), data)

		ret, err := s.pool.Do(data)
		if err != nil {
			s.logger.Println("ssdb do error: ", err)
			ret = []string{"error", err.Error()}
		}

		//	ret := []string{"ok", "value"}
		var bf bytes.Buffer
		for _, str := range ret {
			bf.WriteString(fmt.Sprintf("%d", len(str)))
			bf.WriteByte('\n')
			bf.WriteString(str)
			bf.WriteByte('\n')
		}
		bf.WriteByte('\n')
		_, er := conn.Write(bf.Bytes())
		if er != nil {
			s.logger.Println("response error: ", er)
		}
		s.logger.Println("response ", len(ret), ret)
	}
}

func recv(conn *net.TCPConn) ([]string, error) {
	var buffer bytes.Buffer
	tmp := make([]byte, 102400)
	for {
		data := parse(buffer)
		if data == nil || len(data) > 0 {
			return data, nil
		}
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		buffer.Write(tmp[0:n])
	}
}

func parse(buffer bytes.Buffer) []string {
	data := []string{}
	buf := buffer.Bytes()
	idx, offset := 0, 0
	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(data) == 0 {
				continue
			} else {
				buffer.Next(offset)
				return data
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= buffer.Len() {
			break
		}

		v := buf[offset : offset+size]
		data = append(data, string(v))
		offset += size + 1
	}
	return []string{}
}
