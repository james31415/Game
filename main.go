package main

import (
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
        var flags allegro.DisplayFlags = 0

        // flags |= allegro.WINDOWED
        // flags |= allegro.FULLSCREEN
        // flags |= allegro.OPENGL
        // flags |= allegro.RESIZABLE
        // flags |= allegro.FRAMELESS
        // flags |= allegro.NOFRAME
        // flags |= allegro.GENERATE_EXPOSE_EVENTS
        // flags |= allegro.OPENGL_3_0
        // flags |= allegro.OPENGL_FORWARD_COMPATIBLE
        // flags |= allegro.FULLSCREEN_WINDOW

        if flags == 0 {            println(" = DEFAULT"                  ) } else {
            if flags & 0x001 > 0 { println(" = WINDOWED"                 ) }
            if flags & 0x002 > 0 { println(" = FULLSCREEN"               ) }
            if flags & 0x004 > 0 { println(" = OPENGL"                   ) }
            if flags & 0x010 > 0 { println(" = RESIZABLE"                ) }
            if flags & 0x020 > 0 { println(" = FRAMELESS"                ) }
            if flags & 0x020 > 0 { println(" = NOFRAME"                  ) }
            if flags & 0x040 > 0 { println(" = GENERATE_EXPOSE_EVENTS"   ) }
            if flags & 0x080 > 0 { println(" = OPENGL_3_0"               ) }
            if flags & 0x100 > 0 { println(" = OPENGL_FORWARD_COMPATIBLE") }
            if flags & 0x200 > 0 { println(" = FULLSCREEN_WINDOW"        ) }
        }

        allegro.SetNewDisplayFlags(flags)
        if display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY DISPLAY", display.Destroy)
            display.SetWindowTitle("Game")
            eventQueue.Register(display)
            eventQueue.RegisterEventSource(display.EventSource()) // Redundant?
        }

        /* /////////////////
        ** Register Joystick
        */ /////////////////
        println("INSTALL JOYSTICK")
        if err = allegro.InstallJoystick(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL JOYSTICK", allegro.UninstallJoystick)
            eventQueue.RegisterEventSource(allegro.JoystickEventSource())
        }
        
        /* /////////////////
        ** Register Keyboard
        */ /////////////////
        println("INSTALL KEYBOARD")
        if err = allegro.InstallKeyboard(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL KEYBOARD", allegro.UninstallKeyboard)
            if keyboard, err := allegro.KeyboardEventSource(); err != nil {
                panic(err)
            } else {
                eventQueue.RegisterEventSource(keyboard)
            }
        }
        
        /* //////////////
        ** Register Mouse
        */ //////////////
        println("INSTALL MOUSE")
        if err = allegro.InstallMouse(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL MOUSE", allegro.UninstallMouse)
            if mouse, err := allegro.MouseEventSource(); err != nil {
                panic(err)
            } else {
                eventQueue.RegisterEventSource(mouse)
            }
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
            switch e := eventQueue.WaitForEvent(&event); e.(type) {
            case allegro.DisplayCloseEvent:          println("DisplayCloseEvent"         ); running = false
            case allegro.DisplayExposeEvent:         println("DisplayExposeEvent"        );
            case allegro.DisplayFoundEvent:          println("DisplayFoundEvent"         );
            case allegro.DisplayLostEvent:           println("DisplayLostEvent"          );
            case allegro.DisplayOrientationEvent:    println("DisplayOrientationEvent"   );
            case allegro.DisplayResizeEvent:         println("DisplayResizeEvent"        );
            case allegro.DisplaySwitchInEvent:       println("DisplaySwitchInEvent"      );
            case allegro.DisplaySwitchOutEvent:      println("DisplaySwitchOutEvent"     );
            case allegro.JoystickAxisEvent:          println("JoystickAxisEvent"         );
            case allegro.JoystickButtonDownEvent:    println("JoystickButtonDownEvent"   );
            case allegro.JoystickButtonUpEvent:      println("JoystickButtonUpEvent"     );
            case allegro.JoystickConfigurationEvent: println("JoystickConfigurationEvent");
            case allegro.KeyCharEvent:               println("KeyCharEvent"              );
            case allegro.KeyDownEvent:               println("KeyDownEvent"              );
            case allegro.KeyUpEvent:                 println("KeyUpEvent"                );
            case allegro.MouseAxesEvent:             println("MouseAxesEvent"            );
            case allegro.MouseButtonDownEvent:       println("MouseButtonDownEvent"      );
            case allegro.MouseButtonUpEvent:         println("MouseButtonUpEvent"        );
            case allegro.MouseEnterDisplayEvent:     println("MouseEnterDisplayEvent"    );
            case allegro.MouseLeaveDisplayEvent:     println("MouseLeaveDisplayEvent"    );
            case allegro.MouseWarpedEvent:           println("MouseWarpedEvent"          );
            case allegro.TimerEvent:                 print(  "TimerEvent: "              ); println(timer.Count())
            case allegro.UserEvent:                  println("UserEvent"                 );
            default:                                 print(  "UnknownEvent: "            ); println(e)
            }
        }
    })
}

