package session

import "net/http"

func GetSidFromCookies(cook []*http.Cookie, authName string) string {
	if authName == "" {
		authName = "token"
	}
	for _, c := range cook {
		if c.Name == authName {
			return c.Value
		}
	}

	return ""
}
