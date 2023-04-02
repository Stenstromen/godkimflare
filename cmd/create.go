package cmd

import (
	"godkimflare/resource"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create DKIM Record on Cloudflare",
	RunE:  resource.ResourceCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("domain", "d", "", "Domain Name (Required)")
	createCmd.Flags().StringP("selector", "s", "default", "DKIM Selector")
	createCmd.Flags().StringP("dkimfile", "f", "", "Path to DKIM Private Key File (Required)")
	createCmd.MarkFlagRequired("domain")
	createCmd.MarkFlagRequired("dkimfile")
}
