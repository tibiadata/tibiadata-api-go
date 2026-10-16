package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Child of House
type HousesAuction struct {
	AuctionBid  int    `json:"current_bid"` // The highest bid so far.
	AuctionLeft string `json:"time_left"`   // The number of days or hours left until the bid ends.
	IsFinished  bool   `json:"finished"`    // Whether the auction is finished or not.
}

// Child of HousesHouses
type HousesHouse struct {
	Name        string        `json:"name"`      // The name of the house/guildhall.
	HouseID     int           `json:"house_id"`  // The internal ID of the house/guildhall.
	Size        int           `json:"size"`      // The size in SQM.
	Rent        int           `json:"rent"`      // The monthly cost in gold coins for the house/guildhall.
	IsRented    bool          `json:"rented"`    // Whether the auction is rented or not.
	IsAuctioned bool          `json:"auctioned"` // Whether the auction is auctioned or not.
	Auction     HousesAuction `json:"auction"`   // Details about the auction.
}

// Child of JSONData
type HousesHouses struct {
	World         string        `json:"world"`          // The name of the world the house/guildhall belongs to.
	Town          string        `json:"town"`           // The town where the house/guildhall is located.
	HouseList     []HousesHouse `json:"house_list"`     // List of all houses.
	GuildhallList []HousesHouse `json:"guildhall_list"` // List of all guildhalls.
}

// The base includes two levels: HousesHouses and Information
type HousesOverviewResponse struct {
	Houses      HousesHouses `json:"houses"`
	Information Information  `json:"information"`
}

var (
	houseOverviewDataRegex      = regexp.MustCompile(`<td.*><nobr>(.*)<\/nobr><\/td><td.*><nobr>([0-9]+).sqm<\/nobr><\/td><td.*><nobr>([0-9]+)(k+).gold<\/nobr><\/td><td.*><nobr>(.*)<\/nobr><\/td>.*houseid" value="([0-9]+)"\/><div.*`)
	houseOverviewAuctionedRegex = regexp.MustCompile(`auctioned.\(([0-9]+).gold;.(finished|(.*).left)\)`)
)

// TibiaHousesOverview func
func TibiaHousesOverviewImpl(c *gin.Context, world string, town string, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) (HousesOverviewResponse, error) {
	// Creating empty vars
	var HouseData, GuildhallData []HousesHouse
	var TibiaHouseURLs []string

	// list of different fansite types
	HouseTypes := []string{"houses", "guildhalls"}

	// running over the FansiteTypes array
	for _, HouseType := range HouseTypes {
		houses, houseUrl, err := makeHouseRequest(HouseType, world, town, tibiaDataUserAgentForRequest(c), htmlDataCollector)
		if err != nil {
			return HousesOverviewResponse{}, fmt.Errorf("[error] TibiaHousesOverviewImpl failed at makeHouseRequest, type: %s, err: %s", HouseType, err)
		}

		switch HouseType {
		case "houses":
			HouseData = houses
		case "guildhalls":
			GuildhallData = houses
		}

		TibiaHouseURLs = append(TibiaHouseURLs, houseUrl)
	}

	// Build the data-blob
	return HousesOverviewResponse{
		HousesHouses{
			World:         world,
			Town:          town,
			HouseList:     HouseData,
			GuildhallList: GuildhallData,
		},
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  TibiaHouseURLs,
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}

func makeHouseRequest(HouseType, world, town, userAgent string, htmlDataCollector func(TibiaDataRequestStruct) (string, error)) ([]HousesHouse, string, error) {
	// Creating an empty var
	var output []HousesHouse

	tibiadataRequest := TibiaDataRequestStruct{
		Method:    resty.MethodGet,
		URL:       "https://www.tibia.com/community/?subtopic=houses&world=" + TibiaDataQueryEscapeString(world) + "&town=" + TibiaDataQueryEscapeString(town) + "&type=" + TibiaDataQueryEscapeString(HouseType),
		UserAgent: userAgent,
	}

	BoxContentHTML, err := htmlDataCollector(tibiadataRequest)
	// return error (e.g. for maintenance mode)
	if err != nil {
		return nil, "", err
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, "", err
	}

	var insideError error

	ReaderHTML.Find(".TableContentContainer .TableContent tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
		house := HousesHouse{}

		// Storing HTML into HousesDivHTML
		HousesDivHTML, err := s.Html()
		if err != nil {
			insideError = err
			return false
		}

		// Removing linebreaks from HTML
		HousesDivHTML = TibiaDataHTMLRemoveLinebreaks(HousesDivHTML)
		HousesDivHTML = TibiaDataSanitizeStrings(HousesDivHTML)

		subma1 := houseOverviewDataRegex.FindAllStringSubmatch(HousesDivHTML, -1)

		if len(subma1) > 0 {
			// House details
			house.Name = TibiaDataSanitizeEscapedString(subma1[0][1])
			house.HouseID = TibiaDataStringToInteger(subma1[0][6])
			house.Size = TibiaDataStringToInteger(subma1[0][2])
			house.Rent = TibiaDataConvertValuesWithK(subma1[0][3] + subma1[0][4])

			// HousesAction details
			s := subma1[0][5]
			switch {
			case strings.Contains(s, "rented"):
				house.IsRented = true
			case strings.Contains(s, "auctioned (no bid yet)"):
				house.IsAuctioned = true
			case strings.Contains(s, "auctioned"):
				house.IsAuctioned = true
				subma1b := houseOverviewAuctionedRegex.FindAllStringSubmatch(s, -1)
				house.Auction.AuctionBid = TibiaDataStringToInteger(subma1b[0][1])
				if subma1b[0][2] == "finished" {
					house.Auction.IsFinished = true
				} else {
					house.Auction.AuctionLeft = subma1b[0][3]
				}
			}

			output = append(output, house)
		}

		return true
	})

	return output, tibiadataRequest.URL, insideError
}
