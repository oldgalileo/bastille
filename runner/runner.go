package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	contenders := []string{"../examples/random.py", "../examples/titfortat.py", "../examples/allcoop.py", "../examples/alldefect.py"}
	for i := 0; i < len(contenders); i++ {
		for j := 0; j < len(contenders); j++ {
			if i == j {
				continue
			}
			a, b, rounds := playAgainst(contenders[i], contenders[j])
			fmt.Println(contenders[i], "vs", contenders[j])
			fmt.Println(contenders[i], "avg:", a)
			fmt.Println(contenders[j], "avg:", b)
			fmt.Println("Rounds:", rounds)
			fmt.Println()
		}
	}

}

//takes in two executable paths
//starts up two of these dockers, with port 10000 in the VM bound to two random numbers on the host
//copies in the two executables to /code in both images
//makes a connection to localhost:whatever, which is docker:10000 on both. this makes relay.go start the provided /code, and forwards stdin/stdout over this socket
//runs iterated prisoner's dilemma, and finally returns number of rounds and final socre

func playAgainst(pathToExecA string, pathToExecB string) (float32, float32, int) {

	imagIDRaw, err := exec.Command("docker", "build", "-q", "dock").CombinedOutput()
	if err != nil {
		fmt.Println(string(imagIDRaw))
		panic(err)
	}
	imagID := strings.Trim(string(imagIDRaw), "\n")
	//fmt.Println("Image id:",imagID)

	containerARaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", "5000:10000", imagID).CombinedOutput()
	if err != nil {
		fmt.Println(string(containerARaw))
		panic(err)
	}
	containerA := strings.Trim(string(containerARaw), "\n")
	//fmt.Println(containerA)
	err = exec.Command("docker", "cp", pathToExecA, containerA+":/code").Run()
	if err != nil {
		panic(err)
	}

	containerBRaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", "6000:10000", imagID).CombinedOutput()
	if err != nil {
		fmt.Println(string(containerBRaw))
		panic(err)
	}
	containerB := strings.Trim(string(containerBRaw), "\n")

	err = exec.Command("docker", "cp", pathToExecB, containerB+":/code").Run()
	if err != nil {
		panic(err)
	}

	// cd dock && docker build .
	//dockerImage := "9e983c4697fc"

	// docker run --rm -p 5000:10000 $(dockerImage)
	//containerA := ""
	// docker cp $(pathToExecA) $(containerA):/code

	// docker run --rm -p 6000:10000 $(dockerImage)
	//containerB := ""
	// docker cp $(pathToExecB) $(containerB):/code

	A, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}
	defer A.Close()
	B, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		panic(err)
	}
	defer B.Close()

	AScore := 0
	BScore := 0

	numTurns := 0

	for {
		A.SetDeadline(time.Now().Add(50 * time.Millisecond))
		B.SetDeadline(time.Now().Add(50 * time.Millisecond))
		Am := make([]byte, 1)
		_, err = io.ReadFull(A, Am)
		if err != nil {
			panic(err) // todo dont panic, just disqualify
		}
		Bm := make([]byte, 1)
		_, err = io.ReadFull(B, Bm)
		if err != nil {
			panic(err) // todo dont panic, just disqualify
		}

		if Am[0] != 0 && Am[0] != 1 {
			panic("A made invalid move") // todo dont panic, just disqualify
		}

		if Bm[0] != 0 && Bm[0] != 1 {
			panic("B made invalid move") // todo dont panic, just disqualify
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
	}

	return float32(AScore) / float32(numTurns), float32(BScore) / float32(numTurns), numTurns
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
