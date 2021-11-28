package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	//"log"
)

func comSend(trans chan int, argc int, argv []string) {
	if argc < 3 {
		fmt.Println("usage: send file ip-address port")
		return
	}

	file := argv[1]
	ip := argv[2]
	FInfo, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	size := FInfo.Size()

	sendPtr, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileReader := bufio.NewReader(sendPtr)

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println(err)
		return
	}

	outgoing := bufio.NewWriter(conn)
	outgoing.WriteString(strings.TrimSpace(file) + "\n")
	outgoing.Flush()

	outgoing.WriteString(strconv.FormatInt(size, 10) + "\n")
	outgoing.Flush()

	var transferred int64
	transferred = 0

	for transferred < size {
		out, err := fileReader.ReadByte()
		//fmt.Println(out)
		outArr := make([]byte, 1)
		outArr[0] = out
		if err != nil {
			fmt.Println(err)
			return
		}
		outgoing.Write(outArr)
		size++
	}
}
