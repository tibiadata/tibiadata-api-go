package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Child of JSONData
type Creature struct {
	Name             string   `json:"name"`              // The name of the creature.
	Race             string   `json:"race"`              // The creature's internal name.
	ImageURL         string   `json:"image_url"`         // The URL to this creature's image.
	Description      string   `json:"description"`       // A description of the creature.
	Behaviour        string   `json:"behaviour"`         // The plain description of behaviour of the creature.
	Hitpoints        int      `json:"hitpoints"`         // The number of hitpoints the creature has.
	ImmuneTo         []string `json:"immune"`            // The elements it is immune to.
	StrongAgainst    []string `json:"strong"`            // The elements it is strong against.
	WeaknessAgainst  []string `json:"weakness"`          // The elements it is weak against.
	HealedBy         []string `json:"healed"`            // The elements it is healed when being damaged.
	BeParalysed      bool     `json:"be_paralysed"`      // Whether it can be paralysed or not.
	BeSummoned       bool     `json:"be_summoned"`       // Whether it can be summoned or not.
	SummonMana       int      `json:"summoned_mana"`     // The mana neccessary to summon it.
	BeConvinced      bool     `json:"be_convinced"`      // Whether it can be convinced or not.
	ConvincedMana    int      `json:"convinced_mana"`    // The mana neccessary to convince it.
	SeeInvisible     bool     `json:"see_invisible"`     // Whether it can see even when being invisible or not.
	ExperiencePoints int      `json:"experience_points"` // The number of experience points given for killing it.
	IsLootable       bool     `json:"is_lootable"`       // Whether it can be looted or not.
	LootList         []string `json:"loot_list"`         // Some of the items it drops.
	Featured         bool     `json:"featured"`          // Whether it is featured of not.
}

// The base includes two levels: Creature and Information
type CreatureResponse struct {
	Creature    Creature    `json:"creature"`
	Information Information `json:"information"`
}

var (
	CreatureDataRegex         = regexp.MustCompile(`.*;">(.*)<\/h2>  <img src="(.*)"\/>.*<p>(.*)<\/p>  <p>(.*)<\/p>  <p>(.*)<\/p>.*`)
	CreatureHitpointsRegex    = regexp.MustCompile(`.*have (.*) hitpoints. (.*)`)
	CreatureImmuneRegex       = regexp.MustCompile(`.*are immune to (.*)`)
	CreatureStrongRegex       = regexp.MustCompile(`.*are strong against (.*)`)
	CreatureWeakRegex         = regexp.MustCompile(`.*are weak against (.*)`)
	CreatureHealedRegex       = regexp.MustCompile(`.*are healed when (.*) damage is used.*`)
	CreatureManaRequiredRegex = regexp.MustCompile(`.*It takes (.*) mana to (.*)`)
	CreatureLootRegex         = regexp.MustCompile(`.*yield (.*) experience.*carry (.*)with them.`)
)

