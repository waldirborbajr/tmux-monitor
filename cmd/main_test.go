package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

// Mock SSH client for testing
type mockSSHClient struct{}

func (m *mockSSHClient) Close() error {
	return nil
}

func (m *mockSSHClient) NewSession() (*ssh.Session, error) {
	return &ssh.Session{}, nil
}

// Mock SSH session for testing
type mockSSHSession struct{}

func (m *mockSSHSession) CombinedOutput(cmd string) ([]byte, error) {
	switch cmd {
	case "docker ps -a --format '{{.State}}' | sort | uniq -c":
		return []byte("2 running\n1 exited\n"), nil
	case "docker stats --no-stream --format \"{{.Container}}: {{.CPUPerc}} {{.MemPerc}}\"":
		return []byte("abc123: 5.00% 10.00%\ndef456: 3.00% 8.00%\n"), nil
	default:
		return nil, fmt.Errorf("unknown command")
	}
}

func (m *mockSSHSession) Close() error {
	return nil
}

func TestReadConfig(t *testing.T) {
	// Create a temporary config file
	tmpfile, err := os.CreateTemp("", "test-config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test configuration
	content := `SERVER_ADDRESS=example.com
SERVER_PORT=2222
SERVER_USER=testuser
SERVER_PASSWORD=testpass
UPDATE_INTERVAL=60`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading the config
	config, err := readConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("readConfig() error = %v", err)
	}

	expected := ServerConfig{
		Address:        "example.com",
		User:           "testuser",
		Password:       "testpass",
		Port:           2222,
		UpdateInterval: 60,
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("readConfig() = %v, want %v", config, expected)
	}
}

func TestParseContainerStats(t *testing.T) {
	input := "2 running\n1 exited\n3 stopped"
	expected := map[string]int{
		"running": 2,
		"exited":  1,
		"stopped": 3,
	}

	result := parseContainerStats(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseContainerStats() = %v, want %v", result, expected)
	}
}

func TestParseResourceUsage(t *testing.T) {
	input := "abc123: 5.00% 10.00%\ndef456: 3.00% 8.00%"
	expected := map[string]string{
		"abc123": "CPU: 5.00%, Mem: 10.00%",
		"def456": "CPU: 3.00%, Mem: 8.00%",
	}

	result := parseResourceUsage(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseResourceUsage() = %v, want %v", result, expected)
	}
}

func TestFormatOutput(t *testing.T) {
	stats := map[string]int{
		"running": 2,
		"exited":  1,
		"stopped": 0,
		"failed":  0,
		"dead":    0,
	}
	usage := map[string]string{
		"abc123": "CPU: 5.00%, Mem: 10.00%",
		"def456": "CPU: 3.00%, Mem: 8.00%",
	}

	expected := "🐳 Up: 2, Down: 1, Stopped: 0, Failed: 0, Died: 0 | abc123 (CPU: 5.00%, Mem: 10.00%) def456 (CPU: 3.00%, Mem: 8.00%) "

	result := formatOutput(stats, usage)
	if result != expected {
		t.Errorf("formatOutput() = %v, want %v", result, expected)
	}
}

func TestGetDockerStatus(t *testing.T) {
	// Mock the SSH connection
	originalSSHConnect := sshConnect
	sshConnect = func(config ServerConfig) (ssh.Client, error) {
		return &mockSSHClient{}, nil
	}
	defer func() { sshConnect = originalSSHConnect }()

	config := ServerConfig{
		Address:        "example.com",
		User:           "testuser",
		Password:       "testpass",
		Port:           22,
		UpdateInterval: 30,
	}

	expected := "🐳 Up: 2, Down: 1, Stopped: 0, Failed: 0, Died: 0 | abc12 (CPU: 5.00%, Mem: 10.00%) def45 (CPU: 3.00%, Mem: 8.00%) "

	result := getDockerStatus(config)
	if result != expected {
		t.Errorf("getDockerStatus() = %v, want %v", result, expected)
	}
}

func TestGetHostKeyCallback(t *testing.T) {
	callback, err := getHostKeyCallback()
	if err != nil {
		t.Fatalf("getHostKeyCallback() error = %v", err)
	}

	// Test with a mock host key
	mockHostname := "example.com"
	mockRemoteAddr := &net.TCPAddr{IP: net.ParseIP("192.0.2.1"), Port: 22}
	mockPublicKey := &ssh.PublicKey{}

	err = callback(mockHostname, mockRemoteAddr, *mockPublicKey)
	if err != nil {
		t.Errorf("HostKeyCallback error = %v", err)
	}
}

// TestMain is used to set up any necessary test environment
func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("HOME", "/tmp")

	// Run tests
	exitCode := m.Run()

	// Clean up test environment
	os.Unsetenv("HOME")

	os.Exit(exitCode)
}
