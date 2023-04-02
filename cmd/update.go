package cmd

import (
	"godkimflare/resource"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update DKIM Record on Cloudflare",
	RunE:  resource.ResourceUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("domain", "d", "", "Domain Name (Required)")
	updateCmd.Flags().StringP("selector", "s", "default", "DKIM Selector")
	updateCmd.Flags().StringP("dkimfile", "f", "", "Path to DKIM Private Key File (Required)")
	updateCmd.MarkFlagRequired("domain")
	updateCmd.MarkFlagRequired("dkimfile")
}
