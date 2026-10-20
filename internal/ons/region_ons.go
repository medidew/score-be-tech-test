package ons

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	regionmapping "medidew.name/score-be-tech-test/internal/region_mapping"
)

// Loads the ONS by Region spreadsheet and returns a mapping of Category Index (e.g "1.2.4") -> Regions -> Regional Expenditure (for that category).
//
// For example:
//
//	result["Wales"]["1.1.1"] == 5.00 // OSNByRegion2019.xlsx @Cell_P5
func LoadONSbyRegion(path string) (map[string]map[string]float64, error) {

	/*
		Loading the ONS spreadsheet and performing some basic input sanitization
	*/

	spreadsheet_file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load regional spreadsheet: %v", err)
	}
	defer spreadsheet_file.Close()

	rows, err := spreadsheet_file.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to get spreadsheet rows: %v", err)
	}

	if len(rows) < 1 { // sanity check
		return nil, fmt.Errorf("spreadsheet %v has no rows", path)
	}

	/*
		Parsing the spreadsheet
	*/

	headermap := getHeaders(rows[0])
	result_expenditure := map[string]map[string]float64{}
	//total_categories_seen := 0

	valid_category := regexp.MustCompile(`^\d+$`)
	valid_sub_category := regexp.MustCompile(`^\d+\.\d+$`)
	valid_sub_sub_category := regexp.MustCompile(`^\d+\.\d+\.\d+$`) // yes, really.

	// needed because ONS spreadsheet data separates the category row from float entries in formatting
	last_seen_category := ""

	// spreadsheet parsing
	for row_index := 1; row_index < len(rows); row_index++ {
		row := rows[row_index]

		if len(row) > 1 { // sanity check
			category := row[0]
			sub_category := row[1]

			// tracking the current category
			if valid_category.MatchString(category) {
				last_seen_category = category
				//total_categories_seen++
			} else if valid_sub_category.MatchString(category) {
				last_seen_category = category
				//total_categories_seen++
			} else if category != "" { // one-off to skip certain total/summary lines
				// fmt.Printf("skipping evil summary row: %v\n", row)
				continue
			} else if valid_sub_sub_category.MatchString(sub_category) {
				last_seen_category = sub_category
				//total_categories_seen++
			}

			category_expenditure, err := getCategoryExpenditure(row, headermap)

			if err == nil { // successful row parse

				/* debugging stuff
				if _, exists := result_expenditure[last_seen_category]; exists {
					fmt.Printf("warning: overwriting values for category %s at row %d", last_seen_category, row_index)
				}
				*/
				result_expenditure[last_seen_category] = category_expenditure
			} else { // some parses will fail because of the spreadsheet formatting splitting some entries' descriptions across two rows

				/* debugging stuff
				fmt.Printf("failed to parse row: %v\n", err)
				fmt.Printf("row: %v\n", row)
				fmt.Printf("category_expenditure: %v\n", category_expenditure)
				*/
			}
		}

		// The following line isn't needed anymore, but was used to debug parsing behaviour.
		// In a perfect world I would write a test for this file instead, but considering
		// this will only ever parse a single spreadsheet that produces a massive output, and
		// I would have to write a dummy spreadsheet covering all of the edge cases present in
		// the original to reliably test it otherwise, I've decided that would be
		// overengineering for a take-home test.

		// fmt.Printf("parsed %d categories out of %d seen\n", len(result_expenditure), total_categories_seen) // expected three missing
	}

	return result_expenditure, nil
}

// Given the topmost row from the spreadsheet, returns a map of column indexes to region names.
func getHeaders(header_row []string) map[int]string {
	headermap := map[int]string{}

	for index, cell_value := range header_row {
		if slices.Contains(regionmapping.Regions, cell_value) {
			headermap[index] = cell_value
		}
	}

	return headermap
}

// Parses a row of regional expenses and returns a map of what region is associated with what expense in that row.
func getCategoryExpenditure(row []string, headermap map[int]string) (map[string]float64, error) {
	category_expenditure := map[string]float64{}
	non_zero := false // necessary for the case where an empty row without any special symbols occurs

	for column_index, cell_value := range row {
		column_heading := headermap[column_index]

		if column_heading != "" {
			val, err := parseONSFloat(cell_value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse float value in col %v: %v", column_index, err)
			}
			category_expenditure[column_heading] = val

			if val != 0 {
				non_zero = true
			}
		}
	}

	if !non_zero {
		return nil, fmt.Errorf("row has no values")
	}

	return category_expenditure, nil
}

// Wrapper for strconv.ParseFloat that strips spreadsheet values of estimate/unreliable/non-applicable symbols.
// Throws an error on empty values since ONS entries don't use truly empty values. This helps rule out certain edge cases.
func parseONSFloat(value string) (float64, error) {
	if value == "" {
		return 0, fmt.Errorf("empty value")
	}

	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	value = strings.TrimSuffix(value, "~")

	if value == ".." || value == ":" {
		return 0, nil
	}

	return strconv.ParseFloat(value, 64)
}
