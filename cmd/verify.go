package cmd

import (
	"godkimflare/resource"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify DKIM Record on Cloudflare",
	RunE:  resource.ResourceVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringP("domain", "d", "", "Domain Name (Required)")
	verifyCmd.Flags().StringP("selector", "s", "default", "DKIM Selector")
	verifyCmd.Flags().StringP("dkimfile", "f", "", "Path to DKIM Private Key File (Required)")
	verifyCmd.MarkFlagRequired("domain")
	verifyCmd.MarkFlagRequired("dkimfile")
}
