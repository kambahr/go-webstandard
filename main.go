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
	"time"
)

func main() {

	var portNoArg int
	flag.IntVar(&portNoArg, "portno", 1265, "tcp/ip port no to listen - defaults to port 1265")
	flag.Parse()
	fmt.Println("PortNo: ", portNoArg)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// To handle html pages.
	page := NewPage(dir)

	http.HandleFunc("/assets/", page.WebUtil.ServeStaticFile)

	// Root
	http.HandleFunc("/", page.FromRoot)

	// setup the http server
	svr := http.Server{
		Addr:           fmt.Sprintf(":%v", portNoArg),
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
		ConnState: connState,
	}

	// If you are using https, you can call the following:
	// both files must be in PEM format.
	//   Note that there is no separate arg to apply the
	//   chaine certficate. The one cert file can include the
	//   chain (if any) along with the root cert.
	// tlsCertCertPath:= <full path to the cert file>
	// tlsCertKeyPath:= <full path to the private key file>
	// log.Fatal(svr.ListenAndServeTLS(tlsCertCertPath, tlsCertKeyPath))

	log.Fatal(svr.ListenAndServe())
}

// foo handles monitoring/security tasks. Note that this func
// must run in a separate thread than that of the current (in-process).
// func (ws *website) foo(ipAddr string) {
// 	  TODO:
// }

// connState enables you to monitor callers before their
// requests get to the handlers. You can use this for security
// performance enhancements, or monitoring (i.e. number of
// active connections).
func connState(conn net.Conn, connState http.ConnState) {

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

	// If you want to do any processing here, you'd have to use the go statement
	// so that the processing is done outside of this func (in a speparate
	// thread) -- basically you never want to wait for a func here.
	//
	// TODO: go foo(conn.RemoteAddr().String())

	// These help reduce the too many open files errors.
	maxConnIODeadLine := 20      // Limit the IO to 20 seconds.
	maxConnIOWriteDeadLine := 30 // Extra time for over timeout.
	conn.SetDeadline(time.Now().Add(time.Duration(maxConnIODeadLine) * time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))

	// Uncomment the following, if you want to do anything on the active, idle, and closed states events.
	//
	// csStr := fmt.Sprintf("%v", connState)
	// if csStr == "new" {
	// } else if csStr == "active" {
	// } else if csStr == "idle" {
	// } else if csStr == "closed" {
	// }
}
