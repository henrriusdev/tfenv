package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Model for the CLI
type model struct {
	envFilePath   string
	tfvarsPath    string
	variables     map[string]string
	createVarFile bool
	step          int
}

// Entry point
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Initializes the model with default values
func initialModel() model {
	return model{
		step:      0,
		variables: make(map[string]string),
	}
}

func (m model) Init() tea.Cmd {
	return m.checkForEnvFile()
}

// Update function for managing application state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case 0:
		m.envFilePath = msg.(string)
		m.step++
		return m, m.readEnvFile()
	case 1:
		m.tfvarsPath = msg.(string)
		m.step++
		return m, m.generateTfvarsFile()
	case 2:
		m.createVarFile = msg.(bool)
		m.step++
		if m.createVarFile {
			return m, m.generateVariablesTfFile()
		}
		fmt.Println("\n✅ Process completed successfully.")
		return m, tea.Quit
	case 3:
		fmt.Println("\n✅ `variables.tf` file created successfully.")
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	switch m.step {
	case 0:
		return "\n📂 Looking for `.env` file in the current directory..."
	case 1:
		return "\n💾 Enter the path to save `.tfvars` file: "
	case 2:
		return "\n📌 Do you want to generate a `variables.tf` file? (yes/no): "
	default:
		return ""
	}
}

// Check if the .env file exists in the current directory
func (m model) checkForEnvFile() tea.Cmd {
	return func() tea.Msg {
		if _, err := os.Stat(".env"); err == nil {
			fmt.Println("\n✅ Found `.env` file in the current directory.")
			return ".env"
		}

		var envPath string
		huh.NewInput().
			Title("📂 `.env` file not found").
			Description("Enter the path to your `.env` file").
			Value(&envPath).
			Run()
		return envPath
	}
}

// Reads environment variables from the specified .env file
func parseEnv(r io.Reader) (map[string]string, error) {
	vars := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return vars, nil
}

// Reads environment variables from a file
func parseEnvFromFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseEnv(file) // Now works correctly
}

// Reads environment variables from the file and stores them in the model
func (m model) readEnvFile() tea.Cmd {
	return func() tea.Msg {
		envVars, err := parseEnvFromFile(m.envFilePath)
		if err != nil {
			fmt.Println("❌ Error reading .env file:", err)
			os.Exit(1)
		}
		m.variables = envVars
		return m.askForTfvarsPath()
	}
}

// Asks the user for the `.tfvars` file path
func (m model) askForTfvarsPath(mockValue ...string) tea.Cmd {
	return func() tea.Msg {
		// If a mock value is provided, return it instead of asking the user
		if len(mockValue) > 0 {
			return mockValue[0]
		}

		var tfvarsPath string
		huh.NewInput().
			Title("💾 `.tfvars` File").
			Description("Enter the path to save `.tfvars`").
			Value(&tfvarsPath).
			Run()
		return tfvarsPath
	}
}

// Generates the `.tfvars` file
func (m model) generateTfvarsFile() tea.Cmd {
	return func() tea.Msg {
		file, err := os.Create(m.tfvarsPath)
		if err != nil {
			fmt.Println("❌ Error creating `.tfvars` file:", err)
			os.Exit(1)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for key, value := range m.variables {
			line := fmt.Sprintf("%s = \"%s\"\n", key, value)
			writer.WriteString(line)
		}
		writer.Flush()

		return m.askForVariablesTf()
	}
}

// Asks if the user wants to create `variables.tf`
func (m model) askForVariablesTf(mockValue ...bool) tea.Cmd {
	return func() tea.Msg {
		// If a mock value is provided, return it instead of asking the user
		if len(mockValue) > 0 {
			return mockValue[0]
		}

		var shouldCreate bool
		huh.NewConfirm().
			Title("📌 Generate `variables.tf`").
			Description("Do you want to generate a `variables.tf` file?").
			Affirmative("Yes").
			Negative("No").
			Value(&shouldCreate).
			Run()
		return shouldCreate
	}
}

// Generates the `variables.tf` file
func (m model) generateVariablesTfFile() tea.Cmd {
	return func() tea.Msg {
		file, err := os.Create("variables.tf")
		if err != nil {
			fmt.Println("❌ Error creating `variables.tf`:", err)
			os.Exit(1)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for key := range m.variables {
			line := fmt.Sprintf("variable \"%s\" {\n  description = \"\"\n  type        = string\n}\n\n", key)
			writer.WriteString(line)
		}
		writer.Flush()

		return ""
	}
}
