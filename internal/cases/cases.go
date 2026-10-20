package cases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Case struct {
	CaseID  string   `json:"case_id"`
	Clients []Client `json:"clients"`
}

type Client struct {
	ClientID string        `json:"client_id"`
	Details  ClientDetails `json:"details"`
}

type ClientDetails struct {
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Address          *Address `json:"address"`
	DateOfBirth      string   `json:"dob"`
	Dependents       int      `json:"dependents"`
	EmploymentStatus string   `json:"employment_status"`
}

type Address struct {
	Address1 string `json:"address1"`
	PostTown string `json:"posttown"`
	PostCode string `json:"postcode"`
}

// Takes a json object in the same shape as case.json and unmarshals it into a usable struct.
func ParseCase(case_json []byte) (*Case, error) {
	outcase := &Case{}
	if err := json.Unmarshal(case_json, outcase); err != nil {
		return nil, fmt.Errorf("failed to parse case json: %v", err)
	}
	return outcase, nil
}

// Takes multiple case files and loads them into a map, indexed by caseID.
// Assumes that multiple cases takes the form of multiple files.
func LoadCases(dir string) (map[string]*Case, error) {
	result := map[string]*Case{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cases directory: %v", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		case_json, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read case file %v: %v", entry.Name(), err)
		}

		c, err := ParseCase(case_json)
		if err != nil {
			return nil, fmt.Errorf("failed to parse case file %v: %v", entry.Name(), err)
		}

		result[c.CaseID] = c
	}

	return result, nil
}
