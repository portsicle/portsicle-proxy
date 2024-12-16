package cmd

import (
	"log"
	"strings"

	"github.com/amitsuthar69/attorney-toolkit/internal/database"
	"github.com/spf13/cobra"
)

var (
	addDomain    string
	removeDomain string
)

// blockdomainCmd defines the "block" command
var blockdomainCmd = &cobra.Command{
	Use:   "block",
	Short: "Manage blocked domains",
	Long: `Allows you to add or remove domains from the blocklist.
	--add to block a domain. 
	--remove to unblock a domain.`,
	Run: block,
}

func block(cmd *cobra.Command, args []string) {

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	addDomain, _ := cmd.Flags().GetString("add")
	removeDomain, _ := cmd.Flags().GetString("remove")

	if addDomain != "" && removeDomain != "" {
		log.Println("Error: You cannot use --add and --remove together.")
		return
	}

	if addDomain != "" {
		if err := database.AddBlockedDomain(db, normalizeDomain(addDomain)); err != nil {
			log.Fatalf("Error adding domain to blocklist: %v", err)
		}
		return
	}

	if removeDomain != "" {
		if err := database.RemoveBlockedDomain(db, normalizeDomain(removeDomain)); err != nil {
			log.Fatalf("Error removing domain from blocklist: %v", err)
		}
		return
	}

	log.Println("Error: No domain specified. Use --add or --remove flags.")

}

func init() {
	rootCmd.AddCommand(blockdomainCmd)
	blockdomainCmd.Flags().StringP("add", "a", "", "Add a domain to the blocklist")
	blockdomainCmd.Flags().StringP("remove", "r", "", "Remove a domain from the blocklist")
}

func normalizeDomain(domain string) string {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	parts := strings.Split(domain, ":")
	domain = parts[0]

	domain = strings.ToLower(domain)
	return domain
}
