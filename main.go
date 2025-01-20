package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func handler(w http.ResponseWriter, r *http.Request) {
	redirectURL := strings.TrimPrefix(r.URL.Path, "/")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Unable to parse IP address", http.StatusInternalServerError)
		return
	}

	// Test if IP is already in the set
	testCmd := exec.Command("ipset", "test", "gmad-whitelist", ip)
	if err := testCmd.Run(); err == nil {
		// IP is in the set, delete it first
		delCmd := exec.Command("ipset", "del", "gmad-whitelist", ip)
		if err := delCmd.Run(); err != nil {
			log.Printf("Failed to delete IP from gmad-whitelist: %v", err)
		}
	}

	// Add IP to the set
	addCmd := exec.Command("ipset", "add", "gmad-whitelist", ip, "timeout", "300")
	if err := addCmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("Command %v failed with exit code %d", addCmd, exitError.ExitCode())
		} else {
			http.Error(w, "Unable to add IP to gmad-whitelist", http.StatusInternalServerError)
			return
		}
	}

	if redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		fmt.Fprintf(w, "Your IP address is: %s\n", ip)
	}
}

func runCmds(commands [][]string) error {
	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				log.Printf("Command %v failed with exit code %d", cmd, exitError.ExitCode())
			} else {
				return fmt.Errorf("failed to execute command %v: %v", cmd, err)
			}
		}
	}

	return nil
}

func setupIptables(port string) error {
	commands := [][]string{
		{"ipset", "create", "gmad-whitelist", "hash:ip", "timeout", "300"},
		{"iptables", "-N", "GMAD_PROTECTED"},
		{"iptables", "-A", "GMAD_PROTECTED", "-m", "set", "--match-set", "gmad-whitelist", "src", "-j", "ACCEPT"},
		{"iptables", "-A", "GMAD_PROTECTED", "-m", "hashlimit", "--hashlimit", "35/sec", "--hashlimit-burst", "5", "--hashlimit-mode", "srcip", "--hashlimit-name", "main", "-j", "ACCEPT"},
		{"iptables", "-A", "GMAD_PROTECTED", "-j", "DROP"},
		{"iptables", "-I", "INPUT", "-p", "udp", "--dport", port, "-j", "GMAD_PROTECTED"},
	}

	return runCmds(commands)
}

func cleanIptables(port string) error {
	commands := [][]string{
		{"iptables", "-D", "INPUT", "-p", "udp", "--dport", port, "-j", "GMAD_PROTECTED"},
		{"iptables", "-F", "GMAD_PROTECTED"},
		{"iptables", "-X", "GMAD_PROTECTED"},
		{"ipset", "destroy", "gmad-whitelist"},
	}

	return runCmds(commands)
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	protectedPort := flag.String("protected-port", "9090", "Protected port to listen on")
	flag.Parse()

	if err := setupIptables(*protectedPort); err != nil {
		log.Fatalf("Failed to setup iptables: %v", err)
	}

	defer cleanIptables(*protectedPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		if err := cleanIptables(*protectedPort); err != nil {
			log.Printf("Failed to cleanup iptables: %v", err)
		}
		os.Exit(0)
	}()

	http.HandleFunc("/", handler)

	fmt.Printf("Starting server on :%s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
