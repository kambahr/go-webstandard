package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// logInstanceError logs an instance error to the global
// error var.
func logInstanceError(id int, errTxt string) {

	arryInx := -1
	for j := 0; j < len(allErrLog); j++ {
		if allErrLog[j].ID == id {
			arryInx = j
			break
		}
	}

	// Create the error data
	var lg errLog
	lg.ID = id
	var ev errEvent
	tm := strings.TrimRight(strings.Split(time.Now().String(), "+")[0], " ")
	ev.TimeStamp = tm
	ev.ErrText = errTxt
	lg.ErrEvent = append(lg.ErrEvent, ev)

	// Find this isntance's place in the global array
	// and place it into its position.
	if arryInx < 0 {
		// Add all.
		allErrLog = append(allErrLog, lg)
	} else {
		// Append only the errors.
		allErrLog[arryInx].ErrEvent = append(allErrLog[arryInx].ErrEvent, ev)
	}
}

// callOneSet makes a sequence of http calls to a Go web server.
func callOneSet(id int, wg *sync.WaitGroup) {

	defer wg.Done()
	urlCnt := len(urlx)
	for i := 0; i < urlCnt; i++ {

		// Build the url.
		u := fmt.Sprintf("%s://%s:%d%s", proto, hostName, portNo, urlx[i])

		// Display on screen.
		s := fmt.Sprintf("[%d]", id)
		log.Println(s, u)

		// Make the call.
		req, _ := http.NewRequest("GET", u, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			// we dont' have to wait for this.
			go logInstanceError(id, err.Error())

			// display on screen
			fmt.Println(err)

			// a bit of delay and try again.
			time.Sleep(time.Second)
			continue
		}
		res.Body.Close()
	}
}

// structs to keep the errors organized.
type errEvent struct {
	ErrText   string
	TimeStamp string
}
type errLog struct {
	ID       int
	ErrEvent []errEvent
}

// List of pages in order to hit different urls.
var urlx = []string{
	"/music/theory/counterpoint/what-is-counterpoint",
	"/localhost:1267/gallery",
	"/page_1",
	"/terms-of-use",
	"/about",
	"/privacy-policy",
}

// Holds error for all instances.
var allErrLog []errLog

const proto = "http"
const hostName = "localhost"
const portNo = 1268

func main() {

	var totalInst int
	flag.IntVar(&totalInst, "t", 1000, "number of instances to test")
	flag.Parse()

	if totalInst < 5 {
		totalInst = 5
	}

	fmt.Println("Total isntances:", totalInst)
	time.Sleep(2 * time.Second)

	allErrLog = make([]errLog, totalInst*2)
	var wg sync.WaitGroup

	// Adjust the number of instances, bursts, and delay
	// according to your system's resource and your
	// performance expectations.

	i := 0
	for {

		wg.Add(1)

		go callOneSet(i, &wg)

		if i >= totalInst {
			break
		}

		// This is all needed to create enough gap for disk operations.
		// You may have to increase the delay, depending on your system.
		time.Sleep(time.Millisecond)

		// 2nd burst
		i++
		wg.Add(1)
		go callOneSet(i, &wg)

		if i >= totalInst {
			break
		}
	}
	wg.Wait()

	fmt.Println("done.")

	// Print the results
	errDirty := false
	for i := 0; i < len(allErrLog); i++ {

		if allErrLog[i].ID == 0 || len(allErrLog[i].ErrEvent) == 0 {
			continue
		}

		errDirty = true

		fmt.Println("--- ID", allErrLog[i].ID, "---")

		for k := 0; k < len(allErrLog[i].ErrEvent); k++ {
			fmt.Println("\t", allErrLog[i].ErrEvent[k].TimeStamp, allErrLog[i].ErrEvent[k].ErrText)
		}
	}

	if !errDirty {
		fmt.Println("no errors.")
	}
}
