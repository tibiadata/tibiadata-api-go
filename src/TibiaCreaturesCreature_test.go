package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestDemon(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/demon.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	demonJson, err := TibiaCreaturesCreatureImpl("Demon", string(data), "https://www.tibia.com/library/?subtopic=creatures&race=demon")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	information := demonJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=creatures&race=demon", information.TibiaURLs[0])
	assert.Equal("Demons", demonJson.Creature.Name)
	assert.Equal("Demon", demonJson.Creature.Race)
	assert.Equal("https://static.tibia.com/images/library/demon.gif", demonJson.Creature.ImageURL)
	assert.Equal("The famous knight Apoc once wrote: \"Demons are the most malevolent, powerful, and dangerous creatures in Tibia. In addition to their awesome physical strength, they can also use powerful spells, such as fireballs and fire fields. Especially dangerous is their gaze, which can produce beams of pure energy to annihilate their poor victims. Moreover, they drain mana off their victims, heal themselves and summon fire elementals as their vassals.\nFortunately, Demons usually live in the deepest dungeons near hell but sometimes they appear also on the surface. When they do, they leave a track of death and destruction behind them. Nobody has ever been able to slay even one of these entities and only very few adventurers have survived an encounter with them.\"", demonJson.Creature.Description)
	assert.Equal("They are immune to fire damage and cannot be paralysed. Moreover, they are strong against death, earth, energy and physical damage. On the other hand, they are weak against holy and ice damage. These creatures can neither be summoned nor convinced. In addition, they are able to sense invisible creatures.", demonJson.Creature.Behaviour)
	assert.Equal(8200, demonJson.Creature.Hitpoints)

	assert.Equal(1, len(demonJson.Creature.ImmuneTo))
	assert.Equal("fire", demonJson.Creature.ImmuneTo[0])

	assert.Equal(4, len(demonJson.Creature.StrongAgainst))
	assert.Equal("death", demonJson.Creature.StrongAgainst[0])
	assert.Equal("earth", demonJson.Creature.StrongAgainst[1])
	assert.Equal("energy", demonJson.Creature.StrongAgainst[2])
	assert.Equal("physical", demonJson.Creature.StrongAgainst[3])

	assert.Equal(2, len(demonJson.Creature.WeaknessAgainst))
	assert.Equal("holy", demonJson.Creature.WeaknessAgainst[0])
	assert.Equal("ice", demonJson.Creature.WeaknessAgainst[1])

	assert.Empty(demonJson.Creature.HealedBy)
	assert.False(demonJson.Creature.BeParalysed)
	assert.False(demonJson.Creature.BeSummoned)
	assert.Equal(0, demonJson.Creature.SummonMana)
	assert.False(demonJson.Creature.BeConvinced)
	assert.Equal(0, demonJson.Creature.ConvincedMana)
	assert.True(demonJson.Creature.SeeInvisible)
	assert.Equal(6000, demonJson.Creature.ExperiencePoints)
	assert.True(demonJson.Creature.IsLootable)

	assert.Equal(13, len(demonJson.Creature.LootList))
	assert.Equal("assassin stars", demonJson.Creature.LootList[0])
	assert.Equal("ultimate health potions", demonJson.Creature.LootList[12])

	assert.False(demonJson.Creature.Featured)
}

