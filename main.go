package main

import (
	"fmt"
	"net/http"
	"strings"
)

type cmdData struct {
	command string
	data    string
}

var cmdFromScratch chan cmdData

const (
	scratchListenHost = "127.0.0.1:8001"
)

//fromScratch is a HandlerFunc that checks the URL path for commands from Scratch
// and puts the commands received on a channel to be sent to the Tello drone.
// The drone sends in the format /command/jobID/measure
func fromScratch(w http.ResponseWriter, r *http.Request) {
	u := r.RequestURI
	uSplit := strings.Split(u, "/")

	switch uSplit[1] {
	case "takeoff":
		cmdFromScratch <- cmdData{command: uSplit[1], data: ""}
		fmt.Printf(" * case takeoff detected, uSplit = %#v\n", uSplit)
	case "land":
		cmdFromScratch <- cmdData{command: uSplit[1], data: ""}
		fmt.Printf(" * case land detected uSplit = %#v\n", uSplit)
	case "left":
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case left detected", "uSplit = ", uSplit)
	case "right":
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case right detected", "uSplit = ", uSplit)
	}

}

func handleCommand() {
	for {
		cmd := <-cmdFromScratch

		switch cmd.command {
		case "takeoff":
		case "land":

		}

	}
}

func main() {
	cmdFromScratch = make(chan cmdData, 100)

	http.HandleFunc("/", fromScratch)
	http.ListenAndServe(scratchListenHost, nil)

	//	drone := new(tello.Tello)
	//	err := drone.ControlConnectDefault()
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	drone.TakeOff()
	//	time.Sleep(5 * time.Second)
	//
	//	done, err := drone.AutoTurnByDeg(180)
	//	if err != nil {
	//		log.Println("autoturn failed: ", err)
	//	}
	//	<-done
	//
	//	done, err = drone.AutoTurnByDeg(-180)
	//	if err != nil {
	//		log.Println("autoturn failed: ", err)
	//	}
	//	<-done
	//
	//	drone.
	//		drone.Land()
	//	drone.ControlDisconnect()
}
