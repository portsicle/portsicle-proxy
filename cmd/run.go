package cmd

import (
	"log"

	"github.com/amitsuthar69/portsicle-proxy/internal/proxy"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the proxy server",
	Run:   initServer,
}

func initServer(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		log.Fatalf("Error retrieving port flag: %v", err)
	}

	// Proxy with default config
	p := proxy.NewProxy(proxy.Config{
		ListenAddr: ":" + port,
	})

	log.Printf("proxy running on port=%s", port)

	if err := p.Start(); err != nil {
		log.Fatalf("Proxy startup failed: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("port", "p", "8888", "Port on which proxy is listening")
}
