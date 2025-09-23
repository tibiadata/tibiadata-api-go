package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of World
type OnlinePlayers struct {
	Name     string `json:"name"`     // The name of the character.
	Level    int    `json:"level"`    // The character's level.
	Vocation string `json:"vocation"` // The character's vocation.
}

// Child of JSONData
type World struct {
	Name                string          `json:"name"`                  // The name of the world.
	Status              string          `json:"status"`                // The current status of the world.
	PlayersOnline       int             `json:"players_online"`        // The number of currently online players.
	RecordPlayers       int             `json:"record_players"`        // The world's online players record.
	RecordDate          string          `json:"record_date,omitempty"` // The date when the record was achieved.
	CreationDate        string          `json:"creation_date"`         // The year and month it was created.
	Location            string          `json:"location"`              // The physical location of the servers.
	PvpType             string          `json:"pvp_type"`              // The type of PvP.
	PremiumOnly         bool            `json:"premium_only"`          // Whether only premium account players are allowed to play on it.
	TransferType        string          `json:"transfer_type"`         // The type of transfer restrictions it has. regular / locked / blocked
	WorldsQuestTitles   []string        `json:"world_quest_titles"`    // List of world quest titles the server has achieved.
	BattleyeProtected   bool            `json:"battleye_protected"`    // The type of BattlEye protection. true if protected / false if "Not protected by BattlEye."
	BattleyeDate        string          `json:"battleye_date"`         // The date when BattlEye was added. "" if since release / else show date?
	GameWorldType       string          `json:"game_world_type"`       // The type of world. regular / experimental / tournament (if Tournament World Type exists)
	TournamentWorldType string          `json:"tournament_world_type"` // The type of tournament world. "" (default?) / regular / restricted
	OnlinePlayers       []OnlinePlayers `json:"online_players"`        // List of players being currently online.
}

// The base includes two levels: World and Information
type WorldResponse struct {
	World       World       `json:"world"`
	Information Information `json:"information"`
}

var (
	WorldDataRowRegex           = regexp.MustCompile(`<td class=.*>(.*):<\/td><td>(.*)<\/td>`)
	WorldRecordInformationRegex = regexp.MustCompile(`(.*) players \(on (.*)\)`)
	BattlEyeProtectedSinceRegex = regexp.MustCompile(`Protected by BattlEye since (.*)\.`)
	OnlinePlayerRegex           = regexp.MustCompile(`<td style=.*name=.*">(.*)<\/a>.*">(.*)<\/td>.*">(.*)<\/td>`)
)

