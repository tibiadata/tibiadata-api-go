package main

// InformationV3 stores some API related data
type InformationV3 struct {
	APIversion int    `json:"api_version"` // The API major version currently running.
	Timestamp  string `json:"timestamp"`   // The timestamp from when the data was processed.
}

// BoostableBosses godoc
// @Summary      List of boostable bosses
// @Description  Show all boostable bosses listed
// @Tags         boostable bosses
// @Accept       json
// @Produce      json
// @Success      200  {object}  BoostableBossesOverviewResponseV3
// @Router       /v3/boostablebosses [get]
// @Deprecated
func tibiaBoostableBossesV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type BoostableBossesOverviewResponseV3 struct {
	BoostableBosses BoostableBossesContainer `json:"boostable_bosses"`
	Information     InformationV3            `json:"information"`
}

// Character godoc
// @Summary      Show one character
// @Description  Show all information about one character available
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        name path string true "The character name" extensions(x-example=Trollefar)
// @Success      200  {object}  CharacterResponseV3
// @Router       /v3/character/{name} [get]
// @Deprecated
func tibiaCharactersCharacterV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type CharacterV3 struct {
	CharacterInfo      CharacterInfo      `json:"character"`                     // The character's information.
	AccountBadges      []AccountBadges    `json:"account_badges,omitempty"`      // The account's badges.
	Achievements       []Achievements     `json:"achievements,omitempty"`        // The character's achievements.
	Deaths             []Deaths           `json:"deaths,omitempty"`              // The character's deaths.
	AccountInformation AccountInformation `json:"account_information,omitempty"` // The account information.
	OtherCharacters    []OtherCharacters  `json:"other_characters,omitempty"`    // The account's other characters.
}
type CharacterResponseV3 struct {
	Characters  CharacterV3   `json:"characters"`
	Information InformationV3 `json:"information"`
}

// Creatures godoc
// @Summary      List of creatures
// @Description  Show all creatures listed
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Success      200  {object}  CreaturesOverviewResponseV3
// @Router       /v3/creatures [get]
// @Deprecated
func tibiaCreaturesOverviewV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type CreaturesOverviewResponseV3 struct {
	Creature    Creature      `json:"creature"`
	Information InformationV3 `json:"information"`
}

// Creature godoc
// @Summary      Show one creature
// @Description  Show all information about one creature
// @Tags         creatures
// @Accept       json
// @Produce      json
// @Param        race path string true "The race of creature" extensions(x-example=nightmare)
// @Success      200  {object}  CreatureResponseV3
// @Router       /v3/creature/{race} [get]
// @Deprecated
func tibiaCreaturesCreatureV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type CreatureResponseV3 struct {
	Creatures   CreaturesContainer `json:"creatures"`
	Information InformationV3      `json:"information"`
}

// Fansites godoc
// @Summary      Promoted and supported fansites
// @Description  List of all promoted and supported fansites
// @Tags         fansites
// @Accept       json
// @Produce      json
// @Success      200  {object}  FansitesResponseV3
// @Router       /v3/fansites [get]
// @Deprecated
func tibiaFansitesV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type FansitesResponseV3 struct {
	Fansites    Fansites      `json:"fansites"`
	Information InformationV3 `json:"information"`
}

// Guild godoc
// @Summary      Show one guild
// @Description  Show all information about one guild
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of guild" extensions(x-example=Elysium)
// @Success      200  {object}  GuildResponseV3
// @Router       /v3/guild/{name} [get]
// @Deprecated
func tibiaGuildsGuildV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type GuildResponseV3 struct {
	Guilds      GuildV3       `json:"guilds"`
	Information InformationV3 `json:"information"`
}
type GuildV3 struct {
	Guild Guild `json:"guild"`
}

// Guilds godoc
// @Summary      List all guilds from a world
// @Description  Show all guilds on a certain world
// @Tags         guilds
// @Accept       json
// @Produce      json
// @Param        world path string true "The world" extensions(x-example=Antica)
// @Success      200  {object}  GuildsOverviewResponseV3
// @Router       /v3/guilds/{world} [get]
// @Deprecated
func tibiaGuildsOverviewV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type GuildsOverviewResponseV3 struct {
	Guilds      OverviewGuilds `json:"guilds"`
	Information InformationV3  `json:"information"`
}

// Highscores godoc
// @Summary      Highscores of tibia
// @Description  Show all highscores of tibia
// @Tags         highscores
// @Accept       json
// @Produce      json
// @Param        world    path string true "The world" default(all) extensions(x-example=Antica)
// @Param        category path string true "The category" default(experience) Enums(achievements, axefighting, charmpoints, clubfighting, distancefighting, experience, fishing, fistfighting, goshnarstaint, loyaltypoints, magiclevel, shielding, swordfighting, dromescore, bosspoints) extensions(x-example=fishing)
// @Param        vocation path string true "The vocation" default(all) Enums(all, knights, paladins, sorcerers, druids, monks) extensions(x-example=knights)
// @Param        page     path int    true "The current page" default(1) minimum(1) extensions(x-example=1)
// @Success      200  {object}  HighscoresResponseV3
// @Router       /v3/highscores/{world}/{category}/{vocation}/{page} [get]
// @Deprecated
func tibiaHighscoresV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type HighscoresResponseV3 struct {
	Highscores  Highscores    `json:"highscores"`
	Information InformationV3 `json:"information"`
}

