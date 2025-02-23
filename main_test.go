package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseEnv(t *testing.T) {
	input := "VAR1=value1\nVAR2=value2\nVAR3=value3"
	r := strings.NewReader(input)
	vars, err := parseEnv(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if vars["VAR1"] != "value1" || vars["VAR2"] != "value2" || vars["VAR3"] != "value3" {
		t.Errorf("Parsed values do not match expected output: %v", vars)
	}
}

func TestParseEnvFromFile(t *testing.T) {
	fileContent := "TEST_VAR=hello\nSECOND_VAR=world"
	filePath := "test.env"

	err := os.WriteFile(filePath, []byte(fileContent), 0o644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(filePath)

	vars, err := parseEnvFromFile(filePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if vars["TEST_VAR"] != "hello" || vars["SECOND_VAR"] != "world" {
		t.Errorf("Parsed values do not match expected output: %v", vars)
	}
}

func TestGenerateTfvarsFile(t *testing.T) {
	variables := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	}
	filePath := "test.tfvars"

	err := generateTfvarsFile(filePath, variables)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expected := "VAR1 = \"value1\"\nVAR2 = \"value2\"\n"
	if string(content) != expected {
		t.Errorf("File content mismatch. Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestGenerateVariablesTfFile(t *testing.T) {
	variables := map[string]string{
		"MY_VAR":      "",
		"ANOTHER_VAR": "",
	}
	filePath := "test.variables.tf"

	err := generateVariablesTfFile(filePath, variables)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expected := "variable \"MY_VAR\" {\n  description = \"\"\n  type        = string\n}\n\nvariable \"ANOTHER_VAR\" {\n  description = \"\"\n  type        = string\n}\n\n"
	if string(content) != expected {
		t.Errorf("File content mismatch. Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}
