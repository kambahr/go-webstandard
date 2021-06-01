package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (e *Environment) validateRequest(w http.ResponseWriter, r *http.Request) bool {

	// For performance reasons, keep all logic here (no other functions),
	// even for short lookups that may repeat a couple of times!

	// Note: blocked ips have already been dropped.

	rPath := strings.ToLower(r.URL.Path)

	// Method allowed
	failed := true
	s := r.Method
	for i := 0; i < len(e.Config.HTTP.AllowedMethods); i++ {
		if s == e.Config.HTTP.AllowedMethods[i] {
			failed = false
			break
		}
	}
	if failed {
		e.displayError(w, 405)
		return false
	}

	// This is the order: restrict-paths, exclude-path, forward-paths

	// restrict-paths
	for i := 0; i < len(e.Config.URLPaths.Restrict); i++ {
		sl := strings.ToLower(e.Config.URLPaths.Restrict[i])
		// On-error  the err-text is placed instead of the
		// expected value starting with: ~@error
		if strings.HasPrefix(sl, "~@error") {
			continue
		}
		if sl == rPath {
			e.displayError(w, http.StatusUnauthorized)
			return false
		}
	}

	// exclude-path
	for i := 0; i < len(e.Config.URLPaths.Exclude); i++ {
		sl := strings.ToLower(e.Config.URLPaths.Exclude[i])
		if strings.HasPrefix(sl, "~@error") {
			continue
		}
		if sl == rPath {
			e.displayError(w, http.StatusNotFound)
			return false
		}
	}

	// forward-paths
	for i := 0; i < len(e.Config.URLPaths.Forward); i++ {
		sl := strings.ToLower(e.Config.URLPaths.Forward[i])
		v := strings.Split(sl, "|")
		left := ""
		right := ""
		if len(v) > 1 {
			left = v[0]
			right = v[1]
		}
		if right == "" || strings.HasPrefix(right, "~@error") {
			continue
		}
		if left == rPath {
			http.Redirect(w, r, right, http.StatusTemporaryRedirect)
			return false
		}
	}

	return true
}
func (e *Environment) displayError(w http.ResponseWriter, errCode int) {

	bMaster := e.getRawMaster()

	targetPagePhysPath := fmt.Sprintf("%s/html/errors/%d.html", e.WebRootPath, errCode)

	// Note that the end-use can't get to the /appdata directory.
	if !fileOrDirectoryExists(targetPagePhysPath) {
		// Server flat
		msg := fmt.Sprintf("Error %d - %s", errCode, http.StatusText(errCode))
		fmt.Fprint(w, msg)
		return
	}

	targetPageBytes, _ := ioutil.ReadFile(targetPagePhysPath)

	// This is all contents of the target page into the {{.MainContent}} block inside the master page.
	bFinal := bytes.Replace(bMaster, []byte("{{.MainContent}}"), targetPageBytes, 1)

	bFinal = e.applyPageVars(bFinal, "/", w)

	w.WriteHeader(errCode)

	w.Write(bFinal)
}
