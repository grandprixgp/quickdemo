package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"quickdemo/creation_time"
	"quickdemo/memory_available"
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
	Size    int       `json:"size"`          // size of file
	Created time.Time `json:"creation_date"` // creation time of file
	Demo    demoInfo  `json:"demo"`          // parsed demo object
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
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

	result.Start = creation_time.Get_creation_time(filename)

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

func main() {

	fmt.Println(memory_available.Get_memory_available())
	return

	var results = make(map[string]demoInfo)

	var wg sync.WaitGroup
	demos := parseArgs()
	for demo := range demos {
		demoname := demos[demo]
		wg.Add(1)
		go func() { parseDemo(demoname, results); wg.Done() }()
	}
	wg.Wait()

	resultsJSON, _ := json.MarshalIndent(results, "", "\t")
	fmt.Println(string(resultsJSON))
}
