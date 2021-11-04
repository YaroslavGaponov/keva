package utils

import "os"

func GetEnvVariableOrDefult(name, def_value string) string {
	if value, found := os.LookupEnv(name); found {
		return value
	}
	return def_value
}
