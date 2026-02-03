package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	_ "github.com/mantyr/go-charset/data"
	"github.com/tibiadata/tibiadata-api-go/src/validation"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	// TibiaData app resty vars
	TibiaDataUserAgent, TibiaDataProxyDomain string

	// ErrorNotFound will be returned if the requests ends up in a 404
	ErrorNotFound = errors.New("page not found")
)

// DebugOutInformation wraps OutInformation with some debug info
type DebugOutInformation struct {
	Information Information `json:"information"`
	Debug       Debug       `json:"debug"`
}

// OutInformation wraps Information in other for all json outputs be consistent
type OutInformation struct {
	Information Information `json:"information"`
}

// Information stores some API related data
type Information struct {
	APIDetails APIDetails `json:"api"`        // The API details.
	Timestamp  string     `json:"timestamp"`  // The timestamp from when the data was processed.
	TibiaURLs  []string   `json:"tibia_urls"` // The links to the sources of the data on tibia.com
	Status     Status     `json:"status"`     // The response status information.
}

// API details store information about this API
type APIDetails struct {
	Version int    `json:"version"` // The API major version currently running.
	Release string `json:"release"` // The API release currently running.
	Commit  string `json:"commit"`  // The API GitHub commit sha.
}

// Status stores information about the response
type Status struct {
	HTTPCode int    `json:"http_code"`         // The HTTP response code from the API.
	Error    int    `json:"error,omitempty"`   // The error code thrown by TibiaData API for identification of issue.
	Message  string `json:"message,omitempty"` // The error message thrown by TibiaData API for human readability.
}

// TibiaDataRequest is the struct of request information
type TibiaDataRequestStruct struct {
	Method   string            `json:"method"`    // Request method (default: GET)
	URL      string            `json:"url"`       // Request URL
	FormData map[string]string `json:"form_data"` // Request form content (used when POST)
	RawBody  bool              `json:"raw_body"`  // If set to true the whole content from tibia.com will be passed down
}

