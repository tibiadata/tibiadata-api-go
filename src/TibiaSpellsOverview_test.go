package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestOverviewAll(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/overviewall.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	spellsOverviewJson, err := TibiaSpellsOverviewImpl("", string(data), "https://www.tibia.com/library/?subtopic=spells")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := spellsOverviewJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=spells", information.TibiaURLs[0])

	assert.Equal(152, len(spellsOverviewJson.Spells.Spells))

	firstSpell := spellsOverviewJson.Spells.Spells[0]
	assert.Equal("Animate Dead Rune", firstSpell.Name)
	assert.Equal("animatedeadrune", firstSpell.Spell)
	assert.Equal("adana mort", firstSpell.Formula)
	assert.Equal(27, firstSpell.Level)
	assert.Equal(600, firstSpell.Mana)
	assert.Equal(1200, firstSpell.Price)
	assert.False(firstSpell.GroupAttack)
	assert.False(firstSpell.GroupHealing)
	assert.True(firstSpell.GroupSupport)
	assert.False(firstSpell.TypeInstant)
	assert.True(firstSpell.TypeRune)
	assert.True(firstSpell.PremiumOnly)

	findPersonSpell := spellsOverviewJson.Spells.Spells[60]
	assert.Equal("Find Person", findPersonSpell.Name)
	assert.Equal("findperson", findPersonSpell.Spell)
	assert.Equal("exiva \"name\"", findPersonSpell.Formula)
	assert.Equal(8, findPersonSpell.Level)
	assert.Equal(20, findPersonSpell.Mana)
	assert.Equal(80, findPersonSpell.Price)
	assert.False(findPersonSpell.GroupAttack)
	assert.False(findPersonSpell.GroupHealing)
	assert.True(findPersonSpell.GroupSupport)
	assert.True(findPersonSpell.TypeInstant)
	assert.False(findPersonSpell.TypeRune)
	assert.False(findPersonSpell.PremiumOnly)

	iceBurstSpell := spellsOverviewJson.Spells.Spells[82]
	assert.Equal("Ice Burst", iceBurstSpell.Name)
	assert.Equal(0, iceBurstSpell.Level)
	assert.Equal(0, iceBurstSpell.Price)
}

func TestOverviewDruid(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/overviewdruid.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	spellsOverviewJson, err := TibiaSpellsOverviewImpl("druid", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("druid", spellsOverviewJson.Spells.SpellsVocationFilter)
	assert.Equal(73, len(spellsOverviewJson.Spells.Spells))
}
