package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var fantasywords = []string{
	"alchemy", "astral", "arcane", "assassin", "archer", "armiger", "armorer", "archbishop", "amulet", "arrow", "axe", "abyss", "avian",
	"aegis", "altar", "aureate", "avatar", "artificer", "arboreal", "alchemist", "ancestral", "apparition", "arcanum", "aeon", "aegisguard",
	"archmage", "augury", "augology", "ardent", "aurelia", "aspectral", "aspectual", "acolyte", "bloodmage", "battlemage", "bloodmoon", "blackmage",
	"bestowment", "bestow", "biomagic", "bishop", "barbarian", "bewitch", "blacksmith", "bog", "band", "bard", "book", "bracer", "boots", "buckler",
	"baneblade", "balefire", "bolt", "bane", "brewmaster", "bifrost", "brigand", "bladebound", "barrow", "belltower", "boltflare", "banechant", "blight",
	"briarheart", "baelstrum", "brimstone", "brazen", "barren", "barrenhold", "blessing", "bazaar", "bladecraft", "bladesmith", "berserker", "brightsteel",
	"chronomancer", "chronosmith", "conduit", "conjuration", "conjure", "corrupt", "chalice", "chieftan", "chief", "champion", "chestplate", "cauldron",
	"charm", "chariot", "castle", "cove", "cave", "crown", "crest", "captain", "cloak", "corona", "cryomancer", "cyromancy", "cyrology", "cyrologist",
	"courier", "crypt", "coven", "crag", "cimmerian", "caelum", "catacomb", "chaos", "cairn", "cinderborn", "charmcaster", "cryptic", "caeli", "crimson",
	"crimsonmage", "crownsguard", "druid", "deacon", "darkling", "duke", "dame", "domain", "dungeon", "dungeonmaster", "dagger", "demonslayer", "demonic",
	"disenchant", "dispel", "dawnbreaker", "duskblade", "dreamweaver", "darkmoon", "dirge", "dervan", "divine", "divinity", "destrier", "dreamscape",
	"dreadnought", "diadem", "daemonium", "dragonheart", "swimmer", "eldritch", "elixir", "enchant", "ether", "etheric", "emperor", "empire", "entity",
	"enchantment", "enchanting", "earthmage", "eon", "elder", "echomancer", "echomancy", "ember", "embermage", "embermancy", "embermancer", "enigmancer",
	"enigmancy", "eclipsar", "enclave", "empyrean", "equinox", "ethermancy", "ethermancer", "emberforge", "embersmith", "eldertree", "everflame", "euphoria",
	"folklore", "feast", "forest", "fable", "fairytale", "farsight", "fateweaver", "firemage", "flamecaller", "firemancer", "frostmage", "frostmancer",
	"frostmancy", "firemancy", "frostcaller", "feywild", "fellblade", "fellshadow", "flameforged", "forge", "floralis", "frostshaper", "firebrand", "fulgurite",
	"forsaken", "forsworn", "gauntlets", "goblet", "gorge", "greaves", "guardian", "guard", "gravewalker", "greybeard", "glyph", "glyphsmith", "glypthforge",
	"glade", "guardianship", "gallant", "gatekeeper", "gemstone", "grimoire", "goblinoid", "glacial", "ghostfire", "glypthweaver", "geomancy", "geomancer",
	"gloomcloak", "hexal", "hexalmancer", "hexalsmith", "hexology", "honor", "honorguard", "holy", "hunter", "herald", "hexcaster", "helm", "halcyon", "harmony",
	"harmancy", "harmancer", "harology", "hinterlands", "hellfire", "hallow", "hydromancer", "hydromancy", "hollow", "hypnogem", "hypnomancer", "hypnomancy",
	"hemlock", "hallowsteel", "hexfire", "heavenstone", "horn", "hym", "hymweaver", "havenmist", "hieromagus", "helianth", "hexblade", "hydralith", "horologist",
	"horology", "huntress", "heartwood", "haruspex", "helixstaff", "imbue", "imperium", "icemage", "ivory", "inquisitor", "illusion", "invocation", "invoker",
	"infernal", "infernum", "infernal", "ironclad", "ichor", "incorporeal", "illume", "incarnate", "inlay", "illusory", "ireful", "ivy", "ironwood", "inferno",
	"jewels", "journeyman", "jester", "joker", "jestercraft", "javelin", "jewelcraft", "judgeblade", "journeycraft", "kingdom", "king", "knight", "knightess",
	"knighthood", "kinslayer", "kin", "kindred", "knavish", "kismet", "keybearer", "keystone", "kilnfire", "kylix", "kalonia", "kalonian", "karnelian", "kithara",
	"keyblade", "lair", "lantern", "legend", "lancer", "lord", "lore", "lancer", "lumen", "labryinth", "lunamancy", "lunamancer", "lunasmith", "lunalogy", "luminara",
	"lurk", "lithomancy", "lithomancer", "lithosmith", "leyline", "levitate", "luminary", "lyre", "lorestone", "lunarium", "lambent", "lithoscribe", "mage",
	"magician", "magicsmith", "majesty", "majestic", "maleficent", "medieval", "ministry", "malignant", "morph", "monarch", "mystic", "mythic", "meditate", "marauder",
	"malachite", "monolith", "meadow", "moonrise", "mummy", "mummified", "mirage", "medallion", "moonshadow", "moonlit", "mysticism", "maelstrom", "mandala", "molten",
	"mythos", "monument", "necrosmith", "necromancer", "necromancy", "necralsmith", "necralmancy", "necralmancer", "necrology", "necralology", "nightblade", "noble",
	"necropolis", "nomad", "nightshade", "netherworld", "nethermage", "nethermancer", "nethermancy", "nautical", "nocture", "nightmage", "nightwalker", "novice",
	"nyctophilia", "nex", "nexomancy", "nexomancer", "nexology", "nexer", "oracle", "occult", "ordinator", "ordain", "oracular", "overlord", "obsidian", "omen",
	"oath", "otherworld", "overture", "omnipresent", "oasis", "onyx", "outlandish", "orb", "outpost", "overcast", "opalize", "ornate", "obelisk", "orison", "pantheon",
	"paladin", "paragon", "pauldron", "phantasmal", "potion", "portal", "plain", "prophecy", "provence", "perilous", "pinnacle", "prism", "pyre", "pariah", "parchment",
	"panacea", "pangea", "penumbral", "penumbramancy", "penumbramancer", "queen", "quest", "quiver", "quagmire", "quicksilver", "quicksand", "quell", "quaint",
	"questor", "quill", "quivera", "quellion", "ring", "realm", "rogue", "runes", "runesmith", "runeblade", "runescribe", "relic", "ritual", "ruin", "reaven",
	"revenant", "runic", "rift", "riftwalker", "regal", "runebearer", "ranger", "rustic", "reliquary", "revenant", "rapture", "realmwalker", "rhapsody", "raindancer", "resonant",
	"saint", "scout", "seer", "seersword", "shard", "shadow", "shadowsmith", "shaman", "shield", "shire", "skull", "sin", "scribe", "scroll", "scrollsmith", "shard smith", "sorceress",
	"sorcery", "spell", "spellsword", "spellbook", "staff", "squire", "skirmisher", "sword", "swordmaiden", "swordsman", "sultan", "shah", "smite",
	"tale", "throne", "totem", "totemsmith", "tome", "tomesmith", "thaumaturgy", "tsar", "twilight", "talisman", "tyrant", "tenebria", "talon", "thaumaturge",
	"thundermage", "thundermancy", "thundermancer", "trinket", "transpet", "theocratic", "talismanic", "tornblade", "tranquil", "thaumic", "thaumicacy", "thaumicology",
	"tenebrae", "theurgy", "tincture", "tincturesmith", "unholy", "uprising", "unleash", "undying", "umbral", "umbralmacy", "umbralmancer", "umbralolgy", "ultraessence",
	"utterance", "unabridged", "unyielding", "unruly", "unorthodox", "valor", "valorguard", "valley", "vanguard", "vortex", "veil", "veilwalker", "vagrant", "vitality",
	"voyage", "valiant", "vagabond", "vanquish", "vellichor", "vivarium", "verdigris", "vision", "visonary", "vial", "vestige", "vesper", "vermiform", "voyant", "vellum",
	"velocimancer", "velocimancy", "vestment", "veridical", "vehemence", "vex", "vexsmith", "wand", "weird", "wicked", "wild", "wings", "wisdom", "witchcraft", "wizard",
	"wizardry", "warden", "watchman", "wanderer", "wildheart", "weald", "whirlwind", "witchery", "wayfarer", "wyrd", "warmage", "webweaver", "witching", "wanderlust",
	"warcry", "weave", "windcaller", "warbeast", "weight", "waystone", "xyloid", "xerophyte", "xeromancy", "xeormancer", "xanthism", "xanthist", "xanthilogy",
	"xanthimancer", "xyris", "xylera", "xylanth", "yggdrasil", "yore", "yarn", "yarnwraith", "yggorm", "yewstaff", "yorekeeper", "yorescroll", "yoreologist",
	"yulefire", "yester", "yelldrin", "yarnfrost", "yarnfire", "yestertide", "yonderrealm", "yestermist", "zealot", "zodiac", "zenith", "ziggurat", "zepher",
	"zephyrian", "zircon", "zephyranth", "zeal",
}

