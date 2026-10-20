package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"medidew.name/score-be-tech-test/internal/cases"
	"medidew.name/score-be-tech-test/internal/ons"
	regionmapping "medidew.name/score-be-tech-test/internal/region_mapping"
)

// Holds a mapping of each outgoing category to the ONS categories that compose it.
// If you don't agree with my mappings, feel free to change it!
var outgoingsToCategories map[string][]string = map[string][]string{
	"outgoings_water":                 {"1.2.5", "4.3"}, // "Mineral or spring waters" (1.2.5) might not count, but I am pedantic
	"outgoings_communications":        {"8"},
	"outgoings_mortgage_rent":         {"4.1", "13.1"}, // overlaps with outgoings_council_tax, but 13.1 *does* mention mortgate interest payments
	"outgoings_insurance":             {"12.4", "14.2", "14.4"},
	"outgoings_investments":           {"14.6"},
	"outgoings_council_tax":           {"13.1"},
	"outgoings_food":                  {"1.1"}, // could maybe include drink (1.2), but I did not
	"outgoings_clothing":              {"3"},
	"outgoings_other_living_costs":    {"6", "11.1", "12.1", "12.5"},               // could remove 11.1 as being non-essential
	"outgoings_entertainment":         {"9.3.1", "9.3.2", "9.3.3", "9.4", "9.5.4"}, // deliberately excluded one-off payments or contestable categories like 'pets'
	"outgoings_holidays":              {"9.6", "11.2", "13.3"},
	"outgoings_sports":                {"9.3.3", "9.4.1"}, // overlaps with outgoings_entertainment
	"outgoings_pension":               {"14.1"},
	"outgoings_car_costs":             {"7.2"}, // deliberately excluded one-off payments, such as vehicle purchases
	"outgoings_other_transport_costs": {"7.3"},
	"outgoings_child_care":            {"10", "12.1.4", "13.4.1"}, // tried not to include categories such as clothing where certain households would be penalised twice due to existing overlap
	"outgoings_fuel":                  {"4.4"},
	"outgoings_ground_rent_service_charge_shared_equity_rent": {}, // I have no idea what this is
	"outgoings_television_license":                            {"9.4.3"},
	"outgoings_household_repairs":                             {"4.2"}, // could maybe include 5.6.3, or all of 5.6 if you count 'maintenance', but 5.6.3 isn't solely repair
	"outgoings_additional_details": {
		"1.2.1", "1.2.2", "1.2.3", "1.2.4", "1.2.6", "2", "5", "7.1", "9.1", "9.2", "9.3.4", "9.3.5", "9.5.1", "9.5.2", "9.5.3", "12.2", "12.3", "13.2", "13.4.2", "13.4.3", "14.3", "14.5", "14.7", "14.8",
	}, // I put every category not included thus far into additional_details
}

type App struct {
	Mux             *http.ServeMux
	Cases           map[string]*cases.Case        // Holds a mapping of CaseIDs to each Case
	ONS_Expenditure map[string]map[string]float64 // Mapping of Category Index (e.g "1.2.4") -> Regions (e.g. "London") -> Regional Expenditure (for that category)
	Lookup          *regionmapping.Lookup         // Lookup table to get regions from postcodes
}

// Creates the 'App'. Essentially a setup function that loads case files, reads and preprocesses
// the ONSByRegion spreadsheet, loads the postcode -> region lookup map into memory, and registers
// the handler for the REST api.
func NewApp(cases_dir string, ons_path string) (*App, error) {
	loaded_cases, err := cases.LoadCases(cases_dir)
	if err != nil {
		return nil, fmt.Errorf("error loading cases: %v", err)
	}

	ons_map, err := ons.LoadONSbyRegion(ons_path)
	if err != nil {
		return nil, fmt.Errorf("error loading ons spreadsheet: %v", err)
	}

	app := &App{
		Mux:             http.NewServeMux(),
		Cases:           loaded_cases,
		ONS_Expenditure: ons_map,
		Lookup:          regionmapping.NewLookup(),
	}

	app.Mux.HandleFunc("GET /cases/{caseID}/outgoings", app.handleOutgoings)

	return app, nil
}

// Handler function to return outgoings for a specified caseID in the request's pattern.
func (app *App) handleOutgoings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("caseID")

	c := app.Cases[id]
	if c == nil {
		http.NotFound(w, r)
		return
	}

	// Look at every client's postcode and create a list of every region relevant to this case.
	case_regions := []string{}
	for _, client := range c.Clients {
		reg, err := app.Lookup.GetRegion(client)
		if err == nil {
			case_regions = append(case_regions, reg)
		}
	}

	// If not every client resolves to a region (for whatever reason), we allow this.
	// It only becomes a problem if no client exists.
	if len(case_regions) == 0 {
		err := fmt.Errorf("failed to resolve any client's region for case %v", c.CaseID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outgoings := app.calculateOutgoings(case_regions)

	outgoings_json, err := json.Marshal(outgoings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(outgoings_json)
}

// Calculates each outgoing category using the relevants regions for the case.
// In the event that there are multiple regions involved in a single case, this
// function will calculate outgoings using every ONS category's average value
// across those regions.
func (app *App) calculateOutgoings(case_regions []string) map[string]float64 {
	outgoings := make(map[string]float64, len(outgoingsToCategories))

	for outgoing_var, relevant_categories := range outgoingsToCategories {
		var total float64 = 0

		for _, region := range case_regions {
			for _, ons_category := range relevant_categories {
				total += app.ONS_Expenditure[ons_category][region]
			}
		}
		total /= float64(len(case_regions))
		outgoings[outgoing_var] = math.Round(total*100) / 100
	}

	return outgoings
}
