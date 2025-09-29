package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

// Child of Status
type HouseRental struct {
	Owner            string `json:"owner"`             // The current owner of the house/guildhall.
	OwnerSex         string `json:"owner_sex"`         // The owner's sex.
	PaidUntil        string `json:"paid_until"`        // The date the last paid rent is due.
	MovingDate       string `json:"moving_date"`       // The date when the owner will move out.
	TransferReceiver string `json:"transfer_receiver"` // The character who will receive the house.
	TransferPrice    int    `json:"transfer_price"`    // The price that will be paid from the current owner to the new owner for the transfer.
	TransferAccept   bool   `json:"transfer_accept"`   // Whether the transfer is accepted or not.
}

// Child of Status
type HouseAuction struct {
	CurrentBid     int    `json:"current_bid"`     // The currently highest bid on the house/guildhall.
	CurrentBidder  string `json:"current_bidder"`  // The character that holds the current highest bid.
	AuctionOngoing bool   `json:"auction_ongoing"` // Whether the auction is still ongoing or not.
	AuctionEnd     string `json:"auction_end"`     // The date when the auction will finish.
}

// Child of House
type HouseStatus struct {
	IsAuctioned   bool         `json:"is_auctioned"`   // Whether the house/guildhall is being auctioned.
	IsRented      bool         `json:"is_rented"`      // Wether the house/guildhall is being rented.
	IsMoving      bool         `json:"is_moving"`      // Wether the owner is moving out.
	IsTransfering bool         `json:"is_transfering"` // Wether the house/guildhall is being transfered.
	Auction       HouseAuction `json:"auction"`        // Details about the auction.
	Rental        HouseRental  `json:"rental"`         // Details about the transfer.
	Original      string       `json:"original"`       // Original plain text information.
}

// Child of JSONData
type House struct {
	Houseid int         `json:"houseid"`        // The internal ID of the house/guildhall.
	World   string      `json:"world"`          // The name of the world the house/guildhall belongs to.
	Town    string      `json:"town,omitempty"` // The town where the house/guildhall is located.
	Name    string      `json:"name"`           // The name of the house/guildhall.
	Type    string      `json:"type,omitempty"` // The type of home. (house or guildhall)
	Beds    int         `json:"beds"`           // The number of beds it has.
	Size    int         `json:"size"`           // The number of SQM it has.
	Rent    int         `json:"rent"`           // The monthly cost in gold coins for the house.
	Img     string      `json:"img"`            // The URL to the house's minimap image.
	Status  HouseStatus `json:"status"`         // The current status of the house/guildhall.
}

// The base includes two levels: Houses and Information
type HouseResponse struct {
	House       House       `json:"house"`
	Information Information `json:"information"`
}

var (
	houseDataRegex = regexp.MustCompile(`<td.*src="(.*)" width.*<b>(.*)<\/b>.*This (house|guildhall) can.*to ([0-9]+) beds?..*<b>([0-9]+) square.*<b>([0-9]+)([k]+).gold<\/b>.*on <b>([A-Za-z]+)<\/b>.(.*)<\/td>`)
	// matching for this: and <wants to|will> pass the <HouseType> to <TransferReceiver> for <TransferPrice> gold
	housePassingRegex = regexp.MustCompile(`and (wants to|will) pass the (house|guildhall) to (.*) for ([0-9]+) gold`)
	// matching for this: <OwnerSex> will move out on <MovingDate> (
	moveOutRegex = regexp.MustCompile(`(He|She) will move out on (.*?) \(`)
	// matching for this: The <HouseType> has been rented by <Owner>. <OwnerSex> has paid the rent until <PaidUntil>.
	paidUntilRegex      = regexp.MustCompile(`The (house|guildhall) has been rented by (.*). (He|She) has paid.*until (.*?)\.`)
	houseAuctionedRegex = regexp.MustCompile(`The (house|guildhall) is currently.*The auction (will end|has ended) at (.*)\. The.*is ([0-9]+) gold.*submitted by (.*)\.`)
)

