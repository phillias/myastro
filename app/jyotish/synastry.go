package jyotish

import (
	"fmt"
	"math"
	"strings"
)

const (
	orbConjunction = 8.0
	orbSextile     = 6.0
	orbSquare      = 6.0
	orbTrine       = 8.0
	orbOpposition  = 8.0
)

type Aspect struct {
	PlanetA   string  `json:"planet_a"`
	PlanetB   string  `json:"planet_b"`
	Aspect    string  `json:"aspect"`
	Orb       float64 `json:"orb"`
	PlanetALon float64 `json:"-"`
	PlanetBLon float64 `json:"-"`
}

type VedicDrishti struct {
	Planet       string `json:"planet"`
	FromChart    string `json:"from_chart"`
	AspectsHouse int    `json:"aspects_house"`
	InChart      string `json:"in_chart"`
	Significance string `json:"significance"`
}

type HouseOverlay struct {
	Planet    string `json:"planet"`
	FromChart string `json:"from_chart"`
	House     int    `json:"in_house"`
	OfChart   string `json:"of_chart"`
	HouseLord string `json:"house_lord"`
	Meaning   string `json:"meaning"`
}

type SynastryResult struct {
	ChartA       Chart           `json:"chart_a"`
	ChartB       Chart           `json:"chart_b"`
	Aspects      []Aspect        `json:"aspects"`
	VedicDrishti []VedicDrishti  `json:"vedic_drishti"`
	Overlays     []HouseOverlay  `json:"house_overlays"`
	Summary      string          `json:"summary"`
}

var houseMeanings = []string{
	"Self, identity, physical body",
	"Wealth, family, values",
	"Courage, siblings, communication",
	"Home, mother, emotional roots",
	"Creativity, children, romance",
	"Health, service, daily work",
	"Marriage, partnerships, public",
	"Transformation, occult, shared resources",
	"Philosophy, higher learning, dharma",
	"Career, status, authority",
	"Friendships, aspirations, community",
	"Spirituality, isolation, liberation",
}

// Jupiter→5/7/9, Saturn→3/7/10, Mars→4/7/8 (traditional Vedic)
var vedicDrishtiOffsets = map[string][]int{
	"Jupiter": {5, 7, 9},
	"Saturn":  {3, 7, 10},
	"Mars":    {4, 7, 8},
}

func CalculateSynastry(a, b Chart) SynastryResult {
	aspects := calculateAspects(a, b)
	vedicAspects := calculateVedicDrishtiCross(a, b)
	overlays := calculateOverlays(a, b)
	summary := generateSummary(aspects, vedicAspects, overlays, a, b)

	return SynastryResult{
		ChartA:       a,
		ChartB:       b,
		Aspects:      aspects,
		VedicDrishti: vedicAspects,
		Overlays:     overlays,
		Summary:      summary,
	}
}

func angularDistance(lon1, lon2 float64) float64 {
	d := math.Abs(lon1 - lon2)
	if d > 180 {
		d = 360 - d
	}
	return d
}

func calculateAspects(a, b Chart) []Aspect {
	var as []Aspect

	allA := append([]PlanetPosition{a.Lagna}, a.Planets...)
	allB := append([]PlanetPosition{b.Lagna}, b.Planets...)

	for _, pA := range allA {
		for _, pB := range allB {
			dist := angularDistance(pA.Longitude, pB.Longitude)

			var aspectType string
			var orb float64

			switch {
			case dist <= orbConjunction:
				aspectType = "conjunction"
				orb = dist
			case math.Abs(dist-60) <= orbSextile:
				aspectType = "sextile"
				orb = math.Abs(dist - 60)
			case math.Abs(dist-90) <= orbSquare:
				aspectType = "square"
				orb = math.Abs(dist - 90)
			case math.Abs(dist-120) <= orbTrine:
				aspectType = "trine"
				orb = math.Abs(dist - 120)
			case math.Abs(dist-180) <= orbOpposition:
				aspectType = "opposition"
				orb = math.Abs(dist - 180)
			default:
				continue
			}

			as = append(as, Aspect{
				PlanetA:   pA.Name,
				PlanetB:   pB.Name,
				Aspect:    aspectType,
				Orb:       math.Round(orb*10) / 10,
				PlanetALon: pA.Longitude,
				PlanetBLon: pB.Longitude,
			})
		}
	}

	return as
}

func calculateVedicDrishtiCross(a, b Chart) []VedicDrishti {
	var result []VedicDrishti

	for _, p := range a.Planets {
		offsets, ok := vedicDrishtiOffsets[p.Name]
		if !ok {
			continue
		}
		for _, offset := range offsets {
			targetHouse := p.House + offset
			if targetHouse > 12 {
				targetHouse -= 12
			}
			for _, bp := range b.Planets {
				if bp.House == targetHouse {
					result = append(result, VedicDrishti{
						Planet:       p.Name,
						FromChart:    "A",
						AspectsHouse: targetHouse,
						InChart:      "B",
						Significance: fmt.Sprintf("%s aspects house %d in chart B containing %s",
							p.Name, targetHouse, bp.Name),
					})
				}
			}
		}
	}

	for _, p := range b.Planets {
		offsets, ok := vedicDrishtiOffsets[p.Name]
		if !ok {
			continue
		}
		for _, offset := range offsets {
			targetHouse := p.House + offset
			if targetHouse > 12 {
				targetHouse -= 12
			}
			for _, ap := range a.Planets {
				if ap.House == targetHouse {
					result = append(result, VedicDrishti{
						Planet:       p.Name,
						FromChart:    "B",
						AspectsHouse: targetHouse,
						InChart:      "A",
						Significance: fmt.Sprintf("%s aspects house %d in chart A containing %s",
							p.Name, targetHouse, ap.Name),
					})
				}
			}
		}
	}

	return result
}

