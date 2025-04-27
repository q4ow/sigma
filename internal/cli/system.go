package cli

import (
	"fmt"
	"os"
	"os/exec"
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

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk usage information",
	Long:  `Display disk usage information for all mounted filesystems`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := exec.Command("df", "-h").Output()
		if err != nil {
			return fmt.Errorf("failed to get disk usage: %w", err)
		}

		fmt.Println(formatDiskOutput(string(output)))
		return nil
	},
}

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Show network interfaces",
	Long:  `Display information about network interfaces and their status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := exec.Command("ip", "addr").Output()
		if err != nil {
			return fmt.Errorf("failed to get network info: %w", err)
		}
		fmt.Printf("%s", output)
		return nil
	},
}

func init() {
	systemCmd.AddCommand(fetchCmd)
	systemCmd.AddCommand(diskCmd)
	systemCmd.AddCommand(netCmd)
	rootCmd.AddCommand(systemCmd)
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

func formatDiskOutput(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return output
	}

	var result strings.Builder
	titleColor := "\033[1;36m"
	valueColor := "\033[1;37m"
	headerColor := "\033[1;33m"
	resetColor := "\033[0m"
	separator := strings.Repeat("─", 60)

	fmt.Fprintf(&result, "\n%s%s%s\n", titleColor, separator, resetColor)
	fmt.Fprintf(&result, "%sDisk Usage Information%s\n", titleColor, resetColor)
	fmt.Fprintf(&result, "%s%s%s\n\n", titleColor, separator, resetColor)

	fmt.Fprintf(&result, "%s%-15s %-10s %-10s %-8s%s\n",
		headerColor, "Filesystem", "Size", "Used", "Use%", resetColor)
	fmt.Fprintf(&result, "%s%s%s\n", titleColor, strings.Repeat("─", 45), resetColor)

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		if strings.HasPrefix(fields[0], "tmpfs") ||
			strings.HasPrefix(fields[0], "dev") ||
			strings.HasPrefix(fields[0], "run") ||
			strings.HasPrefix(fields[0], "efivarfs") {
			continue
		}

		fsName := fields[0]
		if strings.Contains(fsName, "/") {
			parts := strings.Split(fsName, "/")
			fsName = parts[len(parts)-1]
		}

		fmt.Fprintf(&result, "%s%-15s %-10s %-10s %-8s%s\n",
			valueColor, fsName, fields[1], fields[2], fields[4], resetColor)
	}

	result.WriteString("\n")
	return result.String()
}