// TibiaHousesHouse func
func TibiaHousesHouseImpl(houseid int, BoxContentHTML string, url string) (HouseResponse, error) {
	// Creating empty vars
	var HouseData House

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return HouseResponse{}, fmt.Errorf("[error] TibiaHousesHouseImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Running query over each div
	HouseHTML, err := ReaderHTML.Find(".BoxContent table tr").First().Html()
	if err != nil {
		return HouseResponse{}, fmt.Errorf("[error] TibiaHousesHouseImpl failed at ReaderHTML.Find, err: %s", err)
	}

	// Regex to get data for house
	subma1 := houseDataRegex.FindAllStringSubmatch(HouseHTML, -1)

	if len(subma1) > 0 {
		HouseData.Houseid = houseid
		HouseData.World = subma1[0][8]

		rawHouse, err := validation.GetHouseRaw(HouseData.Houseid)
		if err != nil {
			return HouseResponse{}, err
		}

		HouseData.Town = rawHouse.Town
		HouseData.Type = rawHouse.Type

		HouseData.Name = TibiaDataSanitizeEscapedString(subma1[0][2])
		HouseData.Img = subma1[0][1]
		HouseData.Beds = TibiaDataStringToInteger(subma1[0][4])
		HouseData.Size = TibiaDataStringToInteger(subma1[0][5])
		HouseData.Rent = TibiaDataConvertValuesWithK(subma1[0][6] + subma1[0][7])

		HouseData.Status.Original = strings.TrimSpace(TibiaDataSanitizeStrings(TibiaDataSanitizeEscapedString(RemoveHtmlTag(subma1[0][9]))))

		switch {
		case strings.Contains(HouseData.Status.Original, "has been rented by"):
			// rented
			switch {
			case strings.Contains(HouseData.Status.Original, " pass the house to "), strings.Contains(HouseData.Status.Original, " pass the guildhall to "):
				HouseData.Status.IsTransfering = true

				subma2 := housePassingRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				if subma2[0][1] == "will" {
					HouseData.Status.Rental.TransferAccept = true
				}
				HouseData.Status.Rental.TransferReceiver = subma2[0][3]
				HouseData.Status.Rental.TransferPrice = TibiaDataStringToInteger(subma2[0][4])
				fallthrough

			case strings.Contains(HouseData.Status.Original, " will move out on "):
				HouseData.Status.IsMoving = true
				subma2 := moveOutRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.MovingDate = TibiaDataDatetime(subma2[0][2])
				fallthrough

			default:
				HouseData.Status.IsRented = true
				subma2 := paidUntilRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Rental.Owner = subma2[0][2]
				HouseData.Status.Rental.PaidUntil = TibiaDataDatetime(subma2[0][4])
				switch subma2[0][3] {
				case "She":
					HouseData.Status.Rental.OwnerSex = "female"
				case "He":
					HouseData.Status.Rental.OwnerSex = "male"
				}
			}
		case strings.Contains(HouseData.Status.Original, "is currently being auctioned"):
			// auctioned
			HouseData.Status.IsAuctioned = true
			// check if bid is going on
			if !strings.Contains(HouseData.Status.Original, "No bid has been submitted so far.") {
				subma2 := houseAuctionedRegex.FindAllStringSubmatch(HouseData.Status.Original, -1)
				// storing values from regex
				HouseData.Status.Auction.AuctionEnd = TibiaDataDatetime(subma2[0][3])
				HouseData.Status.Auction.CurrentBid = TibiaDataStringToInteger(subma2[0][4])
				HouseData.Status.Auction.CurrentBidder = TibiaDataSanitizeStrings(subma2[0][5])
				if subma2[0][2] == "will end" {
					HouseData.Status.Auction.AuctionOngoing = true
				}
			}
		}
	}

	// Build the data-blob
	return HouseResponse{
		HouseData,
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
