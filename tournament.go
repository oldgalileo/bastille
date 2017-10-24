package main

import (
	"io"
	"math/rand"
	"net"
	"os/exec"
	"strings"
	"time"

	"strconv"

	"encoding/hex"

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
		ID:     "1",
		Name:   "allcoop",
		Author: "system",
		Path:   STRATEGIES_DIR + "allcoop.py",
	},
	Strategy{
		ID:     "2",
		Name:   "alldefect",
		Author: "system",
		Path:   STRATEGIES_DIR + "alldefect.py",
	},
	Strategy{
		ID:     "3",
		Name:   "invalidcommand",
		Author: "system",
		Path:   STRATEGIES_DIR + "invalidcommand.py",
	},
	Strategy{
		ID:     "4",
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
	Leaderboard map[StrategyID]float32   `xml:"leaderboard"`
	Strategies  map[StrategyID]*Strategy `xml:"strategies"`
	Matches     map[MatchID]*Match       `xml:"matches"`
}

type MatchID string

type Match struct {
	ID            MatchID    `json:"id",xml:"id",csv:"id"`
	PlayerA       StrategyID `json:"player-a",xml:"player-a",csv:"player-a"`
	PlayerB       StrategyID `json:"player-b",xml:"player-b",csv:"player-b"`
	Rounds        int        `json:"rounds",xml:"rounds",csv:"rounds"`
	ScoreA        float32    `json:"score-a",xml:"score-a",csv:"score-a"`
	ScoreB        float32    `json:"score-b",xml:"score-b",csv:"score-b"`
	DisqualifiedA bool       `json:"dq-a",xml:"dq-a",csv:"dq-a"`
	DisqualifiedB bool       `json:"dq-b",xml:"dq-b",csv:"dq-b"`
	Timestamp     int64      `json:"timestamp",xml:"timestamp",csv:"timestamp"`
}

type StrategyID string

type Strategy struct {
	ID      StrategyID `json:"id"`
	Name    string     `json:"name",xml:"name"`
	Author  string     `json:"author",xml:"author"`
	Path    string     `json:"-",xml:"path"`
	Matches []MatchID  `json:"matches"`
}

func (tm *TournamentManager) playAgainst(aStrat, bStrat Strategy) *Match {
	match := &Match{
		ID:            getMatchID(),
		PlayerA:       aStrat.ID,
		PlayerB:       bStrat.ID,
		Rounds:        0,
		ScoreA:        0,
		ScoreB:        0,
		DisqualifiedA: false,
		DisqualifiedB: false,
		Timestamp:     time.Now().UnixNano(),
	}
	tm.Matches[match.ID] = match
	aStrat.Matches = append(aStrat.Matches, match.ID)
	bStrat.Matches = append(bStrat.Matches, match.ID)
	localLog := trnLog.WithFields(log.Fields{"match": match.ID})
	localLog.Info("Starting round...")

	localLog.Debug("Creating PA container...")
	containerA, portA := createContainer()
	err := exec.Command("docker", "cp", aStrat.Path, containerA+":/code").Run()
	if err != nil {
		match.DisqualifiedA = true
		match.DisqualifiedB = true
		localLog.WithError(err).Panic("Failed to create PA container")
	}

	localLog.Debug("Creating PB container...")
	containerB, portB := createContainer()
	err = exec.Command("docker", "cp", aStrat.Path, containerB+":/code").Run()
	if err != nil {
		match.DisqualifiedA = true
		match.DisqualifiedB = true
		localLog.WithError(err).Panic("Failed to create PB container")
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
		localLog.WithError(err).Panic("Could not establish connection to PA container")
	}
	defer A.Close()
	B, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portB))
	if err != nil {
		localLog.WithError(err).Panic("Could not establish connection to PB container")
	}
	defer B.Close()

	AScore := 0
	BScore := 0

	A.SetDeadline(time.Now().Add(2 * time.Second)) // generous deadline for first move
	B.SetDeadline(time.Now().Add(2 * time.Second))
	Am := make([]byte, 1)
	Bm := make([]byte, 1)
	for {

		_, err = io.ReadFull(A, Am)
		if err != nil {
			localLog.Warn("Disqualified due to stream closed")
			match.DisqualifiedA = true
			return match
		}

		_, err = io.ReadFull(B, Bm)
		if err != nil {
			localLog.Warn("Disqualified due to stream closed")
			match.DisqualifiedB = true
			return match
		}

		if Am[0] != 0 && Am[0] != 1 {
			localLog.Warn("A made invalid move")
			match.DisqualifiedA = true
			return match
		}

		if Bm[0] != 0 && Bm[0] != 1 {
			localLog.Warn("B made invalid move")
			match.DisqualifiedB = true
			return match
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

		match.Rounds++
		if rnd.Float32() < 0.0003 && match.Rounds > 100 {
			break
		}
		A.SetDeadline(time.Now().Add(50 * time.Millisecond)) // short deadline for subsequent moves
		B.SetDeadline(time.Now().Add(50 * time.Millisecond))
	}
	match.ScoreA = float32(AScore) / float32(match.Rounds)
	match.ScoreB = float32(BScore) / float32(match.Rounds)
	return match
}

func createContainer() (string, int) {
	for i := 0; i < 100; i++ {
		port := 5000 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(4000)
		containerRaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", strconv.Itoa(port)+":10000", dockerImageID).CombinedOutput()
		if err != nil {
			if i < 99 {
				continue
			}
			trnLog.WithField("container", string(containerRaw)).Panic("Issue with container")
		}
		container := strings.Trim(string(containerRaw), "\n")
		return container, port

	}
	trnLog.Panic("Should never reach end of function (createContainer)")
	return "", 0
}

func getMatchID() MatchID {
	id := make([]byte, 8)
	rand.Read(id)
	return MatchID(hex.EncodeToString(id))
}

func getStrategyID() StrategyID {
	id := make([]byte, 4)
	rand.Read(id)
	return StrategyID(hex.EncodeToString(id))
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
