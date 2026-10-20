package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type outgoings struct {
	Outgoings_communications float64 `json:"outgoings_communications"`
}

func TestOutgoingCalculation(t *testing.T) {
	app, err := NewApp("../../testcases", "../../ONSByRegion2019.xlsx")
	if err != nil {
		t.Errorf("failed to initialise app: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/cases/56763bf7-d046-4a54-a4a1-0728c724a092/outgoings", nil)
	rec := httptest.NewRecorder()

	app.Mux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("http request to known testcase failed")
	}

	json_response, err := io.ReadAll(res.Body)

	tmp := &outgoings{}
	if err := json.Unmarshal(json_response, tmp); err != nil {
		t.Errorf("json unmarshal failed")
	}

	if tmp.Outgoings_communications != 18.2 {
		t.Errorf("wrong outgoings_communications value")
	}
}
