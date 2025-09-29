package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/PuerkitoBio/goquery"
)

// Child of BazaarAuction
type BazaarAuctionBid struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

// Child of BazaarAuction
type BazaarAuctionDetails struct {
	CharacterName string           `json:"character_name"`
	Level         int              `json:"level"`
	Vocation      string           `json:"vocation"`
	Gender        string           `json:"gender"`
	World         string           `json:"world"`
	AuctionStart  string           `json:"auction_start"`
	AuctionEnd    string           `json:"auction_end"`
	Bid           BazaarAuctionBid `json:"bid"`
}

// Child of BazaarAuction
type General struct {
	HitPoints                 int    `json:"hitPoints"`
	Mana                      int    `json:"mana"`
	Capacity                  int    `json:"capacity"`
	Speed                     int    `json:"speed"`
	Blessings                 int    `json:"blessings"`
	Mounts                    int    `json:"mounts"`
	Outfits                   int    `json:"outfits"`
	Titles                    int    `json:"titles"`
	AxeFighting               int    `json:"axe_fighting"`
	ClubFighting              int    `json:"club_fighting"`
	DistanceFighting          int    `json:"distance_fighting"`
	Fishing                   int    `json:"fishing"`
	FistFighting              int    `json:"fist_fighting"`
	MagicLevel                int    `json:"magic_level"`
	Shielding                 int    `json:"shielding"`
	SwordFighting             int    `json:"sword_fighting"`
	CreationDate              string `json:"creation_date"`
	Experience                int    `json:"experience"`
	Gold                      int    `json:"gold"`
	AchievementPoints         int    `json:"achievement_points"`
	RegularWorldTransfer      string `json:"regular_world_transfer"` // woot
	CharmExpansion            bool   `json:"charm_expansion"`
	AvailableCharmPoints      int    `json:"available_charm_points"`
	SpentCharmPoints          int    `json:"spent_charm_points"`
	DailyRewardStreak         int    `json:"daily_reward_streak"`
	HuntingTaskPoints         int    `json:"hunting_task_points"`
	PermanentHuntingTaskSlots int    `json:"permanent_hunting_task_slots"`
	PermanentPreySlots        int    `json:"permanent_prey_slots"`
	PreyWildCards             int    `json:"prey_wild_cards"`
	Hirelings                 int    `json:"hirelings"`
	HirelingJobs              int    `json:"hireling_jobs"`
	HirelingOutfits           int    `json:"hireling_outfits"`
	ExaltedDust               int    `json:"exalted_dust"`
}

// Child of BazaarAuction
type BazaarAuctionItem struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// Child of BazaarAuction
type BazaarAuctionOutfit struct {
	Name   string `json:"name"`
	Addon1 bool   `json:"addon_1"`
	Addon2 bool   `json:"addon_2"`
}

// Child of BazaarAuction
type BazaarAuctionBlessings struct {
	AdventurersBlessing int `json:"adventurers_blessing"`
	BloodOfTheMountain  int `json:"blood_of_the_mountain"`
	EmbraceOfTibia      int `json:"embrace_of_tibia"`
	FireOfTheSuns       int `json:"fire_of_the_suns"`
	HeartOfTheMountain  int `json:"heart_of_the_mountain"`
	SparkOfThePhoenix   int `json:"spark_of_the_phoenix"`
	SpiritualShielding  int `json:"spiritual_shielding"`
	TwistOfFate         int `json:"twist_of_fate"`
	WisdomOfSolitude    int `json:"wisdom_of_solitude"`
}

