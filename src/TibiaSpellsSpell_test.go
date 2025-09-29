package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestFindPerson(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Find Person.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	findPersonJson, err := TibiaSpellsSpellImpl("Find Person", string(data), "https://www.tibia.com/library/?subtopic=spells&spell=findperson")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := findPersonJson.Spell
	information := findPersonJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=spells&spell=findperson", information.TibiaURLs[0])

	assert.Empty(spell.Description)
	assert.Equal("Find Person", spell.Name)
	assert.Equal("findperson", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("exiva 'name'", spell.SpellInformation.Formula)
	assert.Equal(5, len(spell.SpellInformation.Vocation))
	assert.Equal("Druid", spell.SpellInformation.Vocation[0])
	assert.Equal("Knight", spell.SpellInformation.Vocation[1])
	assert.Equal("Monk", spell.SpellInformation.Vocation[2])
	assert.Equal("Paladin", spell.SpellInformation.Vocation[3])
	assert.Equal("Sorcerer", spell.SpellInformation.Vocation[4])
	assert.False(spell.SpellInformation.GroupAttack)
	assert.False(spell.SpellInformation.GroupHealing)
	assert.True(spell.SpellInformation.GroupSupport)
	assert.True(spell.SpellInformation.TypeInstant)
	assert.False(spell.SpellInformation.TypeRune)
	assert.Equal(2, spell.SpellInformation.CooldownAlone)
	assert.Equal(2, spell.SpellInformation.CooldownGroup)
	assert.Equal(8, spell.SpellInformation.Level)
	assert.Equal(20, spell.SpellInformation.Mana)
	assert.Equal(80, spell.SpellInformation.Price)
	assert.Equal(13, len(spell.SpellInformation.City))
	assert.Equal("Ab'Dendriel", spell.SpellInformation.City[0])
	assert.Equal("Yalahar", spell.SpellInformation.City[12])
	assert.False(spell.SpellInformation.Premium)
	assert.False(spell.HasRuneInformation)
}

func TestHeavyMagicMissileRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Heavy Magic Missile Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	hmmJson, err := TibiaSpellsSpellImpl("Heavy Magic Missile Rune", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := hmmJson.Spell

	assert.Empty(spell.Description)
	assert.Equal("Heavy Magic Missile Rune", spell.Name)
	assert.Equal("heavymagicmissilerune", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("adori vis", spell.SpellInformation.Formula)
	assert.Equal(3, len(spell.SpellInformation.Vocation))
	assert.Equal("Druid", spell.SpellInformation.Vocation[0])
	assert.Equal("Monk", spell.SpellInformation.Vocation[1])
	assert.Equal("Sorcerer", spell.SpellInformation.Vocation[2])
	assert.False(spell.SpellInformation.GroupAttack)
	assert.False(spell.SpellInformation.GroupHealing)
	assert.True(spell.SpellInformation.GroupSupport)
	assert.False(spell.SpellInformation.TypeInstant)
	assert.True(spell.SpellInformation.TypeRune)
	assert.Equal(2, spell.SpellInformation.CooldownAlone)
	assert.Equal(2, spell.SpellInformation.CooldownGroup)
	assert.Equal(2, spell.SpellInformation.SoulPoints)
	assert.Equal(10, spell.SpellInformation.Amount)
	assert.Equal(25, spell.SpellInformation.Level)
	assert.Equal(350, spell.SpellInformation.Mana)
	assert.Equal(1500, spell.SpellInformation.Price)
	assert.Equal(13, len(spell.SpellInformation.City))
	assert.Equal("Ab'Dendriel", spell.SpellInformation.City[0])
	assert.Equal("Venore", spell.SpellInformation.City[11])
	assert.False(spell.SpellInformation.Premium)
	assert.True(spell.HasRuneInformation)
	assert.Equal(4, len(spell.RuneInformation.Vocation))
	assert.Equal("Druid", spell.RuneInformation.Vocation[0])
	assert.Equal("Monk", spell.RuneInformation.Vocation[1])
	assert.Equal("Paladin", spell.RuneInformation.Vocation[2])
	assert.Equal("Sorcerer", spell.RuneInformation.Vocation[3])
	assert.True(spell.RuneInformation.GroupAttack)
	assert.False(spell.RuneInformation.GroupHealing)
	assert.False(spell.RuneInformation.GroupSupport)
	assert.Equal("energy", spell.RuneInformation.DamageType)
	assert.Equal(25, spell.RuneInformation.Level)
	assert.Equal(3, spell.RuneInformation.MagicLevel)
}

func TestAnnihilation(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Annihilation.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	annihilationJson, err := TibiaSpellsSpellImpl("Annihilation", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := annihilationJson.Spell

	assert.Empty(spell.Description)
	assert.Equal("Annihilation", spell.Name)
	assert.Equal("annihilation", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("exori gran ico", spell.SpellInformation.Formula)
	assert.Equal(1, len(spell.SpellInformation.Vocation))
	assert.Equal("Knight", spell.SpellInformation.Vocation[0])
	assert.True(spell.SpellInformation.GroupAttack)
	assert.False(spell.SpellInformation.GroupHealing)
	assert.False(spell.SpellInformation.GroupSupport)
	assert.True(spell.SpellInformation.TypeInstant)
	assert.False(spell.SpellInformation.TypeRune)
	assert.Equal("var.", spell.SpellInformation.DamageType) // weird one..
	assert.Equal(30, spell.SpellInformation.CooldownAlone)
	assert.Equal(4, spell.SpellInformation.CooldownGroup)
	assert.Equal(0, spell.SpellInformation.SoulPoints)
	assert.Equal(0, spell.SpellInformation.Amount)
	assert.Equal(7, len(spell.SpellInformation.City))
	assert.True(spell.SpellInformation.Premium)
	assert.False(spell.HasRuneInformation)
}

func TestBruiseBane(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Bruise Bane.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	bruisebaneJson, err := TibiaSpellsSpellImpl("Bruise Bane", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := bruisebaneJson.Spell

	assert.Empty(spell.Description)
	assert.Equal("Bruise Bane", spell.Name)
	assert.Equal("bruisebane", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("exura infir ico", spell.SpellInformation.Formula)
	assert.Equal(1, len(spell.SpellInformation.Vocation))
	assert.False(spell.SpellInformation.GroupAttack)
	assert.True(spell.SpellInformation.GroupHealing)
	assert.False(spell.SpellInformation.GroupSupport)
	assert.True(spell.SpellInformation.TypeInstant)
	assert.False(spell.SpellInformation.TypeRune)
	assert.Equal(1, spell.SpellInformation.CooldownAlone)
	assert.Equal(1, spell.SpellInformation.CooldownGroup)
	assert.Equal(1, spell.SpellInformation.Level)
	assert.Equal(10, spell.SpellInformation.Mana)
	assert.Equal(0, spell.SpellInformation.Price)
	assert.Equal(1, len(spell.SpellInformation.City))
	assert.Equal("Dawnport", spell.SpellInformation.City[0])
	assert.False(spell.SpellInformation.Premium)
	assert.False(spell.HasRuneInformation)
}

func TestCurePoisonRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Cure Poison Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	curepoisonruneJson, err := TibiaSpellsSpellImpl("Cure Poison Rune", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := curepoisonruneJson.Spell

	assert.Empty(spell.Description)
	assert.Equal("Cure Poison Rune", spell.Name)
	assert.Equal("curepoisonrune", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("adana pox", spell.SpellInformation.Formula)
	assert.Equal(2, len(spell.SpellInformation.Vocation))
	assert.False(spell.SpellInformation.GroupAttack)
	assert.False(spell.SpellInformation.GroupHealing)
	assert.True(spell.SpellInformation.GroupSupport)
	assert.False(spell.SpellInformation.TypeInstant)
	assert.True(spell.SpellInformation.TypeRune)
	assert.False(spell.SpellInformation.Premium)
	assert.True(spell.HasRuneInformation)
	assert.False(spell.RuneInformation.GroupAttack)
	assert.True(spell.RuneInformation.GroupHealing)
	assert.False(spell.RuneInformation.GroupSupport)
}

func TestConvinceCreatureRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Convince Creature Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	convincecreatureruneJson, err := TibiaSpellsSpellImpl("Convince Creature Rune", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	spell := convincecreatureruneJson.Spell

	assert.Empty(spell.Description)
	assert.Equal("Convince Creature Rune", spell.Name)
	assert.Equal("convincecreaturerune", spell.Spell)
	assert.True(spell.HasSpellInformation)
	assert.NotNil(spell.SpellInformation)
	assert.Equal("adeta sio", spell.SpellInformation.Formula)
	assert.True(spell.HasRuneInformation)
	assert.False(spell.RuneInformation.GroupAttack)
	assert.False(spell.RuneInformation.GroupHealing)
	assert.True(spell.RuneInformation.GroupSupport)
}