// RunWebServer starts the gin server
// It blocks the code and will only finish execution on shutdown
func runWebServer() {
	// Setting gin-application to certain mode if GIN_MODE is set to release, test or debug (default is release)
	switch ginMode := getEnv("GIN_MODE", "release"); ginMode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Logging the gin.mode
	log.Printf("[info] TibiaData API gin-mode: %s", gin.Mode())

	// Starting an Engine instance
	router := gin.Default()

	// Gin middleware to enable GZIP support
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Set 404 not found page
	router.NoRoute(func(c *gin.Context) {
		TibiaDataErrorHandler(
			c,
			ErrorNotFound,
			http.StatusNotFound,
		)
	})

	// Set proxy feature of gin
	if isEnvExist("GIN_TRUSTED_PROXIES") {
		trustedProxies := getEnv("GIN_TRUSTED_PROXIES", "")
		_ = router.SetTrustedProxies(strings.Split(trustedProxies, ","))
		log.Printf("[info] TibiaData API gin-trusted-proxies: %s", strings.Split(trustedProxies, ","))
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	// Set the TibiaData restriction mode
	TibiaDataRestrictionMode = getEnvAsBool("TIBIADATA_RESTRICTION_MODE", false)
	log.Printf("[info] TibiaData API restriction-mode: %t", TibiaDataRestrictionMode)

	// Set the ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		data := Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			TibiaURLs:  []string{},
			Status: Status{
				HTTPCode: http.StatusOK,
				Message:  "pong",
			},
		}

		var output OutInformation
		output.Information = data

		c.JSON(http.StatusOK, output)
	})

	// health endpoints for kubernetes
	router.GET("/", rootz)
	router.GET("/health", healthz)
	router.GET("/healthz", healthz)
	router.GET("/readyz", readyz)

	// Set the debug endpoint
	router.GET("/debug", debugHandler)

	// TibiaData API version 3 endpoints
	router.GET("/v3/*action", func(c *gin.Context) {
		c.JSON(299, gin.H{
			"error": "TibiaData v3 is deprecated.",
			"information": InformationV3{
				APIversion: 3,
				Timestamp:  TibiaDataDatetime(""),
			},
		})
	})

	// TibiaData API version 4 endpoints
	v4 := router.Group("/v4")
	{

		v4.GET("/bazaar/auction/:id", tibiaBazaarAuctionV3)

		// Tibia characters
		v4.GET("/boostablebosses", tibiaBoostableBosses)

		// Tibia characters
		v4.GET("/character/:name", tibiaCharactersCharacter)

		// Tibia creatures
		v4.GET("/creature/:race", tibiaCreaturesCreature)
		v4.GET("/creatures", tibiaCreaturesOverview)

		// Tibia fansites
		v4.GET("/fansites", tibiaFansites)

		// Tibia guilds
		v4.GET("/guild/:name", tibiaGuildsGuild)
		// v4.GET("/guild/:name/events",TibiaGuildsGuildEvents)
		// v4.GET("/guild/:name/wars",TibiaGuildsGuildWars)
		v4.GET("/guilds/:world", tibiaGuildsOverview)

		// Tibia highscores
		v4.GET("/highscores/:world", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v4.BasePath()+"/highscores/"+c.Param("world")+"/experience/"+TibiaDataDefaultVoc+"/1")
		})
		v4.GET("/highscores/:world/:category", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, v4.BasePath()+"/highscores/"+c.Param("world")+"/"+c.Param("category")+"/"+TibiaDataDefaultVoc+"/1")
		})
		v4.GET("/highscores/:world/:category/:vocation", tibiaHighscores)
		v4.GET("/highscores/:world/:category/:vocation/:page", tibiaHighscores)

		// Tibia houses
		v4.GET("/house/:world/:house_id", tibiaHousesHouse)
		v4.GET("/houses/:world/:town", tibiaHousesOverview)

		// Tibia killstatistics
		v4.GET("/killstatistics/:world", tibiaKillstatistics)

		// Tibia news
		v4.GET("/news/archive", tibiaNewslist)       // all categories (default 90 days)
		v4.GET("/news/archive/:days", tibiaNewslist) // all categories
		v4.GET("/news/id/:news_id", tibiaNews)       // shows one news entry
		v4.GET("/news/latest", tibiaNewslist)        // only news and articles
		v4.GET("/news/newsticker", tibiaNewslist)    // only news_ticker

		// Tibia spells
		v4.GET("/spell/:spell_id", tibiaSpellsSpell)
		v4.GET("/spells", tibiaSpellsOverview)

		// Tibia worlds
		v4.GET("/world/:name", tibiaWorldsWorld)
		v4.GET("/worlds", tibiaWorldsOverview)
	}

	// Container version details endpoint
	router.GET("/versions", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"release": TibiaDataBuildRelease,
			"build":   TibiaDataBuildBuilder,
			"commit":  TibiaDataBuildCommit,
			"edition": TibiaDataBuildEdition,
		})
	})

	// Build the http server
	server := &http.Server{
		Addr:    ":8080", // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		Handler: router,
	}

	// Prepare for a graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Run a go routine that will receive the shutdown input
	go func() {
		<-quit
		log.Println("[info] TibiaData API received shutdown input")
		if err := server.Close(); err != nil {
			log.Fatal("[error] TibiaData API server close error:", err)
		}
	}()

	// setting readyz endpoint to true
	isReady.Store(true)

	log.Println("[info] TibiaData API starting webserver")

	// Run the server
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("[info] TibiaData API server gracefully shut down")
		} else {
			log.Fatal("[error] TibiaData API server closed unexpectedly")
		}
	}
}

// BoostableBosses godoc
// @Summary      List of boostable bosses
// @Description  Show all boostable bosses listed
// @Tags         boostable bosses
// @Accept       json
// @Produce      json
// @Success      200  {object}  BoostableBossesOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/boostablebosses [get]
func tibiaBoostableBosses(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method:  resty.MethodGet,
		URL:     "https://www.tibia.com/library/?subtopic=boostablebosses",
		RawBody: true,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaBoostableBossesOverviewImpl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaBoostableBosses")
}

// Bazaar auction godoc
// @Summary      Show one bazaar auction
// @Description  Show all information about one bazaar auction
// @Tags         bazaar
// @Accept       json
// @Produce      json
// @Param        id path int true "The ID of auction"
// @Success      200  {object}  BazaarAuctionResponse
// @Router       /v3/bazaar/auction/{id} [get]
func tibiaBazaarAuctionV3(c *gin.Context) {
	// getting params from URL
	id := c.Param("id")

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/charactertrade/?page=details&auctionid=" + id,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaBazaarAuctionV3Impl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaBazaarAuctionV3")
}

