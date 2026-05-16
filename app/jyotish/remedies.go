package jyotish

// Remedy holds the spiritual recommendations for a Sign
type Remedy struct {
	Mantra  string `json:"mantra"`
	Yantra  string `json:"yantra"`
	Color   string `json:"color"`
}

// SignRemedies maps signs to their Vedic recommendations
var SignRemedies = map[string]Remedy{
	"Aries":       {"Om Hram Hrim Hraum Sah Suryaya Namah", "Mangal Yantra", "Red"},
	"Taurus":      {"Om Dram Drim Draum Sah Shukraya Namah", "Shukra Yantra", "White"},
	"Gemini":      {"Om Bram Brim Braum Sah Budhaya Namah", "Budha Yantra", "Green"},
	"Cancer":      {"Om Sram Srim Sraum Sah Chandraya Namah", "Chandra Yantra", "Silver"},
	"Leo":         {"Om Hram Hrim Hraum Sah Suryaya Namah", "Surya Yantra", "Gold"},
	"Virgo":       {"Om Bram Brim Braum Sah Budhaya Namah", "Budha Yantra", "Green"},
	"Libra":       {"Om Dram Drim Draum Sah Shukraya Namah", "Shukra Yantra", "White"},
	"Scorpio":     {"Om Kram Krim Kraum Sah Kuja Namah", "Mangal Yantra", "Deep Red"},
	"Sagittarius": {"Om Gram Grim Graum Sah Gurave Namah", "Guru Yantra", "Yellow"},
	"Capricorn":   {"Om Pram Prim Praum Sah Shanaishvaraya Namah", "Shani Yantra", "Black"},
	"Aquarius":    {"Om Pram Prim Praum Sah Shanaishvaraya Namah", "Shani Yantra", "Dark Blue"},
	"Pisces":      {"Om Gram Grim Graum Sah Gurave Namah", "Guru Yantra", "Saffron"},
}

// GetRemedy fetches remedies for a given sign name
func GetRemedy(signName string) Remedy {
	if remedy, ok := SignRemedies[signName]; ok {
		return remedy
	}
	return Remedy{"Om Namah Shivaya", "General Yantra", "White"}
}
