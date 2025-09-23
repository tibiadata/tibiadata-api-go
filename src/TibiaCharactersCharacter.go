package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
	//"time"
)

// Child of CharacterInfo
type Houses struct {
	Name    string `json:"name"`    // The name of the house.
	Town    string `json:"town"`    // The town where the house is located in.
	Paid    string `json:"paid"`    // The date the last paid rent is due.
	HouseID int    `json:"houseid"` // The internal ID of the house.
}

// Child of CharacterInfo
type CharacterGuild struct {
	GuildName string `json:"name,omitempty"` // The name of the guild.
	Rank      string `json:"rank,omitempty"` // The character's rank in the guild.
}

// Child of Character
type CharacterInfo struct {
	Name              string         `json:"name"`                    // The name of the character.
	FormerNames       []string       `json:"former_names,omitempty"`  // List of former names of the character.
	Traded            bool           `json:"traded,omitempty"`        // Whether the character was traded. (last 6 months)
	DeletionDate      string         `json:"deletion_date,omitempty"` // The date when the character will be deleted. (if scheduled for deletion)
	Sex               string         `json:"sex"`                     // The character's sex.
	Title             string         `json:"title"`                   // The character's selected title.
	UnlockedTitles    int            `json:"unlocked_titles"`         // The number of titles the character has unlocked.
	Vocation          string         `json:"vocation"`                // The character's vocation.
	Level             int            `json:"level"`                   // The character's level.
	AchievementPoints int            `json:"achievement_points"`      // The total of achievement points the character has.
	World             string         `json:"world"`                   // The character's current world.
	FormerWorlds      []string       `json:"former_worlds,omitempty"` // List of former worlds the character was in. (last 6 months)
	Residence         string         `json:"residence"`               // The character's current residence.
	MarriedTo         string         `json:"married_to,omitempty"`    // The name of the character's husband/spouse.
	Houses            []Houses       `json:"houses,omitempty"`        // List of houses the character owns currently.
	Guild             CharacterGuild `json:"guild"`                   // The guild that the character is member of.
	LastLogin         string         `json:"last_login,omitempty"`    // The character's last logged in time.
	Position          string         `json:"position,omitempty"`      // The character's special position.
	AccountStatus     string         `json:"account_status"`          // Whether account is Free or Premium.
	Comment           string         `json:"comment,omitempty"`       // The character's comment.
}

// Child of Character
type AccountBadges struct {
	Name        string `json:"name"`        // The name of the badge.
	IconURL     string `json:"icon_url"`    // The URL to the badge's icon.
	Description string `json:"description"` // The description of the badge.
}

// Child of Character
type Achievements struct {
	Name   string `json:"name"`   // The name of the achievement.
	Grade  int    `json:"grade"`  // The grade/stars of the achievement.
	Secret bool   `json:"secret"` // Whether it is a secret achievement or not.
}

// Child of Deaths
type Killers struct {
	Name   string `json:"name"`   // The name of the killer/assist.
	Player bool   `json:"player"` // Whether it is a player or not.
	Traded bool   `json:"traded"` // If the killer/assist was traded after the death.
	Summon string `json:"summon"` // The name of the summoned creature.
}

// Child of Character
type Deaths struct {
	Time    string    `json:"time"`            // The timestamp when the death occurred.
	Level   int       `json:"level,omitempty"` // The level when the death occurred.
	Killers []Killers `json:"killers"`         // List of killers involved.
	Assists []Killers `json:"assists"`         // List of assists involved.
	Reason  string    `json:"reason"`          // The plain text reason of death.
}

// Child of Character
type AccountInformation struct {
	Position     string `json:"position,omitempty"`      // The account's special position.
	Created      string `json:"created,omitempty"`       // The account's date of creation.
	LoyaltyTitle string `json:"loyalty_title,omitempty"` // The account's loyalty title.
}

