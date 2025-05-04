package main

import (
	"fmt"
)

func main() {
	monitors := getMonitors()

	fmt.Printf("Total monitors: %d\n\n", len(monitors))

	for _, monitor := range monitors {
		primaryStatus := ""
		if monitor.Primary {
			primaryStatus = " (Primary)"
		}
		fmt.Printf("Monitor %d%s:\n", monitor.Index+1, primaryStatus)
		fmt.Printf("  Resolution: %s\n", monitor.Resolution)
		fmt.Printf("  Position: %s\n", monitor.Position)
	}
}
