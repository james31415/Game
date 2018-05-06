package main

import (
    AL "github.com/dradtke/go-allegro/allegro"
)

const (
    WINX = 640
    WINY = 480
    FPS = 60
)

func main() {
    var (
        display    *AL.Display
        eventQueue *AL.EventQueue
        running    bool = true
        err        error
    )

    AL.Run(func() {
        // Create Window
        AL.SetNewDisplayFlags(AL.WINDOWED)
        if display, err = AL.CreateDisplay(WINX,WINY); err == nil {
            defer display.Destroy()
            display.SetWindowTitle("Hello World")
        } else {
            panic(err)
        }

        // Create Event Queue
        if eventQueue, err = AL.CreateEventQueue(); err == nil {
            defer eventQueue.Destroy()
            eventQueue.Register(display)
        } else {
            panic(err)
        }
        
        AL.ClearToColor(AL.MapRGB(32,64,96))
        AL.FlipDisplay()

        timeout := float64(1) / float64(FPS)

        // Main Loop
        var event AL.Event
        for running {
            if e, found := eventQueue.WaitForEventUntil(AL.NewTimeout(timeout), &event); found {
                switch e.(type) {
                case AL.DisplayCloseEvent:
                    running = false
                    break
                }
            }
        }
    })
}

