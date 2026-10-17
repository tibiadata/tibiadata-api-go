package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of Spells
type Spell struct {
	Name         string `json:"name"`          // The name of the spell.
	Spell        string `json:"spell_id"`      // The internal identifier of the spell.
	Formula      string `json:"formula"`       // The formula to cast the spell.
	Level        int    `json:"level"`         // The required level for casting.
	Mana         int    `json:"mana"`          // The required mana for using.
	Price        int    `json:"price"`         // The price in gold coins to learn it.
	GroupAttack  bool   `json:"group_attack"`  // Whether the group is attack.
	GroupHealing bool   `json:"group_healing"` // Whether the group is healing.
	GroupSupport bool   `json:"group_support"` // Whether the group is support.
	TypeInstant  bool   `json:"type_instant"`  // Whether the type is instant.
	TypeRune     bool   `json:"type_rune"`     // Whether the type is rune.
	PremiumOnly  bool   `json:"premium_only"`  // Whether it requires to have premium account to learn and use it.
}

// Child of JSONData
type Spells struct {
	SpellsVocationFilter string  `json:"spells_filter"` // The applied filters on the list
	Spells               []Spell `json:"spell_list"`    // List of spells
}

// The base includes two levels: Spells and Information
type SpellsOverviewResponse struct {
	Spells      Spells      `json:"spells"`
	Information Information `json:"information"`
}

// TibiaSpellsOverview func
func TibiaSpellsOverviewImpl(vocationName, BoxContentHTML, url string) (SpellsOverviewResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return SpellsOverviewResponse{}, fmt.Errorf("[error] TibiaSpellsOverviewImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var SpellsData []Spell

	// Running query over each div
	ReaderHTML.Find(".Table3 table.TableContent tr").Each(func(index int, s *goquery.Selection) {
		// Skip header row
		if index == 0 {
			return
		}

		spellBuilder := Spell{}

		s.Find("td").Each(func(index int, s2 *goquery.Selection) {
			selectionText := s2.Text()
			selectionHtml, _ := s2.Html()

			switch index {
			case 0:
				spellBuilder.Name = selectionText[0:strings.Index(selectionText, " (")]
				spellBuilder.Spell = selectionHtml[strings.Index(selectionHtml, "amp;spell=")+10 : strings.Index(selectionHtml, "&amp;vocation")]
				spellBuilder.Formula = selectionText[strings.Index(selectionText, " (")+2 : strings.Index(selectionText, ")")]
			case 1:
				switch selectionText {
				case "Attack":
					spellBuilder.GroupAttack = true
				case "Healing":
					spellBuilder.GroupHealing = true
				case "Support":
					spellBuilder.GroupSupport = true
				}
			case 2:
				switch selectionText {
				case "Instant":
					spellBuilder.TypeInstant = true
				case "Rune":
					spellBuilder.TypeRune = true
				}
			case 3:
				if selectionText != "-" {
					spellBuilder.Level = TibiaDataStringToInteger(selectionText)
				}
			case 4:
				mana := -1
				if selectionText != "var." {
					mana = TibiaDataStringToInteger(selectionText)
				}

				spellBuilder.Mana = mana
			case 5:
				price := 0
				if selectionText != "free" {
					price = TibiaDataStringToInteger(selectionText)
				}

				spellBuilder.Price = price
			case 6:
				if selectionText == "yes" {
					spellBuilder.PremiumOnly = true
				}
			}
		})

		SpellsData = append(SpellsData, spellBuilder)
	})

	// adding readable SpellsVocationFilter field
	if vocationName == "" {
		vocationName = "all"
	}

	// Build the data-blob
	return SpellsOverviewResponse{
		Spells{
			SpellsVocationFilter: vocationName,
			Spells:               SpellsData,
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