func TibiaCreaturesCreatureImpl(race string, BoxContentHTML string, url string) (CreatureResponse, error) {
	// local strings used in this function
	localDamageString := " damage"

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return CreatureResponse{}, fmt.Errorf("[error] TibiaCreaturesCreatureImp failed at goquery.NewDocumentFromReader, error: %s", err)
	}

	// Getting data
	InnerTableContainerTMP1, err := ReaderHTML.Find(".BoxContent div").First().NextAll().Html()
	if err != nil {
		return CreatureResponse{}, fmt.Errorf("[error] TibiaCreaturesCreatureImp failed at ReaderHTML.Find, error: %s", err)
	}

	// Regex to get data
	subma1 := CreatureDataRegex.FindAllStringSubmatch(InnerTableContainerTMP1, -1)

	// Preparing vars
	var (
		CreatureName, CreatureImageURL, CreatureDescription, CreatureBehaviour                                                    string
		CreatureLootList, CreatureImmuneTo, CreatureStrongAgainst, CreatureWeaknessAgainst, CreatureHealedBy                      []string
		CreatureHitpoints, CreatureSummonedMana, CreatureConvincedMana, CreatureExperiencePoints                                  int
		CreatureBeParalysed, CreatureBeSummoned, CreatureBeConvinced, CreatureSeeInvisible, CreatureIsLootable, CreatureIsBoosted bool
	)

	// Find boosted creature
	boostedMonsterTitle, boostedCreatureFound := ReaderHTML.Find("#Monster").First().Attr("title")

	if boostedCreatureFound {
		boostedCreatureRace := boostedMonsterTitle[strings.Index(boostedMonsterTitle, ": ")+2:]

		CreatureIsBoosted = boostedCreatureRace == race
	}

	// Preparing data for JSONData
	if len(subma1) > 0 {
		CreatureName = TibiaDataSanitizeEscapedString(subma1[0][1])
		CreatureImageURL = subma1[0][2]

		// Description (and unescape hmtl string)
		CreatureDescription = strings.ReplaceAll(subma1[0][3], "<br/>", "\n")
		CreatureDescription = TibiaDataSanitizeEscapedString(CreatureDescription)

		// Behaviour
		// Regex to get data..
		subma2 := CreatureHitpointsRegex.FindAllStringSubmatch(subma1[0][4], -1)
		// Add data to vars
		CreatureHitpoints = TibiaDataStringToInteger(subma2[0][1])
		CreatureBehaviour = subma2[0][2]
		if !strings.Contains(subma1[0][4], "cannot be paralysed") {
			CreatureBeParalysed = true
		}
		if strings.Contains(subma1[0][4], "sense invisible creatures") {
			CreatureSeeInvisible = true
		}
		if strings.Contains(subma1[0][4], " are immune to ") {
			subma21 := CreatureImmuneRegex.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureImmuneToTmp := strings.Split(subma21[0][1], localDamageString)
			CreatureImmuneTo = strings.Split(strings.Replace(CreatureImmuneToTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], " are strong against ") {
			subma22 := CreatureStrongRegex.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureStrongAgainstTmp := strings.Split(subma22[0][1], localDamageString)
			CreatureStrongAgainst = strings.Split(strings.Replace(CreatureStrongAgainstTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], " are weak against ") {
			subma23 := CreatureWeakRegex.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureWeaknessAgainstTmp := strings.Split(subma23[0][1], localDamageString)
			CreatureWeaknessAgainst = strings.Split(strings.Replace(CreatureWeaknessAgainstTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], " are healed when ") {
			subma23 := CreatureHealedRegex.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureHealedByTmp := strings.Split(subma23[0][1], localDamageString)
			CreatureHealedBy = strings.Split(strings.Replace(CreatureHealedByTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], "It takes ") && strings.Contains(subma1[0][4], " mana to ") {
			subma24 := CreatureManaRequiredRegex.FindAllStringSubmatch(subma1[0][4], -1)
			subma2402 := subma24[0][2]
			if strings.Contains(subma2402, "convince these creatures but they cannot be") {
				CreatureBeConvinced = true
				CreatureConvincedMana = TibiaDataStringToInteger(subma24[0][1])
			} else if strings.Contains(subma2402, "summon or convince these creatures") {
				CreatureBeSummoned = true
				CreatureSummonedMana = TibiaDataStringToInteger(subma24[0][1])
				CreatureBeConvinced = true
				CreatureConvincedMana = TibiaDataStringToInteger(subma24[0][1])
			} else if strings.Contains(subma2402, "summon these creatures but they cannot be") {
				CreatureBeSummoned = true
				CreatureSummonedMana = TibiaDataStringToInteger(subma24[0][1])
			}
		}

		// Loot
		// Regex to get loot information
		subma3 := CreatureLootRegex.FindAllStringSubmatch(subma1[0][5], -1)
		// Adding data to vars
		CreatureExperiencePoints = TibiaDataStringToInteger(subma3[0][1])
		if subma3[0][2] != "nothing" {
			CreatureIsLootable = true
			CreatureLootListTmp := strings.SplitSeq(strings.Replace(strings.Replace(subma3[0][2], "items ", "", 1), " and sometimes other ", "", 1), ", ")
			for str := range CreatureLootListTmp {
				if str != "" {
					CreatureLootList = append(CreatureLootList, str)
				}
			}
		}
	} else {
		log.Printf("[warning] TibiaCreaturesCreatureImpl called on invalid creature")
		return CreatureResponse{}, validation.ErrorCreatureNotFound
	}

	// Build the data-blob
	return CreatureResponse{
		Creature{
			Name:             CreatureName,
			Race:             race,
			ImageURL:         CreatureImageURL,
			Description:      CreatureDescription,
			Behaviour:        CreatureBehaviour,
			Hitpoints:        CreatureHitpoints,
			ImmuneTo:         CreatureImmuneTo,
			StrongAgainst:    CreatureStrongAgainst,
			WeaknessAgainst:  CreatureWeaknessAgainst,
			HealedBy:         CreatureHealedBy,
			BeParalysed:      CreatureBeParalysed,
			BeSummoned:       CreatureBeSummoned,
			SummonMana:       CreatureSummonedMana,
			BeConvinced:      CreatureBeConvinced,
			ConvincedMana:    CreatureConvincedMana,
			SeeInvisible:     CreatureSeeInvisible,
			ExperiencePoints: CreatureExperiencePoints,
			IsLootable:       CreatureIsLootable,
			LootList:         CreatureLootList,
			Featured:         CreatureIsBoosted,
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
