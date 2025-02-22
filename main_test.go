package main

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
)

// ✅ Test `parseEnv`
func TestParseEnv(t *testing.T) {
	envContent := `
# Comment line
PORT=8080
STATIC=/app/assets
SESSION_SECRET=supersecret
DB_PATH=/app/data.db
`

	expected := map[string]string{
		"PORT":           "8080",
		"STATIC":         "/app/assets",
		"SESSION_SECRET": "supersecret",
		"DB_PATH":        "/app/data.db",
	}

	reader := strings.NewReader(envContent)
	result, err := parseEnv(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("For key %q, expected %q, got %q", key, expectedValue, result[key])
		}
	}
}

// ✅ Test `parseEnvFromFile`
func TestParseEnvFromFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	envContent := `PORT=8080
STATIC=/app/assets
SESSION_SECRET=supersecret
DB_PATH=/app/data.db`
	tempFile.WriteString(envContent)
	tempFile.Close()

	result, err := parseEnvFromFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := map[string]string{
		"PORT":           "8080",
		"STATIC":         "/app/assets",
		"SESSION_SECRET": "supersecret",
		"DB_PATH":        "/app/data.db",
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("For key %q, expected %q, got %q", key, expectedValue, result[key])
		}
	}
}

// ✅ Test `readEnvFile`
func TestReadEnvFile(t *testing.T) {
	// Create a temporary .env file for testing
	tempFile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Cleanup

	// Write mock environment variables
	envContent := "PORT=8080\n"
	if _, err := tempFile.WriteString(envContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Initialize model as `mdl`
	mdl := initialModel()
	mdl.envFilePath = tempFile.Name()

	// Provide a mock tfvars file path to prevent user input
	mockTfvarsPath := "mock_output.tfvars"

	// Execute readEnvFile() and capture its returned tea.Cmd
	cmd := mdl.readEnvFile(mockTfvarsPath)
	if cmd == nil {
		t.Fatalf("Expected a tea.Cmd but got nil")
	}

	// Execute the command function to obtain the tea.Msg
	msg := cmd()

	// Apply the message to update the model state
	newMdl, _ := mdl.Update(msg)
	updatedMdl := newMdl.(model) // Type assertion

	// Verify that the expected environment variable exists
	expectedValue := "8080"
	gotValue, exists := updatedMdl.variables["PORT"]
	if !exists {
		t.Fatalf("Expected key PORT, but it was not found in mdl.variables")
	}
	if gotValue != expectedValue {
		t.Errorf("Expected PORT=%q, got %q", expectedValue, gotValue)
	}
}

// ✅ Test `askForTfvarsPath`
func TestAskForTfvarsPath(t *testing.T) {
	mdl := initialModel()
	expectedPath := "mock_output.tfvars"

	// Simulate user input by passing the mock value
	cmd := mdl.askForTfvarsPath(expectedPath)
	msg := cmd()

	// Verify that the function returned the simulated input
	if msg != expectedPath {
		t.Errorf("Expected %q, got %q", expectedPath, msg)
	}
}

// ✅ Test `generateTfvarsFile`
func TestGenerateTfvarsFile(t *testing.T) {
	mdl := initialModel()
	mdl.variables = map[string]string{
		"PORT": "8080",
	}
	mdl.tfvarsPath = "test.tfvars"

	defer os.Remove("test.tfvars")

	cmd := mdl.generateTfvarsFile()
	cmd()

	content, err := os.ReadFile("test.tfvars")
	if err != nil {
		t.Fatalf("Error reading tfvars file: %v", err)
	}

	expected := `PORT = "8080"
`
	if string(content) != expected {
		t.Errorf("Expected tfvars content:\n%q\nGot:\n%q", expected, string(content))
	}
}

// ✅ Test `askForVariablesTf`
func TestAskForVariablesTf(t *testing.T) {
	mdl := initialModel()
	expectedResponse := true // Simulating user selecting "Yes"

	// Simulate user input by passing the mock value
	cmd := mdl.askForVariablesTf(expectedResponse)
	msg := cmd()

	// Verify that the function returned the simulated input
	if msg != expectedResponse {
		t.Errorf("Expected %v, got %v", expectedResponse, msg)
	}
}

// ✅ Test `generateVariablesTfFile`
func TestGenerateVariablesTfFile(t *testing.T) {
	mdl := initialModel()
	mdl.variables = map[string]string{
		"PORT": "8080",
	}
	defer os.Remove("variables.tf")

	cmd := mdl.generateVariablesTfFile()
	cmd()

	content, err := os.ReadFile("variables.tf")
	if err != nil {
		t.Fatalf("Error reading variables.tf file: %v", err)
	}

	expected := `variable "PORT" {
  description = ""
  type        = string
}

`
	if string(content) != expected {
		t.Errorf("Expected variables.tf content:\n%q\nGot:\n%q", expected, string(content))
	}
}

// ✅ Test CLI flow to increase coverage
func TestCLIFlow(t *testing.T) {
	tests := []struct {
		name          string
		inputs        []tea.Msg
		expectedSteps int
	}{
		{
			name: "Auto-detect .env, create .tfvars, confirm variables.tf",
			inputs: []tea.Msg{
				".env",          // Auto-detected `.env`
				"output.tfvars", // User input for `.tfvars`
				true,            // User confirms `variables.tf`
			},
			expectedSteps: 3, // Should reach the final step
		},
		{
			name: "Manual input for .env file",
			inputs: []tea.Msg{
				"config/.env",   // User manually enters `.env`
				"output.tfvars", // User input for `.tfvars`
				false,           // User declines `variables.tf`
			},
			expectedSteps: 3, // Should reach the final step
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mdl := initialModel()

			for _, input := range test.inputs {
				newMdl, _ := mdl.Update(input)
				mdl = newMdl.(model)
			}

			if mdl.step != test.expectedSteps {
				t.Errorf("Expected step %d, got %d", test.expectedSteps, mdl.step)
			}
		})
	}
}
