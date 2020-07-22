package main

import (
	"errors"
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

const (
	betray byte = iota
	coop
)

var (
	dockerImageID     string
	dockerHistoryLock sync.Mutex
	dockerHistory     = make(map[string]container)

	trnLog = log.WithFields(log.Fields{
		"prefix": "tournament",
	})
	rnd                  = rand.New(rand.NewSource(time.Now().UnixNano()))
	errFirstStrategyLoad = errors.New("must initialize pairings")
)

var exampleStrategies = []*Strategy{
	{
		ID:      "1",
		Name:    "one-third",
		Author:  "system",
		Path:    STRATEGIES_DIR + "RPS_Strat_1.ipd",
		Matches: []MatchID{},
	},
	{
		ID:      "2",
		Name:    "one-third-two",
		Author:  "system",
		Path:    STRATEGIES_DIR + "RPS_Strat_1.ipd",
		Matches: []MatchID{},
	},
}

var PayoutMatrix = map[byte]map[byte]int{
	coop: {
		coop:   3,
		betray: 0,
	},
	betray: {
		coop:   5,
		betray: 1,
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
	Leaderboard  map[StrategyID]float32
	Strategies   map[StrategyID]*Strategy
	Matches      map[MatchID]*Match
	Pairings     map[[2]StrategyID]int
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
	toggle       chan bool
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
		if s.Disqualified {
			select {
			case ctnr := <-s.containers:
				killContainer(ctnr)
				continue
			default:
				exit <- true
			}
		}
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
	for oldStrat, _ := range tm.Strategies {
		tm.Pairings[[2]StrategyID{strategy.ID, oldStrat}] = 0
		tm.Pairings[[2]StrategyID{oldStrat, strategy.ID}] = 0
	}
	strategy.containers = make(chan container, 1)
	exit := make(chan bool, 1)
	tm.exits = append(tm.exits, exit)
	go strategy.bufferContainers(exit)
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
			trnLog.Info("Ending tournament...")
			return
		default:
			break
		}
		lowestKey := [2]StrategyID{exampleStrategies[0].ID, exampleStrategies[1].ID}
		lowestVal := math.MaxInt32
		firstVal := math.MaxInt32
		tm.RLock()
		for pairing, played := range tm.Pairings {
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
			tm.RUnlock()
			time.Sleep(3 * time.Second)
			continue
		}
		tm.RUnlock()
		match := tm.playAgainst(tm.Strategies[lowestKey[0]], tm.Strategies[lowestKey[1]])
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
		tm.Leaderboard[match.PlayerA] = float32((tm.Leaderboard[match.PlayerA]*float32(len(tm.Strategies[lowestKey[0]].Matches)-1) + float32(match.ScoreA)) / float32(len(tm.Strategies[lowestKey[0]].Matches)))
		tm.Leaderboard[match.PlayerB] = float32((tm.Leaderboard[match.PlayerB]*float32(len(tm.Strategies[lowestKey[1]].Matches)-1) + float32(match.ScoreB)) / float32(len(tm.Strategies[lowestKey[1]].Matches)))
		tm.Pairings[lowestKey] += 1
		tm.Unlock()
	}
}