// Child of BazaarAuction
type BazaarAuctionCharm struct {
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

// Child of BazaarAuction
type BazaarAuctionBestiaryEntry struct {
	Name  string `json:"name"`
	Kills int    `json:"kills"`
	Step  int    `json:"step"`
}

// Child of BazaarAuctionResponse
type BazaarAuction struct {
	Id                          int                          `json:"id"`
	Details                     BazaarAuctionDetails         `json:"details"`
	General                     General                      `json:"general"`
	ItemSummary                 []BazaarAuctionItem          `json:"item_summary"`
	StoreItemSummary            []BazaarAuctionItem          `json:"store_item_summary"`
	Mounts                      []string                     `json:"mounts"`
	StoreMounts                 []string                     `json:"store_mounts"`
	Outfits                     []BazaarAuctionOutfit        `json:"outfits"`
	StoreOutfits                []BazaarAuctionOutfit        `json:"store_outfits"`
	Familiars                   []string                     `json:"familiars"`
	Blessings                   BazaarAuctionBlessings       `json:"blessings"`
	Imbuements                  []string                     `json:"imbuements"`
	Charms                      []BazaarAuctionCharm         `json:"charms"`
	CompletedCyclopediaMapAreas []string                     `json:"completed_cyclopedia_map_areas"`
	CompletedQuestLines         []string                     `json:"completed_quest_lines"`
	Titles                      []string                     `json:"titles"`
	Achievements                []string                     `json:"achievements"`
	BestiaryProgress            []BazaarAuctionBestiaryEntry `json:"bestiary_progress"`
}

// The base includes two levels: Auction and Information
type BazaarAuctionResponse struct {
	Auction     BazaarAuction `json:"auction"`
	Information Information   `json:"information"`
}

const (
	OddEvenSelector         = ".Odd,.Even"
	PageLinkSelector        = ".PageLink"
	CVIconSelector          = ".CVIcon"
	ItemSummarySection      = 0
	StoreItemSummarySection = 1
	MountsSection           = 2
	StoreMountsSection      = 3
	OutfitsSection          = 4
	StoreOutfitsSection     = 5
)

// TibiaBazaarAuctionV3Impl func
func TibiaBazaarAuctionV3Impl(BoxContentHTML string, url string) (BazaarAuctionResponse, error) {

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	var id int
	ReaderHTML.Find("input[name=auctionid]").Each(func(i int, selection *goquery.Selection) {
		// collect the auction ID
		id = TibiaDataStringToInteger(selection.AttrOr("value", ""))
	})

	// Extract details section
	var details BazaarAuctionDetails
	ReaderHTML.Find(".Auction").Each(func(index int, s *goquery.Selection) {
		detailsHeader := strings.Split(s.Find(".AuctionHeader").Text(), "Level: ")
		details.CharacterName = detailsHeader[0]

		detailsHeader = strings.Split(detailsHeader[1], "|")

		level := TibiaDataStringToInteger(strings.TrimSpace(detailsHeader[0]))

		details.Level = level
		details.Vocation = strings.TrimSpace(strings.Split(detailsHeader[1], "Vocation: ")[1])
		details.Gender = strings.TrimSpace(detailsHeader[2])
		details.World = strings.TrimSpace(strings.Split(detailsHeader[3], "World: ")[1])

		s.Find(".ShortAuctionData").Each(func(index int, s *goquery.Selection) {
			nodes := s.Children().Nodes

			lookupIndex := 0
			hasTimer := strings.EqualFold(nodes[0].Attr[0].Val, "AuctionTimer")

			// In case the auction hasn't ended yet, we need to increase the lookup index.
			if hasTimer {
				lookupIndex = 1
			}

			auctionStartDate := TibiaDataSanitizeStrings(nodes[1+lookupIndex].FirstChild.Data)
			auctionStartDate = strings.Split(auctionStartDate, " CET")[0] + ":00 CET"
			auctionStartDate = strings.Split(auctionStartDate, " CEST")[0] + ":00 CEST"

			auctionEndDate := TibiaDataSanitizeStrings(nodes[3+lookupIndex].FirstChild.Data)
			auctionEndDate = strings.Split(auctionEndDate, " CET")[0] + ":00 CET"
			auctionEndDate = strings.Split(auctionEndDate, " CEST")[0] + ":00 CEST"

			details.AuctionStart = TibiaDataDatetime(auctionStartDate)
			details.AuctionEnd = TibiaDataDatetime(auctionEndDate)

			bidType := strings.Split(nodes[4+lookupIndex].FirstChild.FirstChild.Data, " Bid:")[0]
			bidAmount := TibiaDataStringToInteger(nodes[4+lookupIndex].LastChild.FirstChild.FirstChild.Data)

			details.Bid = BazaarAuctionBid{
				Type:   bidType,
				Amount: bidAmount,
			}
		})
	})

	// Extract general section
	var general General
	ReaderHTML.Find("#General").Each(func(index int, s *goquery.Selection) {

		// General
		generalMap := make(map[string]string)
		s.Find(".LabelV").Each(func(index int, s *goquery.Selection) {
			generalMap[strings.Split(s.Nodes[0].FirstChild.Data, ":")[0]] = s.Nodes[0].NextSibling.FirstChild.Data
		})

		// Skills
		skillsMap := make(map[string]int)
		s.Find(".LabelColumn").Each(func(index int, s *goquery.Selection) {
			skillsMap[strings.Split(s.Nodes[0].FirstChild.FirstChild.Data, ":")[0]] =
				TibiaDataStringToInteger(s.Nodes[0].NextSibling.FirstChild.Data)
		})

		general.HitPoints = TibiaDataStringToInteger(generalMap["Hit Points"])
		general.Mana = TibiaDataStringToInteger(generalMap["Mana"])
		general.Capacity = TibiaDataStringToInteger(generalMap["Capacity"])
		general.Speed = TibiaDataStringToInteger(generalMap["Speed"])
		general.Blessings = TibiaDataStringToInteger(strings.Split(generalMap["Blessings"], "/")[0])
		general.Mounts = TibiaDataStringToInteger(generalMap["Mounts"])
		general.Outfits = TibiaDataStringToInteger(generalMap["Outfits"])
		general.Titles = TibiaDataStringToInteger(generalMap["Titles"])
		general.AxeFighting = skillsMap["Axe Fighting"]
		general.ClubFighting = skillsMap["Club Fighting"]
		general.DistanceFighting = skillsMap["Distance Fighting"]
		general.Fishing = skillsMap["Fishing"]
		general.FistFighting = skillsMap["Fist Fighting"]
		general.MagicLevel = skillsMap["Magic Level"]
		general.Shielding = skillsMap["Shielding"]
		general.SwordFighting = skillsMap["Sword Fighting"]
		general.CreationDate = TibiaDataDatetime(generalMap["Creation Date"])
		general.Experience = TibiaDataStringToInteger(generalMap["Experience"])
		general.Gold = TibiaDataStringToInteger(generalMap["Gold"])
		general.AchievementPoints = TibiaDataStringToInteger(generalMap["Achievement Points"])
		general.RegularWorldTransfer = generalMap["Regular World Transfer"]
		general.CharmExpansion = strings.EqualFold(generalMap["Charm Expansion"], "yes")
		general.AvailableCharmPoints = TibiaDataStringToInteger(generalMap["Available Charm Points"])
		general.SpentCharmPoints = TibiaDataStringToInteger(generalMap["Spent Charm Points"])
		general.DailyRewardStreak = TibiaDataStringToInteger(generalMap["Daily Reward Streak"])
		general.HuntingTaskPoints = TibiaDataStringToInteger(generalMap["Hunting Task Points"])
		general.PermanentHuntingTaskSlots = TibiaDataStringToInteger(generalMap["Permanent Hunting Task Slots"])
		general.PermanentPreySlots = TibiaDataStringToInteger(generalMap["Permanent Prey Slots"])
		general.PreyWildCards = TibiaDataStringToInteger(generalMap["Prey Wildcards"])
		general.Hirelings = TibiaDataStringToInteger(generalMap["Hirelings"])
		general.HirelingJobs = TibiaDataStringToInteger(generalMap["Hireling Jobs"])
		general.HirelingOutfits = TibiaDataStringToInteger(generalMap["Hireling Outfits"])
		general.ExaltedDust = TibiaDataStringToInteger(strings.Split(generalMap["Exalted Dust"], "/")[0])
	})

	// Extract items summary
	var itemSummary []BazaarAuctionItem
	ReaderHTML.Find("#ItemSummary").Each(func(index int, s *goquery.Selection) {

		for k, v := range ParseItems(s) {
			itemSummary = append(itemSummary, BazaarAuctionItem{Name: k, Amount: v})
		}

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				itemPage := AjaxJSONDataCollectorV3(id, ItemSummarySection, pageIndex)
				ItemPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(itemPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseItems(ItemPageReaderHTML.Contents()) {
					itemSummary = append(itemSummary, BazaarAuctionItem{Name: k, Amount: v})
				}
			}
		}
	})

	// Extract store items summary
	var storeItemSummary []BazaarAuctionItem
	ReaderHTML.Find("#StoreItemSummary").Each(func(index int, s *goquery.Selection) {

		for k, v := range ParseItems(s) {
			storeItemSummary = append(storeItemSummary, BazaarAuctionItem{Name: k, Amount: v})
		}

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				itemPage := AjaxJSONDataCollectorV3(id, StoreItemSummarySection, pageIndex)
				ItemPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(itemPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseItems(ItemPageReaderHTML.Contents()) {
					storeItemSummary = append(storeItemSummary, BazaarAuctionItem{Name: k, Amount: v})
				}
			}
		}
	})

	// Extract mounts
	var mounts []string
	ReaderHTML.Find("#Mounts").Each(func(index int, s *goquery.Selection) {

		mounts = append(mounts, ParseMounts(s)...)

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				mountsPage := AjaxJSONDataCollectorV3(id, MountsSection, pageIndex)
				MountsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(mountsPage))
				if err != nil {
					log.Fatal(err)
				}
				mounts = append(mounts, ParseMounts(MountsPageReaderHTML.Contents())...)
			}
		}
	})

	// Extract store mounts
	var storeMounts []string
	ReaderHTML.Find("#StoreMounts").Each(func(index int, s *goquery.Selection) {

		storeMounts = append(storeMounts, ParseMounts(s)...)

		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				mountsPage := AjaxJSONDataCollectorV3(id, StoreMountsSection, pageIndex)
				MountsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(mountsPage))
				if err != nil {
					log.Fatal(err)
				}
				storeMounts = append(storeMounts, ParseMounts(MountsPageReaderHTML.Contents())...)
			}
		}
	})

	// Extract outfits
	var outfits []BazaarAuctionOutfit
	ReaderHTML.Find("#Outfits").Each(func(index int, s *goquery.Selection) {
		for k, v := range ParseOutfits(s) {
			outfits = append(outfits, BazaarAuctionOutfit{
				Name:   k,
				Addon1: v[0],
				Addon2: v[1],
			})
		}
		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				outfitsPage := AjaxJSONDataCollectorV3(id, OutfitsSection, pageIndex)
				OutfitsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(outfitsPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseOutfits(OutfitsPageReaderHTML.Contents()) {
					outfits = append(outfits, BazaarAuctionOutfit{
						Name:   k,
						Addon1: v[0],
						Addon2: v[1],
					})
				}
			}
		}
	})

	// Extract store outfits
	var storeOutfits []BazaarAuctionOutfit
	ReaderHTML.Find("#StoreOutfits").Each(func(index int, s *goquery.Selection) {
		for k, v := range ParseOutfits(s) {
			outfits = append(outfits, BazaarAuctionOutfit{
				Name:   k,
				Addon1: v[0],
				Addon2: v[1],
			})
		}
		totalPages := s.Find(PageLinkSelector).Size()
		if totalPages > 1 {
			// Fetch missing pages using links
			for pageIndex := 2; pageIndex <= totalPages; pageIndex++ {
				outfitsPage := AjaxJSONDataCollectorV3(id, StoreOutfitsSection, pageIndex)
				OutfitsPageReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(outfitsPage))
				if err != nil {
					log.Fatal(err)
				}
				for k, v := range ParseOutfits(OutfitsPageReaderHTML.Contents()) {
					outfits = append(outfits, BazaarAuctionOutfit{
						Name:   k,
						Addon1: v[0],
						Addon2: v[1],
					})
				}
			}
		}
	})

	// Extract familiars
	var familiars []string
	ReaderHTML.Find("#Familiars").Each(func(index int, s *goquery.Selection) {
		s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
			familiarName, exists := s.Attr("title")
			if exists {
				familiars = append(familiars, familiarName)
			}
		})
	})

	// Extract blessings
	var blessings BazaarAuctionBlessings
	ReaderHTML.Find("#Blessings").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			blessingsAmount := TibiaDataStringToInteger(strings.Split(node.FirstChild.Data, " x")[0])
			switch blessingName := node.NextSibling.FirstChild.Data; blessingName {
			case "Adventurer's Blessing":
				blessings.AdventurersBlessing = blessingsAmount
			case "Blood of the Mountain":
				blessings.BloodOfTheMountain = blessingsAmount
			case "Embrace of Tibia":
				blessings.EmbraceOfTibia = blessingsAmount
			case "Fire of the Suns":
				blessings.FireOfTheSuns = blessingsAmount
			case "Heart of the Mountain":
				blessings.HeartOfTheMountain = blessingsAmount
			case "Spark of the Phoenix":
				blessings.SparkOfThePhoenix = blessingsAmount
			case "Spiritual Shielding":
				blessings.SpiritualShielding = blessingsAmount
			case "Twist of Fate":
				blessings.TwistOfFate = blessingsAmount
			case "Wisdom of Solitude":
				blessings.WisdomOfSolitude = blessingsAmount
			}
		})
	})

	// Extract Imbuements
	var imbuements []string
	ReaderHTML.Find("#Imbuements").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No imbuements.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				imbuements = append(imbuements, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Charms
	var charms []BazaarAuctionCharm
	ReaderHTML.Find("#Charms").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No charms.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				charms = append(charms, BazaarAuctionCharm{
					Cost: TibiaDataStringToInteger(node.FirstChild.Data),
					Name: node.NextSibling.FirstChild.Data,
				})
			}
		})
	})

	// Extract Completed Cyclopedia Map Areas
	var completedCyclopediaMapAreas []string
	ReaderHTML.Find("#CompletedCyclopediaMapAreas").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if strings.Contains(node.FirstChild.Data, "No areas explored.") {
				return
			}
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				completedCyclopediaMapAreas = append(completedCyclopediaMapAreas, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Completed Quest Lines
	var completedQuestLines []string
	ReaderHTML.Find("#CompletedQuestLines").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				completedQuestLines = append(completedQuestLines, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Titles
	var titles []string
	ReaderHTML.Find("#Titles").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				titles = append(titles, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Achievements
	var achievements []string
	ReaderHTML.Find("#Achievements").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				achievements = append(achievements, strings.TrimSpace(node.FirstChild.Data))
			}
		})
	})

	// Extract Bestiary Progress
	var bestiaryProgress []BazaarAuctionBestiaryEntry
	ReaderHTML.Find("#BestiaryProgress").Each(func(index int, s *goquery.Selection) {
		s.Find(OddEvenSelector).Each(func(index int, s *goquery.Selection) {
			node := s.Nodes[0].FirstChild
			if !strings.Contains(node.Parent.Attr[0].Val, "IndicateMoreEntries") {
				bestiaryProgress = append(bestiaryProgress, BazaarAuctionBestiaryEntry{
					Step:  TibiaDataStringToInteger(node.FirstChild.Data),
					Kills: TibiaDataStringToInteger(strings.Split(node.NextSibling.FirstChild.Data, " x")[0]),
					Name:  node.NextSibling.NextSibling.FirstChild.Data,
				})
			}
		})
	})

	// Build the data-blob
	return BazaarAuctionResponse{
		BazaarAuction{
			Id:                          id,
			Details:                     details,
			General:                     general,
			ItemSummary:                 itemSummary,
			StoreItemSummary:            storeItemSummary,
			Mounts:                      mounts,
			StoreMounts:                 storeMounts,
			Outfits:                     outfits,
			StoreOutfits:                storeOutfits,
			Familiars:                   familiars,
			Blessings:                   blessings,
			Imbuements:                  imbuements,
			Charms:                      charms,
			CompletedCyclopediaMapAreas: completedCyclopediaMapAreas,
			CompletedQuestLines:         completedQuestLines,
			Titles:                      titles,
			Achievements:                achievements,
			BestiaryProgress:            bestiaryProgress,
		},
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{url},
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}

