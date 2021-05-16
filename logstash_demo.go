package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// Logstash的TCP连接
type Logstash struct {
	Hostname   string
	Port       int
	Connection *net.TCPConn
	Timeout    int
}

// 创建一个Logstash连接
func New(hostname string, port int, timeout int) *Logstash {
	l := Logstash{}
	l.Hostname = hostname
	l.Port = port
	l.Connection = nil
	l.Timeout = timeout
	return &l
}

// 设置连接超时
func (l *Logstash) setTimeouts() {
	deadline := time.Now().Add(time.Duration(l.Timeout) * time.Millisecond)
	_ = l.Connection.SetDeadline(deadline)
	_ = l.Connection.SetWriteDeadline(deadline)
	_ = l.Connection.SetReadDeadline(deadline)
}

// 创建TCP连接
func (l *Logstash) Connect() (*net.TCPConn, error) {
	var connection *net.TCPConn
	service := fmt.Sprintf("%s:%d", l.Hostname, l.Port)
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return connection, err
	}
	connection, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return connection, err
	}
	if connection != nil {
		l.Connection = connection
		_ = l.Connection.SetLinger(0) // default -1
		_ = l.Connection.SetNoDelay(true)
		_ = l.Connection.SetKeepAlive(true)
		_ = l.Connection.SetKeepAlivePeriod(time.Duration(5) * time.Second)
		l.setTimeouts()
	}
	return connection, err
}

// 写入数据
func (l *Logstash) Writeln(message string) error {
	var err = errors.New("tpc connection is nil")
	message = fmt.Sprintf("%s\n", message)
	if l.Connection != nil {
		_, err = l.Connection.Write([]byte(message))
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				_ = l.Connection.Close()
				l.Connection = nil
			} else {
				_ = l.Connection.Close()
				l.Connection = nil
				return err
			}
		} else {
			// Successful write! Let's extend the timeout.
			l.setTimeouts()
			return nil
		}
	}
	return err
}

func main() {
	l := New("192.168.24.88", 5044, 5000)
	if _, err := l.Connect(); err != nil {
		panic(err)
	}

	if err := l.Writeln(`{ "haha" : "gaga", "service": "test-service" }`); err != nil {
		panic(err)
	}
}
