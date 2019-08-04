package main

import (
	"C"
	//"encoding/json"
	"flag"
	//"fmt"
	"os"
	//"sync"
	"syscall"
	"time"

	//proto "github.com/gogo/protobuf/proto"
	file_stats "github.com/grandprixgp/file_stats"
	//memory_stats "github.com/grandprixgp/memory_stats"
	dem "github.com/markus-wa/demoinfocs-golang"
	events "github.com/markus-wa/demoinfocs-golang/events"
	msg "github.com/markus-wa/demoinfocs-golang/msg"
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
	Demo    *demoInfo `json:"demo"`          // parsed demo object
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func parseFile(filename string, demoFiles *[]demoFile, totalSize *uint64) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var demoFile demoFile
	demoFile.Name = filename
	demoFile.Size = (file_stats.GetSize(filename) / 1024) / 1024
	demoFile.Created = file_stats.GetCreationTime(filename)
	demoFile.Demo = &demoInfo{}

	*demoFiles = append(*demoFiles, demoFile)
	*totalSize = *totalSize + (demoFile.Size)
}

func parseDemo(filename string, demo *demoInfo, cfg dem.ParserConfig) {
	demoFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer demoFile.Close()

	demo.Valid = false

	demo.File = filename
	demo.Score = make(map[string]int)
	demo.Players = make(map[int64]playerInfo)

	parser := dem.NewParserWithConfig(demoFile, cfg)
	header, err := parser.ParseHeader()

	demo.Start = file_stats.GetCreationTime(demo.File)

	demo.Map = header.MapName
	demo.End = demo.Start.Add(header.PlaybackTime)

	parser.RegisterNetMessageHandler(func(m *msg.CNETMsg_SetConVar) {
		for _, cvar := range m.Convars.Cvars {
			if cvar.Name == "sv_matchstarted" {
				if cvar.Value == "1" {
					demo.State = "live"
					demo.Valid = true
				}
			}
			if cvar.Name == "sv_matchfinished" {
				if cvar.Value == "1" {
					demo.State = "finished"
					demo.Valid = true
				}
			}
		}
	})

	parser.RegisterEventHandler(func(e events.ScoreUpdated) {
		if e.NewScore != 0 {
			if e.TeamState.ID == 3 {
				demo.Score["Counter-Terrorists"] = e.NewScore
			} else {
				demo.Score["Terrorists"] = e.NewScore
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
					demo.Players[player.SteamID] = playerinfo
				}
			}
		}
	})

	parser.ParseToEnd()

	demo.Winner = func() string {
		if demo.Score["Counter-Terrorists"] > demo.Score["Terrorists"] {
			return "Counter-Terrorists"
		}
		return "Terrorists"
	}()
	demo.Loser = func() string {
		if demo.Score["Counter-Terrorists"] < demo.Score["Terrorists"] {
			return "Counter-Terrorists"
		}
		return "Terrorists"
	}()
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

//export Dump
func Dump(a string) string {
	return a
}

func main() {

	//cfg := dem.DefaultParserConfig
	//cfg.AdditionalNetMessageCreators = map[int]dem.NetMessageCreator{
	//	6: func() proto.Message {
	//		return new(msg.CNETMsg_SetConVar)
	//	},
	//}
	//
	//var availableMemory = memory_stats.GetMemoryAvailable()
	//var demosSlice = make([]demoFile, len(flag.Args()))
	//var demosMap = make(map[string]demoFile)
	//var totalSize uint64
	//var waitGroup sync.WaitGroup
	//var filenames = parseArgs()
	//
	//for _, filename := range filenames {
	//	filename := filename
	//	waitGroup.Add(1)
	//	go func() { parseFile(filename, &demosSlice, &totalSize); waitGroup.Done() }()
	//}
	//waitGroup.Wait()
	//
	///* chunk demos to ensure we don't OOM ourselves */
	//
	//var demoChunks [][]demoFile
	//
	//if (totalSize / availableMemory) >= 1 {
	//	chunkCount := totalSize / availableMemory
	//	chunkSize := len(demosSlice) / int(chunkCount)
	//
	//	for i := 0; i < len(demosSlice); i += chunkSize {
	//		chunkEnd := i + chunkSize
	//
	//		if chunkEnd > len(demosSlice) {
	//			chunkEnd = len(demosSlice)
	//		}
	//
	//		demoChunks = append(demoChunks, demosSlice[i:chunkEnd])
	//	}
	//} else {
	//	demoChunks = append(demoChunks, demosSlice)
	//}
	//
	//fmt.Printf("Available Memory: %dMB\nTotal Demos Size: %dMB\nChunks: %d\n", availableMemory, totalSize, len(demoChunks))
	//
	//demosSlice = nil
	//
	//for _, demoChunk := range demoChunks {
	//	for _, demo := range demoChunk {
	//		demo := demo
	//		waitGroup.Add(1)
	//		go func() {
	//			parseDemo(demo.Name, demo.Demo, cfg)
	//			demosMap[demo.Name] = demo
	//			waitGroup.Done()
	//		}()
	//	}
	//	waitGroup.Wait()
	//}
	//
	//demoChunks = nil
	//
	//demosJSON, _ := json.MarshalIndent(demosMap, "", "\t")
	//fmt.Println(string(demosJSON))
}
