package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestWorlds(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/worlds.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldsJson, err := TibiaWorldsOverviewImpl(string(data), "https://www.tibia.com/community/?subtopic=worlds")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := worldsJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=worlds", information.TibiaURLs[0])

	assert.Equal(8720, worldsJson.Worlds.PlayersOnline)
	assert.Equal(64028, worldsJson.Worlds.RecordPlayers)
	assert.Equal("2007-11-28T18:26:00Z", worldsJson.Worlds.RecordDate)
	assert.Equal(76, len(worldsJson.Worlds.RegularWorlds))
	assert.Equal(6, len(worldsJson.Worlds.TournamentWorlds))

	adra := worldsJson.Worlds.RegularWorlds[0]
	assert.Equal("Adra", adra.Name)
	assert.Equal("online", adra.Status)
	assert.Equal(18, adra.PlayersOnline)
	assert.Equal("Europe", adra.Location)
	assert.Equal("Open PvP", adra.PvpType)
	assert.False(adra.PremiumOnly)
	assert.Equal("blocked", adra.TransferType)
	assert.True(adra.BattleyeProtected)
	assert.Equal("release", adra.BattleyeDate)
	assert.Equal("regular", adra.GameWorldType)
	assert.Empty(adra.TournamentWorldType)

	astera := worldsJson.Worlds.RegularWorlds[4]
	assert.Equal("Astera", astera.Name)
	assert.Equal("online", astera.Status)
	assert.Equal(222, astera.PlayersOnline)
	assert.Equal("North America", astera.Location)
	assert.Equal("Optional PvP", astera.PvpType)
	assert.False(astera.PremiumOnly)
	assert.Equal("regular", astera.TransferType)
	assert.True(astera.BattleyeProtected)
	assert.Equal("2017-09-12", astera.BattleyeDate)
	assert.Equal("regular", astera.GameWorldType)
	assert.Empty(astera.TournamentWorldType)

	premia := worldsJson.Worlds.RegularWorlds[50]
	assert.Equal("Premia", premia.Name)
	assert.Equal("offline", premia.Status)
	assert.Equal(0, premia.PlayersOnline)
	assert.Equal("Europe", premia.Location)
	assert.Equal("Open PvP", premia.PvpType)
	assert.True(premia.PremiumOnly)
	assert.Equal("regular", premia.TransferType)
	assert.True(premia.BattleyeProtected)
	assert.Equal("2017-09-05", premia.BattleyeDate)
	assert.Equal("regular", premia.GameWorldType)
	assert.Empty(premia.TournamentWorldType)

	zuna := worldsJson.Worlds.RegularWorlds[74]
	assert.Equal("Zuna", zuna.Name)
	assert.Equal("online", zuna.Status)
	assert.Equal(5, zuna.PlayersOnline)
	assert.Equal("Europe", zuna.Location)
	assert.Equal("Hardcore PvP", zuna.PvpType)
	assert.False(zuna.PremiumOnly)
	assert.Equal("locked", zuna.TransferType)
	assert.False(zuna.BattleyeProtected)
	assert.Empty("", zuna.BattleyeDate)
	assert.Equal("experimental", zuna.GameWorldType)
	assert.Empty(zuna.TournamentWorldType)

	endera := worldsJson.Worlds.TournamentWorlds[1]
	assert.Equal("Endera", endera.Name)
	assert.Equal("unknown", endera.Status)
	assert.Equal(0, endera.PlayersOnline)
	assert.Equal("North America", endera.Location)
	assert.Equal("Optional PvP", endera.PvpType)
	assert.True(endera.PremiumOnly)
	assert.Equal("blocked", endera.TransferType)
	assert.True(endera.BattleyeProtected)
	assert.Equal("release", endera.BattleyeDate)
	assert.Equal("tournament", endera.GameWorldType)
	assert.Equal("restricted", endera.TournamentWorldType)
}
