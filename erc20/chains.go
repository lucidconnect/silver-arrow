package erc20

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

type Chain struct {
	Chain    int64  `json:"chain"`
	Explorer string `json:"explorer"`
}

var chainCache = make(map[int64]string)

func LoadSupportedChains(path string) error {
	var contents string

	log.Info().Msgf("Pulling supported chains from json file: %v", path)

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Send()
		return err
	}
	contents = string(contentBytes)
	initChainsFromJsonString(contents)
	return nil
}

func GetChainExplorer(chain int64) (string, error) {
	explorer, ok := chainCache[chain]
	if !ok {
		err := errors.New("chain not supported")
		log.Err(err).Send()
		return "", err
	}
	return explorer, nil
}

func initChainsFromJsonString(jsonString string) {
	var chains []Chain
	if err := json.Unmarshal([]byte(jsonString), &chains); err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, chain := range chains {
		updateChainCache(chainCache, chain.Chain, chain.Explorer)
	}
}

func updateChainCache(cache map[int64]string, chain int64, explorer string) {
	cache[chain] = explorer
}
