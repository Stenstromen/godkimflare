package resource

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/peterzen/goresolver"
	"github.com/spf13/cobra"
)

func ResourceVerify(cmd *cobra.Command, args []string) error {
	var successCount int

	domain, err := cmd.Flags().GetString("domain")
	if err != nil {
		return err
	}

	selector, err := cmd.Flags().GetString("selector")
	if err != nil {
		return err
	}
	dkimfile, err := cmd.Flags().GetString("dkimfile")
	if err != nil {
		return err
	}

	resolver, err := goresolver.NewResolver("/etc/resolv.conf")
	if err != nil {
		fmt.Printf("Cannot initialize the local resolver: %s\n", err)
		os.Exit(1)
	}

	dkimdnstxt := strings.Split(cloudflareDNSLookupTXT(selector+"._domainkey."+domain), "p=")[1]

	if strings.Contains(cloudflareDNSLookupTXT(selector+"._domainkey."+domain), "v=DKIM1;") {
		fmt.Println("OK - DKIM1 TXT record exists")
		successCount++
	} else {
		fmt.Println("FAIL - No DKMI1 TXT record exists")
	}

	if dkimdnstxt == getb64sum(dkimfile) {
		fmt.Println("OK - DKIM1 TXT record matches DKIM private key")
		successCount++
	} else {
		fmt.Println("FAIL - DKIM1 TXT record does not match DKIM private key")
	}

	result, err := resolver.StrictNSQuery(dns.Fqdn(selector+"._domainkey."+domain), dns.TypeTXT)
	if err != nil {
		fmt.Println("INFO - DKIM1 TXT not DNSSEC signed")
	}
	if result != nil {
		fmt.Println("OK - DKIM1 TXT DNSSEC signed")
	}

	fmt.Println("Success Count: " + fmt.Sprint(successCount) + "/2")

	return nil
}

func cloudflareDNSLookupTXT(domainrecord string) string {
	resolve := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dial := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return dial.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}

	txtresult, _ := resolve.LookupTXT(context.Background(), domainrecord)

	return txtresult[0]
}
