package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// getGalleryCarousel gets the files name in /img/g and concatenates
// html tags that fill the value in the innerHTML of <div class="carousel-inner">
// via {{.Carousel}}.
func (e Environment) getGalleryCarousel() []byte {
	var carousel string

	// gallery images are saved in /img/g dir, so that the
	// filtering for files will be simple.
	p := fmt.Sprintf("%s/assets/img/g", e.WebRootPath)
	f, _ := getFiles(p, ".jpg")
	activeIndx := randInt(0, len(f))

	for i := 0; i < len(f); i++ {
		actv := ""
		if i == activeIndx {
			actv = " active"
		}
		h := fmt.Sprintf(`
		<div class="carousel-item%s">			
		<img src="/assets/img/g/%s" class="d-block w-100" alt="%s">
		</div>	
	`, actv, f[i] /*img src*/, f[i] /*img alt*/)

		// put all the carousel items into one string.
		carousel = fmt.Sprintf("%s%s", carousel, h)
	}
	return []byte(carousel)
}

// getPhysicalPath gets the full path of a page.
func (e Environment) getPhysicalPath(rPath string) (string, bool) {
	targetPagePhysPath := fmt.Sprintf("%s/html%s", e.WebRootPath, rPath)

	// Note that the end-use can't get to the /appdata directory.
	if !fileOrDirectoryExists(targetPagePhysPath) {
		return fmt.Sprintf("%s/html/404.html", e.WebRootPath), false
	}

	return targetPagePhysPath, true
}

// getTargetPageFromCache get the contents of a page (in []byte),
// from the cache object.
func (e Environment) getTargetPageFromCache(rPath string) []byte {
	var targetPageBytes []byte

	targetPagePhysPath, pageFound := e.getPhysicalPath(rPath)

	if !pageFound {
		targetPageBytes, err := ioutil.ReadFile(targetPagePhysPath)
		if err != nil {
			fmt.Println("FromRoot()=>", err)
		}
		return targetPageBytes
	}
	exists := e.WebCach.Exists(rPath)
	if !exists {
		// Read the file from disk.
		targetPageBytes, _ = ioutil.ReadFile(targetPagePhysPath)
		go e.WebCach.AddItem(rPath, targetPageBytes, 5*time.Minute)
	} else {
		// Get the bytes from cache
		targetPageBytes = e.WebCach.GetItem(rPath)
	}
	return targetPageBytes
}

// FromRoot catches all requests.
func (e Environment) FromRoot(w http.ResponseWriter, r *http.Request) {

	rPath := strings.ToLower(r.URL.Path)

	// Here you can process dynamic content or filter requests for securty.
	// If you only want to serve static as-is html file, you can those
	// to WebUtil.ServeStaticFile() -- see main() func.
	// ...
	// Uncomment the Prinln() line -- to see the request path.
	// If you keep requesting images (clicking the next icon on the
	// gallery page), you'll notice that images are displaying accordingly,
	// and no requests are cuaght here. The reason is that those are handled
	// by WebUtil.ServeStaticFile directly; defined in the main func:
	//     http.HandleFunc("/assets/", page.WebUtil.ServeStaticFile)
	//fmt.Println(rPath)

	// The http Method will not have any effect here. But it's
	// good to warn the caller.
	if r.Method != "GET" {
		w.WriteHeader(http.StatusForbidden)
		msg := fmt.Sprintf("error %d - %s\n", http.StatusForbidden, http.StatusText(http.StatusForbidden))
		w.Write([]byte(msg))
		return
	}

	if rPath == "/" || rPath == "/null" {
		rPath = "/index"
	}

	if !strings.HasSuffix(rPath, ".html") {
		rPath = fmt.Sprintf("%s.html", rPath)

		if strings.Contains(rPath, "/.html") {
			rPath = strings.Replace(rPath, "/.html", "/index.html", -1)
		}
	}

	bMaster := e.getRawMaster()

	// This is the page that its content will go into the {{.MainContent}}
	// block in the master.html file before its written to httpResponse.
	var targetPageBytes []byte
	var err error

	if rPath == "/index.html" {
		// The home pag has an image background; it's good to cache it.
		targetPageBytes = e.getTargetPageFromCache(rPath)

	} else if strings.HasPrefix(rPath, "/gallery") {

		// To render the gallery page, we have to query the file names; and this could affect
		// pefromance, so, let's cache the whole page for 5 minutes.
		// Note that the caching occurs in two fold:
		//   --The webutil casches the img files for 5 mint (the length inidicated in this sample).
		//   --The body of the html is chached so that the list of <div> tags are not rebuilt ever time.
		targetPageBytes = e.getTargetPageFromCache(rPath)
		targetPageBytes = bytes.Replace(targetPageBytes, []byte("{{.Carousel}}"), e.getGalleryCarousel(), -1)

	} else {

		// Read the page from disk
		targetPagePhysPath, _ := e.getPhysicalPath(rPath)

		targetPageBytes, err = ioutil.ReadFile(targetPagePhysPath)
		if err != nil {
			fmt.Println("FromRoot()=>", err)
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	// This is all contents of the target page into the {{.MainContent}} block inside the master page.
	bFinal := bytes.Replace(bMaster, []byte("{{.MainContent}}"), targetPageBytes, 1)

	// Replace any other variables.
	bFinal = bytes.Replace(bFinal, []byte("{{.ThisYear}}"), []byte(fmt.Sprintf("%d", time.Now().Year())), -1)

	// Replace the comments last.
	//
	// Comments inside this block {{.COMMENT  }} will be removed before
	// the target page is rendered into a response.
	//
	// The begin and end string for a block could be anything
	// ...you could also use your own (e.g. begin: #### end: ##@).
	// Comments can appear in multiple places anywhere in your
	// html file; see /html/master.html for an example.
	bFinal = e.WebUtil.RemoveCommentsFromByBiteArry(bFinal, "{{.COMMENT ", "}}")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Server master + content file.
	w.Write(bFinal)
}

// getRawMaster gets the master html page as-is from disk.
func (e Environment) getRawMaster() []byte {
	fPath := fmt.Sprintf("%s/html/master.html", e.WebRootPath)
	b, _ := ioutil.ReadFile(fPath)

	return b
}
