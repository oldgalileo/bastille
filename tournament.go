package main

import (
	"io"
	"math/rand"
	"net"
	"os/exec"
	"strings"
	"time"

	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	dockerImageID string
	trnLog        = log.WithFields(log.Fields{
		"prefix": "tournament",
	})
)

var exampleStrategies = []Strategy{
	Strategy{
		Name:   "allcoop",
		Author: "system",
		Path:   STRATEGIES_DIR + "allcoop.py",
	},
	Strategy{
		Name:   "alldefect",
		Author: "system",
		Path:   STRATEGIES_DIR + "alldefect.py",
	},
	Strategy{
		Name:   "invalidcommand",
		Author: "system",
		Path:   STRATEGIES_DIR + "invalidcommand.py",
	},
	Strategy{
		Name:   "titfortat",
		Author: "system",
		Path:   STRATEGIES_DIR + "titfortat.py",
	},
}

var PayoutMatrix = map[bool]map[bool]int{
	true: {
		true:  3,
		false: 0,
	},
	false: {
		true:  5,
		false: 1,
	},
}

func init() {
	imageIDRaw, err := exec.Command("docker", "build", "-q", "dock").CombinedOutput()
	if err != nil {
		log.WithError(err).Panic("Couldn't build docker image with image id: ", string(imageIDRaw))
	}
	dockerImageID = strings.Trim(string(imageIDRaw), "\n")
}

type TournamentManager struct {
	Strategies []Strategy
}

type Strategy struct {
	Name   string
	Author string
	Path   string
}

func (tm *TournamentManager) playAgainst(aStrat, bStrat Strategy) (AavgScore float32, BavgScore float32, numRounds int, aDisqualified bool, bDisqualified bool) {
	localLog := trnLog.WithFields(log.Fields{"pa": aStrat.Name, "pb": bStrat.Name})
	localLog.Info("Starting round...")
	containerA, portA := createContainer()
	err := exec.Command("docker", "cp", aStrat.Path, containerA+":/code").Run()
	if err != nil {
		panic(err)
	}

	containerB, portB := createContainer()

	err = exec.Command("docker", "cp", aStrat.Path, containerB+":/code").Run()
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
			localLog.Warn("Disqualified due to stream closed")
			return 0, 0, 0, true, false
		}

		_, err = io.ReadFull(B, Bm)
		if err != nil {
			localLog.Warn("Disqualified due to stream closed")
			return 0, 0, 0, false, true
		}

		if Am[0] != 0 && Am[0] != 1 {
			localLog.Warn("A made invalid move")
			return 0, 0, 0, true, false
		}

		if Bm[0] != 0 && Bm[0] != 1 {
			localLog.Warn("B made invalid move")
			return 0, 0, 0, false, true
		}

		Amove := Am[0] == 1
		Bmove := Bm[0] == 1

		AScore += PayoutMatrix[Amove][Bmove]
		BScore += PayoutMatrix[Bmove][Amove]

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

func createContainer() (string, int) {
	for i := 0; i < 100; i++ {
		port := 5000 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(4000)
		containerRaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", strconv.Itoa(port)+":10000", dockerImageID).CombinedOutput()
		if err != nil {
			if i < 99 {
				continue
			}
			fmt.Println(string(containerRaw))
			trnLog.WithField("container", string(containerRaw)).Debug("Created container")
			panic(err)
		}
		container := strings.Trim(string(containerRaw), "\n")
		return container, port

	}
	trnLog.Panic("Should never reach end of function (createContainer)")
	return "", 0
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
