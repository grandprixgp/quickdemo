package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	dem "github.com/markus-wa/demoinfocs-golang"
	events "github.com/markus-wa/demoinfocs-golang/events"
	"gopkg.in/djherbis/times.v1"
)

type demo_info struct {
	File   string         `json:"file"`
	Map    string         `json:"map"`
	Winner string         `json:"winner"`
	Loser  string         `json:"loser"`
	Score  map[string]int `json:"score"`
	Start  time.Time      `json:"start_time"`
	End    time.Time      `json:"end_time"`
}

func parse_demo(filename string, info map[string]demo_info) {
	file_stats, err := times.Stat(filename)
	if err != nil {
		panic(err)
	}
	demo_file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer demo_file.Close()

	var result demo_info

	result.File = filename
	result.Score = make(map[string]int)

	parser := dem.NewParser(demo_file)
	header, err := parser.ParseHeader()

	if file_stats.HasBirthTime() {
		result.Start = file_stats.BirthTime()
	}

	result.Map = header.MapName
	result.End = result.Start.Add(header.PlaybackTime)

	parser.RegisterEventHandler(func(e events.ScoreUpdated) {
		if e.NewScore != 0 {
			if e.TeamState.ID == 3 {
				result.Score["Counter-Terrorists"] = e.NewScore
			} else {
				result.Score["Terrorists"] = e.NewScore
			}
		}
	})

	parser.ParseToEnd()

	result.Winner = func() string {
		if result.Score["Counter-Terrorists"] > result.Score["Terrorists"] {
			return "Counter-Terrorists"
		} else {
			return "Terrorists"
		}
	}()
	result.Loser = func() string {
		if result.Score["Counter-Terrorists"] < result.Score["Terrorists"] {
			return "Counter-Terrorists"
		} else {
			return "Terrorists"
		}
	}()

	info[filename] = result
}

func parse_args() []string {
	var first_demo string
	flag.StringVar(&first_demo, "d", "", "A space seperated list of demos")
	flag.Parse()
	demos := flag.Args()
	demos = append([]string{first_demo}, demos...)
	return demos
}

func main() {

	var results = make(map[string]demo_info)

	var wg sync.WaitGroup
	demos := parse_args()
	for demo := range demos {
		demoname := demos[demo]
		wg.Add(1)
		go func() { parse_demo(demoname, results); wg.Done() }()
	}
	wg.Wait()

	results_json, _ := json.MarshalIndent(results, "", "\t")
	fmt.Println(string(results_json))
}
