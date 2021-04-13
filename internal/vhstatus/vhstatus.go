package vhstatus

import (
	"errors"
	"sync"
	"time"
)

type Player struct {
	SteamID   string    `json:"steam_id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"` // last update time of player on the log file.
}

func (p Player) UpdatedAtAsString() string {
	return p.UpdatedAt.Format(time.RFC3339)
}

func (p *Player) update(rhs *Player) {
	if rhs.Status != "" {
		p.Status = rhs.Status
	}
	if rhs.Name != "" {
		p.Name = rhs.Name
	}
	if !rhs.UpdatedAt.IsZero() {
		p.UpdatedAt = rhs.UpdatedAt
	}
}

type Params struct {
	Status            string    `json:"status"`
	UpdatedAt         time.Time `json:"updated_at"`
	ServerID          string    `json:"server_id"`
	ValheimVersion    string    `json:"valheim_version"`
	WorldName         string    `json:"world_name"`
	WorldSeed         string    `json:"world_seed"`
	Day               string    `json:"day"`
	ActivePlayerCount int       `json:"active_player_count"`
	Players           []Player  `json:"players"`
}

func (p Params) UpdatedAtAsString() string {
	return p.UpdatedAt.Format(time.RFC3339)
}

type VHStatus struct {
	// Server Status
	status    string
	updatedAt time.Time // last update time of VHStatus instance.

	// Server Info
	serverID       string
	valheimVersion string

	// World Info
	worldName string
	worldSeed string
	day       string

	// Activity
	players           []Player
	activePlayerCount int

	// internal
	mu sync.Mutex
}

func New() *VHStatus {
	return &VHStatus{
		status:            "init",
		updatedAt:         time.Now(),
		activePlayerCount: 0,
		players:           make([]Player, 0, 10),
	}
}

func (vhs *VHStatus) Params() Params {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	players := make([]Player, len(vhs.players))
	copy(players, vhs.players)

	return Params{
		Status:            vhs.status,
		UpdatedAt:         vhs.updatedAt,
		ServerID:          vhs.serverID,
		ValheimVersion:    vhs.valheimVersion,
		WorldName:         vhs.worldName,
		WorldSeed:         vhs.worldSeed,
		Day:               vhs.day,
		ActivePlayerCount: vhs.activePlayerCount,
		Players:           players,
	}
}

func (vhs *VHStatus) SetStatus(status string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.status = status
	vhs.updatedAt = time.Now()
}

func (vhs *VHStatus) SetServerID(sid string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.serverID = sid
	vhs.updatedAt = time.Now()
}

func (vhs *VHStatus) SetValheimVersion(version string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.valheimVersion = version
	vhs.updatedAt = time.Now()
}

func (vhs *VHStatus) SetWorldName(name string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.worldName = name
	vhs.updatedAt = time.Now()
}

func (vhs *VHStatus) SetWorldSeed(seed string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.worldSeed = seed
	vhs.updatedAt = time.Now()
}

func (vhs *VHStatus) SetDay(day string) {
	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	vhs.day = day
	vhs.updatedAt = time.Now()
}

// UpdatePlayer updates the player information.
// It returns true if the player is a new registration,
// otherwise it returns false.
//
// If player.SteamID is zero-value, UpdatePlayer returns an error.
func (vhs *VHStatus) UpdatePlayer(player Player) (bool, error) {
	if player.SteamID == "" {
		return false, errors.New("player.SteamID is not set.")
	}

	vhs.mu.Lock()
	defer vhs.mu.Unlock()

	new_register := true

	for i, _ := range vhs.players {
		// already registered
		if vhs.players[i].SteamID == player.SteamID {
			vhs.players[i].update(&player)
			new_register = false
			break
		}
	}

	if new_register {
		vhs.players = append(vhs.players, player)
	}

	// count active player
	vhs.activePlayerCount = 0
	for i, _ := range vhs.players {
		if vhs.players[i].Status != "Disconnection" {
			vhs.activePlayerCount += 1
		}
	}

	vhs.updatedAt = time.Now()
	return new_register, nil
}
