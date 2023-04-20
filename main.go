package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var profileName, awsFileName string

	flag.StringVar(&profileName, "profile-name", "", "AWS Profile entries to be extracted")
	flag.StringVar(&awsFileName, "aws-file-name", "~/.aws/credentials", "Name of the AWS file")

	log.SetFlags(0)
	flag.Parse()

	if profileName == "" {
		log.Fatalf("profile-name not provided")
	}
	awsFileName, err := expandTilde(awsFileName)
	if err != nil {
		log.Fatalf("couldn't expand user's home directory")
	}

	fileData, err := os.ReadFile(awsFileName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	fileLines := strings.Split(string(fileData), "\n")
	lines := extractProfileLines(fileLines, profileName)
	lines = getExportableLines(lines)

	fmt.Println(strings.Join(lines, "\n"))
}

func getExportableLines(lines []string) []string {
	result := make([]string, 0, len(lines))

	for _, l := range lines {
		l = strings.ReplaceAll(l, " ", "")
		name, value, _ := strings.Cut(l, "=")

		exportLine := fmt.Sprintf("export %s=%s", strings.ToUpper(name), value)
		result = append(result, exportLine)
	}

	return result
}

// extractProfileLines returns the lines between the profile name and the next empty line
// the profile name is matched with string.Contains
func extractProfileLines(lines []string, profile string) []string {
	result := []string{}
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], profile) {
			for j := i + 1; j < len(lines); j++ {
				if lines[j] == "" {
					return result
				}

				result = append(result, lines[j])
			}
		}
	}

	return result
}

// Replace "~" with user's home directory
func expandTilde(filePath string) (string, error) {
	if !strings.HasPrefix(filePath, "~/") {
		return filePath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath = homeDir + filePath[1:]

	return filePath, nil
}
