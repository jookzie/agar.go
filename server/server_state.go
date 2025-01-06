package main

import (
	"sync"
	
	"github.com/gorilla/websocket"
)

type Player struct {
	Name       string  `json:"name"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Radius     float64 `json:"radius"`
	Color      string  `json:"color"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	Speed      float64 `json:"speed"`
	ClientTime int64   `json:"clientTime"`
	Connection *websocket.Conn
}

type ServerConfig struct {
	MaxX        int `json:"maxX"`
	MaxY        int `json:"maxY"`
	FeedMapSize int `json:"feedmapsize"`
}

type ServerState struct {
	Config  ServerConfig       `json:"config"`
	Players map[string]*Player `json:"players"`
	FeedMap [][2]int           `json:"feedmap"`
	mu      sync.Mutex
}

type ClientMessage struct {
	UID        string  `json:"uid"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	ClientTime int64   `json:"clientTime"`
}

var config ServerConfig = ServerConfig{
	MaxX: 5000,
	MaxY: 5000,
	FeedMapSize: 1000,
}

var state = ServerState{
	Config: config,
	Players: make(map[string]*Player),
	FeedMap: GeneratePoints(config.FeedMapSize, int(config.MaxX), int(config.MaxY)),
}

var nicknames []string = []string{
	"funkymonkey", "pixelpirate", "chilldragon", "sneakyelf", "mysticmuffin", "groovyglitch", "cyberwizard", "shadowzombie", "electricninja", "toastykitten",
	"retroreaper", "blazebutterfly", "glitchyghost", "sugarrush", "neonbat", "wildpirate", "cosmicunicorn", "lunarnightmare", "rainbowrogue", "chocolatechaos",
	"bizarrebird", "hackerhawk", "neonraven", "smokejumper", "invisibleninja", "toxicbutterfly", "scrappywolf", "cosmickitty", "fuzzywarrior", "frostedfury",
	"electricfox", "pixelatedpuma", "mysticwitch", "skywarlock", "toxicshadow", "funkyfury", "sunsetninja", "blazingbear", "crimsonhawk", "glitchghost",
	"ravenousowl", "cyberfalcon", "burntspark", "shadowsurge", "neontiger", "vortexbutterfly", "rainbowqueen", "turbohusky", "fieryrabbit", "warpwarrior",
	"icysparks", "fluffypirate", "lunarshadow", "starryspider", "smokymage", "silenttornado", "jumpybird", "flamerider", "hackerhero", "galaxywolf",
	"spookywizard", "mysticfalcon", "crimsonfox", "retrorebel", "moonlitdragon", "funkyraven", "jumpypanda", "burningkitten", "whirlwindqueen", "sillyelf",
	"candywarrior", "starrywolf", "glimmeringtiger", "spikyninja", "dizzywolf", "neonlegend", "frosteddemon", "psychopirate", "cheekysquirrel", "frostyfalcon",
	"shadowsky", "whirlwindhero", "cosmicwarrior", "stormymage", "flameshadow", "toxicwitch", "electricbutterfly", "neondruid", "crazytiger", "magmaninja",
	"frostqueen", "thunderbear", "electriclion", "bouncybird", "acidicwolf", "starlighthawk", "vampireunicorn", "chaoticfalcon", "twistedraven", "whirlwindmage",
	"psychedelicpirate", "rockybutterfly", "toxiccrow", "fuzzyfox", "crimsonwolf", "burntdemon", "nuclearraven", "neonhusky", "fluffybunny", "wildowl",
	"glitchytiger", "electricdragon", "shadowqueen", "spikemonkey", "thunderghost", "sunshineninja", "glowingraven", "icewitch", "cosmicreaper", "holohawk",
	"funkyfalcon", "neonshadow", "starlightwarrior", "burningpuma", "glitteryowl", "spookyknight", "frostedbunny", "cosmicwizard", "thundercrow", "magmawarrior",
	"crimsonmuffin", "electricrebel", "rockyhawk", "toxyspirit", "icyphoenix", "shadowpanda", "nightmareowl", "glitchreaper", "radiantninja", "solarflame",
	"darkblaze", "mysticshadow", "moonblossom", "frostedfox", "starryninja", "burningphantom", "spikeknight", "crystalwizard", "toxicdemon", "frozenpuma",
	"neonpirate", "tropicalvamp", "cosmicghost", "psychicdragon", "blazinghawk", "icyphoenix", "flamephantom", "neonscorch", "starlightreaper", "skyfirewarrior",
	"smokystar", "acidninja", "burningshadow", "electricfox", "frostybird", "stormrage", "shadowwarrior", "jumpscar", "radioshade", "toxinflare",
	"cosmicfang", "silverbeast", "neonwarlock", "mysticflame", "twilightfang", "icyshade", "lightningstorm", "puffykitten", "icycreeper", "blazingshadow",
	"darklingwolf", "starflame", "glitchyphantom", "stormhacker", "lunarnight", "psychofox", "spicywizard", "burntphoenix", "electroassassin", "stormingbeast",
	"flamedemon", "neonshadow", "bluewarrior", "mysticalwarrior", "crazyshadow", "skyhunts", "bloomstorm", "vortexflame", "frostmagic", "whisperdemon",
	"flamingspirit", "thundershadow", "fluffywarrior", "glitchninja", "frostedcreeper", "darkvenom", "pixelphantom", "blazingwraith", "moonspirit", "acidicmage",
	"neonfang", "spidernight", "glimmeringassassin", "smokewolf", "shiningcrow", "thunderousfox", "wildmage", "vampyrenight", "starlightrogue", "snowfury",
	"phantomfiend", "radiantpuma", "toxicspirit", "flameshade", "sunnyreaper", "shadowwarlock", "firehawk", "mysticwarrior", "icyfiend", "frostbitefox",
	"stormshadow", "neonwraith", "redlightwarrior", "bizarrewarrior", "electricpuma", "spiritualwolf", "icyshadow", "thunderingvamp", "shimmeringelf", "crazywizard",
	"blazingdragon", "phantomwarrior", "mysticowl", "frostreaper", "starlingpuma", "burrytiger", "burningblaze", "blazingreaper", "pulsingwolf", "nightspirit",
	"ragingunicorn", "glowingshadow", "shadowflame", "stormnight", "wildninja", "crystalmage", "phantomninja", "icefirebeast", "stormshadowhunter", "flamelight",
	"vortexbeast", "icyfang", "neonscorch", "blazingknight", "stormwild", "crystalwolf", "neoncreeper", "frostedassassin", "frostbite", "flameghost",
	"darkspirit", "whirlwindknight", "psychoassassin", "lunarshadow", "flamingdragon", "vortexreaper", "crimsonshadow", "electricknight", "neonphoenix", "mysticstorm",
	"flamescorch", "crimsonwarrior", "acidshadow", "radiantbeast", "glowingbeast", "starshadow", "glitchdrake", "frostcreeper", "darktiger", "burnoutmage",
	"icywhirlwind", "bluewarlock", "flaminghunter", "electricshadow", "vortexmage", "blazecreeper", "lunarwarrior", "frostedwarrior", "blazingwitch", "neonhunter",
	"whisperghost", "twinblaze", "mysticbeast", "fluffynight", "stormhunter", "electricfang", "frostglitch", "glowingshadow", "toxicwarlock", "shadowphoenix",
	"lightningrider", "acidicreaper", "icyking", "pixelhunter", "tornadomage", "wildwarrior", "icyblaze", "cosmicshadow", "crystalbeast", "neonpanda",
	"twilightwizard", "toxicspike", "flamingscorch", "nightshadefox", "stormcreeper", "darkknight", "burningknight", "fieryraven", "cosmicdragon", "neonshadowhunter",
	"starbeast", "glimmeringknight", "wildshadow", "shatterknight", "frostbitetiger", "spikebeast", "flamefang", "blazingdemon", "crimsonwitch", "darkcreeper",
	"magmaflame", "cosmicflare", "radiantghost", "neonnight", "toxicwarrior", "crimsonbreeze", "lunarwraith", "spikedemon", "neonphoenixwarrior", "icyassassin",
	"shadowflamer", "flamewizard", "electroscorch", "moonshadow", "frostburn", "burnedrake", "neonflame", "cosmicwarrior", "wildwitch", "fluffyscorch",
	"stormmage", "darkreaper", "fierybeast", "neonphantom", "glitchknight", "frostywarrior", "burningfury", "icyflame", "chaosdragon", "nightwarrior",
	"thunderscorch", "shadowflare", "radiantwarrior", "flamingfiend", "icyglitch", "burningreaper", "cosmicbreeze", "starlightwarrior", "moonflare", "stormwarrior",
	"blazingfox", "darknight", "toxicraven", "neonfire", "glitchraven", "cosmicwarlock", "spikecreeper", "stormdemon", "mooncreeper", "flamedemon",
	"frostshadow", "lunarflame", "spiritualfire", "electricshadowhunter", "demonicfox", "whisperingwarrior", "thunderwraith", "darkfire", "icywarlock", "shadowhunter",
	"glimmeringfire", "fieryhunter", "pixelwarrior", "shimmeringwarrior", "lunarassassin", "whisperingcreeper", "glowbeast", "chaosflame", "stormfiend", "icyphantom",
	"crimsonhunter", "blazingassassin", "darkfury", "electrophantom", "radiantwarlock", "spikewarrior", "lunarhunter", "fieryglitch", "stormspike", "fierywarrior",
	"darkmage", "moonwarrior", "vortexshadow", "stormshadowhunter", "frostinghunter", "glitchknight", "glimmerdemon", "radiantknight", "flamebeast", "mysticphoenix",
}

