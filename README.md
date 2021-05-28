# A Go website template

## Webstandard is a practical website template that lays down the groundwork for a full featured Go website. 

It uses the original Go net/http along with native javascript and bootstrap/jquery to demonstrate rendering dynamic content. It also suggests points on security and mitigation of the too many open files error.

An image gallery is included to further show how html can be rendered via Go. Image-authors are also displayed on the bread-crumb bar via an instance of MutationObserver (*images provided, license-free, by pexels.com*). 

To run the website:

- Start a shell window.
- go build -o webstandard && ./webstandard --portno 1265
- Navigate to http://localhost:1265
