package main

// These documentation-only handlers describe the current v4 API without tying
// its Swagger definitions to the runtime handlers.

// BoostableBosses godoc
//
//	@Summary		List of boostable bosses
//	@Description	Show all boostable bosses listed
//	@Tags			boostable bosses
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	BoostableBossesOverviewResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/boostablebosses [get]
func tibiaBoostableBossesV4() bool { return false }

// Character godoc
//
//	@Summary		Show one character
//	@Description	Show all information about one character available
//	@Tags			characters
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"The character name"	extensions(x-example=Trollefar)
//	@Success		200		{object}	CharacterResponse
//	@Failure		400		{object}	Information
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/character/{name} [get]
func tibiaCharactersCharacterV4() bool { return false }

// Creatures godoc
//
//	@Summary		List of creatures
//	@Description	Show all creatures listed
//	@Tags			creatures
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CreaturesOverviewResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/creatures [get]
func tibiaCreaturesOverviewV4() bool { return false }

// Creature godoc
//
//	@Summary		Show one creature
//	@Description	Show all information about one creature
//	@Tags			creatures
//	@Accept			json
//	@Produce		json
//	@Param			race	path		string	true	"The race of creature"	extensions(x-example=nightmare)
//	@Success		200		{object}	CreatureResponse
//	@Failure		400		{object}	Information
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/creature/{race} [get]
func tibiaCreaturesCreatureV4() bool { return false }

// Fansites godoc
//
//	@Summary		Promoted and supported fansites
//	@Description	List of all promoted and supported fansites
//	@Tags			fansites
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	FansitesResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/fansites [get]
func tibiaFansitesV4() bool { return false }

// Guild godoc
//
//	@Summary		Show one guild
//	@Description	Show all information about one guild
//	@Tags			guilds
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"The name of guild"	extensions(x-example=Elysium)
//	@Success		200		{object}	GuildResponse
//	@Failure		400		{object}	Information
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/guild/{name} [get]
func tibiaGuildsGuildV4() bool { return false }

// Guilds godoc
//
//	@Summary		List all guilds from a world
//	@Description	Show all guilds on a certain world
//	@Tags			guilds
//	@Accept			json
//	@Produce		json
//	@Param			world	path		string	true	"The world"	extensions(x-example=Antica)
//	@Success		200		{object}	GuildsOverviewResponse
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/guilds/{world} [get]
func tibiaGuildsOverviewV4() bool { return false }

// Highscores godoc
//
//	@Summary		Highscores of tibia
//	@Description	Show all highscores of tibia
//	@Tags			highscores
//	@Accept			json
//	@Produce		json
//	@Param			world		path		string	true	"The world"			default(all)		extensions(x-example=Antica)
//	@Param			category	path		string	true	"The category"		default(experience)	Enums(achievements, axefighting, charmpoints, clubfighting, distancefighting, experience, fishing, fistfighting, goshnarstaint, loyaltypoints, magiclevel, shielding, swordfighting, dromescore, bosspoints, bountypoints, weeklytasks, phosphorusrecord)	extensions(x-example=fishing)
//	@Param			vocation	path		string	true	"The vocation"		default(all)		Enums(all, knights, paladins, sorcerers, druids, monks)																																							extensions(x-example=knights)
//	@Param			page		path		int		true	"The current page"	default(1)			minimum(1)																																																		extensions(x-example=1)
//	@Success		200			{object}	HighscoresResponse
//	@Failure		400			{object}	Information
//	@Failure		404			{object}	Information
//	@Failure		502			{object}	Information
//	@Router			/v4/highscores/{world}/{category}/{vocation}/{page} [get]
func tibiaHighscoresV4() bool { return false }

// House godoc
//
//	@Summary		House view
//	@Description	Show all information about one house
//	@Tags			houses
//	@Accept			json
//	@Produce		json
//	@Param			world		path		string	true	"The world to show"		extensions(x-example=Antica)
//	@Param			house_id	path		int		true	"The ID of the house"	extensions(x-example=35019)
//	@Success		200			{object}	HouseResponse
//	@Failure		400			{object}	Information
//	@Failure		404			{object}	Information
//	@Failure		502			{object}	Information
//	@Router			/v4/house/{world}/{house_id} [get]
func tibiaHousesHouseV4() bool { return false }