// Character godoc
// @Summary      Show one character
// @Description  Show all information about one character available
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        name path string true "The character name" extensions(x-example=Trollefar)
// @Success      200  {object}  CharacterResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/character/{name} [get]
func tibiaCharactersCharacter(c *gin.Context) {
	// Getting params from URL
	name := c.Param("name")

	// Validate the name
	err := validation.IsCharacterNameValid(name)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	// Build the request structure
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=characters&name=" + TibiaDataQueryEscapeString(name),
	}

	// Handle the request
	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCharactersCharacterImpl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaCharactersCharacter")
}

// Creatures godoc
// @Summary      List of creatures
// @Description  Show all creatures listed
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Success      200  {object}  CreaturesOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/creatures [get]
func tibiaCreaturesOverview(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCreaturesOverviewImpl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaCreaturesOverview")
}

// Creature godoc
// @Summary      Show one creature
// @Description  Show all information about one creature
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Param        race path string true "The race of creature" extensions(x-example=nightmare)
// @Success      200  {object}  CreatureResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/creature/{race} [get]
func tibiaCreaturesCreature(c *gin.Context) {
	// getting params from URL
	race := c.Param("race")

	// Validate the race
	endpoint, err := validation.IsCreatureNameValid(race)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=creatures&race=" + endpoint,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaCreaturesCreatureImpl(endpoint, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaCreaturesCreature")
}

// Fansites godoc
// @Summary      Promoted and supported fansites
// @Description  List of all promoted and supported fansites
// @Tags         fansites
// @Accept       json
// @Produce      json
// @Success      200  {object}  FansitesResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/fansites [get]
func tibiaFansites(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=fansites",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaFansitesImpl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaFansites")
}

// Guild godoc
// @Summary      Show one guild
// @Description  Show all information about one guild
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of guild" extensions(x-example=Elysium)
// @Success      200  {object}  GuildResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/guild/{name} [get]
func tibiaGuildsGuild(c *gin.Context) {
	// getting params from URL
	guild := c.Param("name")

	// Validate the name
	err := validation.IsGuildNameValid(guild)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=" + TibiaDataQueryEscapeString(guild),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaGuildsGuildImpl(guild, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaGuildsGuild")
}

// Guilds godoc
// @Summary      List all guilds from a world
// @Description  Show all guilds on a certain world
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        world path string true "The world" extensions(x-example=Antica)
// @Success      200  {object}  GuildsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/guilds/{world} [get]
func tibiaGuildsOverview(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitle(world)

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=guilds&world=" + TibiaDataQueryEscapeString(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaGuildsOverviewImpl(world, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaGuildsOverview")
}

