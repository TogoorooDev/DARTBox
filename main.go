package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	// "sync"

	"github.com/BurntSushi/toml"
)

func listenLoop(incoming chan confirmFormat, confirm chan bool, transferred chan int) {
	list, err := net.Listen("tcp", ":"+config.Server.Port)
	if err != nil {
		panic(err)
	}

	defer list.Close()

	for {
		conn, err := list.Accept()
		if err != nil {
			panic(err)
		}

		go handleConn(conn, incoming, confirm)
	}
}

func handleConn(c net.Conn, incoming chan confirmFormat, confirm chan bool) {
	defer c.Close()
	rdr := bufio.NewReader(c)
	fileName, err := rdr.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	if err != nil {
		panic(err)
	}

	out := confirmFormat{c.RemoteAddr().String(), fileName}

	incoming <- out
	c.Close()

	if !(<-confirm) {
		return
	}

	fmt.Printf("Transmission inititated")

	for {
		read, err := rdr.ReadByte()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("%c", read)
	}
}

func comSend(trans chan int, argc int, argv []string) {
	if argc < 3 {
		fmt.Println("usage: send file ip-address")
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
	outgoing.WriteString(strconv.FormatInt(size, 10) + "\n")

	var transferred int64
	transferred = 0

	for transferred < size {
		out, err := fileReader.ReadByte()
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

func comRecv(inc chan confirmFormat, cfrm chan bool, trans chan int) {
	fmt.Println("Awaiting connection...")
	req := <-inc
	var resp string
	fmt.Printf("Accept %s from %s[y/n]? ", req.Filename, req.Ip)
	fmt.Scan(&resp)
	for {
		var resp string
		fmt.Printf("Accept %s from %s[y/n]? ", req.Filename, req.Ip)
		fmt.Scan(&resp)
		if resp == "y" || resp == "Y" {
			fmt.Println("Accepted")
			cfrm <- true
			break
		} else if resp == "n" || resp == "N" {
			fmt.Println("Denied")
			cfrm <- false
			break
		} else {
			fmt.Println("Enter y or n")
			continue
		}
	}
}

func inputLoop(inc chan confirmFormat, cfrm chan bool, transChan chan int) {
	inReader := bufio.NewReader(os.Stdin)
	for {
		var resp string
		fmt.Printf("DARTBox 0.1 %% ")
		resp, err := inReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		comp := strings.ToLower(resp)

		if strings.HasPrefix(comp, "recv") {
			comRecv(inc, cfrm, transChan)
		} else if strings.HasPrefix(comp, "send") {
			fmt.Println("send")
			// Parse
			argv := strings.Split(comp, " ")
			argv[2] = strings.TrimSpace(argv[2])
			comSend(transChan, len(argv), argv)

		}

	}
}

func parseConfig() {
	// do some config stuff here
	configFileData, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(configFileData, &config)
	if err != nil {
		panic(err)
	}

}

func main() {
	writeLock.Lock()
	fmt.Println("DARTBox 0.1")
	fmt.Println("Reading config at config.toml")
	writeLock.Unlock()

	parseConfig()

	writeLock.Lock()
	fmt.Println("Finished reading config.toml")
	fmt.Printf("Options:\n   Port: %s\n", config.Server.Port)

	fmt.Println("Opening up to external peers")
	fmt.Println("Creating message channels")
	incomingChan := make(chan confirmFormat)
	confirmChan := make(chan bool)
	transChan := make(chan int)
	writeLock.Unlock()

	writeLock.Lock()
	fmt.Println("Channels created, starting network thread")
	writeLock.Unlock()
	go listenLoop(incomingChan, confirmChan, transChan)

	writeLock.Lock()
	fmt.Println("Initialization finished. Handing control to you")
	fmt.Println("lets hope this works")
	writeLock.Unlock()

	inputLoop(incomingChan, confirmChan, transChan)
}
