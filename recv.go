package main

import "fmt"

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

	fmt.Println("Continuing")
}
