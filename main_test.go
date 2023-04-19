package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

var (
	//go:embed testdata/opa-coverage.json
	opaCoverage []byte
	//go:embed testdata/codecov-coverage.json
	codecovCoverage []byte
)

func TestProcess(t *testing.T) {
	var (
		expected Out
		actual   = Out{Coverage: map[string]map[string]int{}}
	)
	if err := json.Unmarshal(codecovCoverage, &expected); err != nil {
		t.Fatal(err)
	}
	if err := Process(opaCoverage, &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %#v, but got %#v", expected, actual)
	}
}

func TestProcessFile(t *testing.T) {
	var (
		expected Out
		actual   = Out{Coverage: map[string]map[string]int{}}
	)
	if err := json.Unmarshal(codecovCoverage, &expected); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("testdata/opa-coverage.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := ProcessFile(f, &actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %#v, but got %#v", expected, actual)
	}
}