var colors = []string{
	"#FFB6C1", // LightPink
	"#ADD8E6", // LightBlue
	"#90EE90", // LightGreen
	"#FFD700", // Gold
	"#FFA07A", // LightSalmon
	"#E6E6FA", // Lavender
	"#F0E68C", // Khaki
	"#D8BFD8", // Thistle
}

type Word struct {
	Word  string `json:"word"`
	Color string `json:"color"`
}

type Words struct {
	Words []Word `json:"words"`
}

func (w Words) getWordsString() (wordString string) {
	words := make([]string, 0, len(w.Words))
	for _, word := range w.Words {
		words = append(words, word.Word)
	}

	return strings.Join(words, ",")
}

var colorIdx = 0

func getRandomWords(n int) Words {
	shuffled := make([]string, len(fantasywords))
	copy(shuffled, fantasywords)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	words := make([]Word, 0, n)
	for i := 0; i < n; i += 1 {
		words = append(words, Word{Word: shuffled[i], Color: colors[colorIdx]})
		colorIdx = (colorIdx + 1) % len(colors)
	}

	return Words{Words: words}
}

func getSavedWords(savedWords string) Words {
	splitwords := strings.Split(savedWords, ",")

	words := make([]Word, 0, len(splitwords))
	for i, word := range splitwords {
		words = append(words, Word{Word: word, Color: colors[colorIdx]})
		colorIdx = (colorIdx + 1) % len(colors)
		if i > 2 {
			break
		}
	}

	return Words{Words: words}
}

