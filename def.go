package main

import (
	"github.com/kambahr/go-webconfig"

	"github.com/kambahr/go-webcache"
	"github.com/kambahr/go-webutil"
)

// Environment holds al global object.
type Environment struct {
	WebRootPath    string
	AppDataPath    string
	ConfigFilePath string
	Config         *webconfig.Config
	WebCach        *webcache.Cache
	WebUtil        *webutil.HTTP
}
