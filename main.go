package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/postmannen/tello"
)

type cmdData struct {
	command string
	data    string
}

var cmdFromScratch chan cmdData
var speed int = 100

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
	case "forward":
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case forward detected", "uSplit = ", uSplit)
	case "back":
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case back detected", "uSplit = ", uSplit)
	case "hover":
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case hover detected", "uSplit = ", uSplit)
	}

}

func handleCommand() {
	drone := new(tello.Tello)
	err := drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("*** established connection to the drone ***")

	for {
		cmd := <-cmdFromScratch
		//num1, _ := strconv.ParseInt(cmd.data, 10, 16)
		//num2 := int16(num1)

		switch cmd.command {
		case "takeoff":
			fmt.Println("takeoff")
			time.Sleep(250 * time.Millisecond)
			drone.TakeOff()
			time.Sleep(3 * time.Second)
			fmt.Println("takeoff timer 7 seconds ok")
		case "land":
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("land")
			drone.Land()
			drone.ControlDisconnect()
		case "left":
			fmt.Println("left")
			drone.Left(speed)
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
			time.Sleep(time.Millisecond * 100)
			drone.Left(0)
		case "right":
			fmt.Println("right")
			drone.Right(speed)
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
			time.Sleep(time.Millisecond * 100)
			drone.Right(0)
		case "forward":
			fmt.Println("forward")
			drone.Forward(speed)
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
			time.Sleep(time.Millisecond * 100)
			drone.Forward(0)
		case "back":
			fmt.Println("back")
			drone.Backward(speed)
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
			time.Sleep(time.Millisecond * 100)
			drone.Backward(0)
		case "hover":
			fmt.Println("hover")
			drone.Hover()
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
			time.Sleep(time.Millisecond * 250)
		}

	}
}

func main() {
	cmdFromScratch = make(chan cmdData, 100)

	go handleCommand()

	http.HandleFunc("/", fromScratch)
	http.ListenAndServe(scratchListenHost, nil)

}