type Colors struct {
	Colors []string `json:"Colors"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	savedWords := r.URL.Query().Get("words")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var words Words
	if savedWords != "" {
		words = getSavedWords(savedWords)

	} else {
		words = getRandomWords(2 + rand.Intn(2))
	}

	component := Index(words)

	err := component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template render error: %v", err)
	}
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	words := getRandomWords(2 + rand.Intn(2))
	w.Header().Set("words", words.getWordsString())

	component := WordsDiv(words)

	err := component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Failed to render words", http.StatusInternalServerError)
	}
}

var validWordIdx = regexp.MustCompile("^[0-2]{1}$")

func extractGameId(r *http.Request) (string, error) {
	m := validWordIdx.FindStringSubmatch(r.PathValue("word"))
	if len(m) != 1 {
		return "", errors.New("invalid word idx")
	}
	return m[0], nil
}

func wordHandler(w http.ResponseWriter, r *http.Request) {
	wordIdx, err := extractGameId(r)
	if err != nil {
		http.Error(w, "Failed to retrieve word index", http.StatusInternalServerError)
		return
	}
	savedWords := r.URL.Query().Get("words")
	fmt.Println("savedWords", savedWords)
	savedWordsSplit := strings.Split(savedWords, ",")
	idx, err := strconv.Atoi(wordIdx)
	if err != nil {
		http.Error(w, "Failed to convert word index", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	words := getRandomWords(1)
	savedWordsSplit[idx] = words.Words[0].Word
	w.Header().Set("words", strings.Join(savedWordsSplit, ","))

	component := WordDiv(words.Words[0], wordIdx)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Failed to render word", http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/words/", wordsHandler)
	http.HandleFunc("/word/{word}/", wordHandler)

	log.Println("Starting server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
