package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mitsu-ksgr/vhstatus/internal/vhstatus"
)

func Test_HtmlIndex(t *testing.T) {
	t.Cleanup(cleanup)

	wantCode := http.StatusOK
	wantParams := vhstatus.Params{
		Status:            "Online",
		ServerID:          "1234567890",
		ValheimVersion:    "1.2.3",
		WorldName:         "test-world",
		WorldSeed:         "testseed",
		Day:               "123",
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
	SetTemplateDirPath("/go/src/github.com/mitsu-ksgr/vhstatus/web")
	SetFechVHStatusParamsFunc(func() vhstatus.Params {
		return vhs.Params()
	})
	vhs.SetStatus(wantParams.Status)
	vhs.SetServerID(wantParams.ServerID)
	vhs.SetValheimVersion(wantParams.ValheimVersion)
	vhs.SetWorldName(wantParams.WorldName)
	vhs.SetWorldSeed(wantParams.WorldSeed)
	vhs.SetDay(wantParams.Day)
	for _, p := range wantParams.Players {
		vhs.UpdatePlayer(p)
	}

	// Setup test request
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		bytes.NewBufferString(""),
	)
	resp := httptest.NewRecorder()

	Index(resp, req)
	if resp.Code != wantCode {
		t.Errorf("html#Index response %d, want %d", resp.Code, wantCode)
	}

	strBody := resp.Body.String()
	for _, want := range []string{
		wantParams.Status, wantParams.ValheimVersion,
	} {
		if !strings.Contains(strBody, want) {
			t.Errorf("html#Index response did not contain %q", want)
		}
	}
}

func Test_HtmlIndex_NotSetTemplateDirPath(t *testing.T) {
	t.Cleanup(cleanup)
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		bytes.NewBufferString(""),
	)
	resp := httptest.NewRecorder()

	Index(resp, req)

	wantCode := 500
	if resp.Code != wantCode {
		t.Errorf("html#Index(without index.html) response %d, want %d", resp.Code, wantCode)
	}

	wantBody := "500 Internal Server Error"
	strBody := resp.Body.String()
	if !strings.Contains(strBody, wantBody) {
		t.Errorf("html#Index response did not contain %q\n\tgot: %q", wantBody, strBody)
	}
}
