package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

const (
	characterTradesEndingWithinHours = 24
	maxCharacterTradesPages          = 20
)

// Child of CharacterTrades
type CharacterTradeAuction struct {
	AuctionID    int    `json:"auction_id"`    // The internal ID of the auction.
	Name         string `json:"name"`          // The name of the character.
	Level        int    `json:"level"`         // The character's level.
	Vocation     string `json:"vocation"`      // The character's vocation.
	Sex          string `json:"sex"`           // The character's sex.
	World        string `json:"world"`         // The character's world.
	AuctionStart string `json:"auction_start"` // The timestamp when the auction started.
	AuctionEnd   string `json:"auction_end"`   // The timestamp when the auction ends.
	Bid          int    `json:"bid"`           // The current or minimum bid in Tibia Coins.
	BidType      string `json:"bid_type"`      // The bid type. minimum / current
}

// Child of CharacterTrades
type CharacterTradeFilters struct {
	World     string `json:"world,omitempty"`      // The world filter.
	PvpType   string `json:"pvp_type,omitempty"`   // The PvP type filter.
	BattlEye  string `json:"battleye,omitempty"`   // The BattlEye protection filter.
	Vocation  string `json:"vocation,omitempty"`   // The vocation filter.
	MinLevel  int    `json:"min_level,omitempty"`  // The minimum character level filter.
	MaxLevel  int    `json:"max_level,omitempty"`  // The maximum character level filter.
	Skill     string `json:"skill,omitempty"`      // The skill filter.
	MinSkill  int    `json:"min_skill,omitempty"`  // The minimum skill level filter.
	MaxSkill  int    `json:"max_skill,omitempty"`  // The maximum skill level filter.
	Search    string `json:"search,omitempty"`     // The search string filter.
	SearchType string `json:"search_type,omitempty"` // The search type. item / item_wildcard / character_name
}

// characterTradeFilterQuery stores tibia.com query values for character trades.
type characterTradeFilterQuery struct {
	World               string
	WorldPvpType        string
	WorldBattlEyeState  string
	Profession          string
	LevelRangeFrom      string
	LevelRangeTo        string
	SkillID             string
	SkillRangeFrom      string
	SkillRangeTo        string
	SearchString        string
	SearchType          string
	Display             CharacterTradeFilters
}

// Child of JSONData
type CharacterTrades struct {
	EndingWithinHours int                     `json:"ending_within_hours"` // The ending window in hours used for filtering.
	Filters           CharacterTradeFilters   `json:"filters"`             // The filters applied to the request.
	AuctionList       []CharacterTradeAuction `json:"auction_list"`        // List of character auctions.
}

// The base includes two levels: CharacterTrades and Information
type CharacterTradesResponse struct {
	CharacterTrades CharacterTrades `json:"charactertrades"`
	Information     Information     `json:"information"`
}

var (
	characterTradeAuctionIDRegex = regexp.MustCompile(`auctionid=([0-9]+)`)
	characterTradeHeaderRegex    = regexp.MustCompile(`Level:\s*([0-9]+)\s*\|\s*Vocation:\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*World:\s*(?:<a[^>]*>)?([^|<]+)`)
	characterTradeBidRegex       = regexp.MustCompile(`(?i)(Minimum Bid|Current Bid):.*?<b>([0-9,]+)</b>`)
)

// TibiaCharacterTradesEnding func
func TibiaCharacterTradesEndingImpl(c *gin.Context, filters characterTradeFilterQuery, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) (CharacterTradesResponse, error) {
	return tibiaCharacterTradesEndingImpl(c, filters, htmlDataCollector, time.Now())
}

