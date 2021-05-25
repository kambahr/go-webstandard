package main

import (
	"time"

	"github.com/kambahr/go-webcache"
	"github.com/kambahr/go-webutil"
)

// Page holds the content page and the year.
type Page struct {
	MainContent string // {{.MainContent}}
	ThisYear    string // {{.ThisYear}}
}

// Environment holds al global object.
type Environment struct {
	WebRootPath string
	WebCach     *webcache.Cache
	WebUtil     *webutil.HTTP
}

// NewPage initalizes the NewPage; adds the globals.
func NewPage(webRootPath string) *Environment {
	var e Environment
	e.WebRootPath = webRootPath
	e.WebCach = webcache.NewWebCache(5 * time.Minute)

	// Handles the css, js, and image files via the webutil package.
	e.WebUtil = webutil.NewHTTP(e.WebRootPath, time.Hour /* time to cache js, css,.. files*/)

	return &e
}
