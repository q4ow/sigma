package cli

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/q4ow/sigma/internal/linux"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information utilities",
	Long:  `A collection of useful system information commands that work across different distributions.`,
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Display system information in a pretty format",
	Long:  `Shows detailed system information including distribution details, memory usage, and system specs in a nice format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentUser, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		hostname, err := linux.GetHostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %w", err)
		}

		distro, err := linux.GetDistroInfo()
		if err != nil {
			return fmt.Errorf("failed to get distribution info: %w", err)
		}

		kernel, err := linux.GetKernelVersion()
		if err != nil {
			return fmt.Errorf("failed to get kernel version: %w", err)
		}

		pkgCount, err := linux.GetPackageCount(distro.PackageManager)
		if err != nil {
			pkgCount = 0
		}

		mem, err := linux.GetMemoryInfo()
		if err != nil {
			return fmt.Errorf("failed to get memory info: %w", err)
		}

		uptime := getUptime()

		separator := "─"
		titleColor := "\033[1;36m"
		valueColor := "\033[1;37m"
		resetColor := "\033[0m"

		totalMemGB := float64(mem.Total) / (1024 * 1024)
		usedMemGB := float64(mem.Used) / (1024 * 1024)

		fmt.Printf("\n%s%s@%s%s\n", titleColor, currentUser.Username, hostname, resetColor)
		fmt.Printf("%s%s%s\n", titleColor, strings.Repeat(separator, 40), resetColor)

		fmt.Printf("%sOS%s        │ %s%s %s%s\n", titleColor, resetColor, valueColor, distro.Name, distro.Version, resetColor)
		fmt.Printf("%sKernel%s    │ %s%s%s\n", titleColor, resetColor, valueColor, kernel, resetColor)
		fmt.Printf("%sUptime%s    │ %s%s%s\n", titleColor, resetColor, valueColor, uptime, resetColor)
		fmt.Printf("%sShell%s     │ %s%s%s\n", titleColor, resetColor, valueColor, getShell(), resetColor)

		fmt.Printf("%sCPU%s       │ %s%d cores%s\n", titleColor, resetColor, valueColor, runtime.NumCPU(), resetColor)
		fmt.Printf("%sMemory%s    │ %s%.1f GB / %.1f GB%s\n", titleColor, resetColor, valueColor, usedMemGB, totalMemGB, resetColor)

		fmt.Printf("%sPackages%s  │ %s%d (%s)%s\n", titleColor, resetColor, valueColor, pkgCount, distro.PackageManager, resetColor)

		fmt.Printf("%s%s%s\n", titleColor, strings.Repeat(separator, 40), resetColor)

		return nil
	},
}

func getUptime() string {
	uptimeFile, err := linux.ReadUptimeFile()
	if err != nil {
		return "Unknown"
	}

	duration := time.Duration(uptimeFile.UptimeSeconds) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}

	parts := strings.Split(shell, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}