func tibiaCharacterTradesEndingImpl(c *gin.Context, filters characterTradeFilterQuery, htmlDataCollector func(TibiaDataRequestStruct) (string, error), now time.Time) (CharacterTradesResponse, error) {
	// Creating empty vars
	var AuctionListData []CharacterTradeAuction
	var TibiaURLs []string

	deadline := now.Add(time.Duration(characterTradesEndingWithinHours) * time.Hour)

	for page := 1; page <= maxCharacterTradesPages; page++ {
		auctions, pageURL, stop, err := makeCharacterTradesPageRequest(page, filters, now, deadline, htmlDataCollector)
		if err != nil {
			return CharacterTradesResponse{}, fmt.Errorf("[error] TibiaCharacterTradesEndingImpl failed at makeCharacterTradesPageRequest, page: %d, err: %s", page, err)
		}

		TibiaURLs = append(TibiaURLs, pageURL)
		AuctionListData = append(AuctionListData, auctions...)

		if stop {
			break
		}
	}

	//
	// Build the data-blob
	return CharacterTradesResponse{
		CharacterTrades{
			EndingWithinHours: characterTradesEndingWithinHours,
			Filters:           filters.Display,
			AuctionList:       AuctionListData,
		},
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  TibiaURLs,
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}

func makeCharacterTradesPageRequest(page int, filters characterTradeFilterQuery, now time.Time, deadline time.Time, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) ([]CharacterTradeAuction, string, bool, error) {
	// Creating an empty var
	var output []CharacterTradeAuction
	stop := false

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    buildCharacterTradesURL(page, filters),
	}

	BoxContentHTML, err := htmlDataCollector(tibiadataRequest)
	// return error (e.g. for maintenance mode)
	if err != nil {
		return nil, "", false, err
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, "", false, fmt.Errorf("[error] TibiaCharacterTradesEndingImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var insideError error

	// Running query over each div
	ReaderHTML.Find("div.Auction").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into AuctionDivHTML
		AuctionDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCharacterTradesEndingImpl failed at AuctionDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Removing linebreaks from HTML
		AuctionDivHTML = TibiaDataHTMLRemoveLinebreaks(AuctionDivHTML)
		AuctionDivHTML = TibiaDataSanitizeStrings(AuctionDivHTML)

		OneAuction, auctionEndTime, ok := parseCharacterTradeAuction(s, AuctionDivHTML)
		if !ok {
			return true
		}

		if auctionEndTime.After(deadline) {
			stop = true
			return false
		}

		if characterTradeEndsWithin(auctionEndTime, now, deadline) {
			output = append(output, OneAuction)
		}

		return true
	})

	if insideError != nil {
		return nil, tibiadataRequest.URL, false, insideError
	}

	// Stop when the page had no auctions left to scan.
	if ReaderHTML.Find("div.Auction").Length() == 0 {
		stop = true
	}

	return output, tibiadataRequest.URL, stop, nil
}

func buildCharacterTradesURL(page int, filters characterTradeFilterQuery) string {
	query := url.Values{}
	query.Set("subtopic", "currentcharactertrades")
	query.Set("order_column", "101")
	query.Set("order_direction", "1")
	query.Set("currentpage", strconv.Itoa(page))

	if filters.World != "" {
		query.Set("filter_world", filters.World)
	}
	if filters.WorldPvpType != "" {
		query.Set("filter_worldpvptype", filters.WorldPvpType)
	}
	if filters.WorldBattlEyeState != "" {
		query.Set("filter_worldbattleyestate", filters.WorldBattlEyeState)
	}
	if filters.Profession != "" {
		query.Set("filter_profession", filters.Profession)
	}
	if filters.LevelRangeFrom != "" {
		query.Set("filter_levelrangefrom", filters.LevelRangeFrom)
	}
	if filters.LevelRangeTo != "" {
		query.Set("filter_levelrangeto", filters.LevelRangeTo)
	}
	if filters.SkillID != "" {
		query.Set("filter_skillid", filters.SkillID)
	}
	if filters.SkillRangeFrom != "" {
		query.Set("filter_skillrangefrom", filters.SkillRangeFrom)
	}
	if filters.SkillRangeTo != "" {
		query.Set("filter_skillrangeto", filters.SkillRangeTo)
	}
	if filters.SearchString != "" {
		query.Set("searchstring", filters.SearchString)
		if filters.SearchType != "" {
			query.Set("searchtype", filters.SearchType)
		} else {
			query.Set("searchtype", "1")
		}
	}

	return "https://www.tibia.com/charactertrade/?" + query.Encode()
}

