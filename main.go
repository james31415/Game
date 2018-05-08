package main

import (
    "fmt"
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/audio"
    "github.com/dradtke/go-allegro/allegro/acodec"
    //"github.com/yuin/gopher-lua"
    "os"
)

const (
    // SOUND
    PIANO_WAV = "dat/snd/piano.wav"
    // DISPLAY
    WINX = 800
    WINY = 600
    // MAIN
    FPS  = 1
)

func loadSound(filename string) (instance *audio.SampleInstance, err error) {
    if _, err = os.Stat(filename); os.IsNotExist(err) {
        return nil, err
    }
    if sample, err := audio.LoadSample(filename); err != nil {
        return nil, err
    } else {
        instance = audio.CreateSampleInstance(sample)
    }
    if err = instance.AttachToMixer(audio.DefaultMixer()); err != nil { return nil, err }
    return instance, nil
}

// defer logging utility function while debugging
func logFunc(msg string, function func()) {
    println(msg)
    function()
}

func main() {
    allegro.Run(func() {
        var (
            display    *allegro.Display
            event       allegro.Event
            eventQueue *allegro.EventQueue
            timer      *allegro.Timer
            running    bool = true
            err        error
        )

        /* ////////////////
        ** Install Keyboard
        */ ////////////////
        println("INSTALL KEYBOARD")
        if err = allegro.InstallKeyboard(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL KEYBOARD", allegro.UninstallKeyboard)
        }
        
        /* ////////////////
        ** Install Joystick
        */ ////////////////
        println("INSTALL JOYSTICK")
        if err = allegro.InstallJoystick(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL JOYSTICK", allegro.UninstallJoystick)
        }
        
        /* ///////////////////////////
        ** Install Audio & Audio Codec
        */ ///////////////////////////
        println("INSTALL AUDIO")
        if err = audio.Install(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL AUDIO AND AUDIO CODEC", audio.Uninstall)
            audio.ReserveSamples(1)
            println("INSTALL AUDIO CODEC")
            if err = acodec.Install(); err != nil {
                panic(err)
            }
        }

        /* //////////////////
        ** Create Event Queue
        */ //////////////////
        println("CREATE EVENT QUEUE")
        if eventQueue, err = allegro.CreateEventQueue(); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY EVENT QUEUE", eventQueue.Destroy)
        }
        
        /* ////////////////
        ** Register Display
        */ ////////////////
        println("CREATE DISPLAY")
        allegro.SetNewDisplayFlags(allegro.WINDOWED)
        if display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY DISPLAY", display.Destroy)
            display.SetWindowTitle("Game")
            eventQueue.Register(display)
        }

        /* //////////////
        ** Register Timer
        */ //////////////
        println("CREATE TIMER")
        if timer, err = allegro.CreateTimer(1.0 / FPS); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY TIMER", timer.Destroy)
            eventQueue.Register(timer)
            timer.Start()
        }

        /* ////////////////////
        ** Set Solid Background
        */ ////////////////////
        allegro.ClearToColor(allegro.MapRGB(32,64,96))
        allegro.FlipDisplay()
        allegro.ClearToColor(allegro.MapRGB(96,64,32))

        /* //////////
        ** Sound Test
        */ //////////
        piano, err := loadSound(PIANO_WAV)
        if err != nil { panic(err) }
        piano.Play()

        /* /////////
        ** Main Loop
        */ /////////
        for running {
            switch e := eventQueue.WaitForEvent(&event); e.(type) { // TODO ////////////////////////////////////////////////////////////////////
            //case allegro.JoystickAxisEvent:          println("JoystickAxisEvent")          ;
            //case allegro.JoystickButtonDownEvent:    println("JoystickButtonDownEvent")    ; // Yes, This is a mess right now.
            //case allegro.JoystickButtonUpEvent:      println("JoystickButtonUpEvent")      ; // 
            //case allegro.JoystickConfigurationEvent: println("JoystickConfigurationEvent") ; // Trying to work out why the only events working
            //case allegro.KeyDownEvent:               println("KeyDownEvent")               ; // right now are TimerEvent and DisplayCloseEvent
            //case allegro.KeyUpEvent:                 println("KeyUpEvent")                 ; // 
            //case allegro.KeyCharEvent:               println("KeyCharEvent")               ; // Timer and Display are eventQueue.Register()ed
            //case allegro.MouseAxesEvent:             println("MouseAxesEvent")             ; // Still, other Display Events are not catching
            //case allegro.MouseButtonDownEvent:       println("MouseButtonDownEvent")       ; // 
            //case allegro.MouseButtonUpEvent:         println("MouseButtonUpEvent")         ; // Currently using Ratpoison as my win manager
            //case allegro.MouseWarpedEvent:           println("MouseWarpedEvent")           ; // Will need to test on another wm to make sure
            //case allegro.MouseEnterDisplayEvent:     println("MouseEnterDisplayEvent")     ;
            //case allegro.MouseLeaveDisplayEvent:     println("MouseLeaveDisplayEvent")     ;
            case allegro.TimerEvent:                 print("TimerEvent: ")                 ; println(timer.Count())
            //case allegro.DisplayExposeEvent:         println("DisplayExposeEvent")         ;
            //case allegro.DisplayResizeEvent:         println("DisplayResizeEvent")         ;
            case allegro.DisplayCloseEvent:          println("DisplayCloseEvent")          ; running = false
            //case allegro.DisplayLostEvent:           println("DisplayLostEvent")           ;
            //case allegro.DisplayFoundEvent:          println("DisplayFoundEvent")          ;
            //case allegro.DisplaySwitchOutEvent:      println("DisplaySwitchOutEvent")      ;
            //case allegro.DisplaySwitchInEvent:       println("DisplaySwitchInEvent")       ;
            //case allegro.DisplayOrientationEvent:    println("DisplayOrientationEvent")    ; // Defaulting to printing e.type until I can see
            //case allegro.UserEvent:                  println("UserEvent")                  ; // Proper events triggering
            default:                                 fmt.Printf("T\n", e)                  ;
            }
        }
    })
}