// Child of Character
type OtherCharacters struct {
	Name     string `json:"name"`               // The name of the character.
	World    string `json:"world"`              // The name of the world.
	Status   string `json:"status"`             // The status of the character being online or offline.
	Deleted  bool   `json:"deleted"`            // Whether the character is scheduled for deletion or not.
	Main     bool   `json:"main"`               // Whether this is the main character or not.
	Traded   bool   `json:"traded"`             // Whether the character has been traded last 6 months or not.
	Position string `json:"position,omitempty"` // // The character's special position.
}

// Child of JSONData
type Character struct {
	CharacterInfo      CharacterInfo      `json:"character"`                     // The character's information.
	AccountBadges      []AccountBadges    `json:"account_badges,omitempty"`      // The account's badges.
	Achievements       []Achievements     `json:"achievements,omitempty"`        // The character's achievements.
	Deaths             []Deaths           `json:"deaths,omitempty"`              // The character's deaths.
	DeathsTruncated    bool               `json:"deaths_truncated,omitempty"`    // Whether the character's deaths were truncated or not.
	AccountInformation AccountInformation `json:"account_information,omitempty"` // The account information.
	OtherCharacters    []OtherCharacters  `json:"other_characters,omitempty"`    // The account's other characters.
}

// The base includes two levels, Characters and Information
type CharacterResponse struct {
	Character   Character   `json:"character"`
	Information Information `json:"information"`
}

// From https://pkg.go.dev/golang.org/x/net/html/atom
// This is an Atom. An Atom is an integer code for a string.
// Instead of importing the whole lib, we thought it would be
// best to just simply use the Br constant value.
const Br = 0x202

