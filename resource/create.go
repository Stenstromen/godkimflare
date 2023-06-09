package resource

import (
	"bytes"
	"encoding/json"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func ResourceCreate(cmd *cobra.Command, args []string) error {
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

	content := "v=DKIM1;h=sha256;k=rsa;s=email;p=" + getb64sum(dkimfile)

	postToCloudflare(domain, "TXT", genCloudflareReq("TXT", selector+"._domainkey", content, "Created"))

	return nil
}

func postToCloudflare(nameanddomain string, recordtype string, postBody string) {
	url := "https://api.cloudflare.com/client/v4/zones"
	var bearer = "Bearer " + os.Getenv("TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var res Res
	if err := json.Unmarshal(body, &res); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var did string

	for i := range res.Result {
		if nameanddomain == res.Result[i].Name {
			did = res.Result[i].ID
		}
	}

	posturl := "https://api.cloudflare.com/client/v4/zones/" + did + "/dns_records"
	var jsonStr = []byte(postBody)
	req2, err2 := http.NewRequest("POST", posturl, bytes.NewBuffer(jsonStr))
	if err2 != nil {
		log.Println(err2)
		os.Exit(1)
	}

	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Add("Authorization", bearer)

	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if resp2.StatusCode != 200 {
		log.Println("Error while posting to Cloudflare: Invalid TOKEN och Record Already Exists")
		os.Exit(1)
	}
	if err2 != nil {
		log.Println(err2)
		os.Exit(1)
	}

	defer resp.Body.Close()

	log.Println("Cloudflare Response Status for: "+recordtype, resp2.Status)
}
