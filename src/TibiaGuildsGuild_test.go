package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tibiadata/tibiadata-api-go/src/static"
)

func TestOrderofGlory(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Order of Glory.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	orderOfGloryJson, err := TibiaGuildsGuildImpl("Order of Glory", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := orderOfGloryJson.Guild

	assert.Equal("Order of Glory", guild.Name)
	assert.Equal("Premia", guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Order_of_Glory.gif", guild.LogoURL)
	assert.Equal("We are an English speaking guild of friends and allies from around the world who seek only peaceful questing, exploring, team hunts and a chill place to hang out. Message any of our leaders for an invitation. Contact Zyb with any problems.", guild.Description)
	assert.Nil(guild.Guildhalls)
	assert.True(guild.Active)
	assert.Equal("2020-06-27", guild.Founded)
	assert.True(guild.Applications)
	assert.Empty(guild.Homepage)
	assert.False(guild.InWar)
	assert.Empty(guild.DisbandedDate)
	assert.Empty(guild.DisbandedCondition)
	assert.Equal(1, guild.PlayersOnline)
	assert.Equal(32, guild.PlayersOffline)
	assert.Equal(33, guild.MembersTotal)
	assert.Equal(0, guild.MembersInvited)
	assert.Equal(33, len(guild.Members))

	guildLeader := guild.Members[0]
	assert.Equal("Zyb the Warrior", guildLeader.Name)
	assert.Empty(guildLeader.Title)
	assert.Equal("Leader", guildLeader.Rank)
	assert.Equal("Elite Knight", guildLeader.Vocation)
	assert.Equal(385, guildLeader.Level)
	assert.Equal("2020-10-13", guildLeader.Joined)
	assert.Equal("online", guildLeader.Status)

	assert.Nil(guild.Invited)
}

func TestElysium(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Elysium.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	elysiumJson, err := TibiaGuildsGuildImpl("Elysium", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := elysiumJson.Guild

	assert.Equal("Elysium", guild.Name)
	assert.Equal("Vunira", guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Elysium.gif", guild.LogoURL)
	assert.Equal("The place you want to be...\nIt is the land of peace and harmony, the home of the immortal, the blessed, home of the passed away legends... Hail all defenders of righteousness and the old virtues which shall never be forgotten!\nIf you would like to join us, feel free to contact one of our leaders.", guild.Description)
	assert.NotNil(guild.Guildhalls)
	assert.Equal("Ab'Dendriel Clanhall", guild.Guildhalls[0].Name)
	assert.Equal("2023-02-18", guild.Guildhalls[0].PaidUntil)
	assert.Equal("Vunira", guild.Guildhalls[0].World)
	assert.True(guild.Active)
	assert.Equal("2004-05-26", guild.Founded)
	assert.True(guild.Applications)
	assert.Empty(guild.Homepage)
	assert.False(guild.InWar)
	assert.Empty(guild.DisbandedDate)
	assert.Empty(guild.DisbandedCondition)
	assert.Equal(4, guild.PlayersOnline)
	assert.Equal(154, guild.PlayersOffline)
	assert.Equal(158, guild.MembersTotal)
	assert.Equal(1, guild.MembersInvited)
	assert.Equal(158, len(guild.Members))

	guildFollower := guild.Members[101]
	assert.Equal("Trollefar", guildFollower.Name)
	assert.Equal("Troll Giant", guildFollower.Title)
	assert.Equal("Follower", guildFollower.Rank)
	assert.Equal("Elite Knight", guildFollower.Vocation)
	assert.Equal(202, guildFollower.Level)
	assert.Equal("2013-10-20", guildFollower.Joined)
	assert.Equal("offline", guildFollower.Status)

	assert.NotNil(guild.Invited)
	evelynInvite := guild.Invited[0]
	assert.Equal("Evelyn Earlong", evelynInvite.Name)
	assert.Equal("2023-01-20", evelynInvite.Date)
}

func TestMercenarys(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Mercenarys.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	mercenarysJson, err := TibiaGuildsGuildImpl("Mercenarys", string(data), "https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=Mercenarys")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := mercenarysJson.Guild

	assert.Equal("Mercenarys", guild.Name)
	assert.Equal("Antica", guild.World)
	assert.Equal("https://static.tibia.com/images/guildlogos/Mercenarys.gif", guild.LogoURL)
	assert.NotNil(guild.Guildhalls)
	assert.Equal("Mercenary Tower", guild.Guildhalls[0].Name)
	assert.Equal("2023-01-28", guild.Guildhalls[0].PaidUntil)
	assert.Equal("Antica", guild.Guildhalls[0].World)
	assert.True(guild.Active)
	assert.Equal("2002-02-18", guild.Founded)
	assert.True(guild.Applications)
	assert.Equal("http://www.mercenarys.net", guild.Homepage)
	assert.False(guild.InWar)
	assert.Equal("2023-02-07", guild.DisbandedDate)
	assert.Equal("if there are still less than four vice leaders or an insufficient amount of premium accounts in the leading ranks by then", guild.DisbandedCondition)

	information := mercenarysJson.Information
	assert.Equal("https://www.tibia.com/community/?subtopic=guilds&page=view&GuildName=Mercenarys", information.TibiaURLs[0])
}

func TestKotkiAntica(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Kotki Antica.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	kotkianticaJson, err := TibiaGuildsGuildImpl("Kotki Antica", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := kotkianticaJson.Guild

	assert.Equal("Kotki Antica", guild.Name)
	assert.Equal("Antica", guild.World)
	assert.Empty(guild.Description)
	assert.True(guild.Active)
	assert.Equal("2021-09-22", guild.Founded)
	assert.False(guild.Applications)
}

func TestNightsWatch(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/Nights Watch.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	nightswatchJson, err := TibiaGuildsGuildImpl("Nights Watch", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := nightswatchJson.Guild

	assert.Equal("Nights Watch", guild.Name)
	assert.Equal("Luminera", guild.World)
	assert.Empty(guild.Description)
	assert.True(guild.Active)
	assert.True(guild.InWar)
	assert.Equal("2022-09-25", guild.Founded)
	assert.False(guild.Applications)
}

func TestTruePlayers(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/guilds/guild/True Players.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	trueplayersJson, err := TibiaGuildsGuildImpl("True Players", string(data), "")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	guild := trueplayersJson.Guild

	assert.Equal("True Players", guild.Name)
	assert.Equal("Karmeya", guild.World)
	assert.Empty(guild.Description)
	assert.True(guild.Active)
	assert.Equal("2024-11-24", guild.Founded)
	assert.True(guild.Applications)

	guildLeader := guild.Members[0]
	assert.Equal("Loo Mind Picture", guildLeader.Name)
	assert.Equal("We", guildLeader.Rank)
	assert.Equal("Master Sorcerer", guildLeader.Vocation)
	assert.Equal(606, guildLeader.Level)

	guildViceleader := guild.Members[1]
	assert.Equal("Emres", guildViceleader.Name)
	assert.Equal("Shine", guildViceleader.Rank)
	assert.Equal("Elder Druid", guildViceleader.Vocation)
	assert.Equal(81, guildViceleader.Level)
}
