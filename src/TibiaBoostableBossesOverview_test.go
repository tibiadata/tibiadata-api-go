package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestBoostableBossesOverview(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boostableBossesJson, _ := TibiaBoostableBossesOverviewImpl(string(data), "https://www.tibia.com/library/?subtopic=boostablebosses")
	assert := assert.New(t)
	boosted := boostableBossesJson.BoostableBosses.Boosted
	bosses := boostableBossesJson.BoostableBosses.BoostableBosses
	information := boostableBossesJson.Information

	assert.Equal("https://www.tibia.com/library/?subtopic=boostablebosses", information.TibiaURLs[0])
	assert.Equal(98, len(bosses))
	assert.Equal("Ragiaz", boosted.Name)
	assert.Equal(
		"https://static.tibia.com/images/global/header/monsters/ragiaz.gif",
		boosted.ImageURL,
	)

	for _, tc := range []struct {
		idx      int
		name     string
		featured bool
		imageURL string
	}{
		{
			idx:      21,
			name:     "Gnomevil",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/gnomehorticulist.gif",
		},
		{
			idx:      26,
			name:     "Goshnar's Malice",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/goshnarsmalice.gif",
		},
		{
			idx:      49,
			name:     "Ragiaz",
			featured: true,
			imageURL: "https://static.tibia.com/images/library/ragiaz.gif",
		},
		{
			idx:      57,
			name:     "Sharpclaw",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/sharpclaw.gif",
		},
		{
			idx:      80,
			name:     "The Pale Worm",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/paleworm.gif",
		},
	} {
		boss := bosses[tc.idx]
		assert.Equal(
			tc.name, boss.Name,
			"Wrong name\nidx: %d (%s)\nwant: %s\ngot: %s",
			tc.idx, tc.name, tc.name, boss.Name,
		)
		assert.Equal(
			tc.featured, boss.Featured,
			"Wrong featured status\nidx: %d (%s)\nwant: %v\ngot: %v",
			tc.idx, tc.name, tc.featured, boss.Featured,
		)
		assert.Equal(
			tc.imageURL, boss.ImageURL,
			"Wrong image URL\nidx: %d (%s)\nwant: %s\ngot: %s",
			tc.idx, tc.name, tc.imageURL, boss.ImageURL,
		)
	}
}

var bossSink BoostableBossesOverviewResponse

func BenchmarkTibiaBoostableBossesOverviewImpl(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	rawData, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}
	data := string(rawData)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bossSink, _ = TibiaBoostableBossesOverviewImpl(data, "https://www.tibia.com/library/?subtopic=boostablebosses")
	}
}