func calculateOverlays(a, b Chart) []HouseOverlay {
	var result []HouseOverlay

	for _, p := range a.Planets {
		houseB := houseFromLongitude(p.Longitude, b.Lagna.Longitude)
		lordB := houseLord(b, houseB)

		result = append(result, HouseOverlay{
			Planet:    p.Name,
			FromChart: "A",
			House:     houseB,
			OfChart:   "B",
			HouseLord: lordB,
			Meaning:   fmt.Sprintf("%s in house %d of B: %s", p.Name, houseB, houseMeanings[houseB-1]),
		})
	}

	for _, p := range b.Planets {
		houseA := houseFromLongitude(p.Longitude, a.Lagna.Longitude)
		lordA := houseLord(a, houseA)

		result = append(result, HouseOverlay{
			Planet:    p.Name,
			FromChart: "B",
			House:     houseA,
			OfChart:   "A",
			HouseLord: lordA,
			Meaning:   fmt.Sprintf("%s in house %d of A: %s", p.Name, houseA, houseMeanings[houseA-1]),
		})
	}

	return result
}

func houseLord(c Chart, house int) string {
	if house < 1 || house > 12 {
		return ""
	}
	houseLon := c.Houses[house-1]
	sign := RashiFromLongitude(houseLon)

	signLords := map[string]string{
		"Aries":       "Mars",
		"Taurus":      "Venus",
		"Gemini":      "Mercury",
		"Cancer":      "Moon",
		"Leo":         "Sun",
		"Virgo":       "Mercury",
		"Libra":       "Venus",
		"Scorpio":     "Mars",
		"Sagittarius": "Jupiter",
		"Capricorn":   "Saturn",
		"Aquarius":    "Saturn",
		"Pisces":      "Jupiter",
	}
	if lord, ok := signLords[sign]; ok {
		return lord
	}
	return ""
}

func generateSummary(aspects []Aspect, vedic []VedicDrishti, overlays []HouseOverlay, a, b Chart) string {
	var lines []string

	lines = append(lines, fmt.Sprintf("Chart A: Lagna in %s (%.1f°)", a.Lagna.Sign, a.Lagna.Degree))
	lines = append(lines, fmt.Sprintf("Chart B: Lagna in %s (%.1f°)\n", b.Lagna.Sign, b.Lagna.Degree))

	var conjunctions, sextiles, squares, trines, oppositions int
	for _, asp := range aspects {
		switch asp.Aspect {
		case "conjunction":
			conjunctions++
		case "sextile":
			sextiles++
		case "square":
			squares++
		case "trine":
			trines++
		case "opposition":
			oppositions++
		}
	}

	lines = append(lines, fmt.Sprintf("Ptolemaic Aspects: %d conjunction(s), %d sextile(s), %d square(s), %d trine(s), %d opposition(s)",
		conjunctions, sextiles, squares, trines, oppositions))

	if trines+sextiles > squares+oppositions+conjunctions {
		lines = append(lines, "Harmonious aspects (trines/sextiles) dominate — suggests natural rapport and ease.")
	} else if squares+oppositions > trines+sextiles {
		lines = append(lines, "Challenging aspects (squares/oppositions) dominate — growth through tension and contrast.")
	} else {
		lines = append(lines, "Balanced mix of harmonious and challenging aspects — potential for dynamic growth.")
	}

	for _, asp := range aspects {
		if (asp.PlanetA == "Venus" && asp.PlanetB == "Mars") || (asp.PlanetA == "Mars" && asp.PlanetB == "Venus") {
			lines = append(lines, fmt.Sprintf("\nVenus-%s-Mars connection: %s (orb %.1f°) — romantic chemistry indicator.", asp.Aspect, asp.Aspect, asp.Orb))
		}
	}

	for _, asp := range aspects {
		if (asp.PlanetA == "Moon" && asp.PlanetB == "Lagna") || (asp.PlanetA == "Lagna" && asp.PlanetB == "Moon") {
			if asp.Aspect == "conjunction" || asp.Aspect == "trine" {
				lines = append(lines, "Moon-Lagna: emotional attunement — person B resonates with person A's core self.")
			}
		}
	}

	if len(vedic) > 0 {
		vedicCounts := make(map[string]int)
		for _, v := range vedic {
			vedicCounts[v.Planet]++
		}
		var parts []string
		for p, c := range vedicCounts {
			parts = append(parts, fmt.Sprintf("%s (%d)", p, c))
		}
		lines = append(lines, fmt.Sprintf("\nVedic Graha Drishti (applied across charts): %s", strings.Join(parts, ", ")))
	}

	for _, o := range overlays {
		if o.House == 1 || o.House == 5 || o.House == 7 || o.House == 10 {
			lines = append(lines, fmt.Sprintf("\n%s", o.Meaning))
		}
	}

	if len(lines) == 0 {
		return "No significant synastry connections detected."
	}

	return strings.Join(lines, "\n")
}
