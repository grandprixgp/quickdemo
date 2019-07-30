package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"quickdemo/file_stats"
	"quickdemo/memory_stats"
	"sync"
	"syscall"
	"time"

	dem "github.com/markus-wa/demoinfocs-golang"
	events "github.com/markus-wa/demoinfocs-golang/events"
)

type playerInfo struct {
	Name   string `json:"name"`
	Kills  int    `json:"kills"`
	Deaths int    `json:"deaths"`
}

type demoInfo struct {
	File    string               `json:"file"`
	Valid   bool                 `json:"valid"`
	Map     string               `json:"map"`
	Winner  string               `json:"winner"`
	Loser   string               `json:"loser"`
	Score   map[string]int       `json:"score"`
	Start   time.Time            `json:"start_time"`
	End     time.Time            `json:"end_time"`
	State   string               `json:"state"`
	Players map[int64]playerInfo `json:"players"`
}

type demoFile struct {
	Name    string    `json:"filename"`      // name of file
	Size    uint64    `json:"size"`          // size of file
	Created time.Time `json:"creation_date"` // creation time of file
	Demo    demoInfo  `json:"demo"`          // parsed demo object
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func parseFile(filename string, demos *[]demoFile) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var demo demoFile
	demo.Name = filename
	demo.Size = memory_stats.GetMemoryAvailable()
	demo.Created = file_stats.GetCreationTime(filename)

	*demos = append(*demos, demo)
}

func parseDemo(filename string, info map[string]demoInfo) {
	demoFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer demoFile.Close()

	var result demoInfo
	result.Valid = false

	result.File = filename
	result.Score = make(map[string]int)
	result.Players = make(map[int64]playerInfo)

	parser := dem.NewParser(demoFile)
	header, err := parser.ParseHeader()

	result.Start = file_stats.GetCreationTime(filename)

	result.Map = header.MapName
	result.End = result.Start.Add(header.PlaybackTime)

	if result.Start == result.End {
		result.State = "live"
	} else {
		result.State = "finished"
	}

	parser.RegisterEventHandler(func(e events.ScoreUpdated) {
		if e.NewScore != 0 {
			if e.TeamState.ID == 3 {
				result.Score["Counter-Terrorists"] = e.NewScore
			} else {
				result.Score["Terrorists"] = e.NewScore
			}
		}

		participants := parser.GameState().Participants().Playing()
		if len(participants) >= 1 {
			for player := range participants {
				player := participants[player]
				if player.SteamID != 0 {
					playerinfo := playerInfo{
						Name:   player.Name,
						Kills:  player.AdditionalPlayerInformation.Kills,
						Deaths: player.AdditionalPlayerInformation.Deaths,
					}
					result.Players[player.SteamID] = playerinfo
				}
			}
		}
	})

	parser.RegisterEventHandler(func(e events.BombPickup) {
		result.Valid = true
	})

	parser.ParseToEnd()

	result.Winner = func() string {
		if result.Score["Counter-Terrorists"] > result.Score["Terrorists"] {
			return "Counter-Terrorists"
		}
		return "Terrorists"
	}()
	result.Loser = func() string {
		if result.Score["Counter-Terrorists"] < result.Score["Terrorists"] {
			return "Counter-Terrorists"
		}
		return "Terrorists"
	}()

	info[filename] = result
}

func parseArgs() []string {
	var firstDemo string
	flag.StringVar(&firstDemo, "d", "", "A space seperated list of demos")
	flag.Parse()
	demos := flag.Args()
	demos = append([]string{firstDemo}, demos...)
	return demos
}

func use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func main() {

	var availableMemory = memory_stats.GetMemoryAvailable()
	var demos = make([]demoFile, len(flag.Args()))
	var waitGroup sync.WaitGroup
	var filenames = parseArgs()

	for _, filename := range filenames {
		waitGroup.Add(1)
		go func() { parseFile(filename, &demos); waitGroup.Done() }()
	}
	waitGroup.Wait()

	demosJSON, _ := json.MarshalIndent(demos, "", "\t")
	fmt.Println(string(demosJSON))

	//var results = make(map[string]demoInfo)
	//for demo := range demos {
	//	demoname := demos[demo]
	//	wg.Add(1)
	//	go func() { parseDemo(demoname, results); wg.Done() }()
	//}
	//wg.Wait()

	//resultsJSON, _ := json.MarshalIndent(results, "", "\t")

	use(availableMemory)
	//fmt.Println(string(resultsJSON))
}