// Highscores godoc
// @Summary      Highscores of tibia
// @Description  Show all highscores of tibia
// @Description  In restriction mode, the valid vocation option is all.
// @Tags         highscores
// @Accept       json
// @Produce      json
// @Param        world    path string true "The world" default(all) extensions(x-example=Antica)
// @Param        category path string true "The category" default(experience) Enums(achievements, axefighting, charmpoints, clubfighting, distancefighting, experience, fishing, fistfighting, goshnarstaint, loyaltypoints, magiclevel, shielding, swordfighting, dromescore, bosspoints, bountypoints, weeklytasks) extensions(x-example=fishing)
// @Param        vocation path string true "The vocation" default(all) Enums(all, knights, paladins, sorcerers, druids, monks) extensions(x-example=all)
// @Param        page     path int    true "The current page" default(1) minimum(1) extensions(x-example=1)
// @Success      200  {object}  HighscoresResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/highscores/{world}/{category}/{vocation}/{page} [get]
func tibiaHighscores(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	category := c.Param("category")
	vocation := c.Param("vocation")
	page := c.Param("page")

	// Check if vocation is valid
	err := validation.IsVocationValid(vocation)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	if strings.EqualFold(world, "all") {
		world = ""
	} else {
		world = TibiaDataStringWorldFormatToTitle(world)
	}

	if world != "" {
		// Check if world exists
		exists, err := validation.WorldExists(world)
		if err != nil {
			TibiaDataErrorHandler(c, err, 0)
			return
		}

		if !exists {
			TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
			return
		}
	}

	if category != "" {
		err = validation.IsHighscoreCategoryValid(category)
		if err != nil {
			TibiaDataErrorHandler(c, validation.ErrorHighscoreCategoryDoesNotExist, http.StatusBadRequest)
			return
		}
	}

	highscoreCategory := validation.HighscoreCategoryFromString(category)

	// Sanitize of vocation input
	vocationName, vocationid := TibiaDataVocationValidator(vocation)

	// Check if restriction mode is enabled
	if TibiaDataRestrictionMode && vocationName != "all" {
		TibiaDataErrorHandler(c, validation.ErrorRestrictionMode, http.StatusBadRequest)
		return
	}

	// checking the page provided
	if page == "" {
		page = "1"
	}
	if TibiaDataStringToInteger(page) < 1 {
		TibiaDataErrorHandler(c, validation.ErrorHighscorePageInvalid, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=highscores&world=" + TibiaDataQueryEscapeString(world) + "&category=" + strconv.Itoa(int(highscoreCategory)) + "&profession=" + TibiaDataQueryEscapeString(vocationid) + "&currentpage=" + TibiaDataQueryEscapeString(page),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaHighscoresImpl(world, highscoreCategory, vocationName, TibiaDataStringToInteger(page), BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaHighscores")
}

// House godoc
// @Summary      House view
// @Description  Show all information about one house
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world     path string true "The world to show" extensions(x-example=Antica)
// @Param        house_id  path int    true "The ID of the house" extensions(x-example=35019)
// @Success      200  {object}  HouseResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/house/{world}/{house_id} [get]
func tibiaHousesHouse(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	houseidStr := c.Param("house_id")

	houseid, err := strconv.Atoi(houseidStr)
	if err != nil {
		TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
		return
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitle(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// check if house exists
	exists, err = validation.HouseExistsRaw(houseid)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorHouseDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=houses&page=view&world=" + TibiaDataQueryEscapeString(world) + "&houseid=" + TibiaDataQueryEscapeString(houseidStr),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaHousesHouseImpl(houseid, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaHousesHouse")
}

// Houses godoc
// @Summary      List of houses
// @Description  Show all houses filtered on world and town
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Param        town  path string true "The town to show" extensions(x-example=Venore)
// @Success      200  {object}  HousesOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/houses/{world}/{town} [get]
// TODO: This API needs to be refactored somehow to use tibiaDataRequestHandler
func tibiaHousesOverview(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")
	town := c.Param("town")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitle(world)
	town = strings.ReplaceAll(TibiaDataStringWorldFormatToTitle(town), "+", " ")

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	// Check if town exists
	exists, err = validation.TownExists(town)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorTownDoesNotExist, http.StatusBadRequest)
		return
	}

	// Ab'Dendriel gets formatted as Ab'dendriel by TibiaDataStringWorldFormatToTitle
	// which makes tibia.com not recognize it and return an empty response.
	if strings.EqualFold(town, "ab'dendriel") {
		town = "Ab'Dendriel"
	}

	jsonData, err := TibiaHousesOverviewImpl(c, world, town, TibiaDataHTMLDataCollector)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// return jsonData
	TibiaDataAPIHandleResponse(c, "TibiaHousesOverview", jsonData)
}

