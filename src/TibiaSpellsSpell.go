package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of SpellInformation
type SpellInformation struct {
	Formula       string   `json:"formula"`        // The formula to cast the spell.
	Vocation      []string `json:"vocation"`       // The vocations that can use this spell.
	GroupAttack   bool     `json:"group_attack"`   // Whether the group is attack.
	GroupHealing  bool     `json:"group_healing"`  // Whether the group is healing.
	GroupSupport  bool     `json:"group_support"`  // Whether the group is support.
	TypeInstant   bool     `json:"type_instant"`   // Whether the type is instant.
	TypeRune      bool     `json:"type_rune"`      // Whether the type is rune.
	DamageType    string   `json:"damage_type"`    // The type of damage caused by it.
	CooldownAlone int      `json:"cooldown_alone"` // The individual cooldown of this spell in seconds.
	CooldownGroup int      `json:"cooldown_group"` // The group cooldown of this spell in seconds.
	SoulPoints    int      `json:"soul_points"`    // The number of soul points consumed when casting.
	Amount        int      `json:"amount"`         // The amount of objects created when casting.
	Level         int      `json:"level"`          // The required level for casting.
	Mana          int      `json:"mana"`           // The required mana for using.
	Price         int      `json:"price"`          // The price in gold coins to learn it.
	City          []string `json:"city"`           // The cities where to learn it.
	Premium       bool     `json:"premium_only"`   // Whether it requires a premium account to learn and use it.
}

// Child of RuneInformation
type RuneInformation struct {
	Vocation     []string `json:"vocation"`      // List of vocations that can use the rune.
	GroupAttack  bool     `json:"group_attack"`  // Whether the group is attack.
	GroupHealing bool     `json:"group_healing"` // Whether the group is healing.
	GroupSupport bool     `json:"group_support"` // Whether the group is support.
	DamageType   string   `json:"damage_type"`   // The type of damage caused by it.
	Level        int      `json:"level"`         // The required level for using.
	MagicLevel   int      `json:"magic_level"`   // The required magic level for using.
}

// Child of JSONData
type SpellData struct {
	Name                string           `json:"name"`                  // The name of the spell.
	Spell               string           `json:"spell_id"`              // The internal identifier of the spell.
	ImageURL            string           `json:"image_url"`             // The URL to this spell's image.
	Description         string           `json:"description"`           // A description of it's effect and history.
	HasSpellInformation bool             `json:"has_spell_information"` // Whether the spell has information.
	SpellInformation    SpellInformation `json:"spell_information"`     // Information about the spell.
	HasRuneInformation  bool             `json:"has_rune_information"`  // Whether the spell has rune information.
	RuneInformation     RuneInformation  `json:"rune_information"`      // Information about the spell's rune.
}

// The base includes two levels: Spell and Information
type SpellInformationResponse struct {
	Spell       SpellData   `json:"spell"`
	Information Information `json:"information"`
}

var (
	SpellDataRowRegex      = regexp.MustCompile(`<td.*>(.*):<\/td><td.*>(.*)<\/td>`)
	SpellNameAndImageRegex = regexp.MustCompile(`<td><img src="(.*)" width=.*<h2>(.*)<\/h2>.*`)
	SpellCooldownRegex     = regexp.MustCompile(`([0-9]+)s \(.*:.([0-9]+)s\)`)
	SpellDescriptionRegex  = regexp.MustCompile(`(.*)\.(Spell|Rune) InformationName:.*`)
)

