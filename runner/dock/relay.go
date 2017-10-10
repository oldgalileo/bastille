package main //this is not part of the package system it just gets copied into the docker image
import(
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
	"os/exec"
)


func main(){
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Relay is listening on 10000")
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Println("Relay has received its connection") // one time use

	cmd := exec.Command("/code")
	in, _ := cmd.StdinPipe()
    out, _ := cmd.StdoutPipe()

    go io.Copy(conn,in)
    go io.Copy(out,conn)

    cmd.Start()

}