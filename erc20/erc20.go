package erc20

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type Token struct {
	Name    string `json:"name"`
	Chain   int64  `json:"chain"`
	Address string `json:"address"`
}

var tokenCache = make(map[string]map[int64]string)

func LoadSupportedTokens() error {
	var contents string

	path := "../tokens/tokens.json"
	log.Info().Msgf("Pulling supported tokens from json file: %v", path)

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Send()
		return err
	}
	contents = string(contentBytes)
	initFromJsonString(contents)
	return nil
}

func GetTokenAddress(token string, chain int64) string {
	// var tokenDetails Token
	tokenAddress, ok := tokenCache[token][chain]
	if !ok {
		log.Error().Msgf("%v token details not found", token)
		return ""
	}

	return tokenAddress
}

func initFromJsonString(jsonString string) {
	var tokens []Token
	if err := json.Unmarshal([]byte(jsonString), &tokens); err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, token := range tokens {
		updateTokenCache(tokenCache, token.Name, token.Address, token.Chain)
	}
}

func updateTokenCache(cache map[string]map[int64]string, token, address string, chain int64) {
	// check if the token exists
	if _, ok := cache[token]; !ok {
		// create a new inner map
		cache[token] = make(map[int64]string)
	}
	cache[token][chain] = address
}
