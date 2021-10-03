package main

import (
	"net"
	"fmt"
	"io/ioutil"	
	
	"github.com/BurntSushi/toml"
)

func listen_loop(inc chan string, cfrm chan bool){
	list, err := net.Listen("tcp", ":" + config.Server.Port)
	if err != nil {
		panic(err)
	}

	defer list.Close()
	
	for {
		conn, err := list.Accept()
		if err != nil {
			panic (err)
		}

		inc <- conn.RemoteAddr().String()

		if <-cfrm {
			go handle_conn(conn)	
		}
	}
}

func handle_conn(c net.Conn){
	
}


func input_loop(inc chan string, cfrm chan bool){
	for {
		req_ip := <- inc
		for {
			var resp string
			fmt.Printf("Accept file request from %s? ", req_ip)
			fmt.Scan(&resp)
			if resp == "y" || resp == "Y" {
				fmt.Println("Accepted")	
				cfrm <- true
				break
			}else if resp == "n" || resp == "N" {
				fmt.Println("Denied")
				cfrm <- false
				break
			}else {
				fmt.Println("Enter y or n")

				continue
			}	
		}
	}
}

func parse_config(){
	// do some config stuff here
	config_file_data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(config_file_data, &config)
	if err != nil {
		panic(err)
	}
	
}

func main(){
	fmt.Println("DARTBox 0.1")
	fmt.Println("Reading config at config.toml")

	parse_config()

	fmt.Println("Finished reading config.toml")
	fmt.Printf("Options:\n   Port: %s\n", config.Server.Port)

	fmt.Println("Opening up to external peers")
	fmt.Println("Creating message channels")
	incoming_chan := make(chan string)
	confirm_chan := make(chan bool)

	fmt.Println("Channels created, starting network thread")
	go listen_loop(incoming_chan, confirm_chan)

	fmt.Println("Initialization finished. Handing control to you")
	fmt.Println("lets hope this works")
	
	input_loop(incoming_chan, confirm_chan)
}
