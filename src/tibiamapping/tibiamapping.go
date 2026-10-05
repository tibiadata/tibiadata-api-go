package tibiamapping

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// TibiaMapping stores the values returned by Run()
type TibiaMapping struct {
	RawData   []byte
	Sha256Sum string
	Sha512Sum string
}

const (
	// tibiaAssetsDataMinJsonURL is the endpoint to get the data.min.json file
	tibiaAssetsDataMinJsonURL = "https://assets.tibiadata.com/data.min.json"

	// tibiaAssetsSha256SumURL is the endpoint to get the sha256sum.txt file
	tibiaAssetsSha256SumURL = "https://assets.tibiadata.com/sha256sum.txt"

	// tibiaAssetsSha512SumURL is the endpoint to get the sha512sum.txt file
	tibiaAssetsSha512SumURL = "https://assets.tibiadata.com/sha512sum.txt"
)

// Run is used to load data from the assets JSON file
func Run(userAgent string) (TibiaMapping, error) {
	// Logging the start of tibiamapping
	slog.Info("TibiaData API assets mapping running")

	// Setting up resty client
	client := resty.New()

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   userAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Making the GET request to the data file
	res, err := client.R().Get(tibiaAssetsDataMinJsonURL)
	if err != nil {
		return TibiaMapping{}, err
	}

	// Checking if the response code was OK
	if res.StatusCode() != http.StatusOK {
		return TibiaMapping{}, fmt.Errorf("res status code %d", res.StatusCode())
	}

	// Making the GET request to the sha256 file
	sha256, err := client.R().Get(tibiaAssetsSha256SumURL)
	if err != nil {
		return TibiaMapping{}, err
	}

	// Checking if the response code was OK
	if sha256.StatusCode() != http.StatusOK {
		return TibiaMapping{}, fmt.Errorf("sha256 status code %d", sha256.StatusCode())
	}

	// Making the GET request to the sha512 file
	sha512, err := client.R().Get(tibiaAssetsSha512SumURL)
	if err != nil {
		return TibiaMapping{}, err
	}

	// Checking if the response code was OK
	if sha512.StatusCode() != http.StatusOK {
		return TibiaMapping{}, fmt.Errorf("sha512 status code %d", sha512.StatusCode())
	}

	// Log that Tibia Mapping has been successfully completed
	slog.Info("TibiaData API assets mapping completed")

	return TibiaMapping{
		RawData:   res.Body(),
		Sha256Sum: string(sha256.Body()),
		Sha512Sum: string(sha512.Body()),
	}, nil
}
