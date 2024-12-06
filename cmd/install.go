package cmd

import (
	"fmt"
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

func getServiceContent() string {
	return fmt.Sprintf(`[Unit]
Description=Service for pareto-linux
Requires=pareto-linux.socket

[Service]
ExecStart=%s
User=root
Group=root
StandardInput=socket
Type=oneshot
RemainAfterExit=no

[Install]
WantedBy=multi-user.target`, os.Args[0])
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install root helper",
	Run: func(cmd *cobra.Command, args []string) {
		installSystemdHelper()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
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
	if err := os.WriteFile(servicePath, []byte(getServiceContent()), 0644); err != nil {
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
