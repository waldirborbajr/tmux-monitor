package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	redBold        = "\x1b[1;31m"
	resetColor     = "\x1b[0m"
	configFile     = "$HOME/.tmux-monitor"
	knownHostsFile = "$HOME/.ssh/known_hosts"
)

var version = "[dev build]"

type ServerConfig struct {
	Address  string
	User     string
	Password string
	Port     int
}

func main() {
	app := cli.NewApp()

	app.Version = version
	app.Name = "tmux-monitor"
	app.HelpName = "tmux-monitor"
	app.Usage = "monitor Docker containers running on a remote server"
	app.Description = strings.Join([]string{
		"tmux-monitor connects to remote server and get containers status.",
		"This allows to monitor healthcheck of docker containers.",
	}, " ")
	app.Author = "Waldir \"Borba\" Junior"
	app.Email = "wborbajr@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "short",
			ShortName: "s",
			Usage:     "show short stats version",
			Action:    printShort,
		},
		{
			Name:      "full",
			ShortName: "f",
			Usage:     "show full stats version",
			Action:    printShort,
		},
	}

	app.Run(os.Args)
}

func printShort(*cli.Context) {
	config, err := readConfig(configFile)
	if err != nil {
		fmt.Printf("%sError reading config: %v%s\n", redBold, err, resetColor)
		os.Exit(1)
	}

	status := getDockerStatus(config)
	fmt.Println(status)
}

func readConfig(filename string) (ServerConfig, error) {
	file, err := os.Open(os.ExpandEnv(filename))
	if err != nil {
		return ServerConfig{}, err
	}
	defer file.Close()

	config := ServerConfig{Port: 22} // Default SSH port
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "SERVER_ADDRESS":
			config.Address = value
		case "SERVER_PORT":
			config.Port, _ = strconv.Atoi(value)
		case "SERVER_USER":
			config.User = value
		case "SERVER_PASSWORD":
			config.Password = value
		}
	}

	if config.Address == "" || config.User == "" || config.Password == "" {
		return ServerConfig{}, fmt.Errorf("missing required configuration")
	}

	return config, scanner.Err()
}

func getDockerStatus(config ServerConfig) string {
	client, err := sshConnect(config)
	if err != nil {
		return fmt.Sprintf("%s⚠ Unable to connect to server: %v%s ", redBold, err, resetColor)
	}
	defer client.Close()

	containerStatsCmd := "docker ps -a --format '{{.State}}' | sort | uniq -c"
	resourceUsageCmd := "docker stats --no-stream --format \"{{.Container}}: {{.CPUPerc}} {{.MemPerc}}\""

	containerStatsOutput, err := runCommand(client, containerStatsCmd)
	if err != nil {
		return fmt.Sprintf("%s⚠ Error getting container stats: %v%s ", redBold, err, resetColor)
	}

	resourceUsageOutput, err := runCommand(client, resourceUsageCmd)
	if err != nil {
		return fmt.Sprintf("%s⚠ Error getting resource usage: %v%s ", redBold, err, resetColor)
	}

	containerStats := parseContainerStats(containerStatsOutput)
	resourceUsage := parseResourceUsage(resourceUsageOutput)

	return formatOutput(containerStats, resourceUsage)
}

func sshConnect(config ServerConfig) (*ssh.Client, error) {
	hostKeyCallback, err := getHostKeyCallback()
	if err != nil {
		return nil, fmt.Errorf("failed to get host key callback: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: hostKeyCallback,
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Address, config.Port), sshConfig)
}

func getHostKeyCallback() (ssh.HostKeyCallback, error) {
	knownHosts := os.ExpandEnv(knownHostsFile)
	callback, err := knownhosts.New(knownHosts)
	if err != nil {
		return nil, err
	}
	return ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := callback(hostname, remote, key)
		if err == nil {
			return nil
		}
		if keyErr, ok := err.(*knownhosts.KeyError); ok {
			if len(keyErr.Want) == 0 {
				// Key not in known_hosts, add it
				f, err := os.OpenFile(knownHosts, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = f.WriteString(knownhosts.Line([]string{hostname}, key) + "\n")
				return err
			}
		}
		return err
	}), nil
}

func runCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func parseContainerStats(input string) map[string]int {
	states := make(map[string]int)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			count, _ := strconv.Atoi(fields[0])
			state := fields[1]
			states[state] = count
		}
	}
	return states
}

func parseResourceUsage(input string) map[string]string {
	usage := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			containerID := strings.TrimSuffix(fields[0], ":")
			cpu := fields[1]
			mem := fields[2]
			usage[containerID] = fmt.Sprintf("CPU: %s, Mem: %s", cpu, mem)
		}
	}
	return usage
}

func formatOutput(stats map[string]int, usage map[string]string) string {
	output := fmt.Sprintf("🐳 U: %d, D: %d, S: %d, F: %d, X: %d | ",
		stats["running"], stats["exited"], stats["stopped"], stats["failed"], stats["dead"])

	for container, resources := range usage {
		output += fmt.Sprintf("%s (%s) ", container[:6], resources)
	}

	return output
}
