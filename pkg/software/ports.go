package software

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

// CheckPortAvailable verifies whether a TCP port is free to bind on the host system.
func CheckPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// PromptAndCheckPorts displays default port(s) for a software, asks the tech user
// yes/no to confirm or customize, and checks live port availability in the target system.
func PromptAndCheckPorts(s *Software, reader *bufio.Reader, autoYes bool) (map[string]int, error) {
	chosenPorts := make(map[string]int)

	if len(s.PortConfigs) == 0 {
		return chosenPorts, nil
	}

	fmt.Println()
	ui.PrintBanner(fmt.Sprintf("Target System Port Configuration & Check: %s", s.Name))

	for _, pc := range s.PortConfigs {
		defaultPort := pc.DefaultPort
		fmt.Printf("\n[Port Configuration] %s%s%s (Container Port: %d)\n", ui.ColorCyan, pc.Name, ui.ColorReset, pc.Container)
		fmt.Printf("  • Default Port Number: %s%d%s\n", ui.ColorBold, defaultPort, ui.ColorReset)

		targetPort := defaultPort

		if !autoYes && reader != nil {
			fmt.Printf("  • Do you want to use the default port (%d)? [Y/n]: ", defaultPort)
			input, err := reader.ReadString('\n')
			if err != nil {
				input = "y"
			}
			input = strings.TrimSpace(input)

			if strings.ToLower(input) == "n" || strings.ToLower(input) == "no" {
				for {
					fmt.Printf("  Enter custom port number for %s: ", pc.Name)
					pStr, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					pStr = strings.TrimSpace(pStr)
					pVal, err := strconv.Atoi(pStr)
					if err != nil || pVal < 1 || pVal > 65535 {
						ui.Warning("Invalid port '%s'. Port must be a number between 1 and 65535.", pStr)
						continue
					}
					targetPort = pVal
					break
				}
			}
		}

		// Check availability in the target system
		for {
			ui.Info("Checking availability of port %d in target system...", targetPort)
			if CheckPortAvailable(targetPort) {
				ui.Success("Port %d is AVAILABLE on target system!", targetPort)
				break
			}

			ui.Warning("Port %d is currently IN USE / OCCUPIED on target system!", targetPort)
			if autoYes || reader == nil {
				ui.Warning("Proceeding with occupied port %d at user request.", targetPort)
				break
			}

			fmt.Printf("  Would you like to enter a different port? [Y/n]: ")
			reAns, _ := reader.ReadString('\n')
			reAns = strings.TrimSpace(reAns)
			if strings.ToLower(reAns) == "n" || strings.ToLower(reAns) == "no" {
				ui.Warning("Proceeding with occupied port %d as requested (container may fail to start).", targetPort)
				break
			}

			fmt.Printf("  Enter alternative port number: ")
			pStr, _ := reader.ReadString('\n')
			pStr = strings.TrimSpace(pStr)
			pVal, err := strconv.Atoi(pStr)
			if err != nil || pVal < 1 || pVal > 65535 {
				ui.Warning("Invalid port '%s'. Must be between 1 and 65535.", pStr)
				continue
			}
			targetPort = pVal
		}

		chosenPorts[pc.Name] = targetPort
	}

	return chosenPorts, nil
}
