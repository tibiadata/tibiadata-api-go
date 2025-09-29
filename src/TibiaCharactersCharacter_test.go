package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

func TestNumber1(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Darkside Rafa.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "https://www.tibia.com/community/?subtopic=characters&name=Darkside+Rafa")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo
	information := characterJson.Information

	assert.Equal("Darkside Rafa", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("male", character.Sex)
	assert.Equal("Silencer", character.Title)
	assert.Equal(18, character.UnlockedTitles)
	assert.Equal("Elite Knight", character.Vocation)
	assert.Equal(790, character.Level)
	assert.Equal(596, character.AchievementPoints)
	assert.Equal("Gladera", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Thais", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Nil(character.Houses)
	assert.Equal("Jokerz", character.Guild.GuildName)
	assert.Equal("Trial", character.Guild.Rank)
	assert.Equal("2022-01-05T21:23:32Z", character.LastLogin)
	assert.Equal("Premium Account", character.AccountStatus)
	assert.Empty(character.Comment)

	assert.Equal("https://www.tibia.com/community/?subtopic=characters&name=Darkside+Rafa", information.TibiaURLs[0])
}

func TestNumber2(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Zugspitze Housekeeper.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Zugspitze Housekeeper", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("male", character.Sex)
	assert.Equal("None", character.Title)
	assert.Equal(13, character.UnlockedTitles)
	assert.Equal("Elite Knight", character.Vocation)
	assert.Equal(79, character.Level)
	assert.Equal(262, character.AchievementPoints)
	assert.Equal("Venebra", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Darashia", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Equal(35056, character.Houses[0].HouseID)
	assert.Equal("Loot Lane 1 (Shop)", character.Houses[0].Name)
	assert.Equal("Venore", character.Houses[0].Town)
	assert.Equal("2022-01-16", character.Houses[0].Paid)
	assert.Equal("Magnus Magister", character.Guild.Rank)
	assert.Equal("Lionheart Society", character.Guild.GuildName)
	assert.Equal("2022-01-06T21:38:44Z", character.LastLogin)
	assert.Empty(character.Position)
	assert.Equal("Testa de Ferro do Lejonhjartat ;)", character.Comment)
	assert.Equal("Premium Account", character.AccountStatus)

	// validate other characters
	assert.Equal(7, len(characterJson.Character.OtherCharacters))

	onlineMainCharacter := characterJson.Character.OtherCharacters[3]
	assert.Equal("Lejonhjartat", onlineMainCharacter.Name)
	assert.Equal("Venebra", onlineMainCharacter.World)
	assert.Equal("online", onlineMainCharacter.Status)
	assert.False(onlineMainCharacter.Deleted)
	assert.True(onlineMainCharacter.Main)
	assert.False(onlineMainCharacter.Traded)
	assert.Empty(onlineMainCharacter.Position)

	offlineCharacter := characterJson.Character.OtherCharacters[5]
	assert.Equal("Oak Knight Disruptivo", offlineCharacter.Name)
	assert.Equal("Libertabra", offlineCharacter.World)
	assert.Equal("offline", offlineCharacter.Status)
	assert.False(offlineCharacter.Deleted)
	assert.False(offlineCharacter.Main)
	assert.False(offlineCharacter.Traded)
}

func TestNumber3(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Borttagna Gubben.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Borttagna Gubben", character.Name)
	assert.True(character.Traded)
	assert.Equal(2, len(character.FormerWorlds))
	assert.Equal("Zuna", character.FormerWorlds[0])
	assert.Equal("Zunera", character.FormerWorlds[1])
	assert.Equal("Bubble", character.MarriedTo)
	assert.Equal("2022-03-08T00:09:13Z", character.DeletionDate)
	assert.Empty(character.LastLogin)
	assert.Equal("Free Account", character.AccountStatus)
	assert.Equal("Fansite Admin", characterJson.Character.AccountInformation.Position)
	assert.Empty(characterJson.Character.AccountInformation.LoyaltyTitle)
	assert.True(characterJson.Character.OtherCharacters[0].Deleted)
	assert.False(characterJson.Character.OtherCharacters[1].Deleted)
}

func TestNumber4(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Riley No Hands.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Riley No Hands", character.Name)

	// validate former names
	assert.Equal(4, len(character.FormerNames))
	assert.Equal("Dura Malandro", character.FormerNames[0])
	assert.Equal("Letsgo Brandon", character.FormerNames[1])
	assert.Equal("Letsgo Brandon", character.FormerNames[2]) // yes, this name is listed twice
	assert.Equal("Nataraya Soldrac", character.FormerNames[3])

	// validate death data
	assert.Equal(79, len(characterJson.Character.Deaths))

	firstDeath := characterJson.Character.Deaths[0]
	assert.Equal(28, len(firstDeath.Killers))

	creatureWithOfDeath := characterJson.Character.Deaths[16]
	assert.Equal(2, len(creatureWithOfDeath.Killers))
	assert.Equal(260, creatureWithOfDeath.Level)
	assert.Equal("undead elite gladiator", creatureWithOfDeath.Killers[0].Name)
	assert.False(creatureWithOfDeath.Killers[0].Player)
	assert.False(creatureWithOfDeath.Killers[0].Traded)
	assert.Empty(creatureWithOfDeath.Killers[0].Summon)
	assert.Equal("priestess of the wild sun", creatureWithOfDeath.Killers[1].Name)

	tradedInDeath := characterJson.Character.Deaths[18]
	assert.Equal(3, len(tradedInDeath.Assists))
	assert.Equal(261, tradedInDeath.Level)
	assert.Equal("Vithrann", tradedInDeath.Assists[1].Name)
	assert.True(tradedInDeath.Assists[1].Player)
	assert.Equal("Adam No Hands", tradedInDeath.Assists[2].Name)
	assert.True(tradedInDeath.Assists[2].Traded)

	longDeath := characterJson.Character.Deaths[78]
	assert.Equal(5, len(longDeath.Assists))
	assert.Equal(231, longDeath.Level)
	assert.Equal("Adam No Hands", longDeath.Assists[4].Name)
	assert.Equal("paladin familiar", longDeath.Assists[4].Summon)
	assert.True(longDeath.Assists[4].Traded)

	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jazzim Paralizer", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Murder On Thedancefloor", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Color Chaufa", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "Trend Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
				{Name: "Niikzin", Player: true, Traded: false, Summon: ""},
			},
			Level:  264,
			Reason: "Annihilated at Level 264 by Jupa Infinity, Okiba Kay, Chino Kyle, Jazzim Paralizer, Dont Kill Joustick, Murder On Thedancefloor, Slobansky, Mnnuuuuczzhh, Bravefly Legend, Gorito Fullwar, Ez Josmart, Tiniebla Oddaeri, Rondero Color Chaufa, Jungle Rubi Dominante, Diego Rusher, Psycovzky, Retro Jupa, Nevarez Kyle, King Eteryo, Jupa Traicionado, Retro Demy, Daddy Chiill, Wrathfull Diegoz, Blokzie Prodigy, Pit Haveballs, Trend Fullwar, No Safeword and Niikzin.",
			Time:   "2022-01-05T05:00:55Z",
		},
		{
			Assists: []Killers{
				{Name: "Givsclap", Player: true, Traded: false, Summon: ""},
				{Name: "Richizawer", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Leora Em", Player: true, Traded: false, Summon: ""},
				{Name: "Jaime Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Ek Bombo", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Rick the Bold", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Reptile Stuns You", Player: true, Traded: false, Summon: ""},
				{Name: "Izrehsad Cigam", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jazzim Paralizer", Player: true, Traded: false, Summon: ""},
				{Name: "Murder On Thedancefloor", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Color Chaufa", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Azukiitaa", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
				{Name: "Niikzin", Player: true, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Annihilated at Level 266 by Jupa Infinity, Okiba Kay, Chino Kyle, Jazzim Paralizer, Murder On Thedancefloor, Mnnuuuuczzhh, Bravefly Legend, Demy Mythwar Infinity, Tiniebla Oddaeri, Rondero Color Chaufa, Ztiraleael, Jungle Rubi Dominante, Diego Rusher, Psycovzky, Gallo Kyle, Nevarez Kyle, Jupa Traicionado, Jupa Wezt, Retro Demy, Azukiitaa, Rondero Miguel, Daddy Chiill, Sam Kyle, Blokzie Prodigy, Pit Haveballs, No Safeword and Niikzin. Assisted by Givsclap, Richizawer, Samorbum, Notbrad, Leora Em, Jaime Ardera, Ek Bombo, Rek Bazilha, Tacos Ardera, Nick Pepperoni, Rick the Bold, Netozawer, Hardboss Remix, Hueviin, Reptile Stuns You and Izrehsad Cigam.",
			Time:   "2022-01-05T04:54:54Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "ice golem", Player: false, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Died at Level 266 by an ice golem.",
			Time:   "2022-01-05T04:42:21Z",
		},
		{
			Assists: []Killers{
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Keninho Sinmiedo", Player: true, Traded: false, Summon: ""},
				{Name: "Catwoman Lary", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Annihilated at Level 266 by Jupa Infinity, Okiba Kay, Jahziel Hardcori, Zeus Kyle, Astaroth Kyle, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Chino Kyle, Dont Kill Joustick, Keninho Sinmiedo, Catwoman Lary, Ez Josmart, Tiniebla Oddaeri, Ztiraleael, Psycovzky, Juanjo Infinity, Retro Jupa, Gallo Kyle, King Eteryo, Jupa Traicionado, Millonario Contatinho, Raven Kyle, Street Runner, Retro Demy, Robadob Wallstreet, Rondero Miguel, Wrathfull Diegoz and Renzzo Wallstreet. Assisted by Suprldo.",
			Time:   "2022-01-04T05:01:23Z",
		},
		{
			Assists: []Killers{
				{Name: "May Thirtieth", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Kanabionoia", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Tospa ficha gratis", Player: true, Traded: false, Summon: ""},
				{Name: "Waton Ekoz", Player: true, Traded: false, Summon: ""},
				{Name: "Ret Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekz", Player: true, Traded: false, Summon: ""},
			},

			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Jotade", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Kingsx", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Szpital Psychiatryczny", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Vithransky", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Luisinho Sinmiedo", Player: true, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Eliminated at Level 266 by Jupa Infinity, Dont Kill Trekib, Kaelsia Menardord, Rondero Jotade, Egdark, Von Rokitansky, Retro Kingsx, Demy Mythwar Infinity, Tiniebla Oddaeri, King Peruvian, Retro Jupa, Szpital Psychiatryczny, Millonario Contatinho, Vithransky, Pit Haveballs, Renzzo Wallstreet and Luisinho Sinmiedo. Assisted by May Thirtieth, Guichin Killzejk Boom, Kanabionoia, Suprldo, Tospa ficha gratis, Waton Ekoz, Ret Ardera and Sydekz.",
			Time:   "2022-01-04T04:14:52Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "gazer spectre", Player: false, Traded: false, Summon: ""},
			},
			Level:  264,
			Reason: "Died at Level 264 by a gazer spectre.",
			Time:   "2022-01-04T00:25:34Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Malofur Mangrinder", Player: false, Traded: false, Summon: ""},
			},
			Level:  265,
			Reason: "Died at Level 265 by Malofur Mangrinder.",
			Time:   "2022-01-03T22:39:01Z",
		},
		{
			Assists: []Killers{
				{Name: "Nlliliililiill", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Quali", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Catwoman Lary", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Friendly Smiley Face", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Highclass", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "Vithransky", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Annihilated at Level 266 by Kaelsia Menardord, Nevin kyle, Astaroth Kyle, Infinitywar Sange, Caezors, Dont Kill Joustick, Catwoman Lary, Ez Josmart, Rein is Here, King Peruvian, Touchy Dominante, Retro Jupa, Nevarez Kyle, Friendly Smiley Face, Javo Highclass, Raven Kyle, Rich Contatinho, Robadob Wallstreet, Papi Hydrowar, Sam Kyle, Gustavo Freefrags, Vithransky, Blokzie Prodigy and No Safeword. Assisted by Nlliliililiill, Kanj iro and Librarian Quali.",
			Time:   "2022-01-03T20:56:52Z",
		},
		{
			Assists: []Killers{
				{Name: "Nlliliililiill", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Broccolini", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Raffu", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Sikingg", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Quali", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Niki Salamanca", Player: true, Traded: false, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Highclass", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  267,
			Reason: "Crushed at Level 267 by Dont Kill Chelito, Kaelsia Menardord, Infinitywar Sange, Marchane kee (traded), Caezors, Rein is Here, Viictoryck, Diego Rusher, Retro Jupa, Javo Highclass, Raven Kyle, Street Runner, Blokzie Prodigy and No Safeword. Assisted by Nlliliililiill, Kanj iro, Broccolini, Yuniozawer, Librarian Raffu, Ell Rugalzawer, Sikingg, Librarian Quali, Uanzawer, Niki Salamanca and Pippah.",
			Time:   "2022-01-03T20:52:01Z",
		},
		{
			Assists: []Killers{
				{Name: "Rondero Momosilabo", Player: true, Traded: false, Summon: ""},
				{Name: "Basil Mendoza", Player: true, Traded: false, Summon: ""},
				{Name: "Broccolini", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Tospa ficha gratis", Player: true, Traded: false, Summon: ""},
				{Name: "Crypto Cowboy", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Sikingg", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Quali", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Catwoman Lary", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Highclass", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  268,
			Reason: "Annihilated at Level 268 by Dont Kill Chelito, Kaelsia Menardord, Astaroth Kyle, Infinitywar Sange, Marchane kee (traded), Caezors, Catwoman Lary, Rein is Here, Viictoryck, Touchy Dominante, Jungle Rubi Dominante, Diego Rusher, Retro Jupa, Nevarez Kyle, Javo Highclass, Raven Kyle, Street Runner, Retro Demy, Blokzie Prodigy and No Safeword. Assisted by Rondero Momosilabo, Basil Mendoza, Broccolini, Yuniozawer, Tospa ficha gratis, Crypto Cowboy, Ell Rugalzawer, Sikingg, Librarian Quali, Uanzawer and Pippah.",
			Time:   "2022-01-03T20:50:31Z",
		},
		{
			Assists: []Killers{
				{Name: "Voxiuoz", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
				{Name: "Symexz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Riley No Hands", Player: true, Traded: false, Summon: ""},
			},
			Level:  268,
			Reason: "Killed at Level 268 by Riley No Hands. Assisted by Voxiuoz, Nick Pepperoni, Pippah and Symexz.",
			Time:   "2022-01-02T23:22:53Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Scarlett Etzel", Player: false, Traded: false, Summon: ""},
			},
			Level:  268,
			Reason: "Died at Level 268 by Scarlett Etzel.",
			Time:   "2022-01-02T01:29:30Z",
		},
		{
			Assists: []Killers{
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "unstable spark", Player: false, Traded: false, Summon: ""},
			},
			Level:  267,
			Reason: "Died at Level 267 by an unstable spark. Assisted by Jack Kevorkian, Mapius Akuno and Kanj iro.",
			Time:   "2021-12-30T18:39:30Z",
		},
		{
			Assists: []Killers{
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Niki Salamanca", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "dawnfire asura", Player: false, Traded: false, Summon: ""},
			},
			Level:  261,
			Reason: "Died at Level 261 by a dawnfire asura. Assisted by Samorbum and Niki Salamanca.",
			Time:   "2021-12-30T01:13:18Z",
		},
		{
			Assists: []Killers{
				{Name: "Fuqueta", Player: true, Traded: false, Summon: ""},
				{Name: "Kanabionoia", Player: true, Traded: false, Summon: ""},
				{Name: "Exihva", Player: true, Traded: false, Summon: ""},
				{Name: "Leora Em", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Niki Salamanca", Player: true, Traded: false, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
				{Name: "Zamuxa", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Keninho Sinmiedo", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Kingsx", Player: true, Traded: false, Summon: ""},
				{Name: "Catwoman Lary", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Vyzerion", Player: true, Traded: false, Summon: ""},
				{Name: "Beowulf Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "Trend Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Luisinho Sinmiedo", Player: true, Traded: false, Summon: ""},
			},
			Level:  260,
			Reason: "Annihilated at Level 260 by Jupa Infinity, Okiba Kay, Monarca Mimi, Von Rokitansky, Infinitywar Sange, Chino Kyle, Marchane kee (traded), Rvzor Godslayer, Mnnuuuuczzhh, Keninho Sinmiedo, Retro Kingsx, Catwoman Lary, Utanii Herh, Rein is Here, King Peruvian, Ztiraleael, Monochrome Edowez, Elchico Billonario, Psycovzky, Retro Jupa, Gallo Kyle, King Eteryo, Jupa Traicionado, Vyzerion, Beowulf Kyle, Bello Abreu, Raven Kyle, Niko Tin, Woj Pokashield, Rondero Miguel, Sam Kyle, Wrathfull Diegoz, Pit Haveballs, Trend Fullwar, No Safeword, Renzzo Wallstreet and Luisinho Sinmiedo. Assisted by Fuqueta, Kanabionoia, Exihva, Leora Em, Kanj iro, Rek Bazilha, Niki Salamanca, Pippah and Zamuxa.",
			Time:   "2021-12-29T05:16:27Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Riley No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Dragonking Zyrtarch", Player: false, Traded: false, Summon: ""},
			},
			Level:  260,
			Reason: "Killed at Level 260 by Riley No Hands and Dragonking Zyrtarch.",
			Time:   "2021-12-29T01:50:17Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "undead elite gladiator", Player: false, Traded: false, Summon: ""},
				{Name: "priestess of the wild sun", Player: false, Traded: false, Summon: ""},
			},
			Level:  260,
			Reason: "Died at Level 260 by an undead elite gladiator and a priestess of the wild sun.",
			Time:   "2021-12-28T22:31:46Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Keninho Sinmiedo", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Don Dhimi Hardmode", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Kingsx", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Color Chaufa", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  260,
			Reason: "Annihilated at Level 260 by Jupa Infinity, Zeus Kyle, Keninho Sinmiedo, Bravefly Legend, Don Dhimi Hardmode, Retro Kingsx, Utanii Herh, Nyl Vy, Tiniebla Oddaeri, Rondero Color Chaufa, Ztiraleael, Monochrome Edowez, Diego Rusher, Rondero Tomachon, Millonario Contatinho, Rvzor, Locura Boss, Mnuczhhx, Robadob Wallstreet, Sam Kyle and Pit Haveballs.",
			Time:   "2021-12-28T04:19:21Z",
		},
		{
			Assists: []Killers{
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "midnight asura", Player: false, Traded: false, Summon: ""},
			},
			Level:  261,
			Reason: "Died at Level 261 by a midnight asura. Assisted by Peninsula Boi, Vithrann and Adam No Hands (traded).",
			Time:   "2021-12-28T02:15:39Z",
		},
		{
			Assists: []Killers{
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "midnight asura", Player: false, Traded: false, Summon: ""},
			},
			Level:  262,
			Reason: "Died at Level 262 by a midnight asura. Assisted by Peninsula Boi and Adam No Hands (traded).",
			Time:   "2021-12-28T02:02:19Z",
		},
		{
			Assists: []Killers{
				{Name: "Givsclap", Player: true, Traded: false, Summon: ""},
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Taimati Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Heethard", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Jotade", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  263,
			Reason: "Annihilated at Level 263 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Vrzik, Kaelsia Menardord, Rondero Jotade, Von Rokitansky, Calvo Bebe, Marchane kee (traded), Ez Josmart, Demy Mythwar Infinity, Rein is Here, Licaxn Kyle, Nyl Vy, King Peruvian, Jungle Rubi Dominante, Arthur Heartless, Psycovzky, Retro Jupa, Kyudok, King Eteryo, Rondero Tomachon, Millonario Contatinho, Niko Tin, Rvzor, Jupa Wezt, Rich Contatinho, Spyt Ponzi, Mnuczhhx, Contatinho Ekbomba, Wrathfull Diegoz, Fran Matarindo and Pit Haveballs. Assisted by Givsclap, Rabaab, Jack Kevorkian, Taimati Remix, Heethard, Yuniozawer, Peninsula Boi, Hardboss Remix, Hueviin, Uanzawer, Adam No Hands (traded) and Pippah.",
			Time:   "2021-12-27T23:31:21Z",
		},
		{
			Assists: []Killers{
				{Name: "Bleks Mortem", Player: true, Traded: false, Summon: ""},
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Exihva", Player: true, Traded: false, Summon: ""},
				{Name: "Kedruzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Leora Em", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Reptile Stuns You", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
				{Name: "Fuu Baz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Murder On Thedancefloor", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Rapido Marta", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Luismik", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Beishuu", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Ander Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Samil Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Quick Jamaica", Player: true, Traded: false, Summon: ""},
				{Name: "Savage Papi", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Highclass", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Blokzie Prodigy", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  261,
			Reason: "Annihilated at Level 261 by Jahziel Hardcori, Vrzik, Zeus Kyle, Marchane kee (traded), Orpheus Van Basten, Corker Poker, Murder On Thedancefloor, Bravefly Legend, Rapido Marta, Shady Is Back, Gorito Fullwar, Utanii Herh, Rein is Here, Luismik, King Peruvian, Monochrome Edowez, Beishuu, Arthur Heartless, Psycovzky, Retro Jupa, Tacos Ardera, Gallo Kyle, Ander Wallstreet, Nevarez Kyle, Jupa Traicionado, Samil Kyle, Quick Jamaica, Savage Papi, Javo Highclass, Street Runner, Rvzor, Retro Demy, Sam Kyle, Wrathfull Diegoz, Blokzie Prodigy and No Safeword. Assisted by Bleks Mortem, Notbrad, Exihva, Kedruzawer, Leora Em, Kanj iro, Yuniozawer, Peninsula Boi, Reptile Stuns You, Adam No Hands (traded) and Fuu Baz.",
			Time:   "2021-12-27T03:20:20Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "true dawnfire asura", Player: false, Traded: false, Summon: ""},
			},
			Level:  262,
			Reason: "Killed at Level 262 by Utanii Herh, King Peruvian, Raven Kyle and a true dawnfire asura.",
			Time:   "2021-12-27T01:40:37Z",
		},
		{
			Assists: []Killers{
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Mister Killer Kav", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Don Yoha", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  262,
			Reason: "Annihilated at Level 262 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Chino Kyle, Marchane kee (traded), Caezors, Rvzor Godslayer, Corker Poker, Contatinho Millonario, Demy Mythwar Infinity, Rein is Here, Tiniebla Oddaeri, Touchy Dominante, Ztiraleael, Elchico Billonario, Psycovzky, Retro Jupa, Bello Abreu, Niko Tin, Woj Pokashield, Jupa Wezt, Rich Contatinho, Don Yoha, Retro Demy, Sam Kyle, Gustavo Freefrags and No Safeword. Assisted by Mapius Akuno, Hueviin and Mister Killer Kav.",
			Time:   "2021-12-26T20:14:32Z",
		},
		{
			Assists: []Killers{
				{Name: "Don Brenjun", Player: true, Traded: false, Summon: ""},
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Schalama Rei Delas", Player: true, Traded: false, Summon: ""},
				{Name: "Anthony No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "May Thirtieth", Player: true, Traded: false, Summon: ""},
				{Name: "Daark Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Taimati Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Kana Vys", Player: true, Traded: false, Summon: ""},
				{Name: "Puds", Player: true, Traded: false, Summon: ""},
				{Name: "Exihva", Player: true, Traded: false, Summon: ""},
				{Name: "Kusuko", Player: true, Traded: false, Summon: ""},
				{Name: "Basil Mendoza", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekzu", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Raffu", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Mayj Boater", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Rick the Bold", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Niki Salamanca", Player: true, Traded: false, Summon: ""},
				{Name: "Symexz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Riley No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Savage Papi", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Luka Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  264,
			Reason: "Annihilated at Level 264 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Riley No Hands, Infinitywar Sange, Marchane kee (traded), Rvzor Godslayer, Orpheus Van Basten, Contatinho Millonario, Viictoryck, King Peruvian, Touchy Dominante, Psycovzky, Retro Jupa, Jupa Traicionado, Bello Abreu, Savage Papi, Raven Kyle, Woj Pokashield, Jupa Wezt, King Rexiiruz, Luka Is Back, Sam Kyle, Gustavo Freefrags and No Safeword. Assisted by Don Brenjun, Rabaab, Schalama Rei Delas, Anthony No Hands, May Thirtieth, Daark Remix, Jack Kevorkian, Taimati Remix, Kana Vys, Puds, Exihva, Kusuko, Basil Mendoza, Kanj iro, Sydekzu, Librarian Raffu, Rek Bazilha, Mayj Boater, Nick Pepperoni, Rick the Bold, Hardboss Remix, Hueviin, Niki Salamanca and Symexz.",
			Time:   "2021-12-26T20:10:19Z",
		},
		{
			Assists: []Killers{
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Middle Zocarno", Player: true, Traded: true, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Tospa ficha gratis", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Yisus Jusable", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Beishuu", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Luka Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "Vithransky", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  264,
			Reason: "Annihilated at Level 264 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Zeus Kyle, Yisus Jusable, Infinitywar Sange, Chino Kyle, Marchane kee (traded), Contatinho Millonario, Rein is Here, Retro Mega, Beishuu, Elchico Billonario, Psycovzky, Retro Jupa, Bello Abreu, Raven Kyle, Street Runner, Woj Pokashield, Jupa Wezt, King Rexiiruz, Luka Is Back, Retro Demy, Sam Kyle, Gustavo Freefrags, Vithransky and No Safeword. Assisted by Notbrad, Middle Zocarno (traded), Kanj iro and Tospa ficha gratis.",
			Time:   "2021-12-26T20:06:36Z",
		},
		{
			Assists: []Killers{
				{Name: "Don Brenjun", Player: true, Traded: false, Summon: ""},
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "May Thirtieth", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Nytrander", Player: true, Traded: false, Summon: ""},
				{Name: "Kana Vys", Player: true, Traded: false, Summon: ""},
				{Name: "Mal Victus", Player: true, Traded: false, Summon: ""},
				{Name: "Elpa Tron", Player: true, Traded: false, Summon: ""},
				{Name: "Puds", Player: true, Traded: false, Summon: ""},
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Exihva", Player: true, Traded: false, Summon: ""},
				{Name: "Kusuko", Player: true, Traded: false, Summon: ""},
				{Name: "Middle Zocarno", Player: true, Traded: true, Summon: ""},
				{Name: "Qualitie", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Raffu", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Symexz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Beishuu", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Luka Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Don Yoha", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
				{Name: "Vithransky", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  265,
			Reason: "Annihilated at Level 265 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Zeus Kyle, Infinitywar Sange, Marchane kee (traded), Caezors, Rvzor Godslayer, Orpheus Van Basten, Corker Poker, Mnnuuuuczzhh, Contatinho Millonario, Rein is Here, Viictoryck, Nyl Vy, Tiniebla Oddaeri, Retro Mega, King Peruvian, Beishuu, Elchico Billonario, Psycovzky, Retro Jupa, Nevarez Kyle, Jupa Traicionado, Bello Abreu, Raven Kyle, Street Runner, Niko Tin, Rich Contatinho, Luka Is Back, Don Yoha, Retro Demy, Contatinho Ekbomba, Rondero Miguel, Sam Kyle, Gustavo Freefrags, Wrathfull Diegoz, Fran Matarindo, Vithransky and No Safeword. Assisted by Don Brenjun, Rabaab, May Thirtieth, Jack Kevorkian, Nytrander, Kana Vys, Mal Victus, Elpa Tron, Puds, Notbrad, Exihva, Kusuko, Middle Zocarno (traded), Qualitie, Librarian Raffu, Rek Bazilha, Hardboss Remix and Symexz.",
			Time:   "2021-12-26T20:03:55Z",
		},
		{
			Assists: []Killers{
				{Name: "Librarian Bito", Player: true, Traded: false, Summon: ""},
				{Name: "Qualitie", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Savage Papi", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  266,
			Reason: "Annihilated at Level 266 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Calvo Bebe, Dont Kill Joustick, Contatinho Millonario, Demy Mythwar Infinity, Rein is Here, Nyl Vy, Retro Mega, Ztiraleael, Elchico Billonario, Bello Abreu, Savage Papi, Raven Kyle, King Rexiiruz, Retro Demy, Rondero Miguel, Daddy Chiill, Fran Matarindo and No Safeword. Assisted by Librarian Bito and Qualitie.",
			Time:   "2021-12-26T19:57:32Z",
		},
		{
			Assists: []Killers{
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Duke Krule", Player: false, Traded: false, Summon: ""},
			},
			Level:  259,
			Reason: "Died at Level 259 by Duke Krule. Assisted by Mapius Akuno.",
			Time:   "2021-12-26T06:10:47Z",
		},
		{
			Assists: []Killers{
				{Name: "Mal Victus", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Rick the Bold", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Jotade", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
			},
			Level:  241,
			Reason: "Slain at Level 241 by Jupa Infinity, Rondero Jotade, Zeus Kyle, Nevarez Kyle, Raven Kyle and Rondero Miguel. Assisted by Mal Victus, Yuniozawer, Rick the Bold and Hardboss Remix.",
			Time:   "2021-12-24T02:49:56Z",
		},
		{
			Assists: []Killers{
				{Name: "Toxir Golpista", Player: true, Traded: false, Summon: ""},
				{Name: "Daark Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Taimati Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Tospa ficha gratis", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Raffu", Player: true, Traded: false, Summon: ""},
				{Name: "Magic Dasherzi", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Rapido Marta", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  237,
			Reason: "Annihilated at Level 237 by Jupa Infinity, Okiba Kay, Monarca Mimi, Kaelsia Menardord, Zeus Kyle, Marchane kee (traded), Caezors, Rvzor Godslayer, Corker Poker, Rapido Marta, Nyl Vy, Tiniebla Oddaeri, Monochrome Edowez, Juanjo Infinity, Retro Jupa, Gallo Kyle, Jupa Traicionado, Bello Abreu, Street Runner, Mnuczhhx, Papi Hydrowar, Daddy Chiill and Wrathfull Diegoz. Assisted by Toxir Golpista, Daark Remix, Taimati Remix, Guichin Killzejk Boom, Kanj iro, Tospa ficha gratis, Librarian Raffu and Magic Dasherzi.",
			Time:   "2021-12-23T20:54:45Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  239,
			Reason: "Crushed at Level 239 by Okiba Kay, Dont Kill Chelito, Monarca Mimi, Kaelsia Menardord, Zeus Kyle, Marchane kee (traded), Rvzor Godslayer, Juanjo Infinity, Retro Jupa and Wrathfull Diegoz.",
			Time:   "2021-12-23T20:48:00Z",
		},
		{
			Assists: []Killers{
				{Name: "Taimati Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Magic Dasherzi", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Rapido Marta", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
			},
			Level:  240,
			Reason: "Crushed at Level 240 by Okiba Kay, Dont Kill Chelito, Monarca Mimi, Zeus Kyle, Marchane kee (traded), Caezors, Rapido Marta, Tiniebla Oddaeri, Monochrome Edowez, Juanjo Infinity, Bello Abreu, Street Runner and Daddy Chiill. Assisted by Taimati Remix, Magic Dasherzi and Adam No Hands (traded).",
			Time:   "2021-12-23T20:45:27Z",
		},
		{
			Assists: []Killers{
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Nytrander", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Aereaere", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  241,
			Reason: "Eliminated at Level 241 by Okiba Kay, Dont Kill Chelito, Monarca Mimi, Zeus Kyle, Caezors, Corker Poker, Nyl Vy, Tiniebla Oddaeri, Retro Jupa, Nevarez Kyle, Jupa Traicionado, Bello Abreu, Niko Tin, Rvzor, Spyt Ponzi, Ronald Ardera (traded), Papi Hydrowar, Daddy Chiill and Wrathfull Diegoz. Assisted by Jack Kevorkian, Nytrander, Mapius Akuno and Aereaere.",
			Time:   "2021-12-23T20:31:06Z",
		},
		{
			Assists: []Killers{
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Broccolini", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Finanze", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
			},
			Level:  241,
			Reason: "Crushed at Level 241 by Jupa Infinity, Okiba Kay, Dont Kill Chelito, Calvo Bebe, Caezors, Retro Mega, Monochrome Edowez, Juanjo Infinity, Retro Jupa, Rondero Tomachon, Finanze, Rvzor, Spyt Ponzi and Retro Demy. Assisted by Kanj iro, Broccolini and Adam No Hands (traded).",
			Time:   "2021-12-23T19:20:14Z",
		},
		{
			Assists: []Killers{
				{Name: "Toxir Golpista", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Momosilabo", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Broccolini", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "West Nuukldragor", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Niki Salamanca", Player: true, Traded: false, Summon: ""},
				{Name: "Mister Killer Kav", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Jotade", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Saez Uchiha", Player: true, Traded: false, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Corker Poker", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  241,
			Reason: "Annihilated at Level 241 by Jupa Infinity, Okiba Kay, Dont Kill Chelito, Monarca Mimi, Vrzik, Rondero Jotade, Calvo Bebe, Saez Uchiha, Orpheus Van Basten, Corker Poker, Dont Kill Joustick, Gorito Fullwar, Contatinho Millonario, Rein is Here, Retro Mega, King Peruvian, Arthur Heartless, Retro Jupa, Gallo Kyle, Nevarez Kyle, Bello Abreu, Street Runner, Spyt Ponzi and Renzzo Wallstreet. Assisted by Toxir Golpista, Rondero Momosilabo, Kanj iro, Broccolini, Suprldo, Yuniozawer, West Nuukldragor, Uanzawer, Niki Salamanca and Mister Killer Kav.",
			Time:   "2021-12-23T17:59:56Z",
		},
		{
			Assists: []Killers{
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Tospa ficha gratis", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  242,
			Reason: "Eliminated at Level 242 by Jupa Infinity, Okiba Kay, Monarca Mimi, Vrzik, Caezors, Contatinho Millonario, Rein is Here, Retro Mega, King Peruvian, Arthur Heartless, Juanjo Infinity, Nevarez Kyle, Street Runner, Mnuczhhx and Renzzo Wallstreet. Assisted by Suprldo and Tospa ficha gratis.",
			Time:   "2021-12-23T17:55:47Z",
		},
		{
			Assists: []Killers{
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Daark Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Ek Bombo", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuuucczhh", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Versatilsz", Player: true, Traded: true, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
			},
			Level:  242,
			Reason: "Eliminated at Level 242 by Jupa Infinity, Okiba Kay, Dont Kill Trekib, Vrzik, Chino Kyle, Shady Is Back, Rein is Here, Mnuuucczhh, King Peruvian, Monochrome Edowez, Versatilsz (traded), Retro Jupa, Tacos Ardera, Nevarez Kyle and Street Runner. Assisted by Rabaab, Daark Remix, Guichin Killzejk Boom, Mapius Akuno, Ek Bombo and Sydekz.",
			Time:   "2021-12-23T04:39:07Z",
		},
		{
			Assists: []Killers{
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Cybago", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekzu", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Ek Bombo", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Meiker de Tozir", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuuucczhh", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  243,
			Reason: "Annihilated at Level 243 by Jupa Infinity, Okiba Kay, Jahziel Hardcori, Dont Kill Trekib, Chino Kyle, Slobansky, Shady Is Back, Ez Josmart, Rein is Here, Tiniebla Oddaeri, Mnuuucczhh, King Peruvian, Monochrome Edowez, Diego Rusher, Psycovzky, Retro Jupa, Tacos Ardera, Nevarez Kyle, Street Runner, Robadob Wallstreet, Papi Hydrowar, Wrathfull Diegoz, Pit Haveballs, No Safeword and Renzzo Wallstreet. Assisted by Rabaab, Jack Kevorkian, Cybago, Acid Zero, Kanj iro, Sydekzu, Suprldo, Ek Bombo, Nick Pepperoni, Hardboss Remix and Meiker de Tozir.",
			Time:   "2021-12-23T04:35:35Z",
		},
		{
			Assists: []Killers{
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Highclass", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
				{Name: "Renzzo Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  244,
			Reason: "Eliminated at Level 244 by Jupa Infinity, Okiba Kay, Dont Kill Trekib, Slobansky, Shady Is Back, King Peruvian, Diego Rusher, Psycovzky, Retro Jupa, Tacos Ardera, Nevarez Kyle, Javo Highclass, Street Runner, Retro Demy, Robadob Wallstreet, Papi Hydrowar, Wrathfull Diegoz, Pit Haveballs and Renzzo Wallstreet. Assisted by Rabaab, Jack Kevorkian, Guichin Killzejk Boom, Acid Zero, Suprldo, Yuniozawer and Hardboss Remix.",
			Time:   "2021-12-23T04:30:31Z",
		},
		{
			Assists: []Killers{
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Zurdo Mnuczh", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Ramzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  236,
			Reason: "Annihilated at Level 236 by Jupa Infinity, Astaroth Kyle, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Chino Kyle, Zurdo Mnuczh, Mnnuuuuczzhh, Drunkz Wallstreet (traded), Shady Is Back, Tiniebla Oddaeri, Psycovzky, Juanjo Infinity, King Eteryo, Chapo Kyle, Millonario Contatinho, Street Runner, Mnuczhhx, Retro Demy, Ramzawer and Pit Haveballs. Assisted by Guichin Killzejk Boom, Kanj iro, Rek Bazilha and Hardboss Remix.",
			Time:   "2021-12-21T04:26:30Z",
		},
		{
			Assists: []Killers{
				{Name: "Rabaab", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Shmurdad", Player: true, Traded: false, Summon: ""},
				{Name: "Leora Em", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Fuuba Diretoria", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Ramzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  238,
			Reason: "Eliminated at Level 238 by Jupa Infinity, Astaroth Kyle, Von Rokitansky, Calvo Bebe, Marchane kee (traded), Shady Is Back, Licaxn Kyle, Retro Mega, Psycovzky, Juanjo Infinity, Millonario Contatinho, Street Runner, Locura Boss, Mnuczhhx, Retro Demy, Ramzawer and Pit Haveballs. Assisted by Rabaab, Guichin Killzejk Boom, Mapius Akuno, Shmurdad, Leora Em, Rek Bazilha, Fuuba Diretoria, Netozawer, Hardboss Remix and Adam No Hands (traded).",
			Time:   "2021-12-21T04:24:10Z",
		},
		{
			Assists: []Killers{
				{Name: "Rondero Momosilabo", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Fuuba Diretoria", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Rick the Bold", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Reptile Stuns You", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Ashley Chart", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Juantootwee", Player: true, Traded: true, Summon: ""},
				{Name: "Joao Gomes Tanavoz", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "No Safeword", Player: true, Traded: false, Summon: ""},
			},
			Level:  239,
			Reason: "Annihilated at Level 239 by Jupa Infinity, Zeus Kyle, Astaroth Kyle, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Chino Kyle, Marchane kee (traded), Mnnuuuuczzhh, Drunkz Wallstreet (traded), Gorito Fullwar, Licaxn Kyle, Nyl Vy, Retro Mega, Touchy Dominante, Jungle Rubi Dominante, Psycovzky, Juanjo Infinity, Retro Jupa, Kyudok, King Eteryo, Chapo Kyle, Ashley Chart, Millonario Contatinho, Juantootwee (traded), Joao Gomes Tanavoz, Raven Kyle, Street Runner, Locura Boss, Contatinho Ekbomba and No Safeword. Assisted by Rondero Momosilabo, Mapius Akuno, Ell Rugalzawer, Fuuba Diretoria, Nick Pepperoni, Rick the Bold, Hardboss Remix, Reptile Stuns You and Adam No Hands (traded).",
			Time:   "2021-12-21T04:22:53Z",
		},
		{
			Assists: []Killers{
				{Name: "Richizawer", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Ek Bombo", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
				{Name: "Fuu Baz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Valto Soug", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Senior Marron", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Licus Rubiking Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Versatilsz", Player: true, Traded: true, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Quick Jamaica", Player: true, Traded: false, Summon: ""},
				{Name: "Juantootwee", Player: true, Traded: true, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Don Yoha", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Ramzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  239,
			Reason: "Annihilated at Level 239 by Jupa Infinity, Zeus Kyle, Von Rokitansky, Infinitywar Sange, Valto Soug, Calvo Bebe, Chino Kyle, Marchane kee (traded), Rvzor Godslayer, Mnnuuuuczzhh, Drunkz Wallstreet (traded), Shady Is Back, Gorito Fullwar, Contatinho Millonario, Ez Josmart, Demy Mythwar Infinity, Senior Marron, Nyl Vy, Licus Rubiking Dominante, Touchy Dominante, Ztiraleael, Monochrome Edowez, Versatilsz (traded), Arthur Heartless, Psycovzky, Retro Jupa, Jupa Traicionado, Chapo Kyle, Quick Jamaica, Juantootwee (traded), Raven Kyle, Street Runner, Rich Contatinho, Locura Boss, Mnuczhhx, Don Yoha, Retro Demy, Contatinho Ekbomba, Ramzawer, Wrathfull Diegoz and Pit Haveballs. Assisted by Richizawer, Jack Kevorkian, Ek Bombo, Rek Bazilha, Nick Pepperoni, Vithrann and Fuu Baz.",
			Time:   "2021-12-21T03:25:59Z",
		},
		{
			Assists: []Killers{
				{Name: "Elpa Tron", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Reptile Stuns You", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Valto Soug", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Emma", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Relax Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Ramzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  240,
			Reason: "Eliminated at Level 240 by Jupa Infinity, Von Rokitansky, Valto Soug, Calvo Bebe, Monarca Emma, Mnnuuuuczzhh, Relax Infinity, Tiniebla Oddaeri, Ztiraleael, Monochrome Edowez, Psycovzky, Chapo Kyle, Raven Kyle, Retro Demy, Ramzawer and Pit Haveballs. Assisted by Elpa Tron, Rek Bazilha and Reptile Stuns You.",
			Time:   "2021-12-21T03:18:44Z",
		},
		{
			Assists: []Killers{
				{Name: "Rick the Bold", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Licus Rubiking Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Gallo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  239,
			Reason: "Crushed at Level 239 by Jupa Infinity, Slobansky, Gorito Fullwar, Utanii Herh, Licaxn Kyle, Licus Rubiking Dominante, Monochrome Edowez, Gallo Kyle, King Eteryo, Niko Tin, Retro Demy, Papi Hydrowar, Wrathfull Diegoz and Pit Haveballs. Assisted by Rick the Bold.",
			Time:   "2021-12-20T03:41:50Z",
		},
		{
			Assists: []Killers{
				{Name: "Givsclap", Player: true, Traded: false, Summon: ""},
				{Name: "Veworth Tiva", Player: true, Traded: false, Summon: ""},
				{Name: "Voxiuoz", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  240,
			Reason: "Slain at Level 240 by Jupa Infinity, Marchane kee (traded), Rein is Here, Ztiraleael, Retro Jupa, Niko Tin, Ronald Ardera (traded), Retro Demy and Pit Haveballs. Assisted by Givsclap, Veworth Tiva, Voxiuoz, Peninsula Boi and Sydekz.",
			Time:   "2021-12-20T03:28:22Z",
		},
		{
			Assists: []Killers{
				{Name: "Givsclap", Player: true, Traded: false, Summon: ""},
				{Name: "Don Brenjun", Player: true, Traded: false, Summon: ""},
				{Name: "Schalama Rei Delas", Player: true, Traded: false, Summon: ""},
				{Name: "Daark Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Valdez", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Heste no", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "West Nuukldragor", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Ez Josmart", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Senior Marron", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Licus Rubiking Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Versatilsz", Player: true, Traded: true, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "Don mizzi", Player: true, Traded: false, Summon: ""},
				{Name: "King Eteryo", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Ashley Chart", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Waton Momo", Player: true, Traded: true, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Don Yoha", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Haveballs", Player: true, Traded: false, Summon: ""},
			},
			Level:  241,
			Reason: "Annihilated at Level 241 by Jupa Infinity, Okiba Kay, Kaelsia Menardord, Nevin kyle, Von Rokitansky, Marchane kee (traded), Slobansky, Bravefly Legend, Ez Josmart, Demy Mythwar Infinity, Senior Marron, Utanii Herh, Rein is Here, Licus Rubiking Dominante, Ztiraleael, Monochrome Edowez, Versatilsz (traded), Arthur Heartless, Psycovzky, Retro Jupa, Tacos Ardera, Kyudok, Don mizzi, King Eteryo, Chapo Kyle, Ashley Chart, Niko Tin, Waton Momo (traded), Spyt Ponzi, Don Yoha, Robadob Wallstreet, Contatinho Ekbomba, Sam Kyle, Wrathfull Diegoz and Pit Haveballs. Assisted by Givsclap, Don Brenjun, Schalama Rei Delas, Daark Remix, Jack Kevorkian, Guichin Killzejk Boom, Diego Valdez, Kanj iro, Heste no, Peninsula Boi, Netozawer, Hardboss Remix, West Nuukldragor and Adam No Hands (traded).",
			Time:   "2021-12-20T03:26:21Z",
		},
		{
			Assists: []Killers{
				{Name: "Anthony No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Snurggle", Player: true, Traded: false, Summon: ""},
				{Name: "Ek Bombo", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "Don Yoha", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
				{Name: "Milleerkzs", Player: true, Traded: false, Summon: ""},
				{Name: "Kalecgos Arcane Aspect", Player: true, Traded: true, Summon: ""},
			},
			Level:  241,
			Reason: "Crushed at Level 241 by Von Rokitansky, Slobansky, Bravefly Legend, Shady Is Back, Elchico Billonario, Retro Jupa, Tacos Ardera, Kyudok, Don Yoha, Wrathfull Diegoz, Milleerkzs and Kalecgos Arcane Aspect (traded). Assisted by Anthony No Hands, Jack Kevorkian, Snurggle, Ek Bombo, Ell Rugalzawer and Nick Pepperoni.",
			Time:   "2021-12-19T03:58:58Z",
		},
		{
			Assists: []Killers{
				{Name: "Snurggle", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Burro Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
			},
			Level:  242,
			Reason: "Crushed at Level 242 by Von Rokitansky, Bravefly Legend, Shady Is Back, Monochrome Edowez, Psycovzky, Retro Jupa, Tacos Ardera, Kyudok, Rich Contatinho, Burro Kyle and Retro Demy. Assisted by Snurggle, Ell Rugalzawer, Nick Pepperoni, Netozawer and Hardboss Remix.",
			Time:   "2021-12-19T03:55:11Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Quick Jamaica", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Milleerkzs", Player: true, Traded: false, Summon: ""},
			},
			Level:  236,
			Reason: "Annihilated at Level 236 by Jupa Infinity, Kaelsia Menardord, Nevin kyle, Slobansky, Viictoryck, Monochrome Edowez, Arthur Heartless, Psycovzky, Retro Jupa, Jupa Traicionado, Chapo Kyle, Quick Jamaica, Rvzor, Rich Contatinho, Spyt Ponzi, King Rexiiruz, Retro Demy, Robadob Wallstreet, Sam Kyle and Milleerkzs.",
			Time:   "2021-12-18T03:08:16Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Marta", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: "sorcerer familiar"},
			},
			Level:  237,
			Reason: "Slain at Level 237 by Nevin kyle, Dont Kill Marta, Rein is Here, Arthur Heartless, Psycovzky, Tacos Ardera, Street Runner and a sorcerer familiar of Slobansky.",
			Time:   "2021-12-18T02:53:53Z",
		},
		{
			Assists: []Killers{
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor Godslayer", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Don Dhimi Hardmode", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  238,
			Reason: "Annihilated at Level 238 by Jupa Infinity, Jahziel Hardcori, Nevin kyle, Von Rokitansky, Chino Kyle, Rvzor Godslayer, Slobansky, Don Dhimi Hardmode, Rein is Here, Viictoryck, Tiniebla Oddaeri, Arthur Heartless, Psycovzky, Retro Jupa, Tacos Ardera, Jupa Traicionado, Street Runner, Woj Pokashield, Rich Contatinho, Retro Demy, Robadob Wallstreet, Rondero Miguel and Wrathfull Diegoz. Assisted by Vithrann.",
			Time:   "2021-12-18T02:45:07Z",
		},
		{
			Assists: []Killers{
				{Name: "Fuubaz Ltda", Player: true, Traded: false, Summon: ""},
				{Name: "Adyn Edeus", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Veworth Tiva", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Og Hardmode", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Be Nicee", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Rusher", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Waton Momo", Player: true, Traded: true, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  222,
			Reason: "Annihilated at Level 222 by Dont Kill Trekib, Monarca Mimi, Zeus Kyle, Egdark, Infinitywar Sange, Calvo Bebe, Chino Kyle, Marchane kee (traded), Slobansky, Be Nicee, Utanii Herh, Viictoryck, Tiniebla Oddaeri, King Peruvian, Jungle Rubi Dominante, Monochrome Edowez, Diego Rusher, Retro Jupa, Tacos Ardera, Jupa Traicionado, Chapo Kyle, Millonario Contatinho, Waton Momo (traded), Jupa Wezt, Rich Contatinho, King Rexiiruz, Sam Kyle and Wrathfull Diegoz. Assisted by Fuubaz Ltda, Adyn Edeus, Guichin Killzejk Boom, Mapius Akuno, Veworth Tiva, Rek Bazilha, Og Hardmode, Hardboss Remix and Adam No Hands (traded).",
			Time:   "2021-12-17T04:26:44Z",
		},
		{
			Assists: []Killers{
				{Name: "Cybago", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Leo Madd", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Trekib", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Chiletton", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
			},
			Level:  223,
			Reason: "Annihilated at Level 223 by Jupa Infinity, Dont Kill Chelito, Dont Kill Trekib, Monarca Mimi, Zeus Kyle, Egdark, Infinitywar Sange, Chino Kyle, Marchane kee (traded), Slobansky, Shady Is Back, Utanii Herh, Rein is Here, Chiletton, King Peruvian, Jungle Rubi Dominante, Arthur Heartless, Psycovzky, Retro Jupa, Jupa Traicionado, Millonario Contatinho, Niko Tin, Jupa Wezt, Rich Contatinho, King Rexiiruz, Contatinho Ekbomba, Rondero Miguel and Sam Kyle. Assisted by Cybago, Samorbum, Leo Madd, Yuniozawer, Nick Pepperoni, Peninsula Boi, Hardboss Remix, Vithrann, Adam No Hands (traded) and Pippah.",
			Time:   "2021-12-17T04:18:49Z",
		},
		{
			Assists: []Killers{
				{Name: "Don Brenjun", Player: true, Traded: false, Summon: ""},
				{Name: "May Thirtieth", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Jojo Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Adyn Edeus", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Leo Madd", Player: true, Traded: false, Summon: ""},
				{Name: "Puds", Player: true, Traded: false, Summon: ""},
				{Name: "Kedruzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Reverse Terry", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekzu", Player: true, Traded: false, Summon: ""},
				{Name: "Heste no", Player: true, Traded: false, Summon: ""},
				{Name: "Jaime Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Voxiuoz", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
			},
			Level:  224,
			Reason: "Annihilated at Level 224 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Zeus Kyle, Egdark, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Marchane kee (traded), Shady Is Back, Utanii Herh, Rein is Here, Viictoryck, Licaxn Kyle, Tiniebla Oddaeri, King Peruvian, Jungle Rubi Dominante, Retro Jupa, Tacos Ardera, Jupa Traicionado, Jupa Wezt, Rich Contatinho, King Rexiiruz and Contatinho Ekbomba. Assisted by Don Brenjun, May Thirtieth, Jack Kevorkian, Jojo Ardera, Adyn Edeus, Mapius Akuno, Leo Madd, Puds, Kedruzawer, Reverse Terry, Sydekzu, Heste no, Jaime Ardera, Voxiuoz, Ell Rugalzawer, Peninsula Boi, Adam No Hands (traded) and Pippah.",
			Time:   "2021-12-17T04:14:18Z",
		},
		{
			Assists: []Killers{
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Puds", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Never Forget Loyalty", Player: true, Traded: false, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuuucczhh", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
			},
			Level:  225,
			Reason: "Annihilated at Level 225 by Jupa Infinity, Monarca Mimi, Zeus Kyle, Egdark, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Chino Kyle, Marchane kee (traded), Slobansky, Viictoryck, Tiniebla Oddaeri, Mnuuucczhh, King Peruvian, Monochrome Edowez, Psycovzky, Retro Jupa, Tacos Ardera, Jupa Traicionado, Millonario Contatinho, Woj Pokashield, Rich Contatinho, King Rexiiruz, Robadob Wallstreet, Contatinho Ekbomba and Rondero Miguel. Assisted by Guichin Killzejk Boom, Puds, Netozawer, Hardboss Remix, Never Forget Loyalty and Pippah.",
			Time:   "2021-12-17T04:11:34Z",
		},
		{
			Assists: []Killers{
				{Name: "Don Brenjun", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Mal Victus", Player: true, Traded: false, Summon: ""},
				{Name: "Puds", Player: true, Traded: false, Summon: ""},
				{Name: "Reverse Terry", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
				{Name: "Pippah", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Egdark", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Javo Gryffindor", Player: true, Traded: false, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuuucczhh", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Ztiraleael", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Wezt", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  225,
			Reason: "Annihilated at Level 225 by Jupa Infinity, Dont Kill Chelito, Monarca Mimi, Nevin kyle, Zeus Kyle, Egdark, Von Rokitansky, Infinitywar Sange, Calvo Bebe, Chino Kyle, Marchane kee (traded), Javo Gryffindor, Orpheus Van Basten, Slobansky, Viictoryck, Licaxn Kyle, Tiniebla Oddaeri, Mnuuucczhh, King Peruvian, Ztiraleael, Jungle Rubi Dominante, Monochrome Edowez, Arthur Heartless, Psycovzky, Retro Jupa, Tacos Ardera, Jupa Traicionado, Millonario Contatinho, Niko Tin, Rvzor, Woj Pokashield, Jupa Wezt, Rich Contatinho, King Rexiiruz, Robadob Wallstreet, Contatinho Ekbomba, Rondero Miguel and Wrathfull Diegoz. Assisted by Don Brenjun, Guichin Killzejk Boom, Mal Victus, Puds, Reverse Terry, Hueviin, Adam No Hands (traded) and Pippah.",
			Time:   "2021-12-17T04:07:39Z",
		},
		{
			Assists: []Killers{
				{Name: "Nytrander", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Kusuko", Player: true, Traded: false, Summon: ""},
				{Name: "Valkh Golpista", Player: true, Traded: false, Summon: ""},
				{Name: "Crypto Cowboy", Player: true, Traded: false, Summon: ""},
				{Name: "Don Ballusse", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Orpheus Van Basten", Player: true, Traded: false, Summon: ""},
				{Name: "Rapido Marta", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "No Peleo Solo", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
				{Name: "Wrathfull Diegoz", Player: true, Traded: false, Summon: ""},
			},
			Level:  222,
			Reason: "Crushed at Level 222 by Von Rokitansky, Orpheus Van Basten, Rapido Marta, King Peruvian, No Peleo Solo, Juanjo Infinity, Retro Jupa, Tacos Ardera, Robadob Wallstreet, Contatinho Ekbomba, Rondero Miguel, Gustavo Freefrags and Wrathfull Diegoz. Assisted by Nytrander, Acid Zero, Kusuko, Valkh Golpista, Crypto Cowboy, Don Ballusse, Rek Bazilha, Netozawer and Uanzawer.",
			Time:   "2021-12-17T00:02:36Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Chiletton", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
			},
			Level:  223,
			Reason: "Slain at Level 223 by Von Rokitansky, Chiletton, Juanjo Infinity, Tacos Ardera and Jupa Traicionado.",
			Time:   "2021-12-16T22:57:53Z",
		},
		{
			Assists: []Killers{
				{Name: "Richizawer", Player: true, Traded: false, Summon: ""},
				{Name: "Kedruzawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "fire", Player: false, Traded: false, Summon: ""},
			},
			Level:  223,
			Reason: "Killed at Level 223 by Rein is Here and fire. Assisted by Richizawer and Kedruzawer.",
			Time:   "2021-12-16T01:39:02Z",
		},
		{
			Assists: []Killers{
				{Name: "Eszex Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Ah Pepapipa", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Heethard", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Lilpreben", Player: true, Traded: false, Summon: ""},
				{Name: "Valkh Golpista", Player: true, Traded: false, Summon: ""},
				{Name: "Caladan Bane", Player: true, Traded: false, Summon: ""},
				{Name: "Ander Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Mayj Boater", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Sikingg", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Negaum ardera defender", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Zaka", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Gustavo Freefrags", Player: true, Traded: false, Summon: ""},
			},
			Level:  222,
			Reason: "Eliminated at Level 222 by Jupa Infinity, Jahziel Hardcori, Monarca Mimi, Marchane kee (traded), Caezors, Contatinho Millonario, Jungle Rubi Dominante, Monochrome Edowez, Negaum ardera defender, Juanjo Infinity, Retro Jupa, Rondero Zaka (traded), Niko Tin, Rvzor, Spyt Ponzi, Papi Hydrowar, Contatinho Ekbomba and Gustavo Freefrags. Assisted by Eszex Wallstreet, Ah Pepapipa, Samorbum, Acid Zero, Heethard, Suprldo, Lilpreben, Valkh Golpista, Caladan Bane, Ander Wallstreet, Mayj Boater, Netozawer, Sikingg and Vithrann.",
			Time:   "2021-12-15T21:30:48Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Ardeacion", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Keimithx", Player: true, Traded: false, Summon: ""},
				{Name: "Juantootwee", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
			},
			Level:  223,
			Reason: "Eliminated at Level 223 by Kaelsia Menardord, Nevin kyle, Infinitywar Sange, Calvo Bebe, Marchane kee (traded), Gorito Fullwar, Tiniebla Oddaeri, Jungle Rubi Dominante, Ardeacion, Elchico Billonario, Keimithx, Juantootwee (traded), Niko Tin, Rich Contatinho, Spyt Ponzi, Ronald Ardera (traded), Locura Boss and Fran Matarindo.",
			Time:   "2021-12-15T21:15:55Z",
		},
		{
			Assists: []Killers{
				{Name: "Richizawer", Player: true, Traded: false, Summon: ""},
				{Name: "Daark Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Middle Zocarno", Player: true, Traded: true, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Quali", Player: true, Traded: false, Summon: ""},
				{Name: "Meiker de Tozir", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Fuubaz Ltda", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Keimithx", Player: true, Traded: false, Summon: ""},
				{Name: "Juantootwee", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Luka Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
			},
			Level:  224,
			Reason: "Eliminated at Level 224 by Nevin kyle, Fuubaz Ltda, Infinitywar Sange, Caezors, Gorito Fullwar, Jungle Rubi Dominante, Monochrome Edowez, Juanjo Infinity, Jupa Traicionado, Keimithx, Juantootwee (traded), Niko Tin, Rich Contatinho, Ronald Ardera (traded), Locura Boss, Luka Is Back, Papi Hydrowar, Daddy Chiill and Fran Matarindo. Assisted by Richizawer, Daark Remix, Drunkz Wallstreet (traded), Middle Zocarno (traded), Ell Rugalzawer, Netozawer, Librarian Quali and Meiker de Tozir.",
			Time:   "2021-12-15T21:13:37Z",
		},
		{
			Assists: []Killers{
				{Name: "Schalama Rei Delas", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Negaum ardera defender", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Zaka", Player: true, Traded: true, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "adult goanna", Player: false, Traded: false, Summon: ""},
			},
			Level:  225,
			Reason: "Slain at Level 225 by Nevin kyle, Jungle Rubi Dominante, Negaum ardera defender, Chapo Kyle, Rondero Zaka (traded), Spyt Ponzi, Contatinho Ekbomba, Rondero Miguel and an adult goanna. Assisted by Schalama Rei Delas, Samorbum and Ell Rugalzawer.",
			Time:   "2021-12-15T21:00:57Z",
		},
		{
			Assists: []Killers{
				{Name: "Anthony No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Nytrander", Player: true, Traded: false, Summon: ""},
				{Name: "Mma Axel Mendoza", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Calvo Bebe", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Negaum ardera defender", Player: true, Traded: false, Summon: ""},
				{Name: "Juanjo Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Millonario Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Juantootwee", Player: true, Traded: true, Summon: ""},
				{Name: "Rondero Zaka", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Locura Boss", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Fran Matarindo", Player: true, Traded: false, Summon: ""},
			},
			Level:  226,
			Reason: "Annihilated at Level 226 by Jahziel Hardcori, Nevin kyle, Astaroth Kyle, Infinitywar Sange, Calvo Bebe, Gorito Fullwar, Nyl Vy, Jungle Rubi Dominante, Monochrome Edowez, Negaum ardera defender, Juanjo Infinity, Tacos Ardera, Jupa Traicionado, Chapo Kyle, Millonario Contatinho, Juantootwee (traded), Rondero Zaka (traded), Niko Tin, Rvzor, Rich Contatinho, Spyt Ponzi, Ronald Ardera (traded), Locura Boss, Contatinho Ekbomba and Fran Matarindo. Assisted by Anthony No Hands, Nytrander, Mma Axel Mendoza, Samorbum, Drunkz Wallstreet (traded), Acid Zero and Vithrann.",
			Time:   "2021-12-15T20:56:27Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Jungle Rubi Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
			},
			Level:  227,
			Reason: "Crushed at Level 227 by Astaroth Kyle, Infinitywar Sange, Gorito Fullwar, Tiniebla Oddaeri, Jungle Rubi Dominante, Monochrome Edowez, Rvzor, Rich Contatinho, Ronald Ardera (traded) and Papi Hydrowar.",
			Time:   "2021-12-15T20:47:15Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "skeleton elite warrior", Player: false, Traded: false, Summon: ""},
			},
			Level:  228,
			Reason: "Killed at Level 228 by Kaelsia Menardord and a skeleton elite warrior.",
			Time:   "2021-12-15T20:45:50Z",
		},
		{
			Assists: []Killers{
				{Name: "Mma Axel Mendoza", Player: true, Traded: false, Summon: ""},
				{Name: "Caladan Bane", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "priestess of the wild sun", Player: false, Traded: false, Summon: ""},
			},
			Level:  229,
			Reason: "Died at Level 229 by a priestess of the wild sun. Assisted by Mma Axel Mendoza, Caladan Bane and Netozawer.",
			Time:   "2021-12-15T20:40:47Z",
		},
		{
			Assists: []Killers{
				{Name: "May Thirtieth", Player: true, Traded: false, Summon: ""},
				{Name: "Eszex Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Nytrander", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Reborn Asian", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Chelito", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Roughcut", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Vithransky", Player: true, Traded: false, Summon: ""},
			},
			Level:  229,
			Reason: "Crushed at Level 229 by Jahziel Hardcori, Dont Kill Chelito, Von Rokitansky, Infinitywar Sange, Bravefly Legend, Tiniebla Oddaeri, Elchico Billonario, Retro Jupa, Rondero Tomachon, Roughcut, Ronald Ardera (traded), Contatinho Ekbomba and Vithransky. Assisted by May Thirtieth, Eszex Wallstreet, Nytrander, Samorbum, Guichin Killzejk Boom, Reborn Asian, Uanzawer and Adam No Hands (traded).",
			Time:   "2021-12-15T03:54:08Z",
		},
		{
			Assists: []Killers{
				{Name: "Eszex Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Jack Kevorkian", Player: true, Traded: false, Summon: ""},
				{Name: "Samorbum", Player: true, Traded: false, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Chuy Shelby", Player: true, Traded: false, Summon: ""},
				{Name: "Reborn Asian", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Vithrann", Player: true, Traded: false, Summon: ""},
				{Name: "Vack Vack Vack", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Gryffindor", Player: true, Traded: false, Summon: ""},
				{Name: "Wizard of Awz", Player: true, Traded: true, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Tacos Ardera", Player: true, Traded: false, Summon: ""},
				{Name: "Kyudok", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Madmatheuz", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Gura Druid", Player: true, Traded: false, Summon: ""},
			},
			Level:  230,
			Reason: "Crushed at Level 230 by Jupa Infinity, Von Rokitansky, Chino Kyle, Javo Gryffindor, Wizard of Awz (traded), Retro Jupa, Tacos Ardera, Kyudok, Rondero Tomachon, Madmatheuz (traded), Niko Tin and Gura Druid. Assisted by Eszex Wallstreet, Jack Kevorkian, Samorbum, Guichin Killzejk Boom, Acid Zero, Mapius Akuno, Notbrad, Yuniozawer, Chuy Shelby, Reborn Asian, Nick Pepperoni, Vithrann, Vack Vack Vack, Uanzawer and Adam No Hands (traded).",
			Time:   "2021-12-15T03:47:35Z",
		},
		{
			Assists: []Killers{
				{Name: "Dale Delas", Player: true, Traded: false, Summon: ""},
				{Name: "Kaizer Lobina", Player: true, Traded: false, Summon: ""},
				{Name: "Shmurdad", Player: true, Traded: false, Summon: ""},
				{Name: "Heethard", Player: true, Traded: false, Summon: ""},
				{Name: "Magic Dasherzi", Player: true, Traded: false, Summon: ""},
				{Name: "Hope Dysaster", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Infinitywar Sange", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Slobansky", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Bravefly Legend", Player: true, Traded: false, Summon: ""},
				{Name: "Don Dhimi Hardmode", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Dominante Denovo", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Zaka", Player: true, Traded: true, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
			},
			Level:  226,
			Reason: "Annihilated at Level 226 by Jupa Infinity, Jahziel Hardcori, Astaroth Kyle, Von Rokitansky, Infinitywar Sange, Marchane kee (traded), Slobansky, Mnnuuuuczzhh, Bravefly Legend, Don Dhimi Hardmode, Shady Is Back, Contatinho Millonario, Rein is Here, Nyl Vy, King Peruvian, Elchico Billonario, Contatinho Dominante Denovo, Nevarez Kyle, Jupa Traicionado, Rondero Zaka (traded), King Rexiiruz, Retro Demy and Daddy Chiill. Assisted by Dale Delas, Kaizer Lobina, Shmurdad, Heethard, Magic Dasherzi and Hope Dysaster.",
			Time:   "2021-12-14T03:21:55Z",
		},
		{
			Assists: []Killers{
				{Name: "Bless Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Schalama Rei Delas", Player: true, Traded: false, Summon: ""},
				{Name: "Kaizer Lobina", Player: true, Traded: false, Summon: ""},
				{Name: "Ah Pepapipa", Player: true, Traded: false, Summon: ""},
				{Name: "Kana Vys", Player: true, Traded: false, Summon: ""},
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Guichin Killzejk Boom", Player: true, Traded: false, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
				{Name: "Heethard", Player: true, Traded: false, Summon: ""},
				{Name: "Dexor Laesencia", Player: true, Traded: false, Summon: ""},
				{Name: "Kedruzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Diego Valdez", Player: true, Traded: false, Summon: ""},
				{Name: "Kanj iro", Player: true, Traded: false, Summon: ""},
				{Name: "Yuniozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Najlepzi Przyijeciouly", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Jotaeme Fullwaste", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Uanzawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jahziel Hardcori", Player: true, Traded: false, Summon: ""},
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Riley No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Chino Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Dont Kill Zambrano", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Rein is Here", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Mega", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuuucczhh", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Nevarez Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Zaka", Player: true, Traded: true, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rvzor", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "King Rexiiruz", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Miguel", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
			},
			Level:  227,
			Reason: "Annihilated at Level 227 by Jahziel Hardcori, Vrzik, Zeus Kyle, Astaroth Kyle, Riley No Hands, Chino Kyle, Marchane kee (traded), Dont Kill Zambrano, Mnnuuuuczzhh, Shady Is Back, Contatinho Millonario, Demy Mythwar Infinity, Utanii Herh, Rein is Here, Nyl Vy, Tiniebla Oddaeri, Retro Mega, Mnuuucczhh, Psycovzky, Retro Jupa, Nevarez Kyle, Jupa Traicionado, Rondero Zaka (traded), Niko Tin, Rvzor, Ronald Ardera (traded), King Rexiiruz, Mnuczhhx, Retro Demy, Contatinho Ekbomba, Rondero Miguel, Daddy Chiill and Sam Kyle. Assisted by Bless Wallstreet, Schalama Rei Delas, Kaizer Lobina, Ah Pepapipa, Kana Vys, Drunkz Wallstreet (traded), Guichin Killzejk Boom, Acid Zero, Heethard, Dexor Laesencia, Kedruzawer, Diego Valdez, Kanj iro, Yuniozawer, Najlepzi Przyijeciouly, Ell Rugalzawer, Jotaeme Fullwaste, Nick Pepperoni, Netozawer, Hardboss Remix and Uanzawer.",
			Time:   "2021-12-14T02:04:56Z",
		},
		{
			Assists: []Killers{
				{Name: "Richizawer", Player: true, Traded: false, Summon: ""},
				{Name: "Notbrad", Player: true, Traded: false, Summon: ""},
				{Name: "Kedruzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Reverse Terry", Player: true, Traded: false, Summon: ""},
				{Name: "Suprldo", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Nick Pepperoni", Player: true, Traded: false, Summon: ""},
				{Name: "Peninsula Boi", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Hueviin", Player: true, Traded: false, Summon: ""},
				{Name: "Sydekz", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Vrzik", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Von Rokitansky", Player: true, Traded: false, Summon: ""},
				{Name: "Javo Gryffindor", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Utanii Herh", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Luka Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Retro", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Ekbomba", Player: true, Traded: false, Summon: ""},
				{Name: "Daddy Chiill", Player: true, Traded: false, Summon: ""},
			},
			Level:  228,
			Reason: "Annihilated at Level 228 by Vrzik, Zeus Kyle, Von Rokitansky, Javo Gryffindor, Contatinho Millonario, Demy Mythwar Infinity, Utanii Herh, Nyl Vy, King Peruvian, Arthur Heartless, Psycovzky, Niko Tin, Woj Pokashield, Ronald Ardera (traded), Luka Is Back, Pit Retro, Retro Demy, Robadob Wallstreet, Contatinho Ekbomba and Daddy Chiill. Assisted by Richizawer, Notbrad, Kedruzawer, Reverse Terry, Suprldo, Ell Rugalzawer, Nick Pepperoni, Peninsula Boi, Hardboss Remix, Hueviin and Sydekz.",
			Time:   "2021-12-14T01:50:10Z",
		},
		{
			Assists: []Killers{
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Mal Victus", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Bito", Player: true, Traded: false, Summon: ""},
				{Name: "Rek Bazilha", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Hardboss Remix", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Quali", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Okiba Kay", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Wizard of Awz", Player: true, Traded: true, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
			},
			Level:  227,
			Reason: "Annihilated at Level 227 by Jupa Infinity, Okiba Kay, Monarca Mimi, Zeus Kyle, Astaroth Kyle, Caezors, Wizard of Awz (traded), Contatinho Millonario, Demy Mythwar Infinity, Viictoryck, Monochrome Edowez, Elchico Billonario, Psycovzky, Retro Jupa, Jupa Traicionado, Chapo Kyle, Raven Kyle, Niko Tin, Rich Contatinho, Spyt Ponzi, Ronald Ardera (traded), Retro Demy, Robadob Wallstreet and Papi Hydrowar. Assisted by Mapius Akuno, Mal Victus, Librarian Bito, Rek Bazilha, Netozawer, Hardboss Remix and Librarian Quali.",
			Time:   "2021-12-13T23:27:38Z",
		},
		{
			Assists: []Killers{
				{Name: "Drunkz Wallstreet", Player: true, Traded: true, Summon: ""},
				{Name: "Acid Zero", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Monarca Mimi", Player: true, Traded: false, Summon: ""},
				{Name: "Zeus Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Fuubaz Ltda", Player: true, Traded: false, Summon: ""},
				{Name: "Caezors", Player: true, Traded: false, Summon: ""},
				{Name: "Dont Kill Zambrano", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Izrehsad Cigam", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Ronald Ardera", Player: true, Traded: true, Summon: ""},
				{Name: "Pit Retro", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Papi Hydrowar", Player: true, Traded: false, Summon: ""},
			},
			Level:  228,
			Reason: "Annihilated at Level 228 by Jupa Infinity, Monarca Mimi, Zeus Kyle, Astaroth Kyle, Fuubaz Ltda, Caezors, Dont Kill Zambrano, Contatinho Millonario, Demy Mythwar Infinity, Viictoryck, King Peruvian, Elchico Billonario, Psycovzky, Retro Jupa, Jupa Traicionado, Chapo Kyle, Niko Tin, Izrehsad Cigam, Rich Contatinho, Spyt Ponzi, Ronald Ardera (traded), Pit Retro, Retro Demy, Robadob Wallstreet and Papi Hydrowar. Assisted by Drunkz Wallstreet (traded) and Acid Zero.",
			Time:   "2021-12-13T23:26:45Z",
		},
		{
			Assists: []Killers{
				{Name: "Mal Victus", Player: true, Traded: false, Summon: ""},
				{Name: "Leora Em", Player: true, Traded: false, Summon: ""},
				{Name: "Librarian Bito", Player: true, Traded: false, Summon: ""},
				{Name: "Qualitie", Player: true, Traded: false, Summon: ""},
				{Name: "Ell Rugalzawer", Player: true, Traded: false, Summon: ""},
				{Name: "Ander Wallstreet", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Jupa Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Demy Mythwar Infinity", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Licaxn Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Negaum ardera defender", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Street Runner", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Izrehsad Cigam", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
			},
			Level:  229,
			Reason: "Annihilated at Level 229 by Jupa Infinity, Nevin kyle, Astaroth Kyle, Gorito Fullwar, Contatinho Millonario, Demy Mythwar Infinity, Viictoryck, Licaxn Kyle, Nyl Vy, Tiniebla Oddaeri, King Peruvian, Monochrome Edowez, Negaum ardera defender, Jupa Traicionado, Rondero Tomachon, Raven Kyle, Street Runner, Niko Tin, Izrehsad Cigam, Rich Contatinho, Retro Demy and Sam Kyle. Assisted by Mal Victus, Leora Em, Librarian Bito, Qualitie, Ell Rugalzawer, Ander Wallstreet and Netozawer.",
			Time:   "2021-12-13T22:34:45Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Kaelsia Menardord", Player: true, Traded: false, Summon: ""},
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Fuubaz Ltda", Player: true, Traded: false, Summon: ""},
				{Name: "Gorito Fullwar", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Nyl Vy", Player: true, Traded: false, Summon: ""},
				{Name: "Tiniebla Oddaeri", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Jupa Traicionado", Player: true, Traded: false, Summon: ""},
				{Name: "Rondero Tomachon", Player: true, Traded: false, Summon: ""},
				{Name: "Niko Tin", Player: true, Traded: false, Summon: ""},
				{Name: "Woj Pokashield", Player: true, Traded: false, Summon: ""},
				{Name: "Rich Contatinho", Player: true, Traded: false, Summon: ""},
				{Name: "Spyt Ponzi", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Demy", Player: true, Traded: false, Summon: ""},
				{Name: "Sam Kyle", Player: true, Traded: false, Summon: ""},
			},
			Level:  230,
			Reason: "Eliminated at Level 230 by Kaelsia Menardord, Nevin kyle, Astaroth Kyle, Fuubaz Ltda, Gorito Fullwar, Contatinho Millonario, Nyl Vy, Tiniebla Oddaeri, King Peruvian, Monochrome Edowez, Psycovzky, Jupa Traicionado, Rondero Tomachon, Niko Tin, Woj Pokashield, Rich Contatinho, Spyt Ponzi, Retro Demy and Sam Kyle.",
			Time:   "2021-12-13T22:30:43Z",
		},
		{
			Assists: []Killers{
				{Name: "Anthony No Hands", Player: true, Traded: false, Summon: ""},
				{Name: "Mapius Akuno", Player: true, Traded: false, Summon: ""},
				{Name: "Don Ballusse", Player: true, Traded: false, Summon: ""},
				{Name: "Netozawer", Player: true, Traded: false, Summon: ""},
				{Name: "Adam No Hands", Player: true, Traded: true, Summon: "paladin familiar"},
			},
			Killers: []Killers{
				{Name: "Nevin kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Astaroth Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Marchane kee", Player: true, Traded: true, Summon: ""},
				{Name: "Wizard of Awz", Player: true, Traded: true, Summon: ""},
				{Name: "Dont Kill Joustick", Player: true, Traded: false, Summon: ""},
				{Name: "Mnnuuuuczzhh", Player: true, Traded: false, Summon: ""},
				{Name: "Shady Is Back", Player: true, Traded: false, Summon: ""},
				{Name: "Contatinho Millonario", Player: true, Traded: false, Summon: ""},
				{Name: "Be Nicee", Player: true, Traded: false, Summon: ""},
				{Name: "Viictoryck", Player: true, Traded: false, Summon: ""},
				{Name: "Hiisa", Player: true, Traded: false, Summon: ""},
				{Name: "King Peruvian", Player: true, Traded: false, Summon: ""},
				{Name: "Touchy Dominante", Player: true, Traded: false, Summon: ""},
				{Name: "Monochrome Edowez", Player: true, Traded: false, Summon: ""},
				{Name: "Elchico Billonario", Player: true, Traded: false, Summon: ""},
				{Name: "Arthur Heartless", Player: true, Traded: false, Summon: ""},
				{Name: "Psycovzky", Player: true, Traded: false, Summon: ""},
				{Name: "Retro Jupa", Player: true, Traded: false, Summon: ""},
				{Name: "Chapo Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Beowulf Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Keimithx", Player: true, Traded: false, Summon: ""},
				{Name: "Bello Abreu", Player: true, Traded: false, Summon: ""},
				{Name: "Raven Kyle", Player: true, Traded: false, Summon: ""},
				{Name: "Mnuczhhx", Player: true, Traded: false, Summon: ""},
				{Name: "Pit Retro", Player: true, Traded: false, Summon: ""},
				{Name: "Robadob Wallstreet", Player: true, Traded: false, Summon: ""},
			},
			Level:  231,
			Reason: "Annihilated at Level 231 by Nevin kyle, Astaroth Kyle, Marchane kee (traded), Wizard of Awz (traded), Dont Kill Joustick, Mnnuuuuczzhh, Shady Is Back, Contatinho Millonario, Be Nicee, Viictoryck, Hiisa, King Peruvian, Touchy Dominante, Monochrome Edowez, Elchico Billonario, Arthur Heartless, Psycovzky, Retro Jupa, Chapo Kyle, Beowulf Kyle, Keimithx, Bello Abreu, Raven Kyle, Mnuczhhx, Pit Retro and Robadob Wallstreet. Assisted by Anthony No Hands, Mapius Akuno, Don Ballusse, Netozawer and a paladin familiar of Adam No Hands (traded).",
			Time:   "2021-12-11T02:08:02Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			deaths[idx].Time, tc.Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber5(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Rejana on Fera.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Rejana on Fera", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("female", character.Sex)
	assert.Equal("None", character.Title)
	assert.Equal(12, character.UnlockedTitles)
	assert.Equal("Knight", character.Vocation)
	assert.Equal(8, character.Level)
	assert.Equal(0, character.AchievementPoints)
	assert.Equal("Fera", character.World)
	assert.Nil(character.FormerWorlds)
	assert.Equal("Isle of Solitude", character.Residence)
	assert.Empty(character.MarriedTo)
	assert.Equal("2021-10-25T04:37:46Z", character.LastLogin)
	assert.Equal("CipSoft Member", character.Position)
	assert.Equal("Premium Account", character.AccountStatus)

	assert.Empty(characterJson.Character.Achievements)
	assert.Empty(characterJson.Character.AccountInformation.LoyaltyTitle)

	// validate other characters
	assert.Equal(4, len(characterJson.Character.OtherCharacters))

	positionCharacter := characterJson.Character.OtherCharacters[1]
	assert.Equal("Rejana on Fera", positionCharacter.Name)
	assert.Equal("Fera", positionCharacter.World)
	assert.Equal("offline", positionCharacter.Status)
	assert.False(positionCharacter.Deleted)
	assert.False(positionCharacter.Main)
	assert.False(positionCharacter.Traded)
	assert.Equal("CipSoft Member", positionCharacter.Position)
}

func TestNumber6(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Luminals.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Luminals", character.Name)

	assert.Equal(8, len(characterJson.Character.AccountBadges))
	globalPlayerBadge := characterJson.Character.AccountBadges[4]
	assert.Equal("Global Player (Grade 3)", globalPlayerBadge.Name)
	assert.Equal("https://static.tibia.com/images//badges/badge_globalplayer3.png", globalPlayerBadge.IconURL)
	assert.Equal("Summing up the levels of all characters on the account amounts to at least 2000.", globalPlayerBadge.Description)
	masterClassBadge := characterJson.Character.AccountBadges[7]
	assert.Equal("Master Class (Grade 1)", masterClassBadge.Name)
	assert.Equal("https://static.tibia.com/images//badges/badge_masterclass1.png", masterClassBadge.IconURL)
	assert.Equal("The account has reached at least level 100 with all four vocations.", masterClassBadge.Description)

	assert.Len(characterJson.Character.Achievements, 5)
	assert.Equal(characterJson.Character.Achievements[0].Name, "Alumni")
	assert.Equal(characterJson.Character.Achievements[0].Grade, 2)
	assert.Equal(characterJson.Character.Achievements[0].Secret, false)
	assert.Equal(characterJson.Character.Achievements[1].Name, "Forbidden Fruit")
	assert.Equal(characterJson.Character.Achievements[1].Grade, 1)
	assert.Equal(characterJson.Character.Achievements[1].Secret, true)
	assert.Equal(characterJson.Character.Achievements[2].Name, "Goldhunter")
	assert.Equal(characterJson.Character.Achievements[2].Grade, 1)
	assert.Equal(characterJson.Character.Achievements[2].Secret, true)
	assert.Equal(characterJson.Character.Achievements[3].Name, "Pyromaniac")
	assert.Equal(characterJson.Character.Achievements[3].Grade, 2)
	assert.Equal(characterJson.Character.Achievements[3].Secret, true)
	assert.Equal(characterJson.Character.Achievements[4].Name, "Razing!")
	assert.Equal(characterJson.Character.Achievements[4].Grade, 3)
	assert.Equal(characterJson.Character.Achievements[4].Secret, true)
}

func TestNumber7(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Torbjörn.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Len(characterJson.Character.Achievements, 0)
	assert.Equal("Torbjörn", character.Name)
	assert.Equal("___$$$$$$$$_______$$$$$$$$\n_$$$$$$$$$$$$__$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n_$$$$$$$$$$-Snulliz-$$$$$$$$$$$\n__$$$$$$$$$$$$$$$$$$$$$$$$$$\n____$$$$$$$$$$$$$$$$$$$$$$\n______$$$$$$$$$$$$$$$$$$\n________$$$$$$$$$$$$$$\n___________$$$$$$$$$\n____________$$$$$$\n_____________$$", character.Comment)
}

func TestNumber8(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Mieluffy.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Len(characterJson.Character.Achievements, 0)
	assert.Equal("Mieluffy", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(4, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{
				{Name: "Merlinxd", Player: true, Traded: false, Summon: ""},
				{Name: "Paletero Kriminal", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Pecas Creator", Player: true, Traded: false, Summon: ""},
				{Name: "El Unico", Player: true, Traded: false, Summon: ""},
				{Name: "Supperiore", Player: true, Traded: false, Summon: ""},
				{Name: "Flaco El Pobre", Player: true, Traded: false, Summon: ""},
				{Name: "Rii Rox", Player: true, Traded: false, Summon: ""},
				{Name: "Tiz to", Player: true, Traded: false, Summon: ""},
				{Name: "True Merlinus Druid", Player: true, Traded: false, Summon: ""},
				{Name: "Antii Druida", Player: true, Traded: false, Summon: ""},
				{Name: "El Inestable", Player: true, Traded: false, Summon: ""},
				{Name: "Ga to Relaxsz", Player: true, Traded: false, Summon: ""},
				{Name: "Frodin la Maquina", Player: true, Traded: false, Summon: ""},
				{Name: "Mich Jogadorcaro", Player: true, Traded: false, Summon: ""},
			},
			Level:  508,
			Reason: "Crushed at Level 508 by Pecas Creator, El Unico, Supperiore, Flaco El Pobre, Rii Rox, Tiz to, True Merlinus Druid, Antii Druida, El Inestable, Ga to Relaxsz, Frodin la Maquina and Mich Jogadorcaro. Assisted by Merlinxd and Paletero Kriminal.",
			Time:   "2024-02-01T06:20:46Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "El Unico", Player: true, Traded: false, Summon: ""},
				{Name: "Naireth Sorcerer", Player: true, Traded: false, Summon: ""},
				{Name: "Pedritox Soulfire", Player: true, Traded: false, Summon: ""},
				{Name: "Pepiisho", Player: true, Traded: false, Summon: ""},
				{Name: "Egon Leme", Player: true, Traded: false, Summon: ""},
				{Name: "Queen Marii", Player: true, Traded: false, Summon: ""},
				{Name: "Jicuri Guardian", Player: true, Traded: false, Summon: ""},
				{Name: "Athob Senian", Player: true, Traded: false, Summon: ""},
				{Name: "Noo Friends", Player: true, Traded: false, Summon: ""},
				{Name: "Dhanielz Unstoppable", Player: true, Traded: false, Summon: ""},
				{Name: "Ilysz Sin Tales", Player: true, Traded: false, Summon: ""},
				{Name: "Spy Crusitho Sauvage", Player: true, Traded: false, Summon: ""},
				{Name: "Tiz to", Player: true, Traded: false, Summon: ""},
				{Name: "King Asmiito", Player: true, Traded: false, Summon: ""},
				{Name: "Bolchecoqe", Player: true, Traded: false, Summon: ""},
				{Name: "Jobi", Player: true, Traded: false, Summon: ""},
				{Name: "El Inestable", Player: true, Traded: false, Summon: ""},
				{Name: "Bloomzs", Player: true, Traded: false, Summon: ""},
				{Name: "Natalie Bearskin", Player: true, Traded: false, Summon: ""},
			},
			Level:  508,
			Reason: "Eliminated at Level 508 by El Unico, Naireth Sorcerer, Pedritox Soulfire, Pepiisho, Egon Leme, Queen Marii, Jicuri Guardian, Athob Senian, Noo Friends, Dhanielz Unstoppable, Ilysz Sin Tales, Spy Crusitho Sauvage, Tiz to, King Asmiito, Bolchecoqe, Jobi, El Inestable, Bloomzs and Natalie Bearskin.",
			Time:   "2024-02-01T04:53:23Z",
		},
		{
			Assists: []Killers{
				{Name: "Swiifti", Player: true, Traded: false, Summon: ""},
				{Name: "Duende blanco", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Shantty", Player: true, Traded: false, Summon: ""},
				{Name: "Sin Primalbazzar", Player: true, Traded: false, Summon: ""},
				{Name: "Supperiore", Player: true, Traded: false, Summon: ""},
				{Name: "Righo", Player: true, Traded: false, Summon: ""},
				{Name: "Beltran Tzawayak", Player: true, Traded: false, Summon: ""},
				{Name: "Ilysz Sin Tales", Player: true, Traded: false, Summon: ""},
				{Name: "Flaco El Pobre", Player: true, Traded: false, Summon: ""},
				{Name: "Pecaas", Player: true, Traded: false, Summon: ""},
				{Name: "Knightsitaz", Player: true, Traded: false, Summon: ""},
				{Name: "Dhanielz Acorazado", Player: true, Traded: false, Summon: ""},
				{Name: "Jobi", Player: true, Traded: false, Summon: ""},
				{Name: "Aeronabic", Player: true, Traded: false, Summon: ""},
				{Name: "Side Effectss", Player: true, Traded: false, Summon: ""},
				{Name: "Ekizdd", Player: true, Traded: false, Summon: ""},
				{Name: "Baby Mikoh", Player: true, Traded: false, Summon: ""},
				{Name: "Sneki", Player: true, Traded: false, Summon: ""},
				{Name: "Love and Death", Player: true, Traded: true, Summon: ""},
				{Name: "Next Generation", Player: true, Traded: false, Summon: ""},
			},
			Level:  508,
			Reason: "Eliminated at Level 508 by Shantty, Sin Primalbazzar, Supperiore, Righo, Beltran Tzawayak, Ilysz Sin Tales, Flaco El Pobre, Pecaas, Knightsitaz, Dhanielz Acorazado, Jobi, Aeronabic, Side Effectss, Ekizdd, Baby Mikoh, Sneki, Love and Death (traded) and Next Generation. Assisted by Swiifti and Duende blanco.",
			Time:   "2024-01-12T19:38:09Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Sir Alpha", Player: true, Traded: false, Summon: ""},
				{Name: "Chic Unixzion", Player: true, Traded: false, Summon: ""},
				{Name: "Cadiwax", Player: true, Traded: false, Summon: ""},
				{Name: "Pecaas", Player: true, Traded: false, Summon: ""},
				{Name: "Eveo Zeff", Player: true, Traded: false, Summon: ""},
				{Name: "Monkey Flash", Player: true, Traded: false, Summon: ""},
				{Name: "Flaco Maldito", Player: true, Traded: false, Summon: ""},
				{Name: "Mieluffy", Player: true, Traded: false, Summon: ""},
				{Name: "Geass Evangelion", Player: true, Traded: false, Summon: ""},
			},
			Level:  508,
			Reason: "Slain at Level 508 by Sir Alpha, Chic Unixzion, Cadiwax, Pecaas, Eveo Zeff, Monkey Flash, Flaco Maldito, Mieluffy and Geass Evangelion.",
			Time:   "2024-01-10T21:31:52Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			deaths[idx].Time, tc.Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber9(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Akura Aleus.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character

	assert.Len(character.OtherCharacters, 10)
	assert.Equal(character.OtherCharacters[0].Deleted, false)
	assert.Equal(character.OtherCharacters[0].Main, false)
	assert.Equal(character.OtherCharacters[0].Name, "Akura Aleus")
	assert.Equal(character.OtherCharacters[0].Position, "")
	assert.Equal(character.OtherCharacters[0].Status, "offline")
	assert.Equal(character.OtherCharacters[0].Traded, false)
	assert.Equal(character.OtherCharacters[0].World, "Lobera")

	assert.Equal(character.OtherCharacters[1].Deleted, false)
	assert.Equal(character.OtherCharacters[1].Main, false)
	assert.Equal(character.OtherCharacters[1].Name, "Armnox")
	assert.Equal(character.OtherCharacters[1].Position, "")
	assert.Equal(character.OtherCharacters[1].Status, "offline")
	assert.Equal(character.OtherCharacters[1].Traded, false)
	assert.Equal(character.OtherCharacters[1].World, "Ferobra")

	assert.Equal(character.OtherCharacters[2].Deleted, false)
	assert.Equal(character.OtherCharacters[2].Main, false)
	assert.Equal(character.OtherCharacters[2].Name, "Cheradon")
	assert.Equal(character.OtherCharacters[2].Position, "")
	assert.Equal(character.OtherCharacters[2].Status, "offline")
	assert.Equal(character.OtherCharacters[2].Traded, false)
	assert.Equal(character.OtherCharacters[2].World, "Serdebra")

	assert.Equal(character.OtherCharacters[3].Deleted, false)
	assert.Equal(character.OtherCharacters[3].Main, false)
	assert.Equal(character.OtherCharacters[3].Name, "Dollar Driver")
	assert.Equal(character.OtherCharacters[3].Position, "")
	assert.Equal(character.OtherCharacters[3].Status, "offline")
	assert.Equal(character.OtherCharacters[3].Traded, false)
	assert.Equal(character.OtherCharacters[3].World, "Ousabra")

	assert.Equal(character.OtherCharacters[4].Deleted, false)
	assert.Equal(character.OtherCharacters[4].Main, false)
	assert.Equal(character.OtherCharacters[4].Name, "Goth angel sinner")
	assert.Equal(character.OtherCharacters[4].Position, "")
	assert.Equal(character.OtherCharacters[4].Status, "offline")
	assert.Equal(character.OtherCharacters[4].Traded, true)
	assert.Equal(character.OtherCharacters[4].World, "Ousabra")

	assert.Equal(character.OtherCharacters[5].Deleted, false)
	assert.Equal(character.OtherCharacters[5].Main, false)
	assert.Equal(character.OtherCharacters[5].Name, "Halodrol")
	assert.Equal(character.OtherCharacters[5].Position, "")
	assert.Equal(character.OtherCharacters[5].Status, "offline")
	assert.Equal(character.OtherCharacters[5].Traded, false)
	assert.Equal(character.OtherCharacters[5].World, "Vunira")

	assert.Equal(character.OtherCharacters[6].Deleted, false)
	assert.Equal(character.OtherCharacters[6].Main, false)
	assert.Equal(character.OtherCharacters[6].Name, "Halodrow")
	assert.Equal(character.OtherCharacters[6].Position, "")
	assert.Equal(character.OtherCharacters[6].Status, "offline")
	assert.Equal(character.OtherCharacters[6].Traded, false)
	assert.Equal(character.OtherCharacters[6].World, "Lobera")

	assert.Equal(character.OtherCharacters[7].Deleted, false)
	assert.Equal(character.OtherCharacters[7].Main, false)
	assert.Equal(character.OtherCharacters[7].Name, "Incoggnita")
	assert.Equal(character.OtherCharacters[7].Position, "")
	assert.Equal(character.OtherCharacters[7].Status, "offline")
	assert.Equal(character.OtherCharacters[7].Traded, false)
	assert.Equal(character.OtherCharacters[7].World, "Ferobra")

	assert.Equal(character.OtherCharacters[8].Deleted, false)
	assert.Equal(character.OtherCharacters[8].Main, false)
	assert.Equal(character.OtherCharacters[8].Name, "Lord Kabum")
	assert.Equal(character.OtherCharacters[8].Position, "")
	assert.Equal(character.OtherCharacters[8].Status, "offline")
	assert.Equal(character.OtherCharacters[8].Traded, false)
	assert.Equal(character.OtherCharacters[8].World, "Solidera")

	assert.Equal(character.OtherCharacters[9].Deleted, false)
	assert.Equal(character.OtherCharacters[9].Main, true)
	assert.Equal(character.OtherCharacters[9].Name, "Lord Succubu")
	assert.Equal(character.OtherCharacters[9].Position, "")
	assert.Equal(character.OtherCharacters[9].Status, "offline")
	assert.Equal(character.OtherCharacters[9].Traded, false)
	assert.Equal(character.OtherCharacters[9].World, "Ferobra")
}

func TestNumber10(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Nocna Furia.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Nocna Furia", character.Name)
	assert.Nil(character.FormerNames)
	assert.False(character.Traded)
	assert.Empty(character.DeletionDate)
	assert.Equal("male", character.Sex)
	assert.Equal("Tibia's Topmodel (Grade 1)", character.Title)
	assert.Equal(6, character.UnlockedTitles)
}

func TestNumber11(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Orca Kaoksh.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Orca Kaoksh", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(2, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "fire", Player: false, Traded: false, Summon: ""},
			},
			Level:  10,
			Reason: "Died at Level 10 by fire.",
			Time:   "2023-10-08T16:19:35Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "wasp", Player: false, Traded: false, Summon: ""},
			},
			Level:  8,
			Reason: "Died at Level 8 by wasp.",
			Time:   "2023-10-07T00:27:38Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			deaths[idx].Time, tc.Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber12(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Stalone Matador.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Stalone Matador", character.Name)
	assert.True(characterJson.Character.DeathsTruncated)
	assert.Equal(55, len(characterJson.Character.Deaths))
}

func TestNumber13(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Ninth Dimension.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Ninth Dimension", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(1, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{
				{Name: "Dark Assa", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{
				{Name: "Pess Joeru", Player: true, Traded: false, Summon: ""},
				{Name: "Curly Da Goonx", Player: true, Traded: false, Summon: ""},
				{Name: "Setarehh", Player: true, Traded: false, Summon: ""},
				{Name: "Skkrimz", Player: true, Traded: false, Summon: ""},
				{Name: "Luna Mors", Player: true, Traded: false, Summon: ""},
				{Name: "Micklo", Player: true, Traded: false, Summon: ""},
				{Name: "Kate Morningstar", Player: true, Traded: false, Summon: ""},
				{Name: "Avatar Avatar", Player: true, Traded: false, Summon: ""},
				{Name: "San Bernardino Hoodrat", Player: true, Traded: false, Summon: ""},
				{Name: "Mighty Nitro", Player: true, Traded: false, Summon: ""},
				{Name: "Aiakosz", Player: true, Traded: false, Summon: ""},
				{Name: "Sithaadoz", Player: true, Traded: false, Summon: ""},
				{Name: "Compa Ache", Player: true, Traded: false, Summon: ""},
				{Name: "Cave Stormer", Player: true, Traded: false, Summon: ""},
				{Name: "Doppler and Bankrupt", Player: true, Traded: false, Summon: ""},
			},
			Level:  544,
			Reason: "Eliminated at Level 544 by Pess Joeru, Curly Da Goonx, Setarehh, Skkrimz, Luna Mors, Micklo, Kate Morningstar, Avatar Avatar, San Bernardino Hoodrat, Mighty Nitro, Aiakosz, Sithaadoz, Compa Ache, Cave Stormer and Doppler and Bankrupt. Assisted by Dark Assa.",
			Time:   "2024-03-12T21:02:33Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber14(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Nannorka.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Nannorka", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(35, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-27T16:20:06Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Satashi Xuuu", Player: true, Traded: false, Summon: ""},
			},
			Level:  791,
			Reason: "Killed at Level 791 by Satashi Xuuu.",
			Time:   "2025-02-27T03:46:47Z",
		},
		{
			Assists: []Killers{
				{Name: "Basilicata", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Basilicata.",
			Time:    "2025-02-25T14:11:42Z",
		},
		{
			Assists: []Killers{
				{Name: "Tainerd Ruero", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Tainerd Ruero.",
			Time:    "2025-02-15T10:35:03Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-15T09:45:37Z",
		},
		{
			Assists: []Killers{
				{Name: "Basilicata", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Basilicata.",
			Time:    "2025-02-14T16:56:53Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Satashi Xuuu", Player: true, Traded: false, Summon: ""},
			},
			Level:  826,
			Reason: "Killed at Level 826 by Satashi Xuuu.",
			Time:   "2025-02-14T14:56:33Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-14T13:24:43Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-14T08:19:51Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-12T19:38:28Z",
		},
		{
			Assists: []Killers{
				{Name: "Tikozera To Calmo", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Tikozera To Calmo.",
			Time:    "2025-02-12T17:37:04Z",
		},
		{
			Assists: []Killers{
				{Name: "Basilicata", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Basilicata.",
			Time:    "2025-02-12T17:02:27Z",
		},
		{
			Assists: []Killers{
				{Name: "Du nken", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Du nken.",
			Time:    "2025-02-12T16:56:42Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Satashi Xuuu", Player: true, Traded: false, Summon: ""},
			},
			Level:  878,
			Reason: "Killed at Level 878 by Satashi Xuuu.",
			Time:   "2025-02-11T20:00:44Z",
		},
		{
			Assists: []Killers{
				{Name: "Basilicata", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Basilicata.",
			Time:    "2025-02-11T16:21:46Z",
		},
		{
			Assists: []Killers{
				{Name: "Whiskin", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Whiskin.",
			Time:    "2025-02-11T16:18:28Z",
		},
		{
			Assists: []Killers{
				{Name: "Fjunkes", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Fjunkes.",
			Time:    "2025-02-11T16:16:29Z",
		},
		{
			Assists: []Killers{
				{Name: "Fjunkes", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Fjunkes.",
			Time:    "2025-02-10T14:57:44Z",
		},
		{
			Assists: []Killers{
				{Name: "Whiskin", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Whiskin.",
			Time:    "2025-02-10T14:55:03Z",
		},
		{
			Assists: []Killers{
				{Name: "Kaos Mest", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Kaos Mest.",
			Time:    "2025-02-09T12:13:57Z",
		},
		{
			Assists: []Killers{
				{Name: "Whiskin", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Whiskin.",
			Time:    "2025-02-09T10:26:48Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-09T10:09:39Z",
		},
		{
			Assists: []Killers{
				{Name: "Guzik Szef", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Guzik Szef.",
			Time:    "2025-02-09T10:06:43Z",
		},
		{
			Assists: []Killers{
				{Name: "Whiskin", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Whiskin.",
			Time:    "2025-02-08T18:06:49Z",
		},
		{
			Assists: []Killers{
				{Name: "Shensz Khalifa", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Shensz Khalifa.",
			Time:    "2025-02-08T09:19:19Z",
		},
		{
			Assists: []Killers{
				{Name: "Taiimo", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Taiimo.",
			Time:    "2025-02-07T17:13:05Z",
		},
		{
			Assists: []Killers{
				{Name: "Shensz Khalifa", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Shensz Khalifa.",
			Time:    "2025-02-06T17:51:50Z",
		},
		{
			Assists: []Killers{
				{Name: "Miquudalajarab", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Miquudalajarab.",
			Time:    "2025-02-06T17:50:03Z",
		},
		{
			Assists: []Killers{
				{Name: "Taiimo", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Taiimo.",
			Time:    "2025-02-06T17:15:54Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Sephea", Player: true, Traded: false, Summon: ""},
			},
			Level:  1007,
			Reason: "Killed at Level 1007 by Sephea.",
			Time:   "2025-02-05T17:03:08Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Satashi Xuuu", Player: true, Traded: false, Summon: ""},
			},
			Level:  1016,
			Reason: "Killed at Level 1016 by Satashi Xuuu.",
			Time:   "2025-02-04T03:42:41Z",
		},
		{
			Assists: []Killers{
				{Name: "Rodmago Aesir", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Rodmago Aesir.",
			Time:    "2025-02-03T16:14:53Z",
		},
		{
			Assists: []Killers{
				{Name: "Zain Malek", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Zain Malek.",
			Time:    "2025-02-02T10:19:05Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Itzz Xed", Player: true, Traded: false, Summon: ""},
			},
			Level:  1042,
			Reason: "Killed at Level 1042 by Itzz Xed.",
			Time:   "2025-02-02T10:17:25Z",
		},
		{
			Assists: []Killers{
				{Name: "Rauxzin", Player: true, Traded: false, Summon: ""},
			},
			Killers: []Killers{},
			Reason:  "Assisted by Rauxzin.",
			Time:    "2025-02-02T10:11:07Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber15(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Sir Dumbas.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Sir Dumbas", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(3, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Corruption Toxic", Player: true, Traded: false, Summon: "fire elemental"},
			},
			Level:  43,
			Reason: "Killed at Level 43 by fire elemental of Corruption Toxic.",
			Time:   "2025-02-02T23:08:13Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "cave rat", Player: false, Traded: false, Summon: ""},
			},
			Level:  44,
			Reason: "Died at Level 44 by cave rat.",
			Time:   "2025-02-02T22:47:32Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "cave rat", Player: false, Traded: false, Summon: ""},
			},
			Level:  45,
			Reason: "Died at Level 45 by cave rat.",
			Time:   "2025-02-02T22:46:11Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber16(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Sergiozk.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Sergiozk", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(6, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Phelip On Danera", Player: true, Traded: false, Summon: ""},
				{Name: "fire elemental", Player: false, Traded: false, Summon: ""},
			},
			Level:  597,
			Reason: "Killed at Level 597 by Phelip On Danera and fire elemental.",
			Time:   "2024-03-17T04:47:32Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Phelip On Danera", Player: true, Traded: false, Summon: ""},
				{Name: "Indio Pedibrek", Player: true, Traded: false, Summon: ""},
				{Name: "Guillera", Player: true, Traded: false, Summon: ""},
				{Name: "Sago Dashzik", Player: true, Traded: false, Summon: ""},
				{Name: "Spectrum on Ferobra", Player: true, Traded: false, Summon: ""},
				{Name: "Hawker Destroyer", Player: true, Traded: false, Summon: ""},
				{Name: "Lenuzera Dushow", Player: true, Traded: false, Summon: ""},
				{Name: "Migurina", Player: true, Traded: false, Summon: ""},
				{Name: "Zofia Bosak", Player: true, Traded: false, Summon: ""},
				{Name: "Lekinho Destroyer", Player: true, Traded: false, Summon: ""},
				{Name: "Sergiozk", Player: true, Traded: false, Summon: ""},
				{Name: "Zeqs", Player: true, Traded: false, Summon: ""},
				{Name: "Durantezika", Player: true, Traded: false, Summon: ""},
				{Name: "Maxionz", Player: true, Traded: true, Summon: ""},
				{Name: "Axmc", Player: true, Traded: false, Summon: ""},
				{Name: "Luan Chriistiian", Player: true, Traded: true, Summon: ""},
				{Name: "Utv", Player: true, Traded: false, Summon: ""},
				{Name: "Duzym", Player: true, Traded: false, Summon: ""},
				{Name: "Wizz Gigi", Player: true, Traded: false, Summon: ""},
				{Name: "Kan Zinho", Player: true, Traded: false, Summon: ""},
				{Name: "Yuyo", Player: true, Traded: false, Summon: ""},
				{Name: "skeleton elite warrior", Player: false, Traded: false, Summon: ""},
			},
			Level:  597,
			Reason: "Annihilated at Level 597 by Phelip On Danera, Indio Pedibrek, Guillera, Sago Dashzik, Spectrum on Ferobra, Hawker Destroyer, Lenuzera Dushow, Migurina, Zofia Bosak, Lekinho Destroyer, Sergiozk, Zeqs, Durantezika, Maxionz (traded), Axmc, Luan Chriistiian (traded), Utv, Duzym, Wizz Gigi, Kan Zinho, Yuyo and skeleton elite warrior.",
			Time:   "2024-03-17T04:32:23Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Phelip On Danera", Player: true, Traded: false, Summon: ""},
				{Name: "Indio Pedibrek", Player: true, Traded: false, Summon: ""},
				{Name: "Guillera", Player: true, Traded: false, Summon: ""},
				{Name: "Sago Dashzik", Player: true, Traded: false, Summon: ""},
				{Name: "Hawker Destroyer", Player: true, Traded: false, Summon: ""},
				{Name: "Lenuzera Dushow", Player: true, Traded: false, Summon: ""},
				{Name: "Zofia Bosak", Player: true, Traded: false, Summon: ""},
				{Name: "Zeh da Mannga", Player: true, Traded: false, Summon: ""},
				{Name: "Maxionz", Player: true, Traded: true, Summon: ""},
				{Name: "Luan Chriistiian", Player: true, Traded: true, Summon: ""},
				{Name: "Utv", Player: true, Traded: false, Summon: ""},
				{Name: "Wizz Gigi", Player: true, Traded: false, Summon: ""},
			},
			Level:  597,
			Reason: "Crushed at Level 597 by Phelip On Danera, Indio Pedibrek, Guillera, Sago Dashzik, Hawker Destroyer, Lenuzera Dushow, Zofia Bosak, Zeh da Mannga, Maxionz (traded), Luan Chriistiian (traded), Utv and Wizz Gigi.",
			Time:   "2024-03-17T04:26:44Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Moon Warlock", Player: true, Traded: false, Summon: ""},
				{Name: "Manglu", Player: true, Traded: false, Summon: ""},
				{Name: "Dudodix", Player: true, Traded: false, Summon: ""},
				{Name: "Luchador", Player: true, Traded: false, Summon: ""},
				{Name: "King Madara", Player: true, Traded: false, Summon: ""},
				{Name: "Masstterpali", Player: true, Traded: false, Summon: ""},
				{Name: "Gilbas Concursado", Player: true, Traded: false, Summon: ""},
				{Name: "Cafezinhum", Player: true, Traded: false, Summon: ""},
				{Name: "Antipuush", Player: true, Traded: false, Summon: ""},
			},
			Level:  597,
			Reason: "Slain at Level 597 by Moon Warlock, Manglu, Dudodix, Luchador, King Madara, Masstterpali, Gilbas Concursado, Cafezinhum and Antipuush.",
			Time:   "2024-03-15T13:02:38Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "crazed summer rearguard", Player: false, Traded: false, Summon: ""},
			},
			Level:  598,
			Reason: "Died at Level 598 by crazed summer rearguard.",
			Time:   "2024-03-11T06:21:25Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Irynn", Player: true, Traded: false, Summon: ""},
				{Name: "Indio Pedibrek", Player: true, Traded: false, Summon: ""},
				{Name: "Tenoriio Rex", Player: true, Traded: false, Summon: ""},
				{Name: "Pitucola", Player: true, Traded: false, Summon: ""},
				{Name: "Lekinho Destroyer", Player: true, Traded: false, Summon: ""},
				{Name: "Sergiozk", Player: true, Traded: false, Summon: ""},
				{Name: "Aperninha de grilo", Player: true, Traded: false, Summon: ""},
				{Name: "Bino Pally Rox", Player: true, Traded: false, Summon: ""},
				{Name: "Experience Hunter", Player: true, Traded: false, Summon: ""},
				{Name: "Lohana Darkside", Player: true, Traded: false, Summon: ""},
				{Name: "Flavzk", Player: true, Traded: true, Summon: ""},
				{Name: "Fllockyy", Player: true, Traded: false, Summon: "paladin familiar"},
			},
			Level:  597,
			Reason: "Crushed at Level 597 by Irynn, Indio Pedibrek, Tenoriio Rex, Pitucola, Lekinho Destroyer, Sergiozk, Aperninha de grilo, Bino Pally Rox, Experience Hunter, Lohana Darkside, Flavzk (traded) and paladin familiar of Fllockyy.",
			Time:   "2024-03-10T14:49:27Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber17(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Akayev.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Akayev", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(1, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "Shooter Scotty", Player: true, Traded: false, Summon: ""},
				{Name: "Simpkan", Player: true, Traded: false, Summon: ""},
				{Name: "Spektrozz", Player: true, Traded: false, Summon: ""},
				{Name: "Toxic Raksoo", Player: true, Traded: false, Summon: ""},
				{Name: "Delerioo", Player: true, Traded: false, Summon: ""},
				{Name: "Gugu Lurifax", Player: true, Traded: false, Summon: ""},
				{Name: "Viniciao Novamente", Player: true, Traded: false, Summon: ""},
			},
			Level:  645,
			Reason: "Slain at Level 645 by Shooter Scotty, Simpkan, Spektrozz, Toxic Raksoo, Delerioo, Gugu Lurifax and Viniciao Novamente.",
			Time:   "2025-03-25T11:08:07Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func TestNumber18(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/characters/Tirador Azteca.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	characterJson, err := TibiaCharactersCharacterImpl(string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	character := characterJson.Character.CharacterInfo

	assert.Equal("Tirador Azteca", character.Name)
	assert.False(characterJson.Character.DeathsTruncated)

	// validate death data
	assert.Equal(3, len(characterJson.Character.Deaths))
	deaths := characterJson.Character.Deaths

	for idx, tc := range []struct {
		Assists []Killers
		Killers []Killers
		Level   int
		Reason  string
		Time    string
	}{
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "retainer of Baeloc", Player: false, Traded: false, Summon: ""},
			},
			Level:  500,
			Reason: "Died at Level 500 by retainer of Baeloc.",
			Time:   "2025-09-14T20:29:01Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "terrorsleep", Player: false, Traded: false, Summon: ""},
			},
			Level:  494,
			Reason: "Died at Level 494 by terrorsleep.",
			Time:   "2025-09-07T02:39:55Z",
		},
		{
			Assists: []Killers{},
			Killers: []Killers{
				{Name: "headwalker", Player: false, Traded: false, Summon: ""},
			},
			Level:  487,
			Reason: "Died at Level 487 by headwalker.",
			Time:   "2025-08-29T19:09:51Z",
		},
	} {
		assert.True(
			reflect.DeepEqual(deaths[idx].Assists, tc.Assists),
			"Wrong assists\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Assists, deaths[idx].Assists,
		)
		assert.True(
			reflect.DeepEqual(deaths[idx].Killers, tc.Killers),
			"Wrong killers\nidx: %d\nwant: %#v\n\ngot: %#v",
			idx, tc.Killers, deaths[idx].Killers,
		)
		assert.Equal(
			deaths[idx].Level, tc.Level,
			"Wrong Level\nidx: %d\nwant: %d\n\ngot: %d",
			idx, tc.Level, deaths[idx].Level,
		)
		assert.Equal(
			deaths[idx].Reason, tc.Reason,
			"Wrong Reason\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Reason, deaths[idx].Reason,
		)
		assert.Equal(
			tc.Time, deaths[idx].Time,
			"Wrong Time\nidx: %d\nwant: %s\n\ngot: %s",
			idx, tc.Time, deaths[idx].Time,
		)
	}
}

func BenchmarkNumber1(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/characters/Darkside Rafa.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson, _ := TibiaCharactersCharacterImpl(string(data), "")

		assert.Equal(b, "Darkside Rafa", characterJson.Character.CharacterInfo.Name)
	}
}

func BenchmarkNumber2(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/characters/Riley No Hands.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		characterJson, _ := TibiaCharactersCharacterImpl(string(data), "")

		assert.Equal(b, "Riley No Hands", characterJson.Character.CharacterInfo.Name)
	}
}

func TestEmptyName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "",
		},
	}

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(jerr)

	assert.EqualValues(validation.ErrorCharacterNameEmpty.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameEmpty.Error(), jerr.Information.Status.Message)
}

func TestSmallName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "o",
		},
	}

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(validation.ErrorCharacterNameTooSmall.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameTooSmall.Error(), jerr.Information.Status.Message)
}

func TestInvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert := assert.New(t)

	c.Params = []gin.Param{
		{
			Key:   "name",
			Value: "12",
		},
	}

	tibiaCharactersCharacter(c)
	assert.Equal(http.StatusBadRequest, w.Code)

	var jerr OutInformation
	err := json.Unmarshal(w.Body.Bytes(), &jerr)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(validation.ErrorCharacterNameInvalid.Code(), jerr.Information.Status.Error)
	assert.EqualValues(validation.ErrorCharacterNameInvalid.Error(), jerr.Information.Status.Message)
}
