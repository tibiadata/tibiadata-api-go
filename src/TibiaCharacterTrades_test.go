package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestCharacterTradesEnding(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/charactertrades/current.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	now, err := time.ParseInLocation("Jan 02 2006, 15:04 MST", "Aug 16 2026, 12:00 CEST", berlinLocation)
	if err != nil {
		t.Fatal(err)
	}

	characterTradesJson, err := tibiaCharacterTradesEndingImpl(
		nil,
		characterTradeFilterQuery{},
		func(request TibiaDataRequestStruct) (string, error) {
			return string(data), nil
		},
		now)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := characterTradesJson.Information

	assert.Equal(24, characterTradesJson.CharacterTrades.EndingWithinHours)
	assert.Equal(2, len(characterTradesJson.CharacterTrades.AuctionList))

	firstAuction := characterTradesJson.CharacterTrades.AuctionList[0]
	assert.Equal(1323306, firstAuction.AuctionID)
	assert.Equal("Mundo Carry Eune", firstAuction.Name)
	assert.Equal(8, firstAuction.Level)
	assert.Equal("Knight", firstAuction.Vocation)
	assert.Equal("male", firstAuction.Sex)
	assert.Equal("Thyria", firstAuction.World)
	assert.Equal("2026-08-15T08:05:00Z", firstAuction.AuctionStart)
	assert.Equal("2026-08-16T20:00:00Z", firstAuction.AuctionEnd)
	assert.Equal(454, firstAuction.Bid)
	assert.Equal("current", firstAuction.BidType)

	secondAuction := characterTradesJson.CharacterTrades.AuctionList[1]
	assert.Equal(1323401, secondAuction.AuctionID)
	assert.Equal("Test Druid One", secondAuction.Name)
	assert.Equal(445, secondAuction.Level)
	assert.Equal("Elder Druid", secondAuction.Vocation)
	assert.Equal("male", secondAuction.Sex)
	assert.Equal("Antica", secondAuction.World)
	assert.Equal("2026-08-14T07:00:00Z", secondAuction.AuctionStart)
	assert.Equal("2026-08-17T08:00:00Z", secondAuction.AuctionEnd)
	assert.Equal(1500, secondAuction.Bid)
	assert.Equal("minimum", secondAuction.BidType)

	assert.Equal("https://www.tibia.com/charactertrade/?currentpage=1&order_column=101&order_direction=1&subtopic=currentcharactertrades", information.TibiaURLs[0])
	assert.Equal(1, len(information.TibiaURLs))
}

