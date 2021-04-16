package vhstatus

import (
	"strings"
	"testing"
	"time"
)

func new_vhs_instance() *VHStatus {
	vhs := New()
	vhs.updatedAt = time.Now().Add(time.Minute * -1) // 1 min ago.
	return vhs
}

func Test_New(t *testing.T) {
	vhs := New()

	if vhs == nil {
		t.Error("vhstatus.New() returned nil")
	}

	if vhs.status != "init" {
		t.Error("vhstatus.New() did not initialize vhs.status")
	}

	if vhs.activePlayerCount != 0 {
		t.Error("vhstatus.New() did not initialize vhs.activePlayerCount")
	}

	if vhs.players == nil {
		t.Error("vhstatus.New() did not initialize the player list")
	}
}

//-----------------------------------------------------------------------------
// Setter
//-----------------------------------------------------------------------------

func Test_SetStatus(t *testing.T) {
	cases := []string{"Online", "Offline", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetStatus(c)

		if vhs.status != c {
			t.Errorf("VHStatus#SetStaus did not updated. got %q, want: %q", vhs.status, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetStatus did not update vhs.updatedAt")
		}
	}
}

func Test_SetServerID(t *testing.T) {
	cases := []string{"12345678901234567", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetServerID(c)

		if vhs.serverID != c {
			t.Errorf("VHStatus#SetServerID did not updated. got %q, want: %q", vhs.serverID, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetServerID did not update vhs.updatedAt")
		}
	}
}

func Test_SetValheimVersion(t *testing.T) {
	cases := []string{"0.0.1", "1.2.3", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetValheimVersion(c)

		if vhs.valheimVersion != c {
			t.Errorf("VHStatus#SetValheimVersion did not updated. got %q, want: %q", vhs.valheimVersion, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetValheimVersion did not update vhs.updatedAt")
		}
	}
}

func Test_SetWorldName(t *testing.T) {
	cases := []string{"valheim-world", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetWorldName(c)

		if vhs.worldName != c {
			t.Errorf("VHStatus#SetWorldName did not updated. got %q, want: %q", vhs.worldName, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetWorldName did not update vhs.updatedAt")
		}
	}
}

func Test_SetWorldSeed(t *testing.T) {
	cases := []string{"abcdefg", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetWorldSeed(c)

		if vhs.worldSeed != c {
			t.Errorf("VHStatus#SetWorldSeed did not updated. got %q, want: %q", vhs.worldSeed, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetWorldSeed did not update vhs.updatedAt")
		}
	}
}

func Test_SetDay(t *testing.T) {
	cases := []string{"1", "12", ""}

	vhs := new_vhs_instance()
	for _, c := range cases {
		ua := vhs.updatedAt
		vhs.SetDay(c)

		if vhs.day != c {
			t.Errorf("VHStatus#SetDay did not updated. got %q, want: %q", vhs.day, c)
		}
		if ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#SetWorldSeed did not update vhs.updatedAt")
		}
	}
}

func Test_UpdatePlayer(t *testing.T) {
	cases := []struct {
		player                Player
		want                  bool
		wantErr               string
		wantActivePlayerCount int
	}{
		{Player{SteamID: "1", Status: "Connection", Name: "player1"}, true, "", 1},
		{Player{SteamID: "2", Status: "GotHandshake", Name: "player2"}, true, "", 2},
		{Player{SteamID: "3", Status: "GotCharacter", Name: "player3"}, true, "", 3},
		{Player{SteamID: "4", Status: "Disconnection", Name: "player4"}, true, "", 3},
		{Player{SteamID: "1", Status: "Disconnection", Name: "player1"}, false, "", 2},
		{Player{}, false, "player.SteamID is not set", 2},
	}

	vhs := new_vhs_instance()
	for i, c := range cases {
		ua := vhs.updatedAt
		ret, err := vhs.UpdatePlayer(c.player)

		// if not error is thrown, check updatedAt updates
		if c.wantErr == "" && ua.Equal(vhs.updatedAt) {
			t.Errorf("VHStatus#UpdatePlayer(case[%d]) did not update vhs.updatedAt", i)
		}

		if ret != c.want {
			t.Errorf("VHStatus#UpdatePlayer(case[%d]) = %t, want %t", i, ret, c.want)
		}
		if c.wantErr != "" && !strings.Contains(err.Error(), c.wantErr) {
			t.Errorf("VHStatus#UpdatePlayer(case[%d])\n"+
				"\treturned err: %q\n"+
				"\twant contain: %q",
				i, err.Error(), c.wantErr)
		}
		if vhs.activePlayerCount != c.wantActivePlayerCount {
			t.Errorf("VHStatus#UpdatePlayer(case[%d]) did not update ActivePlayerCount. got %d, want %d",
				i, vhs.activePlayerCount, c.wantActivePlayerCount)
		}
	}
}

//-----------------------------------------------------------------------------
// Getter
//-----------------------------------------------------------------------------

func Test_Player(t *testing.T) {
	srcTime := "2021-04-10T12:34:56Z"
	srcTimeParam, _ := time.Parse(time.RFC3339, srcTime)

	p := Player{"1", "Connection", "player1", srcTimeParam}
	if got := p.UpdatedAtAsString(); got != srcTime {
		t.Errorf("Player#UpdatedAtAsString returns %q, want %q", got, srcTime)
	}
}

func Test_Params(t *testing.T) {
	vhs := New()
	vhs.SetStatus("testing")

	params := vhs.Params()
	if params.Status != vhs.status {
		t.Errorf("VHStatus#Params().status = %q, want = %q", params.Status, vhs.status)
	}
	if params.UpdatedAtAsString() != params.UpdatedAt.Format(time.RFC3339) {
		t.Errorf("VHStatus#Params().UpdatedAtAsString = %q, want = %q",
			params.UpdatedAtAsString(),
			params.UpdatedAt.Format(time.RFC3339))
	}
}

func Test_Params_Players(t *testing.T) {
	vhs := New()
	vhs.UpdatePlayer(Player{SteamID: "1"})
	vhs.UpdatePlayer(Player{SteamID: "2"})

	params := vhs.Params()
	if &params.Players == &vhs.players {
		t.Errorf("VHStatus#Params().Players is same instance as vhs.players\n"+
			"\tparams.Players : %p\n"+
			"\tvhs.players    : %p",
			&params.Players, &vhs.players)
	}

	if params.ActivePlayerCount != vhs.activePlayerCount {
		t.Errorf("VHStatus#Params().ActivePlayerCount = %q, want = %q",
			params.ActivePlayerCount, vhs.activePlayerCount)
	}
	if len(params.Players) != len(vhs.players) {
		t.Errorf("VHStatus#Params().players len = %q, want = %q",
			len(params.Players), len(vhs.players))
	} else {
		for i, _ := range params.Players {
			if params.Players[i].SteamID != vhs.players[i].SteamID {
				t.Errorf("VHStatus#Params().players[%d].SteamID = %q, want = %q",
					i, params.Players[i].SteamID, vhs.players[i].SteamID)
			}

			if &params.Players[i] == &vhs.players[i] {
				t.Errorf("VHStatus#Params().Players[%d] is same instance as vhs.players[%d]\n"+
					"\tparams.Players[%d] : %p\n"+
					"\t   vhs.players[%d] : %p",
					i, i, i, &params.Players[i], i, &vhs.players[i])
			}
		}
	}
}
