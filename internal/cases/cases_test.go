package cases_test

import (
	"reflect"
	"testing"

	"medidew.name/score-be-tech-test/internal/cases"
)

func TestBaseParseCase(t *testing.T) {
	test_case := `{"case_id":"56763bf7-d046-4a54-a4a1-0728c724a092","clients":[{"client_id":"be23f463-0eb9-47b6-9caf-98e9a0865b6d","details":{"first_name":"Joe","last_name":"Bloggs","address":{"address1":"1 High Street","posttown":"Nottingham","postcode":"NG1 2ET"},"dob":"1990/08/11","dependents":2,"employment_status":"EMPLOYED"}},{"client_id":"be23f463-0eb9-47b6-9caf-98e9a0865b6d","details":{"first_name":"Jane","last_name":"Bloggs","address":{"address1":"1 High Street","posttown":"Nottingham","postcode":"NG1 2ET"},"dob":"1992/05/12","dependents":2,"employment_status":"EMPLOYED"}}]}`

	expected_case := cases.Case{CaseID: "56763bf7-d046-4a54-a4a1-0728c724a092",
		Clients: []cases.Client{
			{
				ClientID: "be23f463-0eb9-47b6-9caf-98e9a0865b6d",
				Details: cases.ClientDetails{
					FirstName: "Joe",
					LastName:  "Bloggs",
					Address: &cases.Address{
						Address1: "1 High Street",
						PostTown: "Nottingham",
						PostCode: "NG1 2ET",
					},
					DateOfBirth:      "1990/08/11",
					Dependents:       2,
					EmploymentStatus: "EMPLOYED",
				},
			},
			{
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
			},
		},
	}

	result_case, err := cases.ParseCase([]byte(test_case))

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !reflect.DeepEqual(*result_case, expected_case) {
		t.Errorf("parsed case did not match expected case:\n wanted %v\n got %v\n", expected_case, *result_case)
	}
}

func TestEmptyClientsParseCase(t *testing.T) {
	test_case := `{"case_id":"56763bf7-d046-4a54-a4a1-0728c724a092","clients":[]}`

	expected_case := cases.Case{
		CaseID:  "56763bf7-d046-4a54-a4a1-0728c724a092",
		Clients: []cases.Client{},
	}

	result_case, err := cases.ParseCase([]byte(test_case))

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !reflect.DeepEqual(*result_case, expected_case) {
		t.Errorf("parsed case did not match expected case:\n wanted %v\n got %v\n", expected_case, *result_case)
	}
}
