package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseEnv(t *testing.T) {
	input := "VAR1=value1 # First variable\nVAR2=value2 # Second variable\nVAR3=value3"
	r := strings.NewReader(input)
	vars, descriptions, err := parseEnv(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if vars["VAR1"] != "value1" || vars["VAR2"] != "value2" || vars["VAR3"] != "value3" {
		t.Errorf("Parsed values do not match expected output: %v", vars)
	}

	if descriptions["VAR1"] != "First variable" || descriptions["VAR2"] != "Second variable" {
		t.Errorf("Parsed descriptions do not match expected output: %v", descriptions)
	}
}

func TestParseEnvFromFile(t *testing.T) {
	fileContent := "TEST_VAR=hello # A test variable\nSECOND_VAR=world"
	filePath := "test.env"

	err := os.WriteFile(filePath, []byte(fileContent), 0o644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(filePath)

	vars, descriptions, err := parseEnvFromFile(filePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if vars["TEST_VAR"] != "hello" || vars["SECOND_VAR"] != "world" {
		t.Errorf("Parsed values do not match expected output: %v", vars)
	}

	if descriptions["TEST_VAR"] != "A test variable" {
		t.Errorf("Parsed descriptions do not match expected output: %v", descriptions)
	}
}

func TestGenerateVariablesTfFile(t *testing.T) {
	variables := map[string]string{
		"MY_VAR":      "",
		"ANOTHER_VAR": "",
	}
	descriptions := map[string]string{
		"MY_VAR":      "A sample variable",
		"ANOTHER_VAR": "Another sample variable",
	}
	filePath := "test.variables.tf"

	err := generateVariablesTfFile(filePath, variables, descriptions)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expected := "variable \"MY_VAR\" {\n  description = \"A sample variable\"\n  type        = string\n}\n\nvariable \"ANOTHER_VAR\" {\n  description = \"Another sample variable\"\n  type        = string\n}\n\n"
	if string(content) != expected {
		t.Errorf("File content mismatch. Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}
