package main

import (
	"encoding/json"
	"math/rand"
	"os"
)

// find dublicates -> generate replacement -> apply replacement -> save replacement to file

type nftDublicatesRecord struct {
	TokenID string `json:"token_id"`
	OldURI  string `json:"old_uri"`
	NewURI  string `json:"new_uri"`
}

func extractNFTDublicates(colls []CollectionNew) []nftDublicatesRecord {
	var seenURI = make(map[string]bool)
	var dublURI = make(map[string]bool)
	var records = []nftDublicatesRecord{}
	for _, coll := range colls {
		for _, token := range coll.Tokens {
			if seenURI[token.URI] {
				dublURI[token.URI] = true
			}
			seenURI[token.URI] = true
		}
	}
	for _, coll := range colls {
		for _, token := range coll.Tokens {
			if dublURI[token.URI] {
				records = append(records, nftDublicatesRecord{
					TokenID: token.ID,
					OldURI:  token.URI,
				})
			}
		}
	}
	return records
}

func generateReplacements(uriPrefix string, records *[]nftDublicatesRecord) {
	for i := range *records {
		(*records)[i].NewURI = randomSlug()
	}
}

func fixNFTDublicates(colls *[]CollectionNew, records []nftDublicatesRecord) {
	for i, coll := range *colls {
		for j, token := range coll.Tokens {
			for _, rec := range records {
				if token.ID == rec.TokenID {
					(*colls)[i].Tokens[j].URI = rec.NewURI
					break
				}
			}
		}
	}
}

func exportNFTDublicates(fname string, records []nftDublicatesRecord) error {
	if fname == "" {
		return nil
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	bz, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(bz)
	if err != nil {
		return err
	}
	return nil
}

func randomSlug() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 32)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
