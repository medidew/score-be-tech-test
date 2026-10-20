# ClearScore Backend Tech Test

## Overview

This repository is a solution to ClearScore's [public backend tech test](https://github.com/ClearScore/score-be-tech-test). It imports and parses an ONS Expenditure by Region dataset (`ONSByRegion2019.xlsx`), a directory of json mortgage cases (the `testcases` directory by default), and exposes the average outgoings for the clients listed on a case on a web server endpoint at `/cases/{caseID}/outgoings`. The web server in question can be run using `go run main.go` from the project root, or by building the server binary (`go build`) and running the resulting executable.

By default, the server listens at `localhost`/`127.0.0.1` on port `8082`. These variables are configurable constants in `main.go`, along with the ONS dataset filepath and the json mortgages directory path.

## Project Structure

- `main.go`: starts the web server.
- `internal/`
  - `internal/api`: contains the App struct which calls all file loaders and parsers, creates the postcode -> region lookup table, and registers the `outgoings` HTTP handler. Also contains the `outgoings` estimation logic.
  - `internal/cases`: package responsible for loading and parsing the raw json case files.
  - `internal/ons`: loading and parsing the ONS dataset.
  - `internal/region_mapping`: loads lookup tables from `/internal/region_mapping/sectors.tsv|units.tsv` and handles lookup logic.
- `testcases/`: holds all test json mortgage cases.
- `ONSByRegion2019.xlsx`: ONS dataset.
- `postcode-to-region.py`: python script for creating `internal/region_mapping`'s lookup tables.

The project's only non-standard Go package dependency is `github.com/xuri/excelize/v2` for parsing the ONS dataset. I deliberately chose against using an external API for postcode -> region lookup (such as [postcodes.io](https://postcodes.io/docs/api/lookup-postcode/)) so that this submission's evaluation would not be subject to the uptime of a third party service. This means that in order to perform postcode lookups to determine the region which a case falls under, pre-processed lookup tables are used (`/internal/region_mapping/sectors.tsv|units.tsv`).

These lookup tables are generated using the `postcode-to-region.py` script on the [ONS Postcode Directory from May 2026](https://geoportal.statistics.gov.uk/datasets/6fff67d204fd4f339591ed667a6e3642/about). It parses each row of `ONSPD_MAY_2026_UK.csv` and creates a mapping of postcode sectors to region codes (see Page 35, Field Name RGN25CD of the [ONSPD User Guide](https://geoportal.statistics.gov.uk/datasets/5f31ff4fc7624f9dabcc8148da88ba9e/about)) and outputs this data to `sectors.tsv`. In the event that a postcode sector covers multiple regions, it instead maps every affected postcode onto the region code in `units.tsv`. In the interest of not bloating this repository to over a gigabyte in size, the csv in question has been excluded (and instead linked above).

These files create a two-stage lookup process, whereby we first examine the entries in `units.tsv` to see if the exact postcode we're looking up is present, and if not we instead use the postcode sector in `sectors.tsv` to get the region code. Using a two-tier lookup instead of dumping every postcode into `units.tsv` saves memory by only having ~110,000 entries loaded at once instead of over a million. It also guarantees accuracy in converting postcodes to region codes.

## Design Decisions and Assumptions

The original specification did not specify what format multiple cases would take, whether it be multiple case files or multiple entries in a single json file. This solution assumes the former.

The original specification also did not specify expected outgoing estimation behaviour if two clients on the same case hailed from different regions. This solution lists all regions relevant to a case and averages the expenditure values across regions.

In the event no postcode can be resolved from the case file, the server returns an InternalServerError (500). If even a single postcode can be resolved, clients without resolvable postcodes are excluded from the outgoing estimations.

The ONS dataset parser handles statistical/eligibility symbols ([, ], ~, :, ..) and treats all values as fixed, since this is API is designed to calculate an estimation in the first place.

The mapping table controlling which ONS dataset categories make up the outgoing value estimations is explicitly listed in `internal/api.go`. While working on this solution, I decided that the exact makeup of which categories to use in these estimations was a subjective decision that should and would be reviewed and changed in a true production environment. This is why the design makes it easy to change and adjust which categories are weighted where.

## Known Limitations

A solution to the original spec's extension is not included in this submission.

Certain ONS expenditure categories contribute to more than one outgoing variable. Technically, this will unfairly penalize any household that spends more than other households in these specific categories. Given that there are no clear or obvious divides in some of these areas, and that the mappings can be easily changed, I have left these overlaps in. See `/internal/api.go` comments for more details.

Other expenditure categories (such as `9.4.3`: "TV, video, satellite rental, cable subscriptions and TV licenses") have the issue of being too broad. `9.4.3` is the only category which suits `outgoings_television_license`, but obviously includes these other expenditures as well. It is worth noting that these estimations will therefore be inaccurate.

No structured logging.

## Tests

`go test ./...`

## Disclaimer

Qwen3.8-27B (a Large Langauge Model) was used to generate the `postcode-to-region.py` utility script. Everything else in this repository was typed by hand.