// TibiaCharactersCharacter func
func TibiaCharactersCharacterImpl(BoxContentHTML string, url string) (CharacterResponse, error) {
	var (
		// local strings used in this function
		localDivQueryString = ".TableContentContainer tr"
		localTradedString   = " (traded)"

		// Declaring vars for later use..
		CharacterInfoData      CharacterInfo
		AccountBadgesData      []AccountBadges
		AchievementsData       []Achievements
		DeathsData             []Deaths
		AccountInformationData AccountInformation
		OtherCharactersData    []OtherCharacters
		DeathsTruncated        bool

		// Errors
		characterNotFound bool
		insideError       error
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return CharacterResponse{}, fmt.Errorf("TibiaCharactersCharacterImpl failed at  goquery.NewDocumentFromReader, err: %s", err)
	}

	// Running query on every .TableContainer
	ReaderHTML.Find(".TableContainer").EachWithBreak(func(index int, s *goquery.Selection) bool {
		if insideError != nil {
			return false
		}

		SectionTextQuery := s.Find("div.Text")

		SectionName := SectionTextQuery.Nodes[0].FirstChild.Data

		// query current node with goquery
		CharacterDivQuery := goquery.NewDocumentFromNode(s.Nodes[0])

		switch SectionName {
		case "Could not find character":
			characterNotFound = true
			return false
		case "Character Information", "Account Information":
			// Running query over each tr in character content container
			CharacterDivQuery.Find(localDivQueryString).Each(func(index int, s *goquery.Selection) {
				RowNameQuery := s.Find("td[class^='Label']")

				RowName := RowNameQuery.Nodes[0].FirstChild.Data
				RowData := RowNameQuery.Nodes[0].NextSibling.FirstChild.Data

				switch TibiaDataSanitizeStrings(RowName) {
				case "Name:":
					Tmp := strings.Split(RowData, "<")
					CharacterInfoData.Name = strings.TrimSpace(Tmp[0])
					if strings.Contains(Tmp[0], ", will be deleted at") {
						Tmp2 := strings.Split(Tmp[0], ", will be deleted at ")
						CharacterInfoData.Name = Tmp2[0]
						CharacterInfoData.DeletionDate = TibiaDataDatetime(strings.TrimSpace(Tmp2[1]))
					}
					if strings.Contains(RowData, localTradedString) {
						CharacterInfoData.Traded = true
						CharacterInfoData.Name = strings.Replace(CharacterInfoData.Name, localTradedString, "", -1)
					}
				case "Former Names:":
					CharacterInfoData.FormerNames = strings.Split(RowData, ", ")
				case "Sex:":
					CharacterInfoData.Sex = RowData
				case "Title:":
					leftParenIdx := strings.LastIndex(RowData, "(")
					if leftParenIdx == -1 {
						return
					}

					title := RowData[:leftParenIdx-1]

					spaceIdx := strings.Index(RowData[leftParenIdx:], " ")
					if spaceIdx == -1 {
						return
					}

					unlockedTitles := TibiaDataStringToInteger(
						RowData[leftParenIdx+1 : leftParenIdx+spaceIdx],
					)

					CharacterInfoData.Title = title
					CharacterInfoData.UnlockedTitles = unlockedTitles
				case "Vocation:":
					CharacterInfoData.Vocation = RowData
				case "Level:":
					CharacterInfoData.Level = TibiaDataStringToInteger(RowData)
				case "nobr", "Achievement Points:":
					CharacterInfoData.AchievementPoints = TibiaDataStringToInteger(RowData)
				case "World:":
					CharacterInfoData.World = RowData
				case "Former World:":
					CharacterInfoData.FormerWorlds = strings.Split(RowData, ", ")
				case "Residence:":
					CharacterInfoData.Residence = RowData
				case "Account Status:":
					CharacterInfoData.AccountStatus = RowData
				case "Married To:":
					AnchorQuery := s.Find("a")
					CharacterInfoData.MarriedTo = AnchorQuery.Nodes[0].FirstChild.Data
				case "House:":
					AnchorQuery := s.Find("a")
					HouseName := AnchorQuery.Nodes[0].FirstChild.Data
					HouseHref := AnchorQuery.Nodes[0].Attr[0].Val
					// substring from houseid= to &character in the href for the house
					HouseId := HouseHref[strings.Index(HouseHref, "houseid")+8 : strings.Index(HouseHref, "&character")]
					HouseRawData := RowNameQuery.Nodes[0].NextSibling.LastChild.Data
					HouseTown := HouseRawData[strings.Index(HouseRawData, "(")+1 : strings.Index(HouseRawData, ")")]
					HousePaidUntil := HouseRawData[strings.Index(HouseRawData, "is paid until ")+14:]

					CharacterInfoData.Houses = append(CharacterInfoData.Houses, Houses{
						Name:    HouseName,
						Town:    HouseTown,
						Paid:    TibiaDataDate(HousePaidUntil),
						HouseID: TibiaDataStringToInteger(HouseId),
					})
				case "Guild Membership:":
					CharacterInfoData.Guild.Rank = strings.TrimSuffix(RowData, " of the ")

					// TODO: I don't understand why the unicode nbsp is there...
					CharacterInfoData.Guild.GuildName = TibiaDataSanitizeStrings(RowNameQuery.Nodes[0].NextSibling.LastChild.LastChild.Data)
				case "Last Login:":
					if strings.ToLower(RowData) != "never logged in" {
						CharacterInfoData.LastLogin = TibiaDataDatetime(RowData)
					}
				case "Comment:":
					node := RowNameQuery.Nodes[0].NextSibling.FirstChild

					stringBuilder := strings.Builder{}
					for node != nil {
						if node.DataAtom == Br {
							// It appears we can ignore br because either the encoding or goquery adds an \n for us
							// stringBuilder.WriteString("\n")
						} else {
							stringBuilder.WriteString(node.Data)
						}

						node = node.NextSibling
					}

					CharacterInfoData.Comment = stringBuilder.String()
				case "Loyalty Title:":
					if RowData != "(no title)" {
						AccountInformationData.LoyaltyTitle = RowData
					}
				case "Created:":
					AccountInformationData.Created = TibiaDataDatetime(RowData)
				case "Position:":
					TmpPosition := strings.Split(RowData, "<")
					switch SectionName {
					case "Character Information":
						CharacterInfoData.Position = strings.TrimSpace(TmpPosition[0])
					case "Account Information":
						AccountInformationData.Position = strings.TrimSpace(TmpPosition[0])
					}

				default:
					log.Println("LEFT OVER: `" + RowName + "` = `" + RowData + "`")
				}
			})
		case "Account Badges":
			// Running query over each tr in list
			CharacterDivQuery.Find(".TableContentContainer tr td span[style]").EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Account Badges, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaks(CharacterListHTML)

				// prevent failure of regex that parses account badges
				if CharacterListHTML != "There are no account badges set to be displayed for this character." {
					const (
						nameIndexer = `alt="`
						iconIndexer = `img src="`
						descIndexer = `&#39;, &#39;`
					)

					nameIdx := strings.Index(
						CharacterListHTML, nameIndexer,
					) + len(nameIndexer)
					endNameIdx := strings.Index(
						CharacterListHTML[nameIdx:], `"`,
					) + nameIdx

					iconIdx := strings.Index(
						CharacterListHTML, iconIndexer,
					) + len(iconIndexer)
					endIconIdx := strings.Index(
						CharacterListHTML[iconIdx:], `"`,
					) + iconIdx

					descIdx := strings.Index(
						CharacterListHTML, descIndexer,
					) + len(descIndexer)
					endDescIdx := strings.Index(
						CharacterListHTML[descIdx:], descIndexer,
					) + descIdx

					AccountBadgesData = append(AccountBadgesData, AccountBadges{
						Name:        CharacterListHTML[nameIdx:endNameIdx],
						IconURL:     CharacterListHTML[iconIdx:endIconIdx],
						Description: CharacterListHTML[descIdx:endDescIdx],
					})
				}

				return true
			})
		case "Account Achievements":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Account Achievements, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaks(CharacterListHTML)

				if !strings.Contains(CharacterListHTML, "There are no achievements set to be displayed for this character.") {
					const (
						nameIndexer = `alt="Tibia Achievement"/></td><td>`
					)

					// get the name of the achievement (and ignore the secret image on the right)
					nameIdx := strings.Index(
						CharacterListHTML, nameIndexer,
					) + len(nameIndexer)
					endNameIdx := strings.Index(
						CharacterListHTML[nameIdx:], `<`,
					) + nameIdx

					AchievementsData = append(AchievementsData, Achievements{
						Name:   CharacterListHTML[nameIdx:endNameIdx],
						Grade:  strings.Count(CharacterListHTML, "achievement-grade-symbol"),
						Secret: strings.Contains(CharacterListHTML, "achievement-secret-symbol"),
					})
				}

				return true
			})
		case "Character Deaths", "Character Deaths (Last 30 Days)":
			// Running query over each tr in list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Character Deaths, err: %s", err)
					return false
				}

				// Removing line breaks and sanitizing string
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaks(CharacterListHTML)
				CharacterListHTML = strings.ReplaceAll(strings.ReplaceAll(CharacterListHTML, "<br/>", " "), "  ", " ")
				CharacterListHTML = TibiaDataSanitizeStrings(CharacterListHTML)

				dataNoTags := RemoveHtmlTag(CharacterListHTML)

				if strings.HasPrefix(dataNoTags, "There happened more character deaths in the last 30 days than we can display here.") {
					DeathsTruncated = true
					return false
				}

				// defining responses
				DeathKillers := []Killers{}
				DeathAssists := []Killers{}

				const (
					initCetIndexer  = `CET`
					initCestIndexer = `CEST`
					levelIndexer    = `at Level `
					killersIndexer  = `by `
					assistedIndexer = `Assisted by `
				)

				var initIndexer string
				timeZoneIdentifiers := []string{initCetIndexer, initCestIndexer}
				for _, tz := range timeZoneIdentifiers {
					if strings.Contains(dataNoTags, tz) {
						initIndexer = tz
						break
					}
				}

				initIdx := strings.Index(
					dataNoTags, initIndexer,
				) + len(initIndexer)
				endInitIdx := strings.Index(
					dataNoTags[initIdx:], `by `,
				) + initIdx + len(`by `)

				reasonStart := dataNoTags[initIdx:endInitIdx]
				reasonRest := dataNoTags[endInitIdx:]

				// store for reply later on.. and sanitizing string
				reasonString := reasonStart + reasonRest

				timeIdx := 0
				endTimeIdx := strings.Index(
					dataNoTags[timeIdx:], initIndexer,
				) + timeIdx + len(initIndexer)

				time := TibiaDataDatetime(dataNoTags[timeIdx:endTimeIdx])

				levelIdx := strings.Index(
					dataNoTags, levelIndexer,
				) + len(levelIndexer)
				endLevelIdx := strings.Index(
					dataNoTags[levelIdx:], " ",
				) + levelIdx

				level := TibiaDataStringToInteger(dataNoTags[levelIdx:endLevelIdx])

				// if kill is with assist only (and level is set to 25), then we reset level
				if reasonStart == assistedIndexer && level == 25 {
					level = 0
				}

				killersIdx := strings.Index(
					CharacterListHTML, killersIndexer,
				) + len(killersIndexer)
				endKillersIdx := strings.Index(
					CharacterListHTML[killersIdx:], `.</td>`,
				) + killersIdx

				rawListofKillers := CharacterListHTML[killersIdx:endKillersIdx]

				// A death with assist (or with assist only)
				if strings.Contains(dataNoTags, assistedIndexer) {
					TmpListOfDeath := strings.Split(CharacterListHTML, assistedIndexer)

					// get a list of killers
					if reasonStart != assistedIndexer {
						rawListofKillers = TmpListOfDeath[0][killersIdx:]
						rawListofKillers = strings.TrimSpace(strings.TrimSuffix(rawListofKillers, "."))
					} else {
						rawListofKillers = ""
					}
					TmpAssist := TmpListOfDeath[1]

					// get a list of killers
					ListOfAssists := strings.Split(TmpAssist, ", ")

					// extract if "and" is in last ss1
					ListOfAssistsTmp := strings.Split(ListOfAssists[len(ListOfAssists)-1], " and ")

					// if there is an "and", then we split it..
					if len(ListOfAssistsTmp) > 1 {
						ListOfAssists[len(ListOfAssists)-1] = ListOfAssistsTmp[0]
						ListOfAssists = append(ListOfAssists, ListOfAssistsTmp[1])
					}

					for i := range ListOfAssists {
						name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfAssists[i])
						if name != "" { // Ensure we don't append empty names
							DeathAssists = append(DeathAssists, Killers{
								Name:   name,
								Player: isPlayer,
								Traded: isTraded,
								Summon: theSummon,
							})
						}
					}
				}

				// A death with killers
				if rawListofKillers != "" {

					// get a list of killers
					ListOfKillers := strings.Split(rawListofKillers, ", ")

					const andStr = " and "
					lastItem := ListOfKillers[len(ListOfKillers)-1]
					lastAndIdx := strings.LastIndex(lastItem, andStr)

					if lastAndIdx > -1 {
						if !strings.Contains(lastItem, "<a href=") {
							ListOfKillers[len(ListOfKillers)-1] = lastItem[:lastAndIdx]
							ListOfKillers = append(ListOfKillers, lastItem[lastAndIdx+len(andStr):])
						} else {
							ListOfKillers = ListOfKillers[:len(ListOfKillers)-1]

							const (
								nonTag        = iota // outside of a tag
								openAnchorTag        // inside a <a>
								closeAchorTag        // inside a </a>
							)

							var (
								buffer bytes.Buffer
								state  = nonTag
							)
							buffer.Grow(200) // arbitrary number to avoid allocations

							for i := range lastItem {
								cur := lastItem[i]
								switch state {
								case nonTag:
									if cur == '<' {
										switch lastItem[i+1] {
										case '/':
											state = closeAchorTag
										default:
											state = openAnchorTag
											if buffer.Len() > 0 {
												str := buffer.String()

												str = strings.TrimPrefix(str, " and ")
												str = strings.TrimSuffix(str, " and ")

												if str == "" {
													buffer.Reset()
													buffer.WriteByte(cur)
													continue
												}

												if strings.Contains(str, "of") && !containsCreaturesWithOf(str) {
													// this is a summon
													buffer.WriteByte(cur)
													continue
												}

												buffer.Reset()
												ListOfKillers = append(ListOfKillers, str)
											}
										}
									}
									buffer.WriteByte(cur)

									if cur == ')' {
										str := buffer.String()

										buffer.Reset()
										ListOfKillers = append(ListOfKillers, str)
									}
								case openAnchorTag:
									if cur == '>' {
										state = nonTag
									}
									buffer.WriteByte(cur)
								case closeAchorTag:
									buffer.WriteByte(cur)
									if cur == '>' {
										str := buffer.String()

										str = strings.TrimPrefix(str, " and ")
										str = strings.TrimSuffix(str, " and ")

										ListOfKillers = append(ListOfKillers, str)
										buffer.Reset()
										state = nonTag
									}
								}
							}

							if buffer.Len() > 0 {
								str := buffer.String()
								buffer.Reset()

								str = strings.TrimPrefix(str, " and ")
								str = strings.TrimSuffix(str, " and ")

								if str != "" {
									ListOfKillers = append(ListOfKillers, str)
								}
							}
						}
					}

					// loop through all killers and append to result
					for i := range ListOfKillers {
						name, isPlayer, isTraded, theSummon := TibiaDataParseKiller(ListOfKillers[i])
						if name != "" { // Ensure we don't append empty names
							DeathKillers = append(DeathKillers, Killers{
								Name:   name,
								Player: isPlayer,
								Traded: isTraded,
								Summon: theSummon,
							})
						}
					}
				}

				// append deadentry to death list
				DeathsData = append(DeathsData, Deaths{
					Time:    time,
					Level:   level,
					Killers: DeathKillers,
					Assists: DeathAssists,
					Reason:  reasonString,
				})

				return true
			})
		case "Characters":
			// Running query over each tr in character list
			CharacterDivQuery.Find(localDivQueryString).EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into CharacterListHTML
				CharacterListHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaCharactersCharacterImpl failed at s.Html() inside Characters, err: %s", err)
					return false
				}

				// Removing line breaks
				CharacterListHTML = TibiaDataHTMLRemoveLinebreaks(CharacterListHTML)

				if !strings.Contains(CharacterListHTML, "<td>Name</td><td>World</td><td>Status</td>") {
					const (
						nameIndexer  = `<td style="width: 20%"><nobr>`
						worldIndexer = `<td style="width: 10%"><nobr>`
					)

					nameIdx := strings.Index(
						CharacterListHTML, nameIndexer,
					) + len(nameIndexer)
					nameIdx += strings.Index(
						CharacterListHTML[nameIdx:], " ",
					) + 1
					endNameIdx := strings.Index(
						CharacterListHTML[nameIdx:], `</nobr></td>`,
					) + nameIdx

					tmpCharName := CharacterListHTML[nameIdx:endNameIdx]

					worldIdx := strings.Index(
						CharacterListHTML, worldIndexer,
					) + len(worldIndexer)
					endWorldIdx := strings.Index(
						CharacterListHTML[worldIdx:], `</nobr></td>`,
					) + worldIdx

					world := CharacterListHTML[worldIdx:endWorldIdx]

					var tmpTraded bool
					if strings.Contains(tmpCharName, localTradedString) {
						tmpTraded = true
						tmpCharName = strings.ReplaceAll(tmpCharName, localTradedString, "")
					}

					// If this character is the main character of the account
					var tmpMain bool
					if strings.Contains(tmpCharName, "Main Character") {
						tmpMain = true
						tmp := strings.Split(tmpCharName, "<")
						tmpCharName = strings.TrimSpace(tmp[0])
					}

					// If this character is online or offline
					tmpStatus := "offline"
					if strings.Contains(CharacterListHTML, "<b class=\"green\">online</b>") {
						tmpStatus = "online"
					}

					// Is this character is deleted
					var tmpDeleted bool
					if strings.Contains(CharacterListHTML, "deleted") {
						tmpDeleted = true
					}

					// Is this character having a special position
					var tmpPosition string
					if strings.Contains(CharacterListHTML, "CipSoft Member") {
						tmpPosition = "CipSoft Member"
					}

					// Create the character and append it to the other characters list
					OtherCharactersData = append(OtherCharactersData, OtherCharacters{
						Name:     TibiaDataSanitizeStrings(tmpCharName),
						World:    world,
						Status:   tmpStatus,
						Deleted:  tmpDeleted,
						Main:     tmpMain,
						Traded:   tmpTraded,
						Position: tmpPosition,
					})
				}

				return true
			})
		}

		return true
	})

	// Build the character data
	charData := Character{
		CharacterInfoData,
		AccountBadgesData,
		AchievementsData,
		DeathsData,
		DeathsTruncated,
		AccountInformationData,
		OtherCharactersData,
	}

	// Search for errors
	switch {
	case characterNotFound:
		return CharacterResponse{}, validation.ErrorCharacterNotFound
	case insideError != nil:
		return CharacterResponse{}, insideError
	case reflect.DeepEqual(charData, Character{}):
		// There are some rare cases where a character name would
		// bug out tibia.com (tíbia, for example) and then we would't
		// receive the character not found error, for these edge cases
		// we check if the char structure is empty, if it is, it means
		// the character has not been found
		//
		// Validating those names would also be a pain because of old
		// tibian names such as Kolskägg, which for whatever reason is valid
		return CharacterResponse{}, validation.ErrorCharacterNotFound
	}

	//
	// Build the data-blob
	return CharacterResponse{
		charData,
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

// TibiaDataParseKiller func - insert a html string and get the killers back
func TibiaDataParseKiller(data string) (string, bool, bool, string) {
	var (
		// local strings used in this function
		localTradedString = " (traded)"

		isPlayer, isTraded bool
		theSummon          string
	)

	// check if killer is a traded player
	if strings.Contains(data, localTradedString) {
		isPlayer = true
		isTraded = true
		data = strings.ReplaceAll(data, localTradedString, "")
	}

	// check if killer is a player
	if strings.Contains(data, "https://www.tibia.com") {
		isPlayer = true
		data = RemoveHtmlTag(data)
	}

	// remove htlm, spaces, dots and prefixes from data-string
	data = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(data, "</td>"), "."))
	data = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(data, "and "), "a "), "an ")

	firstRune, _ := utf8.DecodeRuneInString(data)
	// get summon information
	if containsCreaturesWithOf(data) {
		// this is not a summon, since it is a creature with a of in the middle
	} else if unicode.IsLower(firstRune) {
		ofIdx := strings.Index(data, "of")
		if ofIdx != -1 {
			theSummon = data[:ofIdx-1]
			data = data[ofIdx+3:]
		}
	}

	// sanitizing string
	data = TibiaDataSanitizeStrings(data)

	return data, isPlayer, isTraded, theSummon
}

