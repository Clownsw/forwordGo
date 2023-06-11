package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

var targetAddress *net.TCPAddr

func Close[T io.Closer](t T) {
	_ = t.Close()
}

func connectionTarget() (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, targetAddress)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func forword(source net.Conn, target net.Conn, is bool) {
	var buf = make([]byte, 1024)

	for true {
		read, err := source.Read(buf)

		if err != nil {
			println(err)
			return
		}

		if is {
			fmt.Println("is read: ", read)
		} else {
			fmt.Println("read: ", read)
		}

		write, err := target.Write(buf[:read])
		if err != nil {
			println(err)
			return
		}

		if is {
			fmt.Println("write: ", write)
		} else {
			fmt.Println("write: ", write)
		}
	}

	Close(source)
	Close(target)
}

func HandleClient(conn net.Conn) {
	targetConnection, err := connectionTarget()

	if err != nil {
		println(err)
		Close(conn)
	}

	fmt.Println("start forword thread")

	go forword(conn, targetConnection, false)
	go forword(targetConnection, conn, true)
}

func main() {
	port, err := strconv.Atoi(os.Getenv("port"))

	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("0.0.0.0:%d", port)

	targetAddress, err = net.ResolveTCPAddr("tcp", os.Getenv("targetAddress"))

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}

	fmt.Println("listener tcp in", address)

	defer Close(listener)

	for true {
		socket, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go HandleClient(socket)
	}
}
