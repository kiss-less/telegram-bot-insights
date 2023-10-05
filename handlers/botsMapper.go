package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BotsKeysMapping struct {
	ID     int
	APIKey string
}

func CreateBotsKeysMapping(args Arguments) ([]BotsKeysMapping, error) {
	envVar := args.BotsApiEnvVar
	var botsKeysMapping []BotsKeysMapping

	botsAPIKeys := os.Getenv(envVar)

	if botsAPIKeys == "" {
		err := fmt.Errorf("%v env var is not set", envVar)

		return nil, err
	}

	keyValuePairs := strings.Split(botsAPIKeys, " ")

	for _, pair := range keyValuePairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				err := fmt.Errorf("issue with converting %v to int", parts[0])

				return nil, err
			}
			apiKey := parts[1]
			mapping := BotsKeysMapping{
				ID:     id,
				APIKey: apiKey,
			}
			botsKeysMapping = append(botsKeysMapping, mapping)
		}
	}

	return botsKeysMapping, nil
}
