package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

const socketContent = `[Unit]
Description=Socket for pareto-linux

[Socket]
ListenStream=/var/run/pareto-linux.sock
SocketMode=0666
Accept=no

[Install]
WantedBy=sockets.target`

const serviceContent = `[Unit]
Description=Service for pareto-linux
Requires=pareto-linux.socket

[Service]
ExecStart=/usr/bin/paretosecurity helper
User=root
Group=root
StandardInput=socket
Type=oneshot
RemainAfterExit=no

[Install]
WantedBy=multi-user.target`

func runHelper() {
	// Get the socket from file descriptor 0
	file := os.NewFile(0, "socket")
	listener, err := net.FileListener(file)
	if err != nil {
		log.Debugf("Failed to create listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	log.Info("Server is listening on Unix domain socket...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Debugf("Failed to accept connection: %v\n", err)
			continue
		}

		handleConnection(conn)
		break
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info("Connection received")

	// Handle the request
	_, err := conn.Write([]byte("Hello from Go app!\n"))
	if err != nil {
		log.Debugf("Failed to write to connection: %v\n", err)
	}
}

var helperCmd = &cobra.Command{
	Use:   "helper [--install]",
	Short: "A root helper",
	Long:  `A root helper that listens on a Unix domain socket and responds to authenticated requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		installFlag, _ := cmd.Flags().GetBool("install")
		if installFlag {
			installSystemdHelper()
			return
		}
		runHelper()
	},
}

func init() {
	rootCmd.AddCommand(helperCmd)
	helperCmd.Flags().Bool("install", false, "install root helper")
}

func installSystemdHelper() {
	systemdPath := "/etc/systemd/system"

	//ensure the user is root
	if os.Geteuid() != 0 {
		log.Fatal("This command must be run as root")
		return
	}

	// Create socket file
	socketPath := filepath.Join(systemdPath, "pareto-linux.socket")
	if err := os.WriteFile(socketPath, []byte(socketContent), 0644); err != nil {
		log.Infof("Failed to create socket file: %v\n", err)
		return
	}

	// Create service file
	servicePath := filepath.Join(systemdPath, "pareto-linux@.service")
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		fmt.Printf("Failed to create service file: %v\n", err)
		return
	}

	// Execute commands
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		log.Infof("Failed to reload systemd: %v\n", err)
		return
	}
	if err := exec.Command("systemctl", "enable", "pareto-linux.socket").Run(); err != nil {
		log.Infof("Failed to enable socket: %v\n", err)
		return
	}
	if err := exec.Command("systemctl", "start", "pareto-linux.socket").Run(); err != nil {
		log.Infof("Failed to start socket: %v\n", err)
		return
	}
}