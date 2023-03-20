package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type NFT struct {
	Image       string
	Title       string
	Description string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        ethereumAddress := r.FormValue("ethereum_address")
        nfts, err := getNFTs(ethereumAddress)
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        tmpl := template.Must(template.ParseFiles("NFTAPI.html"))
        tmpl.Execute(w, nfts)
        return
    }

    tmpl := template.Must(template.ParseFiles("index.html"))
    tmpl.Execute(w, nil)
}


func getNFTs(ethereumAddress string) ([]NFT, error) {
	alchemyAPIKey := "bkNVadRwumyDz1M1DaEsaEZAsQSz5j7L"
	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s/getNFTs?owner=%s", alchemyAPIKey, ethereumAddress)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type OwnedNFT struct {
		Metadata map[string]interface{} `json:"metadata"`
	}

	var result struct {
		OwnedNfts []OwnedNFT `json:"ownedNfts"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var nfts []NFT

	for _, ownedNFT := range result.OwnedNfts {
		image, _ := ownedNFT.Metadata["image"].(string)
		title, _ := ownedNFT.Metadata["name"].(string)
		description, _ := ownedNFT.Metadata["description"].(string)

		nft := NFT{
			Image:       image,
			Title:       title,
			Description: description,
		}
		nfts = append(nfts, nft)
	}

	return nfts, nil
}

