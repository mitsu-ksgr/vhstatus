package vhlogwatcher

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

type EventType int

const (
	None EventType = iota

	// SERVER
	ValheimVersion
	ServerID
	InitWorldGenSeed
	LoadWorld
	GameServerConnected
	GameServerConnectedFailed
	GameServerDisconnected
	DayHasPassed

	// USER
	Connection
	GotHandshake
	GotCharacter
	Disconnection
)

func (et EventType) String() string {
	switch et {
	case None:
		return "None"

	// SERVER
	case ValheimVersion:
		return "Valheim version"
	case ServerID:
		return "Server ID"
	case InitWorldGenSeed:
		return "Initialize world generator seed"
	case LoadWorld:
		return "Load world"
	case GameServerConnected:
		return "Game server connected"
	case GameServerConnectedFailed:
		return "Game server connected failed"
	case GameServerDisconnected:
		return "Game server disconnected"
	case DayHasPassed:
		return "DayHasPassed"

	// USER
	case Connection:
		return "Connection"
	case GotHandshake:
		return "Got Handshake"
	case GotCharacter:
		return "Got Character"
	case Disconnection:
		return "Disconnection"
	default:
		return ""
	}
}

type VHLogEvent struct {
	Event     EventType
	Timestamp time.Time
	Value     string // generic

	// User event
	SteamID string
	Name    string
}

var reConsoleLog = regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2})`)
var reVHEvents = []struct {
	pattern *regexp.Regexp
	genfunc func([]string) VHLogEvent
}{

	//-------------------------------------------------------------------------
	// Server event
	//-------------------------------------------------------------------------
	{ // Valheim Version
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Valheim version:([0-9\.]+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     ValheimVersion,
				Timestamp: parseLogTime(matches[1]),
				Value:     matches[2], // valheim version
			}
		},
	},
	{ // Server ID
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Server ID (\d+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     ServerID,
				Timestamp: parseLogTime(matches[1]),
				Value:     matches[2], // server id
			}
		},
	},
	{ // Initialize world generator seed
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Initializing world generator seed:([a-zA-Z0-9]{10}) .*`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     InitWorldGenSeed,
				Timestamp: parseLogTime(matches[1]),
				Value:     matches[2], // seed value of world
			}
		},
	},
	{ // Load world
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Load world (.*)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     LoadWorld,
				Timestamp: parseLogTime(matches[1]),
				Value:     matches[2], // world name
			}
		},
	},
	{ // Game server connected
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Game server connected$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     GameServerConnected,
				Timestamp: parseLogTime(matches[1]),
			}
		},
	},
	{ // Game server connected failed
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Game server connected failed$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     GameServerConnectedFailed,
				Timestamp: parseLogTime(matches[1]),
			}
		},
	},
	{ // Game server disconnected
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Game server disconnected$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     GameServerDisconnected,
				Timestamp: parseLogTime(matches[1]),
			}
		},
	},
	{ // Day has passed.
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Time ([0-9\.]{16}), day:([0-9]*)\ *nextm:([0-9\.]{16})  skipspeed:([0-9.]{16})`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     DayHasPassed,
				Timestamp: parseLogTime(matches[1]),
				Value:     matches[3], // day.
			}
		},
	},

	//-------------------------------------------------------------------------
	// User event
	//-------------------------------------------------------------------------
	{ // User connection
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Got connection SteamID (\d+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     Connection,
				Timestamp: parseLogTime(matches[1]),
				SteamID:   matches[2],
			}
		},
	},
	{ // User handshake
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Got handshake from client (\d+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     GotHandshake,
				Timestamp: parseLogTime(matches[1]),
				SteamID:   matches[2],
			}
		},
	},
	{ // User disconnected
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Closing socket (\d+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     Disconnection,
				Timestamp: parseLogTime(matches[1]),
				SteamID:   matches[2],
			}
		},
	},
	{ // User character
		pattern: regexp.MustCompile(`^(\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}): Got character ZDOID from ([A-Za-z]\w*) : (-?\d+:\d+)$`),
		genfunc: func(matches []string) VHLogEvent {
			return VHLogEvent{
				Event:     GotCharacter,
				Timestamp: parseLogTime(matches[1]),
				Name:      matches[2],
			}
		},
	},
}

func parseLogTime(logtime string) time.Time {
	t, _ := time.Parse("01/02/2006 15:04:05", logtime)
	return t
}

func scanLogLine(row string) VHLogEvent {
	if !reConsoleLog.MatchString(row) {
		return VHLogEvent{Event: None}
	}

	for _, ev := range reVHEvents {
		if ret := ev.pattern.FindStringSubmatch(row); ret != nil {
			return ev.genfunc(ret)
		}
	}

	// Unsupported event or not interested event.
	return VHLogEvent{Event: None}
}

func ReadVHLog(logpath string, callback func(VHLogEvent)) {
	file, err := os.Open(logpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lastPlayer := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())
		if event := scanLogLine(row); event.Event != None {
			if event.Event == GotHandshake {
				lastPlayer = event.SteamID
			} else if event.Event == GotCharacter && lastPlayer != "" {
				event.SteamID = lastPlayer
				lastPlayer = ""
			}
			callback(event)
		}
	}
}

func WatchVHLog(logpath string, callback func(VHLogEvent)) {
	t, err := tail.TailFile(logpath, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	lastPlayer := ""
	for line := range t.Lines {
		row := strings.TrimSpace(line.Text)
		if event := scanLogLine(row); event.Event != None {
			if event.Event == GotHandshake {
				lastPlayer = event.SteamID
			} else if event.Event == GotCharacter && lastPlayer != "" {
				event.SteamID = lastPlayer
				lastPlayer = ""
			}
			callback(event)
		}
	}
}
