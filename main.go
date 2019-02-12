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

	if len(uSplit) > 3 {
		fmt.Println("------len was greater than 2")
		cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
		fmt.Println(" * case detected", "uSplit = ", uSplit)
	} else {
		fmt.Println("------len was less than 2")
		cmdFromScratch <- cmdData{command: uSplit[1], data: ""}
		fmt.Println(" * case detected", "uSplit = ", uSplit)
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
			time.Sleep(1000 * time.Millisecond) //let the drone stand still before we land
			fmt.Println("land")
			drone.Land()
			drone.ControlDisconnect()
		case "left":
			fmt.Println("left")
			drone.Left(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Left(0)
		case "right":
			fmt.Println("right")
			drone.Right(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Right(0)
		case "forward":
			fmt.Println("forward")
			drone.Forward(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Forward(0)
		case "back":
			fmt.Println("back")
			drone.Backward(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Backward(0)
		case "up":
			fmt.Println("up")
			drone.Up(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Up(0)
		case "down":
			fmt.Println("down")
			drone.Down(speed)
			time.Sleep(time.Millisecond * 100)
			drone.Down(0)
		case "hover":
			fmt.Println("hover")
			drone.Hover()
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
		case "cw":
			fmt.Println("rotate clockwise")
			drone.TurnRight(speed)
			time.Sleep(time.Millisecond * 100)
			drone.TurnRight(0)
		case "ccw":
			fmt.Println("rotate clockwise")
			drone.TurnLeft(speed)
			time.Sleep(time.Millisecond * 100)
			drone.TurnLeft(0)
		case "flip":
			if cmd.data == "forward" {
				drone.Flip(tello.FlipForward)
			}
			if cmd.data == "backward" {
				drone.Flip(tello.FlipBackward)
			}
			if cmd.data == "left" {
				drone.Flip(tello.FlipLeft)
			}
			if cmd.data == "right" {
				drone.Flip(tello.FlipRight)
			}
			if cmd.data == "forwardleft" {
				drone.Flip(tello.FlipForwardLeft)
			}
			if cmd.data == "forwardright" {
				drone.Flip(tello.FlipForwardRight)
			}
			if cmd.data == "backwarleft" {
				drone.Flip(tello.FlipBackwardLeft)
			}
			if cmd.data == "backwardright" {
				drone.Flip(tello.FlipBackwardRight)
			}
		}
	}
}

func main() {
	cmdFromScratch = make(chan cmdData, 100)

	go handleCommand()

	http.HandleFunc("/", fromScratch)
	http.ListenAndServe(scratchListenHost, nil)

}
