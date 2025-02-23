package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

func main() {
	// Ask for the .env file path
	var envFilePath string
	huh.NewInput().
		Title("📂 `.env` File").
		Description("Enter the path to your `.env` file (Press Enter to use current directory)").
		Value(&envFilePath).
		Run()
	if envFilePath == "" {
		envFilePath = "./.env"
	}

	// Read .env file
	variables, descriptions, err := parseEnvFromFile(envFilePath)
	if err != nil {
		fmt.Println("❌ Error reading .env file:", err)
		os.Exit(1)
	}

	// Ask for the .tfvars file path
	var tfvarsPath string
	huh.NewInput().
		Title("💾 `.tfvars` File").
		Description("Enter the path to save `.tfvars` (Press Enter to use current directory)").
		Value(&tfvarsPath).
		Run()
	if tfvarsPath == "" {
		tfvarsPath = "./terraform.tfvars"
	}

	// Generate `.tfvars` file
	err = generateTfvarsFile(tfvarsPath, variables)
	if err != nil {
		fmt.Println("❌ Error creating `.tfvars` file:", err)
		os.Exit(1)
	}

	// Ask if the user wants to create `variables.tf`
	var createVarFile bool
	huh.NewConfirm().
		Title("📌 Generate `variables.tf`").
		Description("Do you want to generate a `variables.tf` file?").
		Affirmative("Yes").
		Negative("No").
		Value(&createVarFile).
		Run()

	// Generate `variables.tf` if requested
	if createVarFile {
		var variablesTfPath string
		huh.NewInput().
			Title("📌 `variables.tf` File").
			Description("Enter the path to save `variables.tf` (Press Enter to use current directory)").
			Value(&variablesTfPath).
			Run()
		if variablesTfPath == "" {
			variablesTfPath = "./variables.tf"
		}

		err := generateVariablesTfFile(variablesTfPath, variables, descriptions)
		if err != nil {
			fmt.Println("❌ Error creating `variables.tf` file:", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n✅ Process completed successfully.")
}

// Reads and parses a `.env` file with support for comments
func parseEnvFromFile(filePath string) (map[string]string, map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	return parseEnv(file)
}

// Reads environment variables from a file with support for inline comments
func parseEnv(r io.Reader) (map[string]string, map[string]string, error) {
	vars := make(map[string]string)
	descriptions := make(map[string]string) // Store comments as descriptions

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			valueParts := strings.SplitN(parts[1], "#", 2)
			value := strings.TrimSpace(valueParts[0])
			vars[key] = value

			// Store the comment (if present)
			if len(valueParts) > 1 {
				descriptions[key] = strings.TrimSpace(valueParts[1])
			}
		}
	}
	return vars, descriptions, scanner.Err()
}

// Generates the `.tfvars` file
func generateTfvarsFile(filePath string, variables map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range variables {
		line := fmt.Sprintf("%s = %s\n", key, value)
		writer.WriteString(line)
	}
	writer.Flush()

	return nil
}

// Generates the `variables.tf` file with descriptions from `.env`
func generateVariablesTfFile(filePath string, variables map[string]string, descriptions map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key := range variables {
		desc := descriptions[key] // Get the description from `.env`
		if desc == "" {
			desc = "No description available"
		}

		line := fmt.Sprintf("variable \"%s\" {\n  description = \"%s\"\n  type        = string\n}\n\n", key, desc)
		writer.WriteString(line)
	}
	writer.Flush()

	return nil
}