func parseCharacterTradeFilters(c *gin.Context) (characterTradeFilterQuery, error) {
	var filters characterTradeFilterQuery

	if err := applyCharacterTradeWorldFilter(c, &filters); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeMappedFilter(c.Query("pvp_type"), characterTradePvpTypeID, "the provided pvp type is invalid", func(id, name string) {
		filters.WorldPvpType = id
		filters.Display.PvpType = name
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeMappedFilter(c.Query("battleye"), characterTradeBattlEyeID, "the provided battleye filter is invalid", func(id, name string) {
		filters.WorldBattlEyeState = id
		filters.Display.BattlEye = name
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeVocationFilter(c, &filters); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeIntFilter(c.Query("min_level"), func(value int) {
		filters.LevelRangeFrom = strconv.Itoa(value)
		filters.Display.MinLevel = value
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeIntFilter(c.Query("max_level"), func(value int) {
		filters.LevelRangeTo = strconv.Itoa(value)
		filters.Display.MaxLevel = value
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeMappedFilter(c.Query("skill"), characterTradeSkillID, "the provided skill filter is invalid", func(id, name string) {
		filters.SkillID = id
		filters.Display.Skill = name
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeIntFilter(c.Query("min_skill"), func(value int) {
		filters.SkillRangeFrom = strconv.Itoa(value)
		filters.Display.MinSkill = value
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeIntFilter(c.Query("max_skill"), func(value int) {
		filters.SkillRangeTo = strconv.Itoa(value)
		filters.Display.MaxSkill = value
	}); err != nil {
		return filters, err
	}
	if err := applyCharacterTradeSearchFilter(c, &filters); err != nil {
		return filters, err
	}

	return filters, nil
}

func applyCharacterTradeWorldFilter(c *gin.Context, filters *characterTradeFilterQuery) error {
	world := strings.TrimSpace(c.Query("world"))
	if world == "" {
		return nil
	}

	world = TibiaDataStringWorldFormatToTitle(world)
	exists, err := validation.WorldExists(world)
	if err != nil {
		return err
	}
	if !exists {
		return validation.ErrorWorldDoesNotExist
	}

	filters.World = world
	filters.Display.World = world
	return nil
}

func applyCharacterTradeVocationFilter(c *gin.Context, filters *characterTradeFilterQuery) error {
	vocation := strings.TrimSpace(c.Query("vocation"))
	if vocation == "" {
		return nil
	}

	id, name, ok := characterTradeProfessionID(vocation)
	if !ok {
		return validation.ErrorVocationDoesNotExist
	}
	if id == "" {
		return nil
	}

	filters.Profession = id
	filters.Display.Vocation = name
	return nil
}

func applyCharacterTradeSearchFilter(c *gin.Context, filters *characterTradeFilterQuery) error {
	search := strings.TrimSpace(c.Query("search"))
	if search == "" {
		return nil
	}

	id, name, ok := characterTradeSearchTypeID(c.Query("search_type"))
	if !ok {
		return fmt.Errorf("the provided search type is invalid")
	}

	filters.SearchString = search
	filters.Display.Search = search
	filters.SearchType = id
	filters.Display.SearchType = name
	return nil
}

func applyCharacterTradeMappedFilter(raw string, mapper func(string) (string, string, bool), invalidMsg string, apply func(id, name string)) error {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}

	id, name, ok := mapper(value)
	if !ok {
		return fmt.Errorf("%s", invalidMsg)
	}

	apply(id, name)
	return nil
}

func applyCharacterTradeIntFilter(raw string, apply func(value int)) error {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return validation.ErrorStringCanNotBeConvertedToInt
	}

	apply(parsed)
	return nil
}

func characterTradeProfessionID(vocation string) (string, string, bool) {
	// bazaar filter_profession values differ from highscores vocation ids
	switch strings.ToLower(vocation) {
	case "all":
		return "", "", true
	case "none":
		return "1", "none", true
	case "druid", "druids":
		return "2", "druids", true
	case "knight", "knights":
		return "3", "knights", true
	case "paladin", "paladins":
		return "4", "paladins", true
	case "sorcerer", "sorcerers":
		return "5", "sorcerers", true
	case "monk", "monks":
		return "6", "monks", true
	default:
		return "", "", false
	}
}

func characterTradePvpTypeID(pvpType string) (string, string, bool) {
	switch strings.ToLower(strings.ReplaceAll(pvpType, " ", "_")) {
	case "open", "open_pvp":
		return "0", "open", true
	case "optional", "optional_pvp":
		return "1", "optional", true
	case "hardcore", "hardcore_pvp":
		return "2", "hardcore", true
	case "retro_open", "retro_open_pvp":
		return "3", "retro_open", true
	case "retro_hardcore", "retro_hardcore_pvp":
		return "4", "retro_hardcore", true
	default:
		return "", "", false
	}
}

func characterTradeBattlEyeID(battleye string) (string, string, bool) {
	switch strings.ToLower(strings.ReplaceAll(battleye, " ", "_")) {
	case "initially_protected", "initially":
		return "1", "initially_protected", true
	case "protected":
		return "2", "protected", true
	case "not_protected", "unprotected":
		return "3", "not_protected", true
	default:
		return "", "", false
	}
}

func characterTradeSkillID(skill string) (string, string, bool) {
	switch strings.ToLower(strings.ReplaceAll(skill, " ", "_")) {
	case "axe", "axe_fighting":
		return "10", "axe_fighting", true
	case "club", "club_fighting":
		return "9", "club_fighting", true
	case "distance", "distance_fighting":
		return "7", "distance_fighting", true
	case "fishing":
		return "13", "fishing", true
	case "fist", "fist_fighting":
		return "11", "fist_fighting", true
	case "magic", "magic_level", "magiclevel":
		return "1", "magic_level", true
	case "shielding", "shield":
		return "6", "shielding", true
	case "sword", "sword_fighting":
		return "8", "sword_fighting", true
	default:
		return "", "", false
	}
}

func characterTradeSearchTypeID(searchType string) (string, string, bool) {
	switch strings.ToLower(strings.ReplaceAll(searchType, " ", "_")) {
	case "", "item", "item_default", "default":
		return "1", "item", true
	case "item_wildcard", "wildcard":
		return "2", "item_wildcard", true
	case "character", "character_name", "name":
		return "3", "character_name", true
	default:
		return "", "", false
	}
}

func parseCharacterTradeAuction(s *goquery.Selection, AuctionDivHTML string) (CharacterTradeAuction, time.Time, bool) {
	var OneAuction CharacterTradeAuction
	var auctionEndTime time.Time

	nameSelection := s.Find(".AuctionCharacterName a").First()
	OneAuction.Name = strings.TrimSpace(nameSelection.Text())
	if OneAuction.Name == "" {
		return OneAuction, auctionEndTime, false
	}

	href, _ := nameSelection.Attr("href")
	subma1 := characterTradeAuctionIDRegex.FindAllStringSubmatch(href, -1)
	if len(subma1) > 0 {
		OneAuction.AuctionID = TibiaDataStringToInteger(subma1[0][1])
	}

	subma2 := characterTradeHeaderRegex.FindAllStringSubmatch(AuctionDivHTML, -1)
	if len(subma2) > 0 {
		OneAuction.Level = TibiaDataStringToInteger(subma2[0][1])
		OneAuction.Vocation = strings.TrimSpace(subma2[0][2])
		OneAuction.Sex = normalizeCharacterTradeSex(strings.TrimSpace(subma2[0][3]))
		OneAuction.World = strings.TrimSpace(subma2[0][4])
	}

	dateValues := s.Find(".ShortAuctionData > .ShortAuctionDataValue")
	if dateValues.Length() >= 2 {
		startRaw := strings.TrimSpace(dateValues.Eq(0).Text())
		endRaw := strings.TrimSpace(dateValues.Eq(1).Text())
		OneAuction.AuctionStart = TibiaDataDatetime(startRaw)
		OneAuction.AuctionEnd = TibiaDataDatetime(endRaw)

		parsedEnd, err := TibiaDataParseDatetime(endRaw)
		if err != nil {
			return OneAuction, auctionEndTime, false
		}
		auctionEndTime = parsedEnd
	} else {
		return OneAuction, auctionEndTime, false
	}

	bidLabel := strings.TrimSpace(s.Find(".ShortAuctionDataBidRow .ShortAuctionDataLabel").First().Text())
	bidValue := strings.TrimSpace(s.Find(".ShortAuctionDataBidRow .ShortAuctionDataValue b").First().Text())
	if bidValue == "" {
		subma3 := characterTradeBidRegex.FindAllStringSubmatch(AuctionDivHTML, -1)
		if len(subma3) > 0 {
			bidLabel = subma3[0][1]
			bidValue = subma3[0][2]
		}
	}
	OneAuction.Bid = TibiaDataStringToInteger(bidValue)
	OneAuction.BidType = normalizeCharacterTradeBidType(bidLabel)

	return OneAuction, auctionEndTime, true
}

func characterTradeEndsWithin(auctionEnd, now, deadline time.Time) bool {
	return !auctionEnd.Before(now) && !auctionEnd.After(deadline)
}

func normalizeCharacterTradeSex(sex string) string {
	switch strings.ToLower(sex) {
	case "male":
		return "male"
	case "female":
		return "female"
	default:
		return strings.ToLower(sex)
	}
}

func normalizeCharacterTradeBidType(label string) string {
	label = strings.TrimSuffix(strings.TrimSpace(label), ":")
	switch strings.ToLower(label) {
	case "minimum bid":
		return "minimum"
	case "current bid":
		return "current"
	default:
		return strings.ToLower(label)
	}
}
