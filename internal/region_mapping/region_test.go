package regionmapping

import (
	"testing"

	"medidew.name/score-be-tech-test/internal/cases"
)

func TestBasicLookup(t *testing.T) {
	client := cases.Client{
		ClientID: "be23f463-0eb9-47b6-9caf-98e9a0865b6d",
		Details: cases.ClientDetails{
			FirstName: "Jane",
			LastName:  "Bloggs",
			Address: &cases.Address{
				Address1: "1 High Street",
				PostTown: "Nottingham",
				PostCode: "NG1 2ET",
			},
			DateOfBirth:      "1992/05/12",
			Dependents:       2,
			EmploymentStatus: "EMPLOYED",
		},
	}

	lookup := NewLookup()

	region, err := lookup.GetRegion(client)

	if err != nil {
		t.Errorf("failed to lookup region for %v: %v", client.Details.Address.PostCode, err)
	} else if region != East_Midlands {
		t.Errorf("wrong region %v during lookup for %v", region, client.Details.Address.PostCode)
	}

	client = cases.Client{
		ClientID: "be23f463-0eb9-47b6-9caf-98e9a0865b6d",
		Details: cases.ClientDetails{
			FirstName: "Jane",
			LastName:  "Bloggs",
			Address: &cases.Address{
				Address1: "7 Calvert Terrace",
				PostTown: "Swansea",
				PostCode: "SA1 6AY",
			},
			DateOfBirth:      "1992/05/12",
			Dependents:       2,
			EmploymentStatus: "EMPLOYED",
		},
	}

	region, err = lookup.GetRegion(client)

	if err != nil {
		t.Errorf("failed to lookup region for %v: %v", client.Details.Address.PostCode, err)
	} else if region != Wales {
		t.Errorf("wrong region %v during lookup for %v", region, client.Details.Address.PostCode)
	}
}