func ParseItems(s *goquery.Selection) map[string]int {
	m := make(map[string]int)
	s.Find(".CVIconObject").Each(func(index int, s *goquery.Selection) {

		itemTitle, exists := s.Attr("title")

		if exists {
			var (
				itemAmount int
				itemName   string
			)

			nodes := s.Find(".ObjectAmount").First().Nodes
			if nodes == nil {
				itemName = strings.Split(itemTitle, "\n")[0]
				itemAmount = 1
			} else {
				temp := strings.SplitN(itemTitle, "x ", 2)
				itemName = strings.Split(temp[1], "\n")[0]
				itemAmount = TibiaDataStringToInteger(temp[0])
			}
			m[itemName] = itemAmount
		}
	})
	return m
}

func ParseOutfits(s *goquery.Selection) map[string][]bool {
	m := make(map[string][]bool)
	s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
		outfitTitle, exists := s.Attr("title")
		if exists {
			outfitName := strings.Split(outfitTitle, " (")[0]
			hasAddon1 := strings.Contains(outfitTitle, "addon 1")
			hasAddon2 := strings.Contains(outfitTitle, "addon 2")
			m[outfitName] = []bool{hasAddon1, hasAddon2}
		}
	})
	return m
}

func ParseMounts(s *goquery.Selection) []string {
	var mountsList []string
	s.Find(CVIconSelector).Each(func(index int, s *goquery.Selection) {
		mountTitle, exists := s.Attr("title")
		if exists {
			mountsList = append(mountsList, mountTitle)
		}
	})
	return mountsList
}

