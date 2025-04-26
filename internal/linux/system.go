package linux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DistroInfo struct {
	Name           string
	Version        string
	Codename       string
	PackageManager string
}

func GetDistroInfo() (*DistroInfo, error) {
	info, err := readOSRelease()
	if err == nil {
		return info, nil
	}

	return detectWithLSBRelease()
}

func readOSRelease() (*DistroInfo, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &DistroInfo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		value := strings.Trim(parts[1], "\"")
		switch parts[0] {
		case "NAME":
			info.Name = value
		case "VERSION_ID":
			info.Version = value
		case "VERSION_CODENAME":
			info.Codename = value
		}
	}

	info.PackageManager = detectPackageManager()
	return info, nil
}

func detectWithLSBRelease() (*DistroInfo, error) {
	cmd := exec.Command("lsb_release", "-a")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	info := &DistroInfo{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		value := strings.TrimSpace(parts[1])
		switch strings.TrimSpace(parts[0]) {
		case "Distributor ID":
			info.Name = value
		case "Release":
			info.Version = value
		case "Codename":
			info.Codename = value
		}
	}

	info.PackageManager = detectPackageManager()
	return info, nil
}

func detectPackageManager() string {
	packageManagers := []string{"apt", "dnf", "yum", "pacman", "zypper"}
	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm); err == nil {
			return pm
		}
	}
	return "unknown"
}

type SystemMemory struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Cached    uint64
	SwapTotal uint64
	SwapUsed  uint64
}

func GetMemoryInfo() (*SystemMemory, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mem := &SystemMemory{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value := parseMemValue(fields[1])
		switch fields[0] {
		case "MemTotal:":
			mem.Total = value
		case "MemFree:":
			mem.Free = value
		case "Cached:":
			mem.Cached = value
		case "SwapTotal:":
			mem.SwapTotal = value
		case "SwapFree:":
			mem.SwapUsed = mem.SwapTotal - value
		}
	}

	mem.Used = mem.Total - mem.Free - mem.Cached
	return mem, nil
}

func parseMemValue(value string) uint64 {
	var num uint64
	fmt.Sscanf(value, "%d", &num)
	return num
}

type UptimeInfo struct {
	UptimeSeconds float64
	IdleSeconds   float64
}

func ReadUptimeFile() (*UptimeInfo, error) {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return nil, err
	}

	var uptime, idle float64
	_, err = fmt.Sscanf(string(content), "%f %f", &uptime, &idle)
	if err != nil {
		return nil, err
	}

	return &UptimeInfo{
		UptimeSeconds: uptime,
		IdleSeconds:   idle,
	}, nil
}

func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

func GetKernelVersion() (string, error) {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return "", err
	}

	parts := strings.Fields(string(data))
	if len(parts) >= 3 {
		return parts[2], nil
	}
	return "", fmt.Errorf("unable to parse kernel version")
}

func GetPackageCount(packageManager string) (int, error) {
	var cmd *exec.Cmd

	switch packageManager {
	case "pacman":
		cmd = exec.Command("pacman", "-Q")
	case "apt":
		cmd = exec.Command("dpkg", "--get-selections")
	case "dnf", "yum":
		cmd = exec.Command("rpm", "-qa")
	default:
		return 0, fmt.Errorf("unsupported package manager: %s", packageManager)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return len(strings.Split(string(output), "\n")) - 1, nil
}
