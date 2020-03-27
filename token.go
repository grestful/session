package session

import (
	"net/url"
)

func GetSidFromQuery(values url.Values, authName string) string {
	if authName == "" {
		authName = "token"
	}
	for name, value := range values {
		if name == authName {
			return value[0]
		}
	}

	return ""
}
