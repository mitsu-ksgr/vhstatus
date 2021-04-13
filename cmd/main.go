package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/mitsu-ksgr/vhstatus/internal/vhlogwatcher"
	"github.com/mitsu-ksgr/vhstatus/internal/vhstatus"
	"github.com/mitsu-ksgr/vhstatus/internal/web"
)

var vhs = vhstatus.New()

func getenv(key, default_value string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return default_value
	}
	return value
}

func log2store(event vhlogwatcher.VHLogEvent) {
	switch event.Event {
	//---------------------------------------------------------------------
	// Server Event
	case vhlogwatcher.GameServerConnected:
		vhs.SetStatus("Online")

	case vhlogwatcher.GameServerConnectedFailed, vhlogwatcher.GameServerDisconnected:
		vhs.SetStatus("Offile")

	case vhlogwatcher.ValheimVersion:
		vhs.SetValheimVersion(event.Value)

	case vhlogwatcher.ServerID:
		vhs.SetServerID(event.Value)

	case vhlogwatcher.LoadWorld:
		vhs.SetWorldName(event.Value)

	case vhlogwatcher.InitWorldGenSeed:
		vhs.SetWorldSeed(event.Value)

	case vhlogwatcher.DayHasPassed:
		vhs.SetDay(event.Value)

	//---------------------------------------------------------------------
	// User Event
	case vhlogwatcher.Connection,
		vhlogwatcher.GotHandshake,
		vhlogwatcher.Disconnection:
		vhs.UpdatePlayer(vhstatus.Player{
			SteamID:   event.SteamID,
			Status:    event.Event.String(),
			UpdatedAt: event.Timestamp,
		})
	case vhlogwatcher.GotCharacter:
		vhs.UpdatePlayer(vhstatus.Player{
			SteamID:   event.SteamID,
			Status:    event.Event.String(),
			Name:      event.Name,
			UpdatedAt: event.Timestamp,
		})
	}
}

func getPastLogFiles(dirpath string) []string {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	logs := make([]string, 0, 1)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "vhserver-console-") &&
			strings.HasSuffix(file.Name(), ".log") {
			logs = append(logs, dirpath+"/"+file.Name())
		}
	}
	sort.Strings(logs)

	return logs
}

func main() {
	var (
		port            string
		pathLogDir      string
		pathTemplateDir string
	)
	flag.StringVar(&port, "port", "8000", "http port")
	flag.StringVar(&pathLogDir, "log-dir-path", "/home/vhserver/log/console/", "path to direcotry of vhserver-console.log")
	flag.StringVar(&pathTemplateDir, "template-dir-path", "", "path to directory of html templates")
	flag.Parse()

	pathLogDir = strings.TrimSuffix(pathLogDir, "/")
	pathLogFile := pathLogDir + "/vhserver-console.log"

	// Setup data store
	for _, f := range getPastLogFiles(pathLogDir) {
		vhlogwatcher.ReadVHLog(f, log2store)
	}
	go vhlogwatcher.WatchVHLog(pathLogFile, log2store)

	// Setup web server
	web.SetFechVHStatusParamsFunc(func() vhstatus.Params {
		return vhs.Params()
	})
	if pathTemplateDir != "" {
		web.SetTemplateDirPath(pathTemplateDir)
		http.HandleFunc("/", web.Index)
	}
	http.HandleFunc("/api", web.ApiGetStatus)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