// House godoc
// @Summary      House view
// @Description  Show all information about one house
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world     path string true "The world to show" extensions(x-example=Antica)
// @Param        house_id  path int    true "The ID of the house" extensions(x-example=35019)
// @Success      200  {object}  HouseResponseV3
// @Router       /v3/house/{world}/{house_id} [get]
// @Deprecated
func tibiaHousesHouseV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type HouseResponseV3 struct {
	House       House         `json:"house"`
	Information InformationV3 `json:"information"`
}

// Houses godoc
// @Summary      List of houses
// @Description  Show all houses filtered on world and town
// @Tags         houses
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Param        town  path string true "The town to show" extensions(x-example=Venore)
// @Success      200  {object}  HousesOverviewResponseV3
// @Router       /v3/houses/{world}/{town} [get]
// @Deprecated
func tibiaHousesOverviewV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type HousesOverviewResponseV3 struct {
	Houses      HousesHouses  `json:"houses"`
	Information InformationV3 `json:"information"`
}

// Killstatistics godoc
// @Summary      The killstatistics
// @Description  Show all killstatistics filtered on world
// @Tags         killstatistics
// @Accept       json
// @Produce      json
// @Param        world path string true "The world to show" extensions(x-example=Antica)
// @Success      200  {object}  KillStatisticsResponseV3
// @Router       /v3/killstatistics/{world} [get]
// @Deprecated
func tibiaKillstatisticsV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type KillStatisticsResponseV3 struct {
	KillStatistics KillStatistics `json:"killstatistics"`
	Information    InformationV3  `json:"information"`
}

// News archive godoc
// @Summary      Show news archive (90 days)
// @Description  Show news archive with a filtering on 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponseV3
// @Router       /v3/news/archive [get]
// @Deprecated
func tibiaNewslistArchiveV3() bool {
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
// @Success      200  {object}  NewsListResponseV3
// @Router       /v3/news/archive/{days} [get]
// @Deprecated
func tibiaNewslistArchiveDaysV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// Latest news godoc
// @Summary      Show newslist (90 days)
// @Description  Show newslist with filtering on articles and news of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponseV3
// @Router       /v3/news/latest [get]
// @Deprecated
func tibiaNewslistLatestV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

// News ticker godoc
// @Summary      Show news tickers (90 days)
// @Description  Show news of type news tickers of last 90 days
// @Tags         news
// @Accept       json
// @Produce      json
// @Success      200  {object}  NewsListResponseV3
// @Router       /v3/news/newsticker [get]
// @Deprecated
func tibiaNewslistV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type NewsListResponseV3 struct {
	News        News          `json:"news"`
	Information InformationV3 `json:"information"`
}

// News entry godoc
// @Summary      Show one news entry
// @Description  Show one news entry
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        news_id path int true "The ID of news entry" extensions(x-example=6512)
// @Success      200  {object}  NewsResponseV3
// @Router       /v3/news/id/{news_id} [get]
// @Deprecated
func tibiaNewsV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type NewsResponseV3 struct {
	News        []NewsItem    `json:"news"`
	Information InformationV3 `json:"information"`
}

// Spells godoc
// @Summary      List all spells
// @Description  Show all spells
// @Tags         spells
// @Accept       json
// @Produce      json
// @Success      200  {object}  SpellsOverviewResponseV3
// @Router       /v3/spells [get]
// @Deprecated
func tibiaSpellsOverviewV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type SpellsOverviewResponseV3 struct {
	Spells      Spells        `json:"spells"`
	Information InformationV3 `json:"information"`
}

// Spell godoc
// @Summary      Show one spell
// @Description  Show all information about one spell
// @Tags         spells
// @Accept       json
// @Produce      json
// @Param        spell_id path string true "The name of spell" extensions(x-example=stronghaste)
// @Success      200  {object}  SpellInformationResponseV3
// @Router       /v3/spell/{spell_id} [get]
// @Deprecated
func tibiaSpellsSpellV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type SpellInformationResponseV3 struct {
	Spells      SpellV3       `json:"spells"`
	Information InformationV3 `json:"information"`
}
type SpellV3 struct {
	Spell SpellData `json:"spell"`
}

// Worlds godoc
// @Summary      List of all worlds
// @Description  Show all worlds of Tibia
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Success      200  {object}  WorldsOverviewResponseV3
// @Router       /v3/worlds [get]
// @Deprecated
func tibiaWorldsOverviewV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type WorldsOverviewResponseV3 struct {
	Worlds      OverviewWorlds `json:"worlds"`
	Information InformationV3  `json:"information"`
}

// World godoc
// @Summary      Show one world
// @Description  Show all information about one world
// @Tags         worlds
// @Accept       json
// @Produce      json
// @Param        name path string true "The name of world" extensions(x-example=Antica)
// @Success      200  {object}  WorldResponseV3
// @Router       /v3/world/{name} [get]
// @Deprecated
func tibiaWorldsWorldV3() bool {
	// Not used function.. but required for documentation purpose
	return false
}

type WorldResponseV3 struct {
	Worlds      WorldV3       `json:"worlds"`
	Information InformationV3 `json:"information"`
}
type WorldV3 struct {
	World World `json:"world"`
}
