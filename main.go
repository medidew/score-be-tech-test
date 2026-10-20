package main

import (
	"fmt"
	"net/http"

	"medidew.name/score-be-tech-test/internal/api"
)

// Would put all of this in a config/env for an actual production environment.
const address string = "127.0.0.1"
const port int = 8082
const case_dir string = "testcases"
const ons_path string = "ONSByRegion2019.xlsx"

func main() {
	app, err := api.NewApp(case_dir, ons_path)
	if err != nil {
		panic(err)
	}

	listening_address := fmt.Sprintf("%v:%v", address, port)

	fmt.Printf("listening on %v", listening_address)
	http.ListenAndServe(listening_address, app.Mux)
}