// TibiaWorldsWorld func
func TibiaWorldsWorldImpl(world string, BoxContentHTML string, url string) (WorldResponse, error) {
	// TODO: We need to read the world name from the response rather than pass it into this func

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return WorldResponse{}, fmt.Errorf("[error] TibiaWorldsWorldImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty vars
	var (
		WorldsStatus, WorldsRecordDate, WorldsCreationDate, WorldsLocation, WorldsPvpType, WorldsTransferType, WorldsBattleyeDate, WorldsGameWorldType, WorldsTournamentWorldType string
		WorldsQuestTitles                                                                                                                                                         []string
		WorldsPlayersOnline, WorldsRecordPlayers                                                                                                                                  int
		WorldsPremiumOnly, WorldsBattleyeProtected                                                                                                                                bool
		WorldsOnlinePlayers                                                                                                                                                       []OnlinePlayers

		insideError error
	)

	// set default values
	WorldsTransferType = "regular"

	// Running query over each div
	ReaderHTML.Find(".Table1 .InnerTableContainer table tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaWorldsWorldImpl failed at WorldsInformationDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		subma1 := WorldDataRowRegex.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

		if len(subma1) > 0 {
			// Creating easy to use vars (and unescape hmtl right string)
			WorldsInformationLeftColumn := subma1[0][1]
			WorldsInformationRightColumn := TibiaDataSanitizeEscapedString(subma1[0][2])

			if WorldsInformationLeftColumn == "Status" {
				switch {
				case strings.Contains(WorldsInformationRightColumn, "</div>Online"):
					WorldsStatus = "online"
				case strings.Contains(WorldsInformationRightColumn, "</div>Offline"):
					WorldsStatus = "offline"
				default:
					WorldsStatus = "unknown"
				}
			}
			if WorldsInformationLeftColumn == "Players Online" {
				WorldsPlayersOnline = TibiaDataStringToInteger(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Online Record" {
				// Regex to get data for record values
				subma2 := WorldRecordInformationRegex.FindAllStringSubmatch(WorldsInformationRightColumn, -1)

				if len(subma2) > 0 {
					// setting record values
					WorldsRecordPlayers = TibiaDataStringToInteger(subma2[0][1])
					WorldsRecordDate = TibiaDataDatetime(subma2[0][2])
				}
			}

			if WorldsInformationLeftColumn == "Creation Date" {
				WorldsCreationDate = TibiaDataDate(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Location" {
				WorldsLocation = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "PvP Type" {
				WorldsPvpType = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "Premium Type" {
				WorldsPremiumOnly = true
			}

			if WorldsInformationLeftColumn == "Transfer Type" {
				WorldsTransferType = WorldsInformationRightColumn
			}

			if WorldsInformationLeftColumn == "World Quest Titles" {
				if WorldsInformationRightColumn != "This game world currently has no title." {
					WorldsQuestTitlesTmp := strings.SplitSeq(WorldsInformationRightColumn, ", ")
					for str := range WorldsQuestTitlesTmp {
						if str != "" {
							WorldsQuestTitles = append(WorldsQuestTitles, TibiaDataRemoveURLs(str))
						}
					}
				}
			}

			if WorldsInformationLeftColumn == "BattlEye Status" {
				if WorldsInformationRightColumn == "Not protected by BattlEye." {
					WorldsBattleyeProtected = false
				} else {
					WorldsBattleyeProtected = true
					if strings.Contains(WorldsInformationRightColumn, "BattlEye since its release") {
						WorldsBattleyeDate = "release"
					} else {
						subma21 := BattlEyeProtectedSinceRegex.FindAllStringSubmatch(WorldsInformationRightColumn, -1)
						WorldsBattleyeDate = TibiaDataDate(subma21[0][1])
					}
				}
			}

			if WorldsInformationLeftColumn == "Game World Type" {
				WorldsGameWorldType = strings.ToLower(WorldsInformationRightColumn)
			}

			if WorldsInformationLeftColumn == "Tournament World Type" {
				WorldsGameWorldType = "tournament"
				if WorldsInformationRightColumn == "Restricted Store" {
					WorldsTournamentWorldType = "restricted"
				} else {
					WorldsTournamentWorldType = strings.ToLower(WorldsInformationRightColumn)
				}
			}
		}

		return true
	})

	if insideError != nil {
		return WorldResponse{}, insideError
	}

	// Running query over each div
	ReaderHTML.Find(".Table2 .InnerTableContainer table tr").First().NextAll().EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		WorldsInformationDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaWorldsWorldImpl failed at WorldsInformationDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		subma1 := OnlinePlayerRegex.FindAllStringSubmatch(WorldsInformationDivHTML, -1)

		if len(subma1) > 0 {
			WorldsOnlinePlayers = append(WorldsOnlinePlayers, OnlinePlayers{
				Name:     TibiaDataSanitizeStrings(subma1[0][1]),
				Level:    TibiaDataStringToInteger(subma1[0][2]),
				Vocation: TibiaDataSanitizeStrings(subma1[0][3]),
			})
		}

		return true
	})

	if insideError != nil {
		return WorldResponse{}, insideError
	}

	//
	// Build the data-blob
	return WorldResponse{
		World: World{
			Name:                world,
			Status:              WorldsStatus,
			PlayersOnline:       WorldsPlayersOnline,
			RecordPlayers:       WorldsRecordPlayers,
			RecordDate:          WorldsRecordDate,
			CreationDate:        WorldsCreationDate,
			Location:            WorldsLocation,
			PvpType:             WorldsPvpType,
			PremiumOnly:         WorldsPremiumOnly,
			TransferType:        WorldsTransferType,
			WorldsQuestTitles:   WorldsQuestTitles,
			BattleyeProtected:   WorldsBattleyeProtected,
			BattleyeDate:        WorldsBattleyeDate,
			GameWorldType:       WorldsGameWorldType,
			TournamentWorldType: WorldsTournamentWorldType,
			OnlinePlayers:       WorldsOnlinePlayers,
		},
		Information: Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{url},
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
