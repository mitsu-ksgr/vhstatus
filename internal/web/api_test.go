package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mitsu-ksgr/vhstatus/internal/vhstatus"
)

func Test_ApiGetStatus_Params(t *testing.T) {
	t.Cleanup(cleanup)

	wantCode := http.StatusOK
	wantBody := []string{
		"status", "server_id", "world_name", "world_seed", "updated_at",
		"active_player_count", "players", "day",
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/api",
		bytes.NewBufferString(""),
	)
	resp := httptest.NewRecorder()

	ApiGetStatus(resp, req)

	if resp.Code != wantCode {
		t.Errorf("ApiGetStatus response %d, want %d", resp.Code, wantCode)
	}

	strBody := resp.Body.String()
	for _, v := range wantBody {
		if !strings.Contains(strBody, v) {
			t.Errorf("ApiGetStatus response body did not contain '%s'\ngot: '%s'", v, strBody)
		}
	}
}

func Test_ApiGetStatus_NotSetVHStatus(t *testing.T) {
	t.Cleanup(cleanup)

	wantCode := http.StatusOK
	wantBody := "funcFetchVHStatus is nil"

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/api",
		bytes.NewBufferString(""),
	)
	resp := httptest.NewRecorder()

	ApiGetStatus(resp, req)

	if resp.Code != wantCode {
		t.Errorf("ApiGetStatus response %d, want %d", resp.Code, wantCode)
	}

	if got := resp.Body.String(); !strings.Contains(got, wantBody) {
		t.Errorf("ApiGetStatus response body did not contain '%s'\ngot: '%s'", wantBody, got)
	}
}

func Test_ApiGetStatus(t *testing.T) {
	t.Cleanup(cleanup)

	wantCode := http.StatusOK
	wantParams := vhstatus.Params{
		Status:            "Online",
		ServerID:          "1234567890",
		ValheimVersion:    "1.2.3",
		WorldName:         "test-world",
		WorldSeed:         "testseed",
		ActivePlayerCount: 3,
		Players: []vhstatus.Player{
			vhstatus.Player{"1", "Connection", "player1", time.Now()},
			vhstatus.Player{"2", "GotHandshake", "player2", time.Now()},
			vhstatus.Player{"3", "GotCharacter", "player3", time.Now()},
			vhstatus.Player{"4", "Disconnection", "player4", time.Now()},
		},
	}

	// Setup data store
	vhs := vhstatus.New()
	SetFechVHStatusParamsFunc(func() vhstatus.Params {
		return vhs.Params()
	})
	vhs.SetStatus(wantParams.Status)
	vhs.SetServerID(wantParams.ServerID)
	vhs.SetValheimVersion(wantParams.ValheimVersion)
	vhs.SetWorldName(wantParams.WorldName)
	vhs.SetWorldSeed(wantParams.WorldSeed)
	for _, p := range wantParams.Players {
		vhs.UpdatePlayer(p)
	}

	// Setup test request
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/api",
		bytes.NewBufferString(""),
	)
	resp := httptest.NewRecorder()

	ApiGetStatus(resp, req)

	if resp.Code != wantCode {
		t.Errorf("ApiGetStatus response %d, want %d", resp.Code, wantCode)
	}

	var respBody vhstatus.Params
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	if respBody.Status != wantParams.Status {
		t.Errorf("ApiGetStatus response body 'Status' is %q, want %q",
			respBody.Status, wantParams.Status)
	}
	if respBody.ServerID != wantParams.ServerID {
		t.Errorf("ApiGetStatus response body 'ServerID' is %q, want %q",
			respBody.ServerID, wantParams.ServerID)
	}
	if respBody.ValheimVersion != wantParams.ValheimVersion {
		t.Errorf("ApiGetStatus response body 'ValheimVersion' is %q, want %q",
			respBody.ValheimVersion, wantParams.ValheimVersion)
	}
	if respBody.WorldName != wantParams.WorldName {
		t.Errorf("ApiGetStatus response body 'WorldName' is %q, want %q",
			respBody.WorldName, wantParams.WorldName)
	}
	if respBody.WorldSeed != wantParams.WorldSeed {
		t.Errorf("ApiGetStatus response body 'WorldSeed' is %q, want %q",
			respBody.WorldSeed, wantParams.WorldSeed)
	}
	if respBody.ActivePlayerCount != wantParams.ActivePlayerCount {
		t.Errorf("ApiGetStatus response body 'ActivePlayerCount' is %d, want %d",
			respBody.ActivePlayerCount, wantParams.ActivePlayerCount)
	}
}
