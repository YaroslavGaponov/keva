package client

import (
	"strconv"
	"time"

	"github.com/YaroslavGaponov/keva/pkg/utils"
)

const (
	defaultTimeout = 5 * time.Second
)

type ClientOptions struct {
	timeout time.Duration
}

func ClientOptionsFromEnv() *ClientOptions {

	timeout := defaultTimeout

	if value, err := strconv.Atoi(utils.GetEnvVariableOrDefult("TIMEOUT", "5")); err == nil {
		timeout = time.Duration(value) * time.Second
	}
	return &ClientOptions{
		timeout: timeout,
	}
}
