package runner

import (
	"io"
	"math/rand"
	"net"
	"time"
)

//takes in two executable paths
//starts up two of these dockers, with port 10000 in the VM bound to two random numbers on the host
//copies in the two executables to /code in both images
//makes a connection to localhost:whatever, which is docker:10000 on both. this makes relay.go start the provided /code, and forwards stdin/stdout over this socket
//runs iterated prisoner's dilemma, and finally returns number of rounds and final socre

func playAgainst(pathToExecA string, pathToExecB string) (float32, float32, int) {

	// cd dock && docker build .
	dockerImage := "9e983c4697fc"

	// docker run --rm -p 5000:10000 $(dockerImage)
	containerA := ""
	// docker cp $(pathToExecA) $(containerA):/code

	// docker run --rm -p 6000:10000 $(dockerImage)
	containerB := ""
	// docker cp $(pathToExecB) $(containerB):/code

	A, err := net.Dial("tcp", "localhost:5000")
	B, err := net.Dial("tcp", "localhost:6000")

	AScore := 0
	BScore := 0

	numTurns := 0

	for {
		Am := make([]byte, 1)
		io.ReadFull(A, Am)
		Bm := make([]byte, 1)
		io.ReadFull(B, Bm)

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

func value(my bool, their bool) float32 {
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