func (tm *TournamentManager) init() {
	tm.load()
	for _, strategy := range tm.Strategies {
		strategy.containers = make(chan container, 1)
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
	defer tm.Unlock()
	for aStratID, aStrat := range tm.Strategies {
		if aStrat.Disqualified {
			continue
		}
		for bStratID, bStrat := range tm.Strategies {
			if bStrat.Disqualified {
				continue
			}
			if aStrat.ID == bStrat.ID {
				continue
			}
			tm.Pairings[[2]StrategyID{aStratID, bStratID}] = 0
		}
	}
}

func (tm *TournamentManager) cleanup() {
	for _, exit := range tm.exits {
		exit <- true
	}
	tm.save()
	//for _, strat := range tm.Strategies {
	//	for {
	//		select {
	//		case cntr := <-strat.containers:
	//			killContainer(cntr)
	//			continue
	//		default:
	//			break
	//		}
	//		return
	//	}
	//}
	dockerHistoryLock.Lock()
	for _, cntr := range dockerHistory {
		killContainer(cntr)
	}
	dockerHistoryLock.Unlock()
}

func (tm *TournamentManager) load() {
	trnLog.Info("Loading tournament data...")
	defer trnLog.Info("Successfully loaded!")

	var (
		leadErr     error
		stratErr    error
		matchErr    error
		pairingsErr error
	)

	tm.Leaderboard, leadErr = loadLeaderboard()
	tm.Strategies, stratErr = loadStrategies()
	tm.Matches, matchErr = loadMatches()
	tm.Pairings, pairingsErr = loadPairings()
	if pairingsErr != nil && pairingsErr == errFirstStrategyLoad {
		tm.buildPairs()
		pairingsErr = nil
	}

	if leadErr == nil && stratErr == nil && matchErr == nil && pairingsErr == nil {
		trnLog.Info("Loaded without errors!")
	} else {
		trnLog.WithFields(log.Fields{
			"lead":  leadErr,
			"strat": stratErr,
			"match": matchErr,
			"pair":  pairingsErr,
		}).Info("Errors occurred while loading...")
	}

	tm.exits = []chan bool{}
}

func loadLeaderboard() (map[StrategyID]float32, error) {
	trnLog.Info("Loading leaderboard...")
	defer trnLog.Info("Finished loaded leaderboard!")
	if _, err := os.Stat(TOURNAMENT_DIR + "leaderboard.json"); os.IsNotExist(err) {
		return make(map[StrategyID]float32), nil
	} else {
		leaderboard := make(map[StrategyID]float32)
		raw, rawErr := ioutil.ReadFile(TOURNAMENT_DIR + "leaderboard.json")
		if rawErr != nil {
			return leaderboard, rawErr
		}
		json.Unmarshal(raw, &leaderboard)
		return leaderboard, nil
	}
}

func loadStrategies() (map[StrategyID]*Strategy, error) {
	trnLog.Info("Loading strategies...")
	defer trnLog.Info("Finished loaded strategies!")
	if _, err := os.Stat(STRATEGIES_DIR + "core.json"); os.IsNotExist(err) {
		strategies := make(map[StrategyID]*Strategy)
		for _, strategy := range exampleStrategies {
			strategies[strategy.ID] = strategy
		}
		return strategies, nil
	} else {
		strategies := make(map[StrategyID]*Strategy)
		raw, rawErr := ioutil.ReadFile(STRATEGIES_DIR + "core.json")
		if rawErr != nil {
			return strategies, rawErr
		}
		json.Unmarshal(raw, &strategies)
		return strategies, nil
	}
}

func loadMatches() (map[MatchID]*Match, error) {
	trnLog.Info("Loading matches...")
	defer trnLog.Info("Finished loaded matches!")
	if _, err := os.Stat(TOURNAMENT_DIR + "matches.json"); os.IsNotExist(err) {
		return make(map[MatchID]*Match), nil
	} else {
		matches := make(map[MatchID]*Match)
		raw, rawErr := ioutil.ReadFile(TOURNAMENT_DIR + "matches.json")
		if rawErr != nil {
			return matches, rawErr
		}
		json.Unmarshal(raw, &matches)
		return matches, nil
	}
}

func loadPairings() (map[[2]StrategyID]int, error) {
	trnLog.Info("Loading pairings...")
	defer trnLog.Info("Finished loaded pairings!")
	if _, err := os.Stat(TOURNAMENT_DIR + "pairings.json"); os.IsNotExist(err) {
		return make(map[[2]StrategyID]int), errFirstStrategyLoad
	} else {
		tempPairings := make(map[string]int)
		raw, rawErr := ioutil.ReadFile(TOURNAMENT_DIR + "pairings.json")
		if rawErr != nil {
			return nil, rawErr
		}
		json.Unmarshal(raw, &tempPairings)
		pairings := make(map[[2]StrategyID]int)
		for key, value := range tempPairings {
			split := strings.Split(key, "-")
			pairings[[2]StrategyID{StrategyID(split[0]), StrategyID(split[1])}] = value
		}
		return pairings, nil
	}
}

func (tm *TournamentManager) periodic() {
	for range time.Tick(1 * time.Minute) {
		tm.save()
	}
}

func (tm *TournamentManager) save() {
	trn.RLock()
	defer trn.RUnlock()
	trnLog.Info("Saving tournament data...")
	tm.saveLeaderboard()
	tm.saveStrategies()
	tm.saveMatches()
	tm.savePairings()
	trnLog.Info("Finished saving tournament data...")
}

func (tm *TournamentManager) saveLeaderboard() {
	raw, rawErr := xml.MarshalIndent(tm.Leaderboard, "", "    ")
	if rawErr != nil {
		var jsonErr error
		raw, jsonErr = json.Marshal(tm.Leaderboard)
		if jsonErr != nil {
			trnLog.WithError(jsonErr).Error("Could not marshal data in both XML and JSON")
			return
		}
		trnLog.WithError(rawErr).Warning("Could not marshal data in XML... (this is safe to ignore)")
	}
	writeErr := ioutil.WriteFile(TOURNAMENT_DIR+"leaderboard.json", raw, 0644)
	if writeErr != nil {
		trnLog.WithError(writeErr).Error("Could not write core data")
	}
}

func (tm *TournamentManager) saveStrategies() {
	raw, rawErr := xml.MarshalIndent(tm.Strategies, "", "    ")
	if rawErr != nil {
		var jsonErr error
		raw, jsonErr = json.Marshal(tm.Strategies)
		if jsonErr != nil {
			trnLog.WithError(jsonErr).Error("Could not marshal data in both XML and JSON")
			return
		}
		trnLog.WithError(rawErr).Warning("Could not marshal data in XML... (this is safe to ignore)")
	}
	writeErr := ioutil.WriteFile(STRATEGIES_DIR+"core.json", raw, 0644)
	if writeErr != nil {
		trnLog.WithError(writeErr).Error("Could not write core data")
	}
}

func (tm *TournamentManager) saveMatches() {
	raw, rawErr := xml.MarshalIndent(tm.Matches, "", "    ")
	if rawErr != nil {
		var jsonErr error
		raw, jsonErr = json.Marshal(tm.Matches)
		if jsonErr != nil {
			trnLog.WithError(jsonErr).Error("Could not marshal data in both XML and JSON")
			return
		}
		trnLog.WithError(rawErr).Warning("Could not marshal data in XML... (this is safe to ignore)")
	}
	writeErr := ioutil.WriteFile(TOURNAMENT_DIR+"matches.json", raw, 0644)
	if writeErr != nil {
		trnLog.WithError(writeErr).Error("Could not write core data")
	}
}

func (tm *TournamentManager) savePairings() {
	tempPairings := make(map[string]int)
	for key, value := range tm.Pairings {
		tempPairings[string(key[0])+"-"+string(key[1])] = value
	}
	raw, rawErr := xml.MarshalIndent(tempPairings, "", "    ")
	if rawErr != nil {
		var jsonErr error
		raw, jsonErr = json.Marshal(tempPairings)
		if jsonErr != nil {
			trnLog.WithError(jsonErr).Error("Could not marshal data in both XML and JSON")
			return
		}
		trnLog.WithError(rawErr).Warning("Could not marshal data in XML... (this is safe to ignore)")
	}
	writeErr := ioutil.WriteFile(TOURNAMENT_DIR+"pairings.json", raw, 0644)
	if writeErr != nil {
		trnLog.WithError(writeErr).Error("Could not write core data")
	}
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
	for key := range tm.Pairings {
		if key[0] == id || key[1] == id {
			delete(tm.Pairings, key)
		}
	}
	strategy.Disqualified = true
	trnLog.Warn("Successfully disqualified Strategy #" + id)
}

func (tm *TournamentManager) playAgainst(aStrat, bStrat *Strategy) *Match {
	tm.Lock()
	defer tm.Unlock()
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

	containerA := <-aStrat.containers
	containerB := <-bStrat.containers
	defer killContainer(containerA)
	defer killContainer(containerB)

	time.Sleep(75 * time.Millisecond) // go run relay.go may take time to start

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
			localLog.WithError(err).WithFields(log.Fields{"a-name": aStrat.Name, "a-id": match.PlayerA}).Warn("A Disqualified due to stream closed")
			match.DisqualifiedA = true
			return match
		}

		_, err = io.ReadFull(B, Bm)
		if err != nil {
			localLog.WithError(err).WithFields(log.Fields{"b-name": bStrat.Name, "b-id": match.PlayerB}).Warn("B Disqualified due to stream closed")
			match.DisqualifiedB = true
			return match
		}

		if _, ok := PayoutMatrix[Am[0]]; !ok {
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

		if _, ok := PayoutMatrix[Bm[0]]; !ok {
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

		AScore += PayoutMatrix[Am[0]][Bm[0]]
		BScore += PayoutMatrix[Bm[0]][Am[0]]

		B.Write([]byte{Am[0]})
		A.Write([]byte{Bm[0]})

		match.Rounds++
		if rnd.Float32() < 0.0003 && match.Rounds > 100 {
			break
		}
		A.SetDeadline(time.Now().Add(50 * time.Millisecond)) // short deadline for subsequent moves
		B.SetDeadline(time.Now().Add(50 * time.Millisecond))
	}
	exec.Command("docker", "cp", containerA.id+":/stuff/history.txt", "history/"+string(match.ID)+"-"+string(aStrat.ID)).Run()
	exec.Command("docker", "cp", containerB.id+":/stuff/history.txt", "history/"+string(match.ID)+"-"+string(bStrat.ID)).Run()
	A.Write([]byte{})
	B.Write([]byte{})
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
		trnLog.Debug("Created container " + string(containerRaw)[:8])
		containerID := strings.Trim(string(containerRaw), "\n")
		container := container{containerID, strconv.Itoa(port)}
		dockerHistoryLock.Lock()
		dockerHistory[container.id] = container
		dockerHistoryLock.Unlock()
		return container

	}
	trnLog.Panic("Should never reach end of function (createContainer)")
	return container{"", ""}
}

func killContainer(cntr container) {
	err := exec.Command("docker", "kill", cntr.id).Run()
	if err != nil {
		trnLog.WithError(err).WithField("container", cntr.id).Warning("Could not kill container...")
	}
	dockerHistoryLock.Lock()
	delete(dockerHistory, cntr.id)
	dockerHistoryLock.Unlock()
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
