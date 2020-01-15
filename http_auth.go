package session

import "net/http"

func GetSidFromHeader(h *http.Request) (string,string) {
	u,p,ok := h.BasicAuth()
	if !ok {
		return "",""
	}
	return u,p
}
