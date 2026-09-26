package software

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

// PromptSelectSoftware prompts the user interactively to pick a software from the registry.
func PromptSelectSoftware() (*Software, error) {
	fmt.Println()
	ui.PrintBanner("PiltiSmart Interactive Software Selector")
	fmt.Println("Please choose a software from the list:")
	fmt.Println(strings.Repeat("-", 60))

	for i, s := range Registry {
		status := CheckStatus(&s)
		statusColor := ui.ColorGreen
		if strings.Contains(status, "STOPPED") {
			statusColor = ui.ColorYellow
		} else if strings.Contains(status, "NOT") {
			statusColor = ui.ColorRed
		}
		fmt.Printf(" [%d] %-14s : %-30s [%s%s%s]\n", i+1, s.ID, s.Name, statusColor, status, ui.ColorReset)
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Select an option [1-%d] (or 'q' to quit): ", len(Registry))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "q" || input == "" {
		return nil, fmt.Errorf("operation cancelled by user")
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(Registry) {
		return nil, fmt.Errorf("invalid selection '%s'", input)
	}

	selected := Registry[choice-1]
	return &selected, nil
}

// PromptSelectAction prompts the user for an action on the selected software.
func PromptSelectAction(s *Software) (string, error) {
	fmt.Printf("\nSelected Software: %s%s (%s)%s\n", ui.ColorCyan, s.Name, s.ID, ui.ColorReset)
	fmt.Println("Available Actions:")
	fmt.Println(" [1] Install / Deploy (docker compose up -d)")
	fmt.Println(" [2] Check Status")
	fmt.Println(" [3] View Container Logs")
	fmt.Println(" [4] Restart Container")
	fmt.Println(" [5] Stop Container")
	fmt.Println(" [6] Remove Container")
	fmt.Print("Choose action [1-6] (or 'q' to cancel): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		return "install", nil
	case "2":
		return "status", nil
	case "3":
		return "logs", nil
	case "4":
		return "restart", nil
	case "5":
		return "stop", nil
	case "6":
		return "remove", nil
	case "q", "Q", "":
		return "", fmt.Errorf("action cancelled")
	default:
		return "", fmt.Errorf("invalid action choice: %s", input)
	}
}