// Killstatistics godoc
// @Summary      The killstatistics
// @Description  Show all killstatistics filtered on world
// @Tags         killstatistics
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Success      200  {object}  KillStatisticsResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/killstatistics/{world} [get]
func tibiaKillstatistics(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitle(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiaDataQueryEscapeString(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaKillstatisticsImpl(world, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaKillstatistics")
}

// News archive godoc
// @Summary      Show news archive (90 days)
// @Description  Show news archive with a filtering on 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/news/archive [get]
func tibiaNewslistArchive() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// News archive (with day filter) godoc
// @Summary      Show news archive (with days filter)
// @Description  Show news archive with a filtering option on days
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        days path int true "The number of days to show" default(90) minimum(1) extensions(x-example=30)
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/news/archive/{days} [get]
func tibiaNewslistArchiveDays() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// Latest news godoc
// @Summary      Show newslist (90 days)
// @Description  Show newslist with filtering on articles and news of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/news/latest [get]
func tibiaNewslistLatest() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// News ticker godoc
// @Summary      Show news tickers (90 days)
// @Description  Show news of type news tickers of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/news/newsticker [get]
func tibiaNewslist(c *gin.Context) {
	// getting params from URL
	daysStr := c.Param("days")

	var (
		days int
		err  error
	)

	if daysStr != "" {
		// convert param to int
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
			return
		}
	}

	if days == 0 {
		days = 90 // default for recent posts
	}

	// generating dates to pass to FormData
	DateBegin := time.Now().AddDate(0, 0, -days)
	DateEnd := time.Now()

	tibiadataRequest := TibiaDataRequestStruct{
		Method: http.MethodPost,
		URL:    "https://www.tibia.com/news/?subtopic=newsarchive",
		FormData: map[string]string{
			"filter_begin_day":   strconv.Itoa(DateBegin.UTC().Day()),        // period
			"filter_begin_month": strconv.Itoa(int(DateBegin.UTC().Month())), // period
			"filter_begin_year":  strconv.Itoa(DateBegin.UTC().Year()),       // period
			"filter_end_day":     strconv.Itoa(DateEnd.UTC().Day()),          // period
			"filter_end_month":   strconv.Itoa(int(DateEnd.UTC().Month())),   // period
			"filter_end_year":    strconv.Itoa(DateEnd.UTC().Year()),         // period
			"filter_cipsoft":     "cipsoft",                                  // category
			"filter_community":   "community",                                // category
			"filter_development": "development",                              // category
			"filter_support":     "support",                                  // category
			"filter_technical":   "technical",                                // category
		},
	}

	if c.Request != nil {
		// getting type of news list
		switch tmp := strings.Split(c.Request.URL.Path, "/"); tmp[3] {
		case "newsticker":
			tibiadataRequest.FormData["filter_ticker"] = "ticker"
		case "latest":
			tibiadataRequest.FormData["filter_article"] = "article"
			tibiadataRequest.FormData["filter_news"] = "news"
		case "archive":
			tibiadataRequest.FormData["filter_ticker"] = "ticker"
			tibiadataRequest.FormData["filter_article"] = "article"
			tibiadataRequest.FormData["filter_news"] = "news"
		}
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaNewslistImpl(days, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaNewslist")
}

// News entry godoc
// @Summary      Show one news entry
// @Description  Show one news entry
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        news_id path int true "The ID of news entry" extensions(x-example=6512)
// @Success      200  {object}  NewsResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/news/id/{news_id} [get]
func tibiaNews(c *gin.Context) {
	// getting params from URL
	newsIDStr := c.Param("news_id")

	// convert param to int
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		TibiaDataErrorHandler(c, validation.ErrorStringCanNotBeConvertedToInt, http.StatusBadRequest)
		return
	}

	// checking the NewsID provided
	err = validation.IsNewsIDValid(newsID)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/news/?subtopic=newsarchive&id=" + newsIDStr,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaNewsImpl(newsID, tibiadataRequest.URL, BoxContentHTML)
		},
		"TibiaNews")
}

// Spells godoc
// @Summary      List all spells
// @Description  Show all spells
// @Tags         spells
// @Accept       json
// @Produce      json
// @Success      200  {object}  SpellsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/spells [get]
func tibiaSpellsOverview(c *gin.Context) {
	// getting params from URL
	vocation := c.Param("vocation")
	if vocation == "" {
		vocation = TibiaDataDefaultVoc
	}

	err := validation.IsVocationValid(vocation)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// Sanitize of vocation input
	vocationName, _ := TibiaDataVocationValidator(vocation)
	if vocationName == "all" || vocationName == "none" {
		vocationName = ""
	} else {
		// removes the last letter (s) from the string (required for spells page)
		vocationName = strings.TrimSuffix(vocationName, "s")
		vocationName = cases.Title(language.English).String(vocationName)
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiaDataQueryEscapeString(vocationName),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaSpellsOverviewImpl(vocationName, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaSpellsOverview")
}

// Spell godoc
// @Summary      Show one spell
// @Description  Show all information about one spell
// @Tags         spells
// @Accept       json
// @Produce      json
// @Param        spell_id path string true "The name of spell" extensions(x-example=stronghaste)
// @Success      200  {object}  SpellInformationResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/spell/{spell_id} [get]
func tibiaSpellsSpell(c *gin.Context) {
	// getting params from URL
	spellRaw := c.Param("spell_id")

	spell, err := validation.IsSpellNameOrFormulaValid(spellRaw)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/library/?subtopic=spells&spell=" + spell,
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaSpellsSpellImpl(spell, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaSpellsSpell")
}

// Worlds godoc
// @Summary      List of all worlds
// @Description  Show all worlds of Tibia
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Success      200  {object}  WorldsOverviewResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/worlds [get]
func tibiaWorldsOverview(c *gin.Context) {
	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds",
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaWorldsOverviewImpl(BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaWorldsOverview")
}

// World godoc
// @Summary      Show one world
// @Description  Show all information about one world
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of world" extensions(x-example=Antica)
// @Success      200  {object}  WorldResponse
// @Failure      400  {object}  Information
// @Failure      404  {object}  Information
// @Failure      503  {object}  Information
// @Router       /v4/world/{name} [get]
func tibiaWorldsWorld(c *gin.Context) {
	// getting params from URL
	world := c.Param("name")

	// Adding fix for First letter to be upper and rest lower
	world = TibiaDataStringWorldFormatToTitle(world)

	// Check if world exists
	exists, err := validation.WorldExists(world)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	if !exists {
		TibiaDataErrorHandler(c, validation.ErrorWorldDoesNotExist, http.StatusBadRequest)
		return
	}

	tibiadataRequest := TibiaDataRequestStruct{
		Method: resty.MethodGet,
		URL:    "https://www.tibia.com/community/?subtopic=worlds&world=" + TibiaDataQueryEscapeString(world),
	}

	tibiaDataRequestHandler(
		c,
		tibiadataRequest,
		func(BoxContentHTML string) (interface{}, error) {
			return TibiaWorldsWorldImpl(world, BoxContentHTML, tibiadataRequest.URL)
		},
		"TibiaWorldsWorld")
}

func TibiaDataErrorHandler(c *gin.Context, err error, httpCode int) {
	if err == nil {
		panic(errors.New("TibiaDataErrorHandler called with nil err"))
	}

	info := Information{
		APIDetails: TibiaDataAPIDetails,
		Timestamp:  TibiaDataDatetime(""),
		TibiaURLs:  []string{},
		Status: Status{
			HTTPCode: httpCode,
		},
	}

	switch t := err.(type) {
	case validation.Error:
		if httpCode == 0 {
			if t.Code() == 10 || t.Code() == 11 {
				httpCode = http.StatusInternalServerError
			} else {
				httpCode = http.StatusBadRequest
			}
		}

		// An error occurred at tibia.com
		if t.Code() > 20000 {
			httpCode = http.StatusBadGateway
		}

		info.Status.HTTPCode = httpCode
		info.Status.Error = t.Code()
		info.Status.Message = t.Error()
	case error:
		if httpCode == 0 {
			httpCode = http.StatusBadGateway
			info.Status.HTTPCode = httpCode
		}

		info.Status.Message = err.Error()

		log.Printf("[TibiaDataErrorHandler] HTTPCode: %d], Message: %s", info.Status.HTTPCode, info.Status.Message)
	}

	var output OutInformation
	output.Information = info

	c.JSON(httpCode, output)
}

func tibiaDataRequestHandler(c *gin.Context, tibiaDataRequest TibiaDataRequestStruct, requestHandler func(string) (interface{}, error), handlerName string) {
	BoxContentHTML, err := TibiaDataHTMLDataCollector(tibiaDataRequest)
	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataErrorHandler(c, err, http.StatusBadGateway)
		return
	}

	jsonData, err := requestHandler(BoxContentHTML)
	if err != nil {
		TibiaDataErrorHandler(c, err, 0)
		return
	}

	// return jsonData
	TibiaDataAPIHandleResponse(c, handlerName, jsonData)
}

// TibiaDataAPIHandleResponse func - handling of responses..
// This should NOT be invoked if an error occured
func TibiaDataAPIHandleResponse(c *gin.Context, s string, j interface{}) {
	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] " + s + " - (" + c.Request.RequestURI + ") returned data:")
		js, err := json.Marshal(j)
		log.Printf("[debug] %s\n", js)
		if err != nil {
			log.Printf("[debug] the above had an error: %s\n", err)
		}
	}

	if TibiaDataDebug {
		log.Println("[info] " + s + " - (" + c.Request.RequestURI + ") executed successfully.")
	}

	// return successful response
	c.JSON(http.StatusOK, j)
}

// TibiaDataUserAgentGenerator func - creates User-Agent for requests
func TibiaDataUserAgentGenerator(version int) string {
	// setting product name
	useragent := "TibiaData-API/v" + strconv.Itoa(version)

	// adding details in parenthesis
	useragentDetails := []string{
		"release/" + TibiaDataBuildRelease,
		"build/" + TibiaDataBuildBuilder,
		"commit/" + TibiaDataBuildCommit,
		"edition/" + TibiaDataBuildEdition,
	}

	// adding TibiaDataHost if set
	if TibiaDataHost != "" {
		useragentDetails = append(useragentDetails, "+"+TibiaDataProtocol+"://"+TibiaDataHost)
	}

	useragent += " (" + strings.Join(useragentDetails, "; ") + ")"

	return useragent
}

// TibiaDataHTMLDataCollector func
func TibiaDataHTMLDataCollector(TibiaDataRequest TibiaDataRequestStruct) (string, error) {
	// Setting up resty client
	client := resty.New()

	// Set Debug if enabled by TibiaDataDebug var
	if TibiaDataDebug {
		client.SetDebug(true)
		client.EnableTrace()
	}

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   TibiaDataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Replace domain with proxy if env TIBIADATA_PROXY set
	if TibiaDataProxyDomain != "" {
		TibiaDataRequest.URL = strings.ReplaceAll(TibiaDataRequest.URL, "https://www.tibia.com/", TibiaDataProxyDomain)
	}

	// defining values for request
	var (
		res        *resty.Response
		err        error
		LogMessage string
	)

	switch TibiaDataRequest.Method {
	case resty.MethodPost:
		res, err = client.R().
			SetFormData(TibiaDataRequest.FormData).
			Post(TibiaDataRequest.URL)
	default:
		res, err = client.R().Get(TibiaDataRequest.URL)
	}

	if TibiaDataDebug {
		// logging trace information for resty
		TibiaDataRequestTraceLogger(res, err)
	}

	if err != nil {
		log.Printf("[error] TibiaDataHTMLDataCollector (Status: %s, URL: %s) in resp1: %s", res.Status(), res.Request.URL, err)
		return "", err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		// ok request, nothing to be done
	case http.StatusForbidden:
		// throttled request
		LogMessage = "request throttled due to rate-limitation on tibia.com"
		log.Printf("[warning] TibiaDataHTMLDataCollector: %s!", LogMessage)
		return "", validation.ErrStatusForbidden
	case http.StatusFound:
		// Check if page is in maintenance mode
		location, _ := res.RawResponse.Location()
		if location != nil && location.Host == "maintenance.tibia.com" {
			LogMessage := "maintenance mode detected on tibia.com"
			log.Printf("[info] TibiaDataHTMLDataCollector: %s!", LogMessage)
			return "", validation.ErrorMaintenanceMode
		}

		LogMessage = fmt.Sprintf(
			"unknown error occurred on tibia.com (Status: %d, RequestURL: %s)",
			http.StatusFound, res.Request.URL,
		)
		log.Printf("[error] TibiaDataHTMLDataCollector: %s!", LogMessage)
		return "", validation.ErrStatusFound
	default:
		LogMessage = fmt.Sprintf(
			"unknown error and status occurred on tibia.com (Status: %d, RequestURL: %s)",
			res.StatusCode(), res.Request.URL,
		)
		log.Printf("[error] TibiaDataHTMLDataCollector: %s!", LogMessage)
		return "", validation.ErrStatusUnknown
	}

	if TibiaDataRequest.RawBody {
		return string(res.Body()), nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		log.Printf("[error] TibiaDataHTMLDataCollector (URL: %s) error: %s", res.Request.URL, err)
		return "", err
	}

	data, err := doc.Find(".Border_2 .Border_3").Html()
	if err != nil {
		return "", err
	}

	// Return of extracted html to functions..
	return data, nil
}

// rootz is a root path (for helm-testing)
func rootz(c *gin.Context) {
	if !isReady.Load() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": http.StatusText(http.StatusServiceUnavailable)})
		return
	}
	TibiaDataAPIHandleResponse(c, "rootz", gin.H{"message": "TibiaData API up and running."})
}

// healthz is a k8s liveness probe
func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusText(http.StatusOK)})
}

// readyz is a k8s readiness probe
func readyz(c *gin.Context) {
	if !isReady.Load() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": http.StatusText(http.StatusServiceUnavailable)})
		return
	}
	TibiaDataAPIHandleResponse(c, "readyz", gin.H{"status": http.StatusText(http.StatusOK)})
}
