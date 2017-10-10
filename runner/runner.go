package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	contenders := []string{"../examples/random.py", "../examples/titfortat.py", "../examples/allcoop.py", "../examples/alldefect.py", "../examples/invalidcommand.py", "../examples/invalidquit.py"}
	for i := 0; i < len(contenders); i++ {
		for j := 0; j < len(contenders); j++ {
			if i == j {
				continue
			}
			a, b, rounds, adisqual, bdisqual := playAgainst(contenders[i], contenders[j])
			if adisqual {
				fmt.Println("WARNING contender A: ", contenders[i], " IS DISQUALIFIED")
				continue
			}
			if bdisqual {
				fmt.Println("WARNING contender B: ", contenders[j], " IS DISQUALIFIED")
				continue
			}
			fmt.Println(contenders[i], "vs", contenders[j])
			fmt.Println(contenders[i], "avg:", a)
			fmt.Println(contenders[j], "avg:", b)
			fmt.Println("Rounds:", rounds)
			fmt.Println()
		}
	}

}

var imagID string

func init() {
	imagIDRaw, err := exec.Command("docker", "build", "-q", "dock").CombinedOutput()
	if err != nil {
		fmt.Println(string(imagIDRaw))
		panic(err)
	}
	imagID = strings.Trim(string(imagIDRaw), "\n")
}

func makeMeAContainerBaby() (string, int) {

	for i := 0; i < 100; i++ {
		port := 5000 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(4000)
		containerRaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", strconv.Itoa(port)+":10000", imagID).CombinedOutput()
		if err != nil {

			if i < 99 {
				continue
			}
			fmt.Println(string(containerRaw))
			panic(err)
		}
		container := strings.Trim(string(containerRaw), "\n")
		return container, port

	}
	panic("shouldn't be able to get to here")

}

//takes in two executable paths
//starts up two of these dockers, with port 10000 in the VM bound to two random numbers on the host
//copies in the two executables to /code in both images
//makes a connection to localhost:whatever, which is docker:10000 on both. this makes relay.go start the provided /code, and forwards stdin/stdout over this socket
//runs iterated prisoner's dilemma, and finally returns number of rounds and final socre

func playAgainst(pathToExecA string, pathToExecB string) (AavgScore float32, BavgScore float32, numRounds int, aDisqualified bool, bDisqualified bool) {
	//fmt.Println("Image id:",imagID)

	containerA, portA := makeMeAContainerBaby()
	//fmt.Println(containerA)
	err := exec.Command("docker", "cp", pathToExecA, containerA+":/code").Run()
	if err != nil {
		panic(err)
	}

	containerB, portB := makeMeAContainerBaby()

	err = exec.Command("docker", "cp", pathToExecB, containerB+":/code").Run()
	if err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Millisecond) // go run relay.go may take time to start

	// cd dock && docker build .
	//dockerImage := "9e983c4697fc"

	// docker run --rm -p 5000:10000 $(dockerImage)
	//containerA := ""
	// docker cp $(pathToExecA) $(containerA):/code

	// docker run --rm -p 6000:10000 $(dockerImage)
	//containerB := ""
	// docker cp $(pathToExecB) $(containerB):/code

	A, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portA))
	if err != nil {
		panic(err)
	}
	defer A.Close()
	B, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portB))
	if err != nil {
		panic(err)
	}
	defer B.Close()

	AScore := 0
	BScore := 0

	numTurns := 0
	A.SetDeadline(time.Now().Add(2 * time.Second)) // generous deadline for first move
	B.SetDeadline(time.Now().Add(2 * time.Second))
	Am := make([]byte, 1)
	Bm := make([]byte, 1)
	for {

		_, err = io.ReadFull(A, Am)
		if err != nil {
			fmt.Println("Disqualified due to stream closed")
			return 0, 0, 0, true, false
		}

		_, err = io.ReadFull(B, Bm)
		if err != nil {
			fmt.Println("Disqualified due to stream closed")
			return 0, 0, 0, false, true
		}

		if Am[0] != 0 && Am[0] != 1 {
			fmt.Println("A made invalid move")
			return 0, 0, 0, true, false
		}

		if Bm[0] != 0 && Bm[0] != 1 {
			fmt.Println("B made invalid move")
			return 0, 0, 0, false, true
		}

		Amove := Am[0] == 1
		Bmove := Bm[0] == 1

		AScore += value(Amove, Bmove)
		BScore += value(Bmove, Amove)

		if Amove { // let the other one know what just happened
			B.Write([]byte{1})
		} else {
			B.Write([]byte{0})
		}
		if Bmove {
			A.Write([]byte{1})
		} else {
			A.Write([]byte{0})
		}

		numTurns++
		if rnd.Float32() < 0.0003 && numTurns > 100 {
			break
		}
		A.SetDeadline(time.Now().Add(50 * time.Millisecond)) // short deadline for subsequent moves
		B.SetDeadline(time.Now().Add(50 * time.Millisecond))
	}

	return float32(AScore) / float32(numTurns), float32(BScore) / float32(numTurns), numTurns, false, false
}

const R = 3

const S = 0
const T = 5

const P = 1

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func value(my bool, their bool) int {
	if my {
		if their {
			return R
		} else {
			return S
		}
	} else {
		if their {
			return T
		} else {
			return P
		}
	}
}
