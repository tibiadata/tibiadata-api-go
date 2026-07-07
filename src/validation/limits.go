package validation

const (
	MaxRunesAllowedInACharacterName     = 29 // Larggest character name length possible
	MinRunesAllowedInACharacterName     = 2  // Smallest character name length possible
	MaxRunesAllowedInACharacterNameWord = 20 // Larggest character name word length possible (new names only 14, but older ones are longer)
	MinRunesAllowedInACharacterNameWord = 1  // Smallest character name word length possible

	MaxRunesAllowedInAGuildName     = 29 // Larggest guild name length possible
	MinRunesAllowedInAGuildName     = 2  // Smallest guild name length possible
	MaxRunesAllowedInAGuildNameWord = 14 // Larggest guild name word length possible
	MinRunesAllowedInAGuildNameWord = 2  // Smallest character name word length possible

	// AmountOfBoostableBosses is the amount of boostable bosses.
	// Last updated: Feb 03 2025
	AmountOfBoostableBosses = 99
)
