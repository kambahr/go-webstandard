package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	const defaultPortFlag int = -32956103397293
	var portNoArg int
	flag.IntVar(&portNoArg, "portno", defaultPortFlag, "tcp/ip port no to listen")
	flag.Parse()

	if portNoArg < 1 && portNoArg != defaultPortFlag {
		fmt.Println("Invalid port number")
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// To handle requests for html file.
	page := NewPage(dir)

	// First look at the config.
	portNo := page.Config.PortNo

	// If port no is passed from command line, it will take priority.
	if portNoArg > 0 {
		portNo = portNoArg
	}

	// Handler for static files.
	http.HandleFunc("/assets/", page.WebUtil.ServeStaticFile)

	// Root
	http.HandleFunc("/", page.FromRoot)

	// setup the http server.
	svr := http.Server{
		Addr:           fmt.Sprintf(":%d", portNo),
		MaxHeaderBytes: 20480,

		// This seems to prevent http2 calls by the client browser
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),

		// No checks on the client cert
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},

		// These should reflect what you expect from your site
		// also see comments in connState().
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,

		// More options before the request gets to the handlers.
		ConnState: page.connState,
	}
	fmt.Println("Listening to port:", portNo, "for host:", page.Config.HostName)
	if page.Config.Proto == "HTTP" {
		log.Fatal(svr.ListenAndServe())
	} else {
		if fileOrDirectoryExists(page.Config.TLS.CertFilePath) &&
			fileOrDirectoryExists(page.Config.TLS.KeyFilePath) {

			// both files must be in PEM format.
			//   Note that there is no separate arg to apply the
			//   chaine certficate. The one cert file can include the
			//   chain (if any) along with the root cert.

			// Assume that they are valid and run
			log.Fatal(svr.ListenAndServeTLS(page.Config.TLS.CertFilePath, page.Config.TLS.KeyFilePath))
		}
	}
}

// connState enables you to monitor callers before their
// requests get to the handlers. You can use this for security
// performance enhancements, or monitoring (i.e. number of
// active connections).
func (e *Environment) connState(conn net.Conn, connState http.ConnState) {

	// You can block target IP addresses here by closing connection of
	// the offenders or the ones that are not be in your authorized group.

	// conn.RemoteAddr().String() is in the form: <IP address>:<identifier> e.g. 127.0.0.1:35692.
	// The numbers after the colon is a unique id given to each request.
	//    For example:
	//    http://localhost:8000/mycss.css........127.0.0.1:35692
	//    http://localhost:8000/myjs.js..........127.0.0.1:35693
	//
	// If you use XMLHttpRequest to get a response back. The identifier will remain
	// the same even though you may get a different content.
	//    For example:
	//    http://localhost:8000/page1.html........127.0.0.1:35692
	//    http://localhost:8000/page2.html........127.0.0.1:35692
	//
	// For the most part, you can ignore the identifier, however, you can make use
	// of it -- to build on security, content/performance-smart concepts to enhance
	// your website.

	// Check the blocked ips
	// TODO: log further to revoke permanently or re-instate after a period.
	ip := strings.Split(strings.Replace(strings.Replace(conn.RemoteAddr().String(), "[", "", -1), "]", "", -1), ":")[0]

	// blank means its ::1 (ipv6 loopback ip)
	if ip != "" {
		for i := 0; i < len(e.Config.BlockedIP); i++ {
			if e.Config.BlockedIP[i] == ip {
				conn.Close()
				return
			}
		}
	}

	// These help reduce the too many open files errors.
	maxConnIODeadLine := 20      // Limit the IO to 20 seconds.
	maxConnIOWriteDeadLine := 30 // Extra time for over timeout.
	conn.SetDeadline(time.Now().Add(time.Duration(maxConnIODeadLine) * time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))

	// Uncomment the following, if you want to do anything on the active, idle, and closed states events.
	// csStr := fmt.Sprintf("%v", connState)
	// if csStr == "new" {
	// } else if csStr == "active" {
	// } else if csStr == "idle" {
	// } else if csStr == "closed" {
	// }
}
