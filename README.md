# A Go website template

## Webstandard is a practical website template that lays down the groundwork for a full featured Go website. 

It uses the original Go net/http along with native javascript and bootstrap/jquery to demonstrate rendering dynamic content. It also suggests points on security and mitigation of the *too many open files* error.

An image gallery is included to further show how html can be rendered via Go. Image-authors are also displayed on the bread-crumb bar via an instance of MutationObserver (*images provided, license-free, by pexels.com*). 

#### Handling Requests
There are only two handlers

- http.HandleFunc("/assets/", page.WebUtil.ServeStaticFile) ... all asset (.css, .js, img,...) files.
- http.HandleFunc("/", page.FromRoot) .................................... all other paths (.htm files). 

#### Master Page
On each request, *master.html* is first read; and then the content of the target page is
loaded into the {{.MainContent}} block, before http response is written. 

##### Search Engines
Although there is one master and many partial html files, search engines will still receive
the entire page as one whole page. You can add a &lt;meta name="robots" content="noindex"&gt;
(in the &lt;header&gt; block of the master.html) for the pages that are not be indexed by search engines.

#### Security
There are two places that you can apply security
- connState(conn net.Conn, connState http.ConnState)
  * You can close (drop) connections by IP address.
- http.HandleFunc("/", page.FromRoot)
  * You can send 403 (Forbidden) or 401 (Unauthorized) error back to callers based
    on URL path (.env, .config, appdata, _private,...), http method (POST, DELETE,...), or headers (bad characters, unwanted cookie values,... ). Note that requests handled by page.WebUtil.ServeStaticFile()
    will not be caught by page.FromRoot.
    
##### javascript code in the URL
script code in the URL path will not be interpreted and/or cause run-time error; the caller will only receive a 404 (Not Found) error.

#### Dealing with the *too many open files* error
The Deadline, ReadDeadline, and WriteDeadline are set to a default value in the the connState() func. This reduces the risk of getting the *too many open files* error. You can flood the website with requests to test this.

Run a few instances of the following, simultaneously, in shell windows.
- Linux
  * for i in {1..1000};do curl http:&#47;&#47;localhost:1265/gallery;echo ... $i ...;done;
- Windows
  * for /l %f in (1,1,1000) do curl http:&#47;&#47;localhost:1265/gallery

#### Caching
- Assets
  * All asset files are cached via page.WebUtil, although browsers do not make repeated calls for static files (i.e. css, js,...).
- Image Gallery
  * The image gallery page (gallery.html) is explicitly cached for a few minutes. Note that the images will still be handled directly by page.WebUtil and it is only the body of the rendered html (from gallery.html) that is cached, so that the operations of getting file names from disk, and redering html tags are not repeated on every request. 

#### Comments
All text inside the {{.COMMENT <text goes here> }} blocks (on the html files) are removed before http response is written, so the 
caller will not be able to see those by viewing the page-code.

#### Running the website

- Start a shell window.
- go build -o webstandard && ./webstandard --portno 1265
- Navigate to http:&#47;&#47;localhost:1265

See https://go-webstandard.githubsamples.com for a live demo.
