package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestCormaya10(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/houses/Premia/Edron/Cormaya10.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	houseJson, err := TibiaHousesHouseImpl(54025, string(data), "https://www.tibia.com/community/?subtopic=houses&page=view&world=Premia&houseid=54025")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := houseJson.Information

	assert.Equal("https://www.tibia.com/community/?subtopic=houses&page=view&world=Premia&houseid=54025", information.TibiaURLs[0])

	assert.Equal(54025, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("Edron", houseJson.House.Town)
	assert.Equal("Cormaya 10", houseJson.House.Name)
	assert.Equal("house", houseJson.House.Type)
	assert.Equal(3, houseJson.House.Beds)
	assert.Equal(80, houseJson.House.Size)
	assert.Equal(300000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_54025.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.False(houseStatus.IsAuctioned)
	assert.True(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseAuction{CurrentBid: 0, CurrentBidder: "", AuctionOngoing: false, AuctionEnd: ""}, houseStatus.Auction)
	assert.Equal("The house has been rented by Xendor of Askara. He has paid the rent until Feb 02 2022, 10:05:26 CET.", houseStatus.Original)

	houseRental := houseJson.House.Status.Rental
	assert.NotNil(houseRental)
	assert.Equal("Xendor of Askara", houseRental.Owner)
	assert.Equal("male", houseRental.OwnerSex)
	assert.Equal("2022-02-02T09:05:26Z", houseRental.PaidUntil)
	assert.Empty(houseRental.MovingDate)
	assert.Empty(houseRental.TransferReceiver)
	assert.Equal(0, houseRental.TransferPrice)
	assert.False(houseRental.TransferAccept)
}

func TestCormaya11(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/houses/Premia/Edron/Cormaya11.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	houseJson, err := TibiaHousesHouseImpl(54026, string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(54026, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("Edron", houseJson.House.Town)
	assert.Equal("Cormaya 11", houseJson.House.Name)
	assert.Equal("house", houseJson.House.Type)
	assert.Equal(2, houseJson.House.Beds)
	assert.Equal(43, houseJson.House.Size)
	assert.Equal(150000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_54026.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.True(houseStatus.IsAuctioned)
	assert.False(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseRental{Owner: "", OwnerSex: "", PaidUntil: "", MovingDate: "", TransferReceiver: "", TransferPrice: 0, TransferAccept: false}, houseStatus.Rental)
	assert.Equal("The house is currently being auctioned. The auction will end at Jan 21 2022, 10:00:00 CET. The highest bid so far is 200000 gold and has been submitted by Ciuchy Szajba.", houseStatus.Original)

	houseAuction := houseJson.House.Status.Auction
	assert.NotNil(houseAuction)
	assert.Equal(200000, houseAuction.CurrentBid)
	assert.Equal("Ciuchy Szajba", houseAuction.CurrentBidder)
	assert.True(houseAuction.AuctionOngoing)
	assert.Equal("2022-01-21T09:00:00Z", houseAuction.AuctionEnd)
}

func TestCormaya9c(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/houses/Premia/Edron/Cormaya9c.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	houseJson, _ := TibiaHousesHouseImpl(54023, string(data), "")
	assert := assert.New(t)

	assert.Equal(54023, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("Edron", houseJson.House.Town)
	assert.Equal("Cormaya 9c", houseJson.House.Name)
	assert.Equal("house", houseJson.House.Type)
	assert.Equal(2, houseJson.House.Beds)
	assert.Equal(25, houseJson.House.Size)
	assert.Equal(80000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_54023.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.True(houseStatus.IsAuctioned)
	assert.False(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseRental{Owner: "", OwnerSex: "", PaidUntil: "", MovingDate: "", TransferReceiver: "", TransferPrice: 0, TransferAccept: false}, houseStatus.Rental)
	assert.Equal("The house is currently being auctioned. The auction has ended at Jan 21 2022, 10:00:00 CET. The highest bid so far is 12345 gold and has been submitted by Ciuchy Szajba.", houseStatus.Original)

	houseAuction := houseJson.House.Status.Auction
	assert.NotNil(houseAuction)
	assert.Equal(12345, houseAuction.CurrentBid)
	assert.Equal("Ciuchy Szajba", houseAuction.CurrentBidder)
	assert.False(houseAuction.AuctionOngoing)
	assert.Equal("2022-01-21T09:00:00Z", houseAuction.AuctionEnd)
}

func TestBeachHomeApartmentsFlat14(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/houses/Premia/Thais/BeachHomeApartmentsFlat14.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	houseJson, err := TibiaHousesHouseImpl(10214, string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(10214, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("Thais", houseJson.House.Town)
	assert.Equal("Beach Home Apartments, Flat 14", houseJson.House.Name)
	assert.Equal("house", houseJson.House.Type)
	assert.Equal(1, houseJson.House.Beds)
	assert.Equal(7, houseJson.House.Size)
	assert.Equal(25000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_10214.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.True(houseStatus.IsAuctioned)
	assert.False(houseStatus.IsRented)
	assert.False(houseStatus.IsMoving)
	assert.False(houseStatus.IsTransfering)
	assert.Equal(HouseRental{Owner: "", OwnerSex: "", PaidUntil: "", MovingDate: "", TransferReceiver: "", TransferPrice: 0, TransferAccept: false}, houseStatus.Rental)
	assert.Equal("The house is currently being auctioned. No bid has been submitted so far.", houseStatus.Original)

	houseAuction := houseJson.House.Status.Auction
	assert.NotNil(houseAuction)
	assert.Equal(0, houseAuction.CurrentBid)
	assert.Empty(houseAuction.CurrentBidder)
	assert.False(houseAuction.AuctionOngoing)
	assert.Empty(houseAuction.AuctionEnd)
}

func TestBeachHomeApartmentsFlat15(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/houses/Premia/Thais/BeachHomeApartmentsFlat15.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	houseJson, err := TibiaHousesHouseImpl(10215, string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(10215, houseJson.House.Houseid)
	assert.Equal("Premia", houseJson.House.World)
	assert.Equal("Thais", houseJson.House.Town)
	assert.Equal("Beach Home Apartments, Flat 15", houseJson.House.Name)
	assert.Equal("house", houseJson.House.Type)
	assert.Equal(1, houseJson.House.Beds)
	assert.Equal(7, houseJson.House.Size)
	assert.Equal(25000, houseJson.House.Rent)
	assert.Equal("https://static.tibia.com/images/houses/house_10215.png", houseJson.House.Img)

	houseStatus := houseJson.House.Status
	assert.NotNil(houseStatus)
	assert.False(houseStatus.IsAuctioned)
	assert.True(houseStatus.IsRented)
	assert.True(houseStatus.IsMoving)
	assert.True(houseStatus.IsTransfering)
	assert.Equal(HouseRental{Owner: "Xenaris mag", OwnerSex: "female", PaidUntil: "2019-01-10T09:20:52Z", MovingDate: "2018-12-12T09:00:00Z", TransferReceiver: "Ivarr Bezkosci", TransferPrice: 850000, TransferAccept: true}, houseStatus.Rental)
	assert.Equal("The house has been rented by Xenaris mag. She has paid the rent until Jan 10 2019, 10:20:52 CET. She will move out on Dec 12 2018, 10:00:00 CET (time of daily server save) and will pass the house to Ivarr Bezkosci for 850000 gold coins.", houseStatus.Original)

	houseAuction := houseJson.House.Status.Auction
	assert.NotNil(houseAuction)
	assert.Empty(houseAuction.CurrentBid)
	assert.Empty(houseAuction.CurrentBidder)
	assert.False(houseAuction.AuctionOngoing)
	assert.Empty(houseAuction.AuctionEnd)
}