// TibiaSpellsSpell func
func TibiaSpellsSpellImpl(spell string, BoxContentHTML string, url string) (SpellInformationResponse, error) {
	// TODO: There is currently a bug with description, it always comes back empty

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return SpellInformationResponse{}, fmt.Errorf("[error] TibiaSpellsSpellImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var (
		// creating empty vars for later use
		SpellsInfoVocation, SpellsInfoCity, RuneInfoVocation []string
		// var SpellsInfoName, RuneInfoName string
		SpellInformationSection, SpellName, SpellID, SpellImageURL, SpellDescription, SpellsInfoFormula, SpellsInfoDamageType, RuneInfoDamageType                                                                                                         string
		SpellsInfoCooldownAlone, SpellsInfoCooldownGroup, SpellsInfoSoulPoints, SpellsInfoAmount, SpellsInfoLevel, SpellsInfoMana, SpellsInfoPrice, RuneInfoLevel, RuneInfoMagicLevel                                                                     int
		SpellsInfoGroupAttack, SpellsInfoGroupHealing, SpellsInfoGroupSupport, SpellsInfoTypeInstant, SpellsInfoTypeRune, RuneInfoGroupAttack, RuneInfoGroupHealing, RuneInfoGroupSupport, SpellsInfoPremium, SpellsHasSpellSection, SpellsHasRuneSection bool

		insideError error
	)

	ReaderHTML.Find(".BoxContent").EachWithBreak(func(index int, s *goquery.Selection) bool {
		NameAndImageSection, err := s.Find("table tr").First().Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaSpellsSpellImpl failed at NameAndImageSection, err := s.Find, err: %s", err)
			return false
		}

		// Get the name, spell_id and image
		subma2 := SpellNameAndImageRegex.FindAllStringSubmatch(NameAndImageSection, -1)
		if len(subma2) > 0 {
			SpellName = TibiaDataSanitize0026String(subma2[0][2])
			SpellImageURL = subma2[0][1]
			SpellID = SpellImageURL[strings.Index(SpellImageURL, "library/")+8 : strings.Index(SpellImageURL, ".png")]
		}

		s.Find(".TableContainer").Each(func(index int, s *goquery.Selection) {
			SectionName := s.Find(".CaptionInnerContainer div.Text").Text()

			// Determine if this is the spell or rune section
			switch SectionName {
			case "Spell Information":
				SpellInformationSection = "spell"
				SpellsHasSpellSection = true
			case "Rune Information":
				SpellInformationSection = "rune"
				SpellsHasRuneSection = true
			}

			// Running query over each div
			s.Find("table.Table2 .TableContentContainer tbody tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into SpellDivHTML
				SpellDivHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaSpellsSpellImpl failed at SpellDivHTML, err := s.Html(), err: %s", err)
					return false
				}

				subma1 := SpellDataRowRegex.FindAllStringSubmatch(SpellDivHTML, -1)

				// check if regex return length is over 0 and the match of name is over 1
				if len(subma1) > 0 {
					// Creating easy to use vars (and unescape hmtl right string)
					LeftColumn := subma1[0][1]
					RightColumn := TibiaDataSanitizeEscapedString(subma1[0][2])

					// Formula
					if LeftColumn == "Formula" {
						SpellsInfoFormula = TibiaDataSanitizeDoubleQuoteString(RightColumn)
					}

					// Vocation
					if LeftColumn == "Vocation" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoVocation = strings.Split(RightColumn, ", ")
						case "rune":
							RuneInfoVocation = strings.Split(RightColumn, ", ")
						}
					}

					// Group information
					if LeftColumn == "Group" {
						switch SpellInformationSection {
						case "spell":
							switch RightColumn {
							case "Attack":
								SpellsInfoGroupAttack = true
							case "Healing":
								SpellsInfoGroupHealing = true
							case "Support":
								SpellsInfoGroupSupport = true
							}
						case "rune":
							switch RightColumn {
							case "Attack":
								RuneInfoGroupAttack = true
							case "Healing":
								RuneInfoGroupHealing = true
							case "Support":
								RuneInfoGroupSupport = true
							}
						}
					}

					// Spell type
					if LeftColumn == "Type" {
						switch RightColumn {
						case "Instant":
							SpellsInfoTypeInstant = true
						case "Rune":
							SpellsInfoTypeRune = true
						}
					}

					// Damage
					if LeftColumn == "Damage Type" || LeftColumn == "Magic Type" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoDamageType = strings.ToLower(RightColumn)
						case "rune":
							RuneInfoDamageType = strings.ToLower(RightColumn)
						}
					}

					// Cooldown
					if LeftColumn == "Cooldown" {
						subma3 := SpellCooldownRegex.FindAllStringSubmatch(SpellDivHTML, -1)
						if len(subma3) > 0 {
							SpellsInfoCooldownAlone = TibiaDataStringToInteger(subma3[0][1])
							SpellsInfoCooldownGroup = TibiaDataStringToInteger(subma3[0][2])
						}
					}

					// Soul Points
					if LeftColumn == "Soul Points" {
						SpellsInfoSoulPoints = TibiaDataStringToInteger(RightColumn)
					}

					// Amount
					if LeftColumn == "Amount" {
						SpellsInfoAmount = TibiaDataStringToInteger(RightColumn)
					}

					// Experience Level
					if LeftColumn == "Exp Lvl" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoLevel = TibiaDataStringToInteger(RightColumn)
						case "rune":
							RuneInfoLevel = TibiaDataStringToInteger(RightColumn)
						}
					}

					// Mana
					if LeftColumn == "Mana" {
						SpellsInfoMana = TibiaDataStringToInteger(RightColumn)
					}

					// Price
					if LeftColumn == "Price" {
						if RightColumn == "free" {
							SpellsInfoPrice = 0
						} else {
							SpellsInfoPrice = TibiaDataStringToInteger(RightColumn)
						}
					}

					// City
					if LeftColumn == "City" {
						SpellsInfoCity = strings.Split(RightColumn, ", ")
					}

					// Premium
					if LeftColumn == "Premium" {
						if RightColumn == "yes" {
							SpellsInfoPremium = true
						} else {
							SpellsInfoPremium = false
						}
					}

					// Magic level
					if LeftColumn == "Mag Lvl" {
						RuneInfoMagicLevel = TibiaDataStringToInteger(RightColumn)
					}
				}

				return true
			})
		})

		return true
	})

	if insideError != nil {
		return SpellInformationResponse{}, insideError
	}

	// Getting the description
	InnerTableContainerTMPB := ReaderHTML.Find(".BoxContent").Text()
	subma4 := SpellDescriptionRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)
	if len(subma4) > 0 {
		SpellDescription = subma4[0][1] + "."
	}

	//
	// Build the data-blob
	return SpellInformationResponse{
		SpellData{
			Name:                SpellName,
			Spell:               SpellID,
			ImageURL:            SpellImageURL,
			Description:         SpellDescription,
			HasSpellInformation: SpellsHasSpellSection,
			SpellInformation: SpellInformation{
				Formula:       SpellsInfoFormula,
				Vocation:      SpellsInfoVocation,
				GroupAttack:   SpellsInfoGroupAttack,
				GroupHealing:  SpellsInfoGroupHealing,
				GroupSupport:  SpellsInfoGroupSupport,
				TypeInstant:   SpellsInfoTypeInstant,
				TypeRune:      SpellsInfoTypeRune,
				DamageType:    SpellsInfoDamageType,
				CooldownAlone: SpellsInfoCooldownAlone,
				CooldownGroup: SpellsInfoCooldownGroup,
				SoulPoints:    SpellsInfoSoulPoints,
				Amount:        SpellsInfoAmount,
				Level:         SpellsInfoLevel,
				Mana:          SpellsInfoMana,
				Price:         SpellsInfoPrice,
				City:          SpellsInfoCity,
				Premium:       SpellsInfoPremium,
			},
			HasRuneInformation: SpellsHasRuneSection,
			RuneInformation: RuneInformation{
				Vocation:     RuneInfoVocation,
				GroupAttack:  RuneInfoGroupAttack,
				GroupHealing: RuneInfoGroupHealing,
				GroupSupport: RuneInfoGroupSupport,
				DamageType:   RuneInfoDamageType,
				Level:        RuneInfoLevel,
				MagicLevel:   RuneInfoMagicLevel,
			},
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
