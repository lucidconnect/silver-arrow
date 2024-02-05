package erc20

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Token struct {
	Name        string `json:"name"`
	Chain       int64  `json:"chain"`
	Address     string `json:"address"`
	MinorFactor int64  `json:"minorFactor"`
}

var tokenCache = make(map[string]map[int64]string)

func LoadSupportedTokens(tokenPath string, db *gorm.DB) error {
	var contents string

	log.Info().Msgf("Pulling supported tokens from json file: %v", tokenPath)

	contentBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Fatal().Err(err).Send()
		return err
	}
	contents = string(contentBytes)
	tokenConfigs := initFromJsonString(contents)

	if db != nil {
		var existing models.Token
		for _, tokenConfig := range tokenConfigs {
			token := &models.Token{
				Name:        tokenConfig.Name,
				Chain:       tokenConfig.Chain,
				Address:     tokenConfig.Address,
				MinorFactor: tokenConfig.MinorFactor,
			}
			if errors.Is(db.Where(&models.Token{Name: tokenConfig.Name, Chain: tokenConfig.Chain}).First(&existing).Error, gorm.ErrRecordNotFound) {
				db.Create(&token)
			} else {
				db.Where(&models.Token{Name: tokenConfig.Name, Chain: tokenConfig.Chain}).Updates(&tokenConfig)
			}
		}
	}
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

func GetNativeToken(chain int64) string {
	var token string

	switch chain {
	case 10:
		token = "ETH"
	case 80001:
		token = "MATIC"
	}
	return token
}

func initFromJsonString(jsonString string) []Token {
	var tokens []Token
	if err := json.Unmarshal([]byte(jsonString), &tokens); err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, token := range tokens {
		updateTokenCache(tokenCache, token.Name, token.Address, token.Chain)
	}
	return tokens
}

func updateTokenCache(cache map[string]map[int64]string, token, address string, chain int64) {
	// check if the token exists
	if _, ok := cache[token]; !ok {
		// create a new inner map
		cache[token] = make(map[int64]string)
	}
	cache[token][chain] = address
}
