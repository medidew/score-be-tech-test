package regionmapping

import (
	_ "embed"
	"fmt"
	"strings"

	"medidew.name/score-be-tech-test/internal/cases"
)

const ( // Nice and neat standardisation of region names and codes. No typos allowed.
	North_East       string = "North East"
	North_West       string = "North West"
	Yorkshire_Humber string = "Yorkshire and the Humber"
	East_Midlands    string = "East Midlands"
	West_Midlands    string = "West Midlands"
	East_England     string = "East"
	London           string = "London"
	South_East       string = "South East"
	South_West       string = "South West"
	Wales            string = "Wales"
	Scotland         string = "Scotland"
	Northern_Ireland string = "Northern Ireland"

	North_East_Code       string = "E12000001"
	North_West_Code       string = "E12000002"
	Yorkshire_Humber_Code string = "E12000003"
	East_Midlands_Code    string = "E12000004"
	West_Midlands_Code    string = "E12000005"
	East_England_Code     string = "E12000006"
	London_Code           string = "E12000007"
	South_East_Code       string = "E12000008"
	South_West_Code       string = "E12000009"
	Wales_Code            string = "W99999999"
	Scotland_Code         string = "S99999999"
	Northern_Ireland_Code string = "N99999999"
)

// Slice of all valid regions in the ONS spreadsheet.
var Regions = []string{North_East, North_West, Yorkshire_Humber, East_Midlands, West_Midlands, East_England, London, South_East, South_West, Wales, Scotland, Northern_Ireland}

var codeToRegion map[string]string = map[string]string{
	North_East_Code:       North_East,
	North_West_Code:       North_West,
	Yorkshire_Humber_Code: Yorkshire_Humber,
	East_Midlands_Code:    East_Midlands,
	West_Midlands_Code:    West_Midlands,
	East_England_Code:     East_England,
	London_Code:           London,
	South_East_Code:       South_East,
	South_West_Code:       South_West,
	Wales_Code:            Wales,
	Scotland_Code:         Scotland,
	Northern_Ireland_Code: Northern_Ireland,
}

//go:embed sectors.tsv
var sectorData string

//go:embed units.tsv
var unitData string

type Lookup struct {
	sectors   map[string]string
	postcodes map[string]string
}

func NewLookup() *Lookup {
	lookup := &Lookup{
		sectors:   make(map[string]string, 12428), // number of entries in sectors.tsv
		postcodes: make(map[string]string, 96703), // number of entries in units.tsv
	}

	for _, line := range strings.Split(sectorData, "\n") {
		entries := strings.Split(line, "\t")
		if len(entries) == 2 {
			lookup.sectors[entries[0]] = entries[1]
		}
	}

	for _, line := range strings.Split(unitData, "\n") {
		entries := strings.Split(line, "\t")
		if len(entries) == 2 {
			lookup.postcodes[entries[0]] = entries[1]
		}
	}

	return lookup
}

func (lookup *Lookup) GetRegion(client cases.Client) (string, error) {
	if client.Details.Address == nil {
		return "", fmt.Errorf("client %v has no address", client.ClientID)
	}

	postcode := strings.ReplaceAll(client.Details.Address.PostCode, " ", "") // strip whitespace

	if len(postcode) < 6 {
		return "", fmt.Errorf("postcode '%v' is too small/invalid", postcode)
	}

	sector := postcode[:len(postcode)-2]

	var res string

	if res = lookup.postcodes[postcode]; res != "" {
		return codeToRegion[res], nil
	} else if res = lookup.sectors[sector]; res != "" {
		return codeToRegion[res], nil
	}

	return "", fmt.Errorf("region lookup failed: postcode not found")
}