// Houses godoc
//
//	@Summary		List of houses
//	@Description	Show all houses filtered on world and town
//	@Tags			houses
//	@Accept			json
//	@Produce		json
//	@Param			world	path		string	true	"The world to show"	extensions(x-example=Antica)
//	@Param			town	path		string	true	"The town to show"	extensions(x-example=Venore)
//	@Success		200		{object}	HousesOverviewResponse
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/houses/{world}/{town} [get]
func tibiaHousesOverviewV4() bool { return false }

// Killstatistics godoc
//
//	@Summary		The killstatistics
//	@Description	Show all killstatistics filtered on world
//	@Tags			killstatistics
//	@Accept			json
//	@Produce		json
//	@Param			world	path		string	true	"The world to show"	extensions(x-example=Antica)
//	@Success		200		{object}	KillStatisticsResponse
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/killstatistics/{world} [get]
func tibiaKillstatisticsV4() bool { return false }

// News archive godoc
//
//	@Summary		Show news archive (90 days)
//	@Description	Show news archive with a filtering on 90 days
//	@Tags			news
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	NewsListResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/news/archive [get]
func tibiaNewslistArchiveV4() bool { return false }

// News archive (with day filter) godoc
//
//	@Summary		Show news archive (with days filter)
//	@Description	Show news archive with a filtering option on days
//	@Tags			news
//	@Accept			json
//	@Produce		json
//	@Param			days	path		int	true	"The number of days to show"	default(90)	minimum(1)	extensions(x-example=30)
//	@Success		200		{object}	NewsListResponse
//	@Failure		400		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/news/archive/{days} [get]
func tibiaNewslistArchiveDaysV4() bool { return false }

// Latest news godoc
//
//	@Summary		Show newslist (90 days)
//	@Description	Show newslist with filtering on articles and news of last 90 days
//	@Tags			news
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	NewsListResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/news/latest [get]
func tibiaNewslistLatestV4() bool { return false }

// News ticker godoc
//
//	@Summary		Show news tickers (90 days)
//	@Description	Show news of type news tickers of last 90 days
//	@Tags			news
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	NewsListResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/news/newsticker [get]
func tibiaNewslistV4() bool { return false }

// News entry godoc
//
//	@Summary		Show one news entry
//	@Description	Show one news entry
//	@Tags			news
//	@Accept			json
//	@Produce		json
//	@Param			news_id	path		int	true	"The ID of news entry"	extensions(x-example=6512)
//	@Success		200		{object}	NewsResponse
//	@Failure		400		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/news/id/{news_id} [get]
func tibiaNewsV4() bool { return false }

// Spells godoc
//
//	@Summary		List all spells
//	@Description	Show all spells
//	@Tags			spells
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SpellsOverviewResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/spells [get]
func tibiaSpellsOverviewV4() bool { return false }

// Spell godoc
//
//	@Summary		Show one spell
//	@Description	Show all information about one spell
//	@Tags			spells
//	@Accept			json
//	@Produce		json
//	@Param			spell_id	path		string	true	"The name of spell"	extensions(x-example=stronghaste)
//	@Success		200			{object}	SpellInformationResponse
//	@Failure		400			{object}	Information
//	@Failure		404			{object}	Information
//	@Failure		502			{object}	Information
//	@Router			/v4/spell/{spell_id} [get]
func tibiaSpellsSpellV4() bool { return false }

// Worlds godoc
//
//	@Summary		List of all worlds
//	@Description	Show all worlds of Tibia
//	@Tags			worlds
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	WorldsOverviewResponse
//	@Failure		502	{object}	Information
//	@Router			/v4/worlds [get]
func tibiaWorldsOverviewV4() bool { return false }

// World godoc
//
//	@Summary		Show one world
//	@Description	Show all information about one world
//	@Tags			worlds
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"The name of world"	extensions(x-example=Antica)
//	@Success		200		{object}	WorldResponse
//	@Failure		404		{object}	Information
//	@Failure		502		{object}	Information
//	@Router			/v4/world/{name} [get]
func tibiaWorldsWorldV4() bool { return false }
