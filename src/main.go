package main

import (
	"log/slog"
	"sync/atomic"

	"github.com/tibiadata/tibiadata-api-go/src/validation"
)

var (
	// application readyz endpoint value for k8s
	isReady atomic.Bool

	// TibiaDataDefaultVoc - default vocation when not specified in request
	TibiaDataDefaultVoc string = "all"

	// TibiaData app flags for running
	TibiaDataAPIversion      int = 4
	TibiaDataDebug           bool
	TibiaDataRestrictionMode bool
	TibiaDataCacheControl    bool = true

	// TibiaData app settings
	TibiaDataAPIDetails APIDetails // containing information from build
	TibiaDataHost       string     // set through env TIBIADATA_HOST
	TibiaDataProtocol   = "https"  // can be overridden by env TIBIADATA_PROTOCOL

	// TibiaData app details set to release/build on GitHub
	TibiaDataBuildRelease = "unknown"     // will be set by GitHub Actions (to release number)
	TibiaDataBuildBuilder = "manual"      // will be set by GitHub Actions
	TibiaDataBuildCommit  = "-"           // will be set by GitHub Actions (to git commit)
	TibiaDataBuildEdition = "open-source" //

	// Tibia Fansite API
	TibiaFansiteAPI   bool   // indicates if the Tibia Fansite API is used
	TibiaFansiteToken string // the token used for accessing the Tibia Fansite API
)

// @title           TibiaData API
// @version         edge
// @description     This is the API documentation for the TibiaData API.
// @description     The documentation contains version 3 and above.
// @termsOfService  https://tibiadata.com/terms/

// @contact.name   TibiaData
// @contact.url    https://tibiadata.com/contact/
// @contact.email  tobias@tibiadata.com

// @license.name  MIT
// @license.url   https://github.com/tibiadata/tibiadata-api-go/blob/main/LICENSE

// @schemes   http
// @host      localhost:8080
// @BasePath  /

func init() {
	initTibiaDataLogging()

	slog.Info("TibiaData API initializing")

	slog.Info("TibiaData API release", "release", TibiaDataBuildRelease)
	slog.Info("TibiaData API build", "build", TibiaDataBuildBuilder)
	slog.Info("TibiaData API commit", "commit", TibiaDataBuildCommit)
	slog.Info("TibiaData API edition", "edition", TibiaDataBuildEdition)

	TibiaDataAPIDetails = APIDetails{
		Version: TibiaDataAPIversion,
		Release: TibiaDataBuildRelease,
		Commit:  TibiaDataBuildCommit,
	}

	if getEnvAsBool("DEBUG_MODE", false) {
		TibiaDataDebug = true
	}
	slog.Info("TibiaData API debug-mode", "enabled", TibiaDataDebug)

	TibiaDataInitializer()

	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)

	if TibiaDataDebug {
		slog.Debug("TibiaData API User-Agent", "user_agent", TibiaDataUserAgent)
	}

	initTibiaDataClient()

	if TibiaFansiteAPI {
		checkTibiaFansiteAPIStatus()
	}

	err := validation.Initiate(TibiaDataUserAgent)
	if err != nil {
		panic(err)
	}

}

func main() {
	slog.Info("TibiaData API starting")

	runWebServer()
}

// TibiaDataInitializer set the background for the webserver
func TibiaDataInitializer() {
	if isEnvExist("TIBIADATA_EDITION") {
		TibiaDataBuildEdition = getEnv("TIBIADATA_EDITION", "open-source")
	}

	if isEnvExist("TIBIA_FANSITEAPI_TOKEN") {
		TibiaFansiteToken = getEnv("TIBIA_FANSITEAPI_TOKEN", "")
		if err := validateTibiaFansiteToken(TibiaFansiteToken); err == nil {
			TibiaFansiteAPI = true
			slog.Info("TibiaData API fansiteapi enabled")
		} else {
			TibiaFansiteToken = ""
			TibiaFansiteAPI = false
			slog.Warn("TibiaData API fansiteapi token is invalid", "error", err)
		}
	}

	TibiaDataCacheControl = getEnvAsBool("TIBIADATA_CACHE_CONTROL_HEADERS", true)
	slog.Info("TibiaData API cache-control headers", "enabled", TibiaDataCacheControl)

	if isEnvExist("TIBIADATA_HOST") {
		TibiaDataHost = getEnv("TIBIADATA_HOST", "")
		slog.Info("TibiaData API hostname", "host", TibiaDataHost)
	}
	if isEnvExist("TIBIADATA_PROTOCOL") {
		TibiaDataProtocol = getEnv("TIBIADATA_PROTOCOL", "https")
		slog.Info("TibiaData API protocol", "protocol", TibiaDataProtocol)
	}

	if isEnvExist("TIBIADATA_PROXY") {

		TibiaDataProxyProtocol := getEnv("TIBIADATA_PROXY_PROTOCOL", "https")
		switch TibiaDataProxyProtocol {
		case "http":
			TibiaDataProxyProtocol = "http"
		}

		TibiaDataProxyDomain = TibiaDataProxyProtocol + "://" + getEnv("TIBIADATA_PROXY", "www.tibia.com") + "/"
		slog.Info("TibiaData API proxy", "proxy", TibiaDataProxyDomain)
	}

	_ = tibiaNewslistArchive()
	_ = tibiaNewslistArchiveDays()
	_ = tibiaNewslistLatest()

	_ = tibiaBoostableBossesV3()
	_ = tibiaCharactersCharacterV3()
	_ = tibiaCreaturesOverviewV3()
	_ = tibiaCreaturesCreatureV3()
	_ = tibiaFansitesV3()
	_ = tibiaGuildsGuildV3()
	_ = tibiaGuildsOverviewV3()
	_ = tibiaHighscoresV3()
	_ = tibiaHousesHouseV3()
	_ = tibiaHousesOverviewV3()
	_ = tibiaKillstatisticsV3()
	_ = tibiaNewslistArchiveV3()
	_ = tibiaNewslistArchiveDaysV3()
	_ = tibiaNewslistLatestV3()
	_ = tibiaNewslistV3()
	_ = tibiaNewsV3()
	_ = tibiaSpellsOverviewV3()
	_ = tibiaSpellsSpellV3()
	_ = tibiaWorldsOverviewV3()
	_ = tibiaWorldsWorldV3()

}