func AjaxJSONDataCollectorV3(AuctionId int, SectionType int, PageIndex int) string {
	TibiaURL := "https://www.tibia.com/websiteservices/handle_charactertrades.php?" +
		"auctionid=" + strconv.Itoa(AuctionId) +
		"&type=" + strconv.Itoa(SectionType) +
		"&currentpage=" + strconv.Itoa(PageIndex)

	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiaDataDebug var
	if TibiaDataDebug {
		client.SetDebug(true)
	}

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Content-Type":     "application/json",
		"User-Agent":       TibiaDataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	res, err := client.R().Get(TibiaURL)
	if err != nil {
		log.Printf("[error] AjaxJSONDataCollectorV3 (URL: %s) in resp1: %s", TibiaURL, err)
	}

	// Checking if status is something else than 200
	if res.StatusCode() != 200 {
		log.Printf("[warni] AjaxJSONDataCollectorV3 (URL: %s) status code: %s", TibiaURL, res.Status())

		// Check if page is in maintenance mode
		if res.StatusCode() == 302 {
			log.Printf("[info] AjaxJSONDataCollectorV3 (URL: %s): Page tibia.com returns 302, probably maintenance mode enabled?", TibiaURL)
		}
	}

	var result AjaxJSONData
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		log.Printf("[error] AjaxJSONDataCollectorV3 (URL: %s) in deserialization process: %s", TibiaURL, err)
	}

	// Return of extracted html to functions
	return result.AjaxObjects[0].Data
}

// AjaxResponseObject - child of AjaxJSONData
type AjaxResponseObject struct {
	Data     string `json:"Data"`
	DataType string `json:"DataType"`
	Target   string `json:"Target"`
}

// AjaxJSONData - base response for auction items page links
type AjaxJSONData struct {
	AjaxObjects []AjaxResponseObject `json:"AjaxObjects"`
}
