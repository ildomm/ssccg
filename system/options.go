package system

import (
	"log"
	"os"
	"strconv"
)

const (
	ListenAddressEnvVar = "SERVER_PORT"
)

// ExtractServerPort extracts the server port from the environment variable SERVER_PORT.
func ExtractServerPort() *int {
	if env, found := os.LookupEnv(ListenAddressEnvVar); found {
		value, err := strconv.Atoi(env)

		if err != nil {
			log.Println("Could not parse server port address from environment variable ", ListenAddressEnvVar)
			return nil
		}

		return &value
	}

	return nil
}
