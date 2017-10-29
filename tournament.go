package main

import (
	"io"
	"math/rand"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"strconv"

	"encoding/hex"

	"io/ioutil"

	"encoding/json"

	"encoding/xml"

	"os"

	"math"

	log "github.com/sirupsen/logrus"
)

const STRATEGIES_DIR = "strategies/"
const EXAMPLES_DIR = "examples/"
const TOURNAMENT_DIR = "tournament/"

var (
	dockerImageID string
	trnLog        = log.WithFields(log.Fields{
		"prefix": "tournament",
	})
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var exampleStrategies = []*Strategy{
	&Strategy{
		ID:      "1",
		Name:    "allcoop",
		Author:  "system",
		Path:    EXAMPLES_DIR + "allcoop.py",
		Matches: []MatchID{},
	},
	&Strategy{
		ID:      "2",
		Name:    "alldefect",
		Author:  "system",
		Path:    EXAMPLES_DIR + "alldefect.py",
		Matches: []MatchID{},
	},
	&Strategy{
		ID:      "3",
		Name:    "invalidcommand",
		Author:  "system",
		Path:    EXAMPLES_DIR + "invalidcommand.py",
		Matches: []MatchID{},
	},
	&Strategy{
		ID:      "4",
		Name:    "titfortat",
		Author:  "system",
		Path:    EXAMPLES_DIR + "titfortat.py",
		Matches: []MatchID{},
	},
	&Strategy{
		ID:      "5",
		Name:    "random",
		Author:  "system",
		Path:    EXAMPLES_DIR + "random.py",
		Matches: []MatchID{},
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
	sync.RWMutex `json:"-",xml:"-"`
	Leaderboard  map[StrategyID]float32   `json:"leaderboard",xml:"leaderboard"`
	Strategies   map[StrategyID]*Strategy `json:"strategies",xml:"strategies"`
	Matches      map[MatchID]*Match       `json:"matches",xml:"matches"`
	pairings     map[[2]*Strategy]int
	exits        []chan bool
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
	sync.RWMutex `json:"-",xml:"-"`
	ID           StrategyID `json:"id",xml:"id"`
	Name         string     `json:"name",xml:"name"`
	Author       string     `json:"author",xml:"author"`
	Description  string     `json:"desc",xml:"desc"`
	Path         string     `json:"path",xml:"path"`
	Disqualified bool       `json:"disqualified",xml:"disqualified"`
	Matches      []MatchID  `json:"matches",xml:"matches"`
	containers   chan container
}

type container struct {
	id   string
	port string
}

func (s *Strategy) bufferContainers(exit chan bool) {
	for {
		select {
		case <-exit:
			return
		default:
			break
		}
		//log.WithFields(log.Fields{
		//	"prefix": "strategy",
		//	"name":   s.Name,
		//	"id":     s.ID,
		//}).Debug("Creating strategy...")
		tempContainer := createContainer()
		err := exec.Command("docker", "cp", s.Path, tempContainer.id+":/code").Run()
		if err != nil {
			log.WithError(err).WithField("path", s.Path).Error("Could not copy strategy #" + string(s.ID) + " to container #" + tempContainer.id)
			continue
		}
		s.containers <- tempContainer
	}
}

func (tm *TournamentManager) add(strategy *Strategy) {
	tm.Lock()
	defer tm.Unlock()
	tm.Strategies[strategy.ID] = strategy
	for _, oldStrat := range tm.Strategies {
		tm.pairings[[2]*Strategy{strategy, oldStrat}] = 0
		tm.pairings[[2]*Strategy{oldStrat, strategy}] = 0
	}
}

func (tm *TournamentManager) run() {
	trnLog.Info("Starting tournament...")
	tm.Lock()
	exit := make(chan bool, 1)
	tm.exits = append(tm.exits, exit)
	tm.Unlock()
	for {
		select {
		case <-exit:
			return
		default:
			break
		}
		lowestKey := [2]*Strategy{exampleStrategies[0], exampleStrategies[1]}
		lowestVal := math.MaxInt32
		firstVal := math.MaxInt32
		tm.RLock()
		for pairing, played := range tm.pairings {
			if lowestVal == math.MaxInt32 {
				firstVal = played
			}
			if lowestVal >= played {
				lowestKey = pairing
				lowestVal = played
			}
		}
		if firstVal == lowestVal && lowestVal == 100 {
			trnLog.Info("No new strategies... Skipping")
			time.Sleep(3 * time.Second)
			continue
		}
		tm.RUnlock()
		match := tm.playAgainst(lowestKey[0], lowestKey[1])
		if match.DisqualifiedA && !tm.validateStrategy(match.PlayerA) {
			tm.disqualifyStrategy(match.PlayerA)
			continue
		}
		if match.DisqualifiedB && !tm.validateStrategy(match.PlayerB) {
			tm.disqualifyStrategy(match.PlayerB)
			continue
		}
		// heck yeah
		tm.Lock()
		tm.Leaderboard[match.PlayerA] = float32((tm.Leaderboard[match.PlayerA]*float32(len(lowestKey[0].Matches)-1) + float32(match.ScoreA)) / float32(len(lowestKey[0].Matches)))
		tm.Leaderboard[match.PlayerB] = float32((tm.Leaderboard[match.PlayerB]*float32(len(lowestKey[1].Matches)-1) + float32(match.ScoreB)) / float32(len(lowestKey[1].Matches)))
		tm.pairings[lowestKey] += 1
		tm.Unlock()
	}
	defer trnLog.Info("Ending tournament...")
}

func (tm *TournamentManager) init() {
	tm.load()
	tm.buildPairs()
	for _, strategy := range tm.Strategies {
		strategy.containers = make(chan container, 5)
		tm.Lock()
		exit := make(chan bool, 1)
		tm.exits = append(tm.exits, exit)
		tm.Unlock()
		go strategy.bufferContainers(exit)
	}
	go tm.periodic()
}

func (tm *TournamentManager) buildPairs() {
	tm.Lock()
	for _, aStrat := range tm.Strategies {
		if aStrat.Disqualified {
			continue
		}
		for _, bStrat := range tm.Strategies {
			if bStrat.Disqualified {
				continue
			}
			if aStrat.ID == bStrat.ID {
				continue
			}
			tm.pairings[[2]*Strategy{aStrat, bStrat}] = 0
		}
	}
	tm.Unlock()
}

func (tm *TournamentManager) cleanup() {
	for _, exit := range tm.exits {
		exit <- true
	}
	tm.save()
}

func (tm *TournamentManager) load() {
	trnLog.Info("Loading core data...")
	defer trnLog.Info("Successfully loaded!")
	tm.pairings = make(map[[2]*Strategy]int)
	if _, err := os.Stat(TOURNAMENT_DIR + "core.json"); os.IsNotExist(err) {
		trnLog.Info("Could not find core data. Initializing...")
		tm.Leaderboard = make(map[StrategyID]float32)
		tm.Strategies = make(map[StrategyID]*Strategy)
		tm.Matches = make(map[MatchID]*Match)
		for _, strategy := range exampleStrategies {
			tm.Strategies[strategy.ID] = strategy
		}
		return
	}
	raw, rawErr := ioutil.ReadFile(TOURNAMENT_DIR + "core.json")
	if rawErr != nil {
		trnLog.WithError(rawErr).Panic("Could not read core data")
	}
	json.Unmarshal(raw, tm)
}

func (tm *TournamentManager) periodic() {
	for range time.Tick(1 * time.Minute) {
		tm.save()
	}
}

func (tm *TournamentManager) save() {
	trn.RLock()
	trnLog.Info("Saving core data...")
	raw, rawErr := xml.MarshalIndent(tm, "", "    ")
	if rawErr != nil {
		var jsonErr error
		raw, jsonErr = json.Marshal(tm)
		if jsonErr != nil {
			trnLog.WithError(jsonErr).Error("Could not marshal data in both XML and JSON")
			return
		}
		trnLog.WithError(rawErr).Error("Could not marshal data in XML")
	}
	writeErr := ioutil.WriteFile(TOURNAMENT_DIR+"core.json", raw, 0644)
	if writeErr != nil {
		trnLog.WithError(writeErr).Error("Could not write core data")
	}
	trn.RUnlock()
}

func (tm *TournamentManager) validateStrategy(id StrategyID) bool {
	strategy := tm.Strategies[id]
	trnLog.WithFields(log.Fields{
		"name":    strategy.Name,
		"matches": len(strategy.Matches),
		"id":      id,
	}).Debug("Validating")
	if len(strategy.Matches) < 5 {
		return true
	}
	disqualifyCount := 0
	for ; disqualifyCount < 5; disqualifyCount++ {
		tempMatchID := strategy.Matches[len(strategy.Matches)-1-disqualifyCount]
		tempMatch := tm.Matches[tempMatchID]
		trnLog.WithFields(log.Fields{
			"rounds":    tempMatch.Rounds,
			"timestamp": tempMatch.Timestamp,
			"dq-a":      tempMatch.DisqualifiedA,
			"dq-b":      tempMatch.DisqualifiedB,
			"id-a":      tempMatch.PlayerA,
			"id-b":      tempMatch.PlayerB,
			"player":    tempMatch.PlayerA == id,
			"count":     disqualifyCount,
		}).Debug("Checking match #" + tempMatchID)
		if tempMatch.PlayerA == id && !tempMatch.DisqualifiedA {
			return true
		}
		if tempMatch.PlayerB == id && !tempMatch.DisqualifiedB {
			return true
		}
	}
	return false
}

func (tm *TournamentManager) disqualifyStrategy(id StrategyID) {
	trnLog.Warn("Disqualifying Strategy #" + id)
	strategy := tm.Strategies[id]
	for _, matchID := range strategy.Matches {
		trnLog.Warn("Undoing match #" + matchID)
		match := tm.Matches[matchID]
		var enemyID StrategyID
		var enemyScore float32
		if match.PlayerA != id {
			enemyID = match.PlayerA
			enemyScore = match.ScoreA
		} else {
			enemyID = match.PlayerB
			enemyScore = match.ScoreB
		}
		enemy := tm.Strategies[enemyID]
		tm.Lock()
		tm.Leaderboard[enemyID] = ((tm.Leaderboard[enemyID] * float32(len(enemy.Matches))) - enemyScore) / float32(len(enemy.Matches)-1)
		tm.Unlock()
	}
	for key := range tm.pairings {
		if key[0] == strategy || key[1] == strategy {
			delete(tm.pairings, key)
		}
	}
	strategy.Disqualified = true
	trnLog.Warn("Successfully disqualified Strategy #" + id)
}

func (tm *TournamentManager) playAgainst(aStrat, bStrat *Strategy) *Match {
	tm.Lock()
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
	tm.Unlock()
	aStrat.Matches = append(aStrat.Matches, match.ID)
	bStrat.Matches = append(bStrat.Matches, match.ID)
	localLog := trnLog.WithFields(log.Fields{"match": match.ID})
	localLog.Info("Starting round...")

	containerA := <-aStrat.containers
	containerB := <-bStrat.containers

	//time.Sleep(50 * time.Millisecond) // go run relay.go may take time to start

	// cd dock && docker build .
	//dockerImage := "9e983c4697fc"

	// docker run --rm -p 5000:10000 $(dockerImage)
	//containerA := ""
	// docker cp $(pathToExecA) $(containerA):/code

	// docker run --rm -p 6000:10000 $(dockerImage)
	//containerB := ""
	// docker cp $(pathToExecB) $(containerB):/code

	A, err := net.Dial("tcp", net.JoinHostPort("localhost", containerA.port))
	if err != nil {
		localLog.WithError(err).Panic("Could not establish connection to PA container")
	}
	defer A.Close()
	B, err := net.Dial("tcp", net.JoinHostPort("localhost", containerB.port))
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
			localLog.WithFields(log.Fields{"a-name": aStrat.Name, "a-id": match.PlayerA}).Warn("A Disqualified due to stream closed")
			match.DisqualifiedA = true
			return match
		}

		_, err = io.ReadFull(B, Bm)
		if err != nil {
			localLog.WithFields(log.Fields{"b-name": bStrat.Name, "b-id": match.PlayerB}).Warn("B Disqualified due to stream closed")
			match.DisqualifiedB = true
			return match
		}

		if Am[0] != 0 && Am[0] != 1 {
			localLog.WithFields(log.Fields{"a-name": aStrat.Name, "a-id": match.PlayerA}).Warn("A made invalid move")
			match.DisqualifiedA = true
			localLog.WithFields(log.Fields{
				"id":    match.PlayerA,
				"name":  aStrat.Name,
				"score": match.ScoreA,
				"dq":    match.DisqualifiedA,
			}).Debug("Player A information")
			localLog.WithFields(log.Fields{
				"id":    match.PlayerB,
				"name":  bStrat.Name,
				"score": match.ScoreB,
				"dq":    match.DisqualifiedB,
			}).Debug("Player B information")
			return match
		}

		if Bm[0] != 0 && Bm[0] != 1 {
			localLog.WithFields(log.Fields{"b-name": bStrat.Name, "b-id": match.PlayerB}).Warn("B made invalid move")
			match.DisqualifiedB = true
			localLog.WithFields(log.Fields{
				"id":    match.PlayerA,
				"name":  aStrat.Name,
				"score": match.ScoreA,
				"dq":    match.DisqualifiedA,
			}).Debug("Player A information")
			localLog.WithFields(log.Fields{
				"id":    match.PlayerB,
				"name":  bStrat.Name,
				"score": match.ScoreB,
				"dq":    match.DisqualifiedB,
			}).Debug("Player B information")
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
	localLog.WithFields(log.Fields{
		"rounds":    match.Rounds,
		"timestamp": match.Timestamp,
	}).Debug("Round information")
	localLog.WithFields(log.Fields{
		"id":    match.PlayerA,
		"name":  aStrat.Name,
		"score": match.ScoreA,
		"dq":    match.DisqualifiedA,
	}).Debug("Player A information")
	localLog.WithFields(log.Fields{
		"id":    match.PlayerB,
		"name":  bStrat.Name,
		"score": match.ScoreB,
		"dq":    match.DisqualifiedB,
	}).Debug("Player B information")
	localLog.Info("Finished round...")
	return match
}

func createContainer() container {
	for i := 0; i < 100; i++ {
		port := 5000 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(4000)
		containerRaw, err := exec.Command("docker", "run", "--rm", "-d", "-p", strconv.Itoa(port)+":10000", dockerImageID).CombinedOutput()
		if err != nil {
			if i < 99 {
				continue
			}
			trnLog.WithField("container", string(containerRaw)).Panic("Issue with container")
		}
		containerID := strings.Trim(string(containerRaw), "\n")
		return container{containerID, strconv.Itoa(port)}

	}
	trnLog.Panic("Should never reach end of function (createContainer)")
	return container{"", ""}
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
