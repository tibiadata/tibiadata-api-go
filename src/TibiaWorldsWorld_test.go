package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestWorldAntica(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Antica.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Antica", string(data), "https://www.tibia.com/community/?subtopic=worlds&world=Antica")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World
	information := worldJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=worlds&world=Antica", information.TibiaURLs[0])

	umlautPlayer := world.OnlinePlayers[171]
	assert.Equal("Nöber", umlautPlayer.Name)
	assert.Equal(455, umlautPlayer.Level)
	assert.Equal("Elder Druid", umlautPlayer.Vocation)
}

func TestWorldEndebra(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Endebra.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Endebra", string(data), "https://www.tibia.com/community/?subtopic=worlds&world=Endebra")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World
	information := worldJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=worlds&world=Endebra", information.TibiaURLs[0])

	assert.Equal("Endebra", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(0, world.PlayersOnline)
	assert.Equal(89, world.RecordPlayers)
	assert.Equal("2020-04-23T01:30:30Z", world.RecordDate)
	assert.Equal("2019-05", world.CreationDate)
	assert.Equal("South America", world.Location)
	assert.Equal("Optional PvP", world.PvpType)
	assert.True(world.PremiumOnly)
	assert.Equal("blocked", world.TransferType)
	assert.Equal(0, len(world.WorldsQuestTitles))
	assert.True(world.BattleyeProtected)
	assert.Equal("release", world.BattleyeDate)
	assert.Equal("tournament", world.GameWorldType)
	assert.Equal("restricted", world.TournamentWorldType)
	assert.Equal(0, len(world.OnlinePlayers))
}

func TestWorldPremia(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Premia.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Premia", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World

	assert.Equal("Premia", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(53, world.PlayersOnline)
	assert.Equal(531, world.RecordPlayers)
	assert.Equal("2013-08-08T15:30:30Z", world.RecordDate)
	assert.Equal("2002-04", world.CreationDate)
	assert.Equal("Europe", world.Location)
	assert.Equal("Open PvP", world.PvpType)
	assert.True(world.PremiumOnly)
	assert.Equal("regular", world.TransferType)
	assert.Equal(7, len(world.WorldsQuestTitles))
	assert.Equal("The Lightbearer", world.WorldsQuestTitles[0])
	assert.Equal("Rise of Devovorga", world.WorldsQuestTitles[1])
	assert.Equal("Annual Autumn Vintage", world.WorldsQuestTitles[2])
	assert.Equal("A Piece of Cake", world.WorldsQuestTitles[3])
	assert.True(world.BattleyeProtected)
	assert.Equal("2017-09-05", world.BattleyeDate)
	assert.Equal("regular", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(53, len(world.OnlinePlayers))
}

func TestWorldWintera(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Wintera.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Wintera", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World

	assert.Equal("Wintera", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(281, world.PlayersOnline)
	assert.Equal(1023, world.RecordPlayers)
	assert.Equal("2020-05-04T01:25:30Z", world.RecordDate)
	assert.Equal("2018-04", world.CreationDate)
	assert.Equal("North America", world.Location)
	assert.Equal("Open PvP", world.PvpType)
	assert.False(world.PremiumOnly)
	assert.Equal("regular", world.TransferType)
	assert.Equal(4, len(world.WorldsQuestTitles))
	assert.Equal("A Piece of Cake", world.WorldsQuestTitles[0])
	assert.Equal("Rise of Devovorga", world.WorldsQuestTitles[1])
	assert.Equal("Bewitched", world.WorldsQuestTitles[2])
	assert.Equal("The Colours of Magic", world.WorldsQuestTitles[3])
	assert.True(world.BattleyeProtected)
	assert.Equal("2018-04-19", world.BattleyeDate)
	assert.Equal("regular", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(281, len(world.OnlinePlayers))

	firstPlayer := world.OnlinePlayers[0]
	assert.Equal("Akiles Boy", firstPlayer.Name)
	assert.Equal(281, firstPlayer.Level)
	assert.Equal("Royal Paladin", firstPlayer.Vocation)
}

func TestWorldZuna(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Zuna.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Zuna", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World

	assert.Equal("Zuna", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(15, world.PlayersOnline)
	assert.Equal(174, world.RecordPlayers)
	assert.Equal("2019-07-29T16:55:30Z", world.RecordDate)
	assert.Equal("2017-10", world.CreationDate)
	assert.Equal("Europe", world.Location)
	assert.Equal("Hardcore PvP", world.PvpType)
	assert.False(world.PremiumOnly)
	assert.Equal("locked", world.TransferType)
	assert.Equal(2, len(world.WorldsQuestTitles))
	assert.Equal("The Colours of Magic", world.WorldsQuestTitles[0])
	assert.Equal("A Piece of Cake", world.WorldsQuestTitles[1])
	assert.False(world.BattleyeProtected)
	assert.Empty(world.BattleyeDate)
	assert.Equal("experimental", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(15, len(world.OnlinePlayers))

	firstPlayer := world.OnlinePlayers[0]
	assert.Equal("Bright soul", firstPlayer.Name)
	assert.Equal(20, firstPlayer.Level)
	assert.Equal("Paladin", firstPlayer.Vocation)
}

func TestWorldOceanis(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Oceanis.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Oceanis", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World

	assert.Equal("Oceanis", world.Name)
	assert.Equal("offline", world.Status)
	assert.Equal(0, world.PlayersOnline)
	assert.Equal(0, world.RecordPlayers)
	assert.Empty(world.RecordDate)
	assert.Equal("2024-04", world.CreationDate)
	assert.Equal("Oceania", world.Location)
	assert.Equal("Optional PvP", world.PvpType)
	assert.False(world.PremiumOnly)
	assert.Equal("regular", world.TransferType)
	assert.Equal(0, len(world.WorldsQuestTitles))
	assert.True(world.BattleyeProtected)
	assert.Equal("2024-04-10", world.BattleyeDate)
	assert.Equal("regular", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(0, len(world.OnlinePlayers))
}

func TestWorldTesta(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/worlds/world/Testa.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	worldJson, err := TibiaWorldsWorldImpl("Testa", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	world := worldJson.World

	assert.Equal("Testa", world.Name)
	assert.Equal("online", world.Status)
	assert.Equal(87, world.PlayersOnline)
	assert.Equal(781, world.RecordPlayers)
	assert.Equal("2012-12-04T15:42:47Z", world.RecordDate)
	assert.Equal("2013-01", world.CreationDate)
	assert.Equal("Europe", world.Location)
	assert.Equal("Optional PvP", world.PvpType)
	assert.False(world.PremiumOnly)
	assert.Equal("regular", world.TransferType)
	assert.Equal(0, len(world.WorldsQuestTitles))
	assert.True(world.BattleyeProtected)
	assert.Equal("release", world.BattleyeDate)
	assert.Equal("regular", world.GameWorldType)
	assert.Empty(world.TournamentWorldType)
	assert.Equal(87, len(world.OnlinePlayers))

	firstPlayer := world.OnlinePlayers[0]
	assert.Equal("Alius Harg E1", firstPlayer.Name)
	assert.Equal(8, firstPlayer.Level)
	assert.Equal("Monk", firstPlayer.Vocation)

	thirdPlayer := world.OnlinePlayers[2]
	assert.Equal("Andris Hun paladin E1", thirdPlayer.Name)
	assert.Equal(23, thirdPlayer.Level)
	assert.Equal("Exalted Monk", thirdPlayer.Vocation)
}
