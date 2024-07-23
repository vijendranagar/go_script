package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	rmAndMkdir("autogen_output")
	var path = "./code-solidity-vault-guardians/"
	var environment = detectEnvironment(path)
	fmt.Println("Detected environment : ", environment)
	installDependencies(environment, path)
	runTool()
	extractSolidityDetails()
}

func rmAndMkdir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Printf("Error removing %s: %v\n", dir, err)
		return
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", dir, err)
	}
}

func detectEnvironment(path string) string {
	var DIRECTORY = path
	var environment string

	// Check for existence of hardhat config files
	_, err := os.Stat(filepath.Join(DIRECTORY, "hardhat.config.js"))
	_, err1 := os.Stat(filepath.Join(DIRECTORY, "hardhat.config.ts"))
	if err == nil || err1 == nil {
		environment = "hardhat"
	} else {
		// Only check for foundry.toml if hardhat config JS doesn't exist
		_, err = os.Stat(filepath.Join(DIRECTORY, "foundry.toml"))
		if err == nil {
			environment = "foundry"
		} else {
			// Check for existence of brownie-config.yaml
			_, err = os.Stat(filepath.Join(DIRECTORY, "brownie-config.yaml"))
			if err == nil {
				environment = "brownie"
			} else {
				environment = "unknown"
			}
		}
	}

	fmt.Println(environment)
	return environment
}

func installDependencies(env string, directory string) {
	switch env {
	case "hardhat":
		fmt.Println("hardhat environment")
		if _, err := os.Stat(filepath.Join(directory, "package.json")); err == nil {
			cmd := exec.Command("npm", "install")
			cmd.Dir = directory
			err := cmd.Run()
			if err != nil {
				fmt.Printf("npm install failed: %v\n", err)
			}
		} else if _, err := os.Stat(filepath.Join(directory, "yarn.lock")); err == nil {
			cmd := exec.Command("yarn", "install")
			cmd.Dir = directory
			err := cmd.Run()
			if err != nil {
				fmt.Printf("yarn install failed: %v\n", err)
			}
		}
	case "foundry":
		fmt.Println("foundry environment")
		if _, err := os.Stat(filepath.Join(directory, "foundry.toml")); err == nil {
			cmd := exec.Command("forge", "install")
			cmd.Dir = directory
			err := cmd.Run()
			if err != nil {
				fmt.Printf("forge install failed: %v\n", err)
			}
		}
	case "brownie":
		fmt.Println("brownie environment")
		if _, err := os.Stat(filepath.Join(directory, "requirements.txt")); err == nil {
			cmd := exec.Command("pip", "install", "-r", filepath.Join(directory, "requirements.txt"))
			cmd.Dir = directory
			err := cmd.Run()
			if err != nil {
				fmt.Printf("pip install failed: %v\n", err)
			}
		}
	default:
		fmt.Println("Unknown environment:", env)
	}
}

func runTool() {
	startTime := getCurrentTimeMillis()
	// Define the command to run
	cmd := exec.Command("falcon", ".", "--checklist", "--json", "autogen_output/output.json")

	fmt.Printf("🚀 ~ funcrunTool ~ cmd:", cmd)
	// Redirect stderr to stdout
	cmd.Stderr = cmd.Stdout
	endTime := getCurrentTimeMillis()
	executionTime := endTime - startTime
	fmt.Printf("Execution time: %d milliseconds\n", executionTime)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Output: %v", err)
	}

	// Write the output to a file
	err = logToFile(string(output))
	if err != nil {
		log.Fatalf("Writing to file: %v", err)
	}

}

// logToFile writes the given text to a file, appending to it if it exists.
func logToFile(text string) error {
	file, err := os.OpenFile("autogen_output/results.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text + "\n")
	return err
}

// func runTool() {
// 	// Execute the tool command here

// }

func extractSolidityDetails() {
	// Implement the logic to extract details about Solidity files here
	// This would involve reading the solidity_files.txt file, processing its contents,
	// and writing the results to the results.md file
}

func getCurrentTimeMillis() int64 {
	return int64(time.Now().UnixNano()) / int64(time.Millisecond)
}
