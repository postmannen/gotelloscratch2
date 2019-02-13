package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/postmannen/tello"
)

type cmdData struct {
	command string
	data    string
}

var cmdFromScratch chan cmdData
var flightInfo chan tello.FlightData

var speed int = 100

const (
	scratchListenHost = "127.0.0.1:8001"
)

//fromScratch is a HandlerFunc that checks the URL path for commands from Scratch
// and puts the commands received on a channel to be sent to the Tello drone.
// The drone sends in the format /command/jobID/measure
func (d *droneData) fromScratch(w http.ResponseWriter, r *http.Request) {
	u := r.RequestURI
	uSplit := strings.Split(u, "/")

	if uSplit[1] != "poll" {
		if len(uSplit) > 3 {
			cmdFromScratch <- cmdData{command: uSplit[1], data: uSplit[3]}
			//fmt.Printf(" * len was greater than 2, case detected, uSplit = %#v\n", uSplit)
		} else {
			cmdFromScratch <- cmdData{command: uSplit[1], data: ""}
			//fmt.Printf(" * len was less than 2, case detected, uSplit = %#v\n", uSplit)
		}
	} else {
		fmt.Fprintf(w, "%v %v\n", "battery?", d.drone.GetFlightData().BatteryPercentage)
		fmt.Fprintf(w, "%v %v\n", "speed?", d.drone.GetFlightData().GroundSpeed)
		fmt.Fprintf(w, "%v %v\n", "time?", d.drone.GetFlightData().DroneFlyTimeLeft)
		fmt.Fprintf(w, "%v %v\n", "height?", d.drone.GetFlightData().Height)
		fmt.Fprintf(w, "%v %v\n", "wind?", d.drone.GetFlightData().WindState)
		fmt.Fprintf(w, "%v %v\n", "ssid?", d.drone.GetFlightData().SSID)

	}
}

func (d *droneData) handleCommand() {
	err := d.drone.ControlConnectDefault()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("*** established connection to the drone ***")

	for {
		//To be used for specifying degrees etc in case's below.
		var num2 int16

		cmd := <-cmdFromScratch
		//Convert to int16 length, if it fails its a string and not a number
		num1, err := strconv.ParseInt(cmd.data, 10, 16)
		if err == nil {
			num2 = int16(num1)
			fmt.Println("cmd.Data is a number", num2)
		}

		switch cmd.command {
		case "disconnect":
			d.drone.ControlDisconnect()
		case "takeoff":
			fmt.Println("takeoff")
			time.Sleep(250 * time.Millisecond)
			d.drone.TakeOff()
			time.Sleep(5 * time.Second)
			fmt.Println("takeoff timer 7 seconds ok")
		case "land":
			time.Sleep(1000 * time.Millisecond) //let the drone stand still before we land
			fmt.Println("land")
			d.drone.Land()
			//d.drone.ControlDisconnect()
		case "left":
			fmt.Println("left")
			d.drone.Left(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Left(0)
		case "right":
			fmt.Println("right")
			d.drone.Right(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Right(0)
		case "forward":
			fmt.Println("forward")
			d.drone.Forward(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Forward(0)
		case "back":
			fmt.Println("back")
			d.drone.Backward(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Backward(0)
		case "up":
			fmt.Println("up")
			d.drone.Up(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Up(0)
		case "down":
			fmt.Println("down")
			d.drone.Down(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.Down(0)
		case "setspeed":
			fmt.Println("down")
			speed = int(num2)
		case "hover":
			fmt.Println("hover")
			d.drone.Hover()
			if err != nil {
				log.Println("autoturn failed: ", err)
			}
		case "cw":
			fmt.Println("rotate clockwise")
			d.drone.TurnRight(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.TurnRight(0)
		case "ccw":
			fmt.Println("rotate clockwise")
			d.drone.TurnLeft(speed)
			time.Sleep(time.Millisecond * 100)
			d.drone.TurnLeft(0)
		case "flip":
			if cmd.data == "forward" {
				d.drone.Flip(tello.FlipForward)
			}
			if cmd.data == "backward" {
				d.drone.Flip(tello.FlipBackward)
			}
			if cmd.data == "left" {
				d.drone.Flip(tello.FlipLeft)
			}
			if cmd.data == "right" {
				d.drone.Flip(tello.FlipRight)
			}
			if cmd.data == "forwardleft" {
				d.drone.Flip(tello.FlipForwardLeft)
			}
			if cmd.data == "forwardright" {
				d.drone.Flip(tello.FlipForwardRight)
			}
			if cmd.data == "backwarleft" {
				d.drone.Flip(tello.FlipBackwardLeft)
			}
			if cmd.data == "backwardright" {
				d.drone.Flip(tello.FlipBackwardRight)
			}
			//case "poll":
			//	fmt.Println("found a poll request")
			//	select {
			//	case flightInfo <- drone.GetFlightData():
			//	default:
			//		fmt.Println("FlightData channel full")
			//	}
		}
	}
}

type droneData struct {
	drone *tello.Tello
}

func main() {
	d := &droneData{
		drone: new(tello.Tello),
	}
	cmdFromScratch = make(chan cmdData, 100)
	flightInfo = make(chan tello.FlightData)

	go d.handleCommand()

	http.HandleFunc("/", d.fromScratch)
	http.ListenAndServe(scratchListenHost, nil)

}