func TestQuaraPredator(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/quara predator.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	quaraPredatorJson, err := TibiaCreaturesCreatureImpl("Quara Predator", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Quara Predators", quaraPredatorJson.Creature.Name)
	assert.Equal("Quara Predator", quaraPredatorJson.Creature.Race)
	assert.Equal("https://static.tibia.com/images/library/quarapredator.gif", quaraPredatorJson.Creature.ImageURL)
	assert.Equal("Even more bloodthirsty than the other quara and even feared by their own kind for their murderous frenzy, Quara Predators are usually found in the first lines of a battlefield whenever the quara decide to fight their sworn enemies (which basically means all non-quara). Obviously a mixture of quara and sharks, they share the most disgusting and most lethal attributes of both. Ruthless and vicious they know no mercy and no surrender. Quara Predators are known to fight until the last spark of life leaves their bodies. Sometimes they even fight against each other after slaying all enemies in sight.On shore they are somewhat weaker, losing their speed and agility, but their huge jaws are still fatal.", quaraPredatorJson.Creature.Description)
	assert.Equal("They are immune to fire and ice damage. On the other hand, they are weak against earth and energy damage. These creatures can neither be summoned nor convinced. In addition, they are able to sense invisible creatures.", quaraPredatorJson.Creature.Behaviour)
	assert.Equal(2200, quaraPredatorJson.Creature.Hitpoints)

	assert.Equal(2, len(quaraPredatorJson.Creature.ImmuneTo))
	assert.Equal("fire", quaraPredatorJson.Creature.ImmuneTo[0])
	assert.Equal("ice", quaraPredatorJson.Creature.ImmuneTo[1])

	assert.Equal(0, len(quaraPredatorJson.Creature.StrongAgainst))

	assert.Equal(2, len(quaraPredatorJson.Creature.WeaknessAgainst))
	assert.Equal("earth", quaraPredatorJson.Creature.WeaknessAgainst[0])
	assert.Equal("energy", quaraPredatorJson.Creature.WeaknessAgainst[1])

	assert.True(quaraPredatorJson.Creature.BeParalysed)
	assert.False(quaraPredatorJson.Creature.BeSummoned)
	assert.Equal(0, quaraPredatorJson.Creature.SummonMana)
	assert.False(quaraPredatorJson.Creature.BeConvinced)
	assert.Equal(0, quaraPredatorJson.Creature.ConvincedMana)
	assert.True(quaraPredatorJson.Creature.SeeInvisible)
	assert.Equal(1850, quaraPredatorJson.Creature.ExperiencePoints)
	assert.True(quaraPredatorJson.Creature.IsLootable)

	assert.Equal(2, len(quaraPredatorJson.Creature.LootList))
	assert.Equal("gold coins", quaraPredatorJson.Creature.LootList[0])
	assert.Equal("quara bones", quaraPredatorJson.Creature.LootList[1])

	assert.False(quaraPredatorJson.Creature.Featured)
}

func TestCentipede(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/centipede.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	centipedeJson, _ := TibiaCreaturesCreatureImpl("Centipede", string(data), "")
	assert := assert.New(t)

	assert.Equal("Centipedes", centipedeJson.Creature.Name)
	assert.Equal("Centipede", centipedeJson.Creature.Race)

	assert.True(centipedeJson.Creature.BeSummoned)
	assert.Equal(335, centipedeJson.Creature.SummonMana)
	assert.False(centipedeJson.Creature.BeConvinced)
	assert.Equal(0, centipedeJson.Creature.ConvincedMana)
}

func TestHunter(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/hunter.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	hunterJson, _ := TibiaCreaturesCreatureImpl("Hunter", string(data), "")
	assert := assert.New(t)

	assert.Equal("Hunters", hunterJson.Creature.Name)
	assert.Equal("Hunter", hunterJson.Creature.Race)

	assert.False(hunterJson.Creature.BeSummoned)
	assert.Equal(0, hunterJson.Creature.SummonMana)
	assert.True(hunterJson.Creature.BeConvinced)
	assert.Equal(530, hunterJson.Creature.ConvincedMana)
}

func TestSkunk(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/skunk.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	skunkJson, _ := TibiaCreaturesCreatureImpl("Skunk", string(data), "")
	assert := assert.New(t)

	assert.Equal("Skunks", skunkJson.Creature.Name)
	assert.Equal("Skunk", skunkJson.Creature.Race)

	assert.True(skunkJson.Creature.BeSummoned)
	assert.Equal(200, skunkJson.Creature.SummonMana)
	assert.True(skunkJson.Creature.BeConvinced)
	assert.Equal(200, skunkJson.Creature.ConvincedMana)
}

func TestLavaLurkers(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/lava lurkers.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	lavalurkersJson, _ := TibiaCreaturesCreatureImpl("Lava Lurkers", string(data), "")
	assert := assert.New(t)

	assert.Equal("Lava Lurkers", lavalurkersJson.Creature.Name)
	assert.Equal("Lava Lurkers", lavalurkersJson.Creature.Race)

	assert.Equal(1, len(lavalurkersJson.Creature.HealedBy))
	assert.Equal("fire", lavalurkersJson.Creature.HealedBy[0])
}

func TestFeralWerecrocodilesFeatured(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/creatures/creature/feral werecrocodile.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	feralwerecrocodileJson, _ := TibiaCreaturesCreatureImpl("Feral Werecrocodiles", string(data), "")
	assert := assert.New(t)

	assert.Equal("Feral Werecrocodiles", feralwerecrocodileJson.Creature.Name)
	assert.Equal("Feral Werecrocodiles", feralwerecrocodileJson.Creature.Race)

	assert.False(feralwerecrocodileJson.Creature.Featured)
}
