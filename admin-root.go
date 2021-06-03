package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (e *Environment) getRemoteIP(r *http.Request) string {
	return strings.Split(strings.Replace(strings.Replace(r.RemoteAddr, "[", "", -1), "]", "", -1), ":")[0]
}
func (e *Environment) validateAdminRequest(w http.ResponseWriter, r *http.Request) bool {

	ip := e.getRemoteIP(r)

	// if blank it's local ip.
	if ip != "" {
		ipFound := false
		for i := 0; i < len(e.Config.Admin.AllowedIP); i++ {
			if ip == e.Config.Admin.AllowedIP[i] {
				ipFound = true
				break
			}
		}
		if !ipFound {
			// Not authorized, but return not found,
			// using the public master page.
			e.displayError(w, http.StatusNotFound)
			return false
		}
	}

	// rPath := strings.ToLower(r.URL.Path)

	return true
}
func (e *Environment) adminRoot(w http.ResponseWriter, r *http.Request) {

	if !e.validateAdminRequest(w, r) {
		return
	}

	fPath := fmt.Sprintf("%s/html/admin/master.html", e.WebRootPath)
	master, _ := ioutil.ReadFile(fPath)

	fPath = fmt.Sprintf("%s/html/admin/index.html", e.WebRootPath)
	bPage, _ := ioutil.ReadFile(fPath)

	bFinal := bytes.Replace(master, []byte("{{.AdminMainContent}}"), bPage, -1)

	bFinal = e.WebUtil.RemoveCommentsFromByBiteArry(bFinal, "{{.COMMENT ", "}}")

	w.Write(bFinal)
}
