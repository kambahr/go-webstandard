package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kambahr/go-webcache"
	"github.com/kambahr/go-webconfig"
	"github.com/kambahr/go-webutil"
)

// NewPage initalizes the NewPage; adds the globals.
func NewPage(webRootPath string) *Environment {
	//openssl req -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out webstandard.crt -keyout webstandard.key
	var e Environment
	e.WebRootPath = webRootPath
	e.WebCach = webcache.NewWebCache(5 * time.Minute)

	// Handles the css, js, and image files via the webutil package.
	e.WebUtil = webutil.NewHTTP(e.WebRootPath, time.Hour /* time to cache js, css,.. files*/)

	// Create the appdata if it does not exist
	e.AppDataPath = fmt.Sprintf("%s/appdata", e.WebRootPath)
	if !fileOrDirectoryExists(e.AppDataPath) {
		os.Mkdir(e.AppDataPath, os.ModePerm)
	}

	// Initializes an instance of webconfig. e.Config is refreshed
	// automatically. To update values:
	// e.Config.UpdateConfigValue(<key>, <value>).
	// See https://github.com/kambahr/go-webconfig.
	e.Config = webconfig.NewWebConfig(e.WebRootPath)

	return &e
}
