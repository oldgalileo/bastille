package main //this is not part of the package system it just gets copied into the docker image
import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Relay is listening on 10000")

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Relay has received its connection") // one time use

	os.Chmod("/code", 0777)

	signalChan := make(chan os.Signal, 1)

	timeToExit := make(chan bool)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...\n")
			timeToExit <- true
		}
	}()

	cmd := exec.Command("/code")

	cmd.Stderr = os.Stdout

	in, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer out.Close()

	go func() {
		io.Copy(in, conn)
		fmt.Println("Done copying from command to socket")
		in.Close() // we're done reading input, not quite done with the rest yet
	}()

	go func() {
		io.Copy(conn, out)
		fmt.Println("Done copying from socket to command")
		timeToExit <- true
	}()

	cmd.Start()
	go func() {
		cmd.Wait()
		fmt.Println("Command done waiting")
		timeToExit <- true
	}()

	<-timeToExit
	fmt.Println("Exiting")

}