// containsCreaturesWithOf checks if creature is present in special creatures list
func containsCreaturesWithOf(str string) bool {
	// trim away "an " and "a "
	str = strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(str, "an "), "a "))

	switch str {
	case "acolyte of darkness",
		"acolyte of the cult",
		"adept of the cult",
		"ancient spawn of morgathla",
		"aspect of power",
		"baby pet of chayenne",
		"bane of light",
		"bloom of doom",
		"bride of night",
		"cloak of terror",
		"energuardian of tales",
		"enlightened of the cult",
		"eruption of destruction",
		"essence of darkness",
		"essence of malice",
		"eye of the seven",
		"flame of omrafir",
		"fury of the emperor",
		"ghastly pet of chayenne",
		"ghost of a planegazer",
		"greater splinter of madness",
		"groupie of skyrr",
		"guardian of tales",
		"gust of wind",
		"hand of cursed fate",
		"harbinger of darkness",
		"herald of gloom",
		"izcandar champion of summer",
		"izcandar champion of winter",
		"lesser splinter of madness",
		"lord of the elements",
		"lost ghost of a planegazer",
		"memory of a banshee",
		"memory of a book",
		"memory of a carnisylvan",
		"memory of a dwarf",
		"memory of a faun",
		"memory of a frazzlemaw",
		"memory of a fungus",
		"memory of a golem",
		"memory of a hero",
		"memory of a hydra",
		"memory of a lizard",
		"memory of a mammoth",
		"memory of a manticore",
		"memory of a pirate",
		"memory of a scarab",
		"memory of a shaper",
		"memory of a vampire",
		"memory of a werelion",
		"memory of a wolf",
		"memory of a yalahari",
		"memory of an amazon",
		"memory of an elf",
		"memory of an insectoid",
		"memory of an ogre",
		"mighty splinter of madness",
		"minion of gaz'haragoth",
		"minion of versperoth",
		"monk of the order",
		"muse of penciljack",
		"nightmare of gaz'haragoth",
		"noble pet of chayenne",
		"novice of the cult",
		"pillar of death",
		"pillar of draining",
		"pillar of healing",
		"pillar of protection",
		"pillar of summoning",
		"priestess of the wild sun",
		"rage of mazoran",
		"reflection of a mage",
		"reflection of mawhawk",
		"reflection of obujos",
		"retainer of baeloc",
		"scorn of the emperor",
		"servant of tentugly",
		"shadow of boreth",
		"shadow of lersatio",
		"shadow of marziel",
		"shard of corruption",
		"shard of magnor",
		"sight of surrender",
		"son of verminor",
		"soul of dragonking zyrtarch",
		"spark of destruction",
		"spawn of despair",
		"spawn of devovorga",
		"spawn of havoc",
		"spawn of the schnitzel",
		"spawn of the welter",
		"sphere of wrath",
		"spirit of earth",
		"spirit of fertility",
		"spirit of fire",
		"spirit of light",
		"spirit of water",
		"spite of the emperor",
		"squire of nictros",
		"stolen knowledge of armor",
		"stolen knowledge of healing",
		"stolen knowledge of lifesteal",
		"stolen knowledge of spells",
		"stolen knowledge of summoning",
		"stolen tome of portals",
		"sword of vengeance",
		"symbol of fear",
		"symbol of hatred",
		"tentacle of the deep terror",
		"the book of death",
		"the book of secrets",
		"the cold of winter",
		"the corruptor of souls",
		"the count of the core",
		"the devourer of secrets",
		"the duke of the depths",
		"the heat of summer",
		"the lily of night",
		"the lord of the lice",
		"the scion of havoc",
		"the scourge of oblivion",
		"the source of corruption",
		"the voice of ruin",
		"tin lizzard of lyxoph",
		"undead pet of chayenne",
		"weak harbinger of darkness",
		"weak spawn of despair",
		"wildness of urmahlullu",
		"wisdom of urmahlullu",
		"wrath of the emperor",
		"zarcorix of yalahar":
		return true
	default:
		return false
	}
}
