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

func ResourceUpdate(cmd *cobra.Command, args []string) error {
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

	putToCloudflare(selector+"._domainkey."+domain, "TXT", genCloudflareReq("TXT", selector+"._domainkey.", content, "Updated"))

	return nil
}

func putToCloudflare(nameanddomain string, recordtype string, putBody string) {
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

	searchurl := "https://api.cloudflare.com/client/v4/zones/" + res.Result[0].ID + "/dns_records"
	req2, err2 := http.NewRequest("GET", searchurl, nil)
	if err2 != nil {
		log.Println(err2)
		os.Exit(1)
	}

	req2.Header.Add("Authorization", bearer)
	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp2.Body.Close()
	body2, err2 := io.ReadAll(resp2.Body)
	if err2 != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var recordsres RecordsRes
	if err2 := json.Unmarshal(body2, &recordsres); err2 != nil {
		log.Println(err2)
		os.Exit(1)
	}

	var did string

	for i := range recordsres.Result {
		if nameanddomain == recordsres.Result[i].Name {
			did = recordsres.Result[i].ID
		}
	}

	puturl := "https://api.cloudflare.com/client/v4/zones/" + res.Result[0].ID + "/dns_records/" + did
	var jsonStr = []byte(putBody)
	req3, err3 := http.NewRequest("PUT", puturl, bytes.NewBuffer(jsonStr))
	if err3 != nil {
		log.Println(err3)
		os.Exit(1)
	}

	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Add("Authorization", bearer)

	client3 := &http.Client{}
	resp3, err3 := client3.Do(req3)
	if err3 != nil {
		log.Println(err3)
		os.Exit(1)
	}
	defer resp.Body.Close()

	log.Println("Cloudflare Response Status for: "+recordtype, resp3.Status)
}
