package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Because we are getting random values early tests were flaky.
// Repeating test many times allowed me to see intermittent
// errors
var trials = 1000

func TestHexColorString(t *testing.T) {

	for i := 0; i < trials; i++ {

		actual := hexColorString()

		if err := validateColor(actual); err != nil {
			t.Errorf("%v", err)
		}

	}

}

func validateColor(s string) error {
	if string(s[0]) != "#" {
		return fmt.Errorf("Hex color show begin with '#' got %s", string(s[0]))
	}

	if len(s) != 7 {
		return fmt.Errorf("Hex color should be 7 characters, got %d %s", len(s), s)
	}

	num, err := strconv.ParseInt(s[1:6], 16, 32)
	if err != nil {
		return fmt.Errorf("Error converting result to number: %v", err)
	}

	if num < 0 || num > 16777215 {
		return fmt.Errorf("Hex should be > 0 (000000) and  16777215 (ffffff) got: %d", num)
	}
	return nil
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `ok`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestColorHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/api/color", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(color)
	handler.ServeHTTP(rr, req)

	if statusGot := rr.Code; statusGot != http.StatusOK {
		t.Errorf("wrong status code %s: got %d want %d", "/api/color", statusGot, http.StatusOK)
	}

	colorGot := rr.Body.String()
	if err := validateColor(colorGot); err != nil {
		t.Errorf("Invalid color got: %s, err: %v", colorGot, err)
	}

}

type result struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

func TestColorCompleteHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/color-complete", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(colorComplete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var got result
	actual := rr.Body.Bytes()
	if err := json.Unmarshal(actual, &got); err != nil {
		t.Errorf("could not parse server response: %v", err)
	}

	if err := validateColor(got.Color); err != nil {
		t.Errorf("Invalid color got: %s, err: %v", got.Color, err)
	}
}