func TestCharacterTradesEndingPagination(t *testing.T) {
	page1File, err := static.TestFiles.Open("testdata/charactertrades/current_page1.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer page1File.Close()

	page1Data, err := io.ReadAll(page1File)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	page2File, err := static.TestFiles.Open("testdata/charactertrades/current_page2.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer page2File.Close()

	page2Data, err := io.ReadAll(page2File)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	now, err := time.ParseInLocation("Jan 02 2006, 15:04 MST", "Aug 16 2026, 12:00 CEST", berlinLocation)
	if err != nil {
		t.Fatal(err)
	}

	characterTradesJson, err := tibiaCharacterTradesEndingImpl(
		nil,
		characterTradeFilterQuery{},
		func(request TibiaDataRequestStruct) (string, error) {
			if strings.Contains(request.URL, "currentpage=2") {
				return string(page2Data), nil
			}

			return string(page1Data), nil
		},
		now)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := characterTradesJson.Information

	assert.Equal(2, len(characterTradesJson.CharacterTrades.AuctionList))
	assert.Equal(2, len(information.TibiaURLs))
	assert.Equal("https://www.tibia.com/charactertrade/?currentpage=1&order_column=101&order_direction=1&subtopic=currentcharactertrades", information.TibiaURLs[0])
	assert.Equal("https://www.tibia.com/charactertrade/?currentpage=2&order_column=101&order_direction=1&subtopic=currentcharactertrades", information.TibiaURLs[1])
}

func TestCharacterTradesEndingFilter(t *testing.T) {
	assert := assert.New(t)

	now, err := time.ParseInLocation("Jan 02 2006, 15:04 MST", "Aug 16 2026, 12:00 CEST", berlinLocation)
	if err != nil {
		t.Fatal(err)
	}
	deadline := now.Add(24 * time.Hour)

	within, err := TibiaDataParseDatetime("Aug 16 2026, 22:00 CEST")
	if err != nil {
		t.Fatal(err)
	}
	outside, err := TibiaDataParseDatetime("Aug 18 2026, 10:00 CEST")
	if err != nil {
		t.Fatal(err)
	}
	past, err := TibiaDataParseDatetime("Aug 15 2026, 10:00 CEST")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(characterTradeEndsWithin(within, now, deadline))
	assert.False(characterTradeEndsWithin(outside, now, deadline))
	assert.False(characterTradeEndsWithin(past, now, deadline))
}

func TestCharacterTradesEndingWithFilters(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/charactertrades/current.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	now, err := time.ParseInLocation("Jan 02 2006, 15:04 MST", "Aug 16 2026, 12:00 CEST", berlinLocation)
	if err != nil {
		t.Fatal(err)
	}

	filters := characterTradeFilterQuery{
		World:              "Antica",
		WorldPvpType:       "1",
		WorldBattlEyeState: "2",
		Profession:         "3",
		LevelRangeFrom:     "100",
		LevelRangeTo:       "500",
		SkillID:            "1",
		SkillRangeFrom:     "80",
		SkillRangeTo:       "120",
		SearchString:       "Ferumbras",
		SearchType:         "1",
		Display: CharacterTradeFilters{
			World:      "Antica",
			PvpType:    "optional",
			BattlEye:   "protected",
			Vocation:   "knights",
			MinLevel:   100,
			MaxLevel:   500,
			Skill:      "magic_level",
			MinSkill:   80,
			MaxSkill:   120,
			Search:     "Ferumbras",
			SearchType: "item",
		},
	}

	var requestedURL string
	characterTradesJson, err := tibiaCharacterTradesEndingImpl(
		nil,
		filters,
		func(request TibiaDataRequestStruct) (string, error) {
			requestedURL = request.URL
			return string(data), nil
		},
		now)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Antica", characterTradesJson.CharacterTrades.Filters.World)
	assert.Equal("optional", characterTradesJson.CharacterTrades.Filters.PvpType)
	assert.Equal("protected", characterTradesJson.CharacterTrades.Filters.BattlEye)
	assert.Equal("knights", characterTradesJson.CharacterTrades.Filters.Vocation)
	assert.Equal(100, characterTradesJson.CharacterTrades.Filters.MinLevel)
	assert.Equal(500, characterTradesJson.CharacterTrades.Filters.MaxLevel)
	assert.Equal("magic_level", characterTradesJson.CharacterTrades.Filters.Skill)
	assert.Equal(80, characterTradesJson.CharacterTrades.Filters.MinSkill)
	assert.Equal(120, characterTradesJson.CharacterTrades.Filters.MaxSkill)
	assert.Equal("Ferumbras", characterTradesJson.CharacterTrades.Filters.Search)
	assert.Equal("item", characterTradesJson.CharacterTrades.Filters.SearchType)

	assert.Contains(requestedURL, "filter_world=Antica")
	assert.Contains(requestedURL, "filter_worldpvptype=1")
	assert.Contains(requestedURL, "filter_worldbattleyestate=2")
	assert.Contains(requestedURL, "filter_profession=3")
	assert.Contains(requestedURL, "filter_levelrangefrom=100")
	assert.Contains(requestedURL, "filter_levelrangeto=500")
	assert.Contains(requestedURL, "filter_skillid=1")
	assert.Contains(requestedURL, "filter_skillrangefrom=80")
	assert.Contains(requestedURL, "filter_skillrangeto=120")
	assert.Contains(requestedURL, "searchstring=Ferumbras")
	assert.Contains(requestedURL, "searchtype=1")
	assert.Contains(requestedURL, "order_column=101")
	assert.Contains(requestedURL, "order_direction=1")
}

func TestParseCharacterTradeFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/v4/charactertrades/ending?vocation=knights&pvp_type=optional&battleye=protected&min_level=100&max_level=300&skill=magic_level&min_skill=70&max_skill=90&search=boots&search_type=item_wildcard", nil)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	filters, err := parseCharacterTradeFilters(c)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal("3", filters.Profession)
	assert.Equal("knights", filters.Display.Vocation)
	assert.Equal("1", filters.WorldPvpType)
	assert.Equal("optional", filters.Display.PvpType)
	assert.Equal("2", filters.WorldBattlEyeState)
	assert.Equal("protected", filters.Display.BattlEye)
	assert.Equal("100", filters.LevelRangeFrom)
	assert.Equal("300", filters.LevelRangeTo)
	assert.Equal("1", filters.SkillID)
	assert.Equal("magic_level", filters.Display.Skill)
	assert.Equal("70", filters.SkillRangeFrom)
	assert.Equal("90", filters.SkillRangeTo)
	assert.Equal("boots", filters.SearchString)
	assert.Equal("2", filters.SearchType)
	assert.Equal("item_wildcard", filters.Display.SearchType)
}

func TestBuildCharacterTradesURL(t *testing.T) {
	assert := assert.New(t)

	url := buildCharacterTradesURL(1, characterTradeFilterQuery{
		World:      "Premia",
		Profession: "2",
	})

	assert.Contains(url, "subtopic=currentcharactertrades")
	assert.Contains(url, "currentpage=1")
	assert.Contains(url, "filter_world=Premia")
	assert.Contains(url, "filter_profession=2")
	assert.Contains(url, "order_column=101")
	assert.Contains(url, "order_direction=1")
}
