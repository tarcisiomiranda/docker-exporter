package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the Docker metrics exporter",
	Run: func(cmd *cobra.Command, args []string) {
		restartServer()
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func restartServer() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), executable, "serve", "--port", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server restarted successfully")
	os.Exit(0)
}
