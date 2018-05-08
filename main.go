package main

import (
    "fmt" // for log
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/audio"
    "github.com/dradtke/go-allegro/allegro/acodec"
    //"github.com/yuin/gopher-lua"
    "os" // for loadSound check
)

const (
    LOGGING = false
    LOOPING = true
    // SOUND
    PIANO_WAV = "dat/snd/piano.wav"
    // DISPLAY
    WINX = 800
    WINY = 600
    // MAIN
    FPS  = 30
)

// dump Joystick information, if present
func checkJoysticks() {
    log("CONFIGURING JOYSTICKS")
    if allegro.ReconfigureJoysticks() {
        log(" = PASSED")
    } else {
        log(" = FAILED")
        return
    }
    // TODO // Too Deep
    if joys := allegro.NumJoysticks(); joys > 0 {
        log(" = JOYS:", joys)
        for joy := 0; joy < joys; joy++ {
            if joystick, err := allegro.GetJoystick(joy); err != nil {
                log(err)
            } else {
                log(" = JOYSTICK:", joy)
                log(" = NAME:", joystick.Name())
                if sticks := joystick.NumSticks(); sticks > 0 {
                    log(" = STICKS:", sticks)
                    for stick := 0; stick < sticks; stick++ {
                        if stickName, err := joystick.StickName(stick); err != nil {
                            log(err)
                        } else {
                            log(" = STICK:", stick)
                            log(" = NAME:", stickName)
                            if axes := joystick.NumAxes(stick); axes > 0 {
                                log(" = AXES:", axes)
                                for axis := 0; axis < axes; axis++ {
                                    if axisName, err := joystick.AxisName(stick, axis); err != nil {
                                        log(err)
                                    } else {
                                        log(" = AXIS:", axis)
                                        log(" = NAME:", axisName)
                                    }
                                }
                            }
                            if nbutt := joystick.NumButtons(); nbutt > 0 {
                                log(" = NBUTT:", nbutt)
                                for butt := 0; butt < nbutt; butt++ {
                                    if buttName, err := joystick.ButtonName(stick); err != nil {
                                        log(err)
                                    } else {
                                        log(" = BUTT:", butt)
                                        log(" = NAME:", buttName)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// temporary, until sprites are implemented and can load their own assets
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

func log(msg ...interface{}) { if LOGGING { fmt.Println(msg...) } }
func logFunc(msg string, function func()) { log(msg); function() }

func main() {
    log("STARTING..."); defer log("...FINISHED")
    allegro.Run(func() {
        var (
            display    *allegro.Display
            event       allegro.Event
            eventQueue *allegro.EventQueue
            joyState    allegro.JoystickState
            keyState    allegro.KeyboardState
            timer      *allegro.Timer
            running    bool = LOOPING
            err        error
        )

        /* /////////////////
        ** Configure Display
        */ /////////////////
        log("CONFIG DISPLAY")
        var flags allegro.DisplayFlags = 0

        // Manually set DisplayFlags for now
        // planning on setting with a config

           flags |= allegro.WINDOWED
        // flags |= allegro.FULLSCREEN
        // flags |= allegro.OPENGL
        // flags |= allegro.RESIZABLE
        // flags |= allegro.FRAMELESS
        // flags |= allegro.NOFRAME
           flags |= allegro.GENERATE_EXPOSE_EVENTS
        // flags |= allegro.OPENGL_3_0
        // flags |= allegro.OPENGL_FORWARD_COMPATIBLE
        // flags |= allegro.FULLSCREEN_WINDOW

        if flags == 0 {            log(" = DEFAULT"                  ) } else {
            if flags & 0x001 > 0 { log(" = WINDOWED"                 ) }
            if flags & 0x002 > 0 { log(" = FULLSCREEN"               ) }
            if flags & 0x004 > 0 { log(" = OPENGL"                   ) }
            if flags & 0x010 > 0 { log(" = RESIZABLE"                ) }
            if flags & 0x020 > 0 { log(" = FRAMELESS"                ) } // FRAMELESS == NOFRAME
            if flags & 0x020 > 0 { log(" = NOFRAME"                  ) } // NOFRAME == FRAMELESS
            if flags & 0x040 > 0 { log(" = GENERATE_EXPOSE_EVENTS"   ) }
            if flags & 0x080 > 0 { log(" = OPENGL_3_0"               ) }
            if flags & 0x100 > 0 { log(" = OPENGL_FORWARD_COMPATIBLE") }
            if flags & 0x200 > 0 { log(" = FULLSCREEN_WINDOW"        ) }
        }

        /* ///////////////////////////
        ** Install Audio & Audio Codec
        */ ///////////////////////////
        log("INSTALL AUDIO")
        if err = audio.Install(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL AUDIO AND AUDIO CODEC", audio.Uninstall)
            audio.ReserveSamples(1)
            log("INSTALL AUDIO CODEC")
            if err = acodec.Install(); err != nil {
                panic(err)
            }
        }

        /* //////////////////
        ** Create Event Queue
        */ //////////////////
        log("CREATE EVENT QUEUE")
        if eventQueue, err = allegro.CreateEventQueue(); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY EVENT QUEUE", eventQueue.Destroy)
        }
        
        /* ////////////////
        ** Register Display
        */ ////////////////
        log("CREATE DISPLAY")
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
        log("INSTALL JOYSTICK")
        if err = allegro.InstallJoystick(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL JOYSTICK", allegro.UninstallJoystick)
            eventQueue.RegisterEventSource(allegro.JoystickEventSource())
            checkJoysticks()
        }
        
        /* /////////////////
        ** Register Keyboard
        */ /////////////////
        log("INSTALL KEYBOARD")
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
        log("INSTALL MOUSE")
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
        log("CREATE TIMER")
        if timer, err = allegro.CreateTimer(1.0 / FPS); err != nil {
            panic(err)
        } else {
            defer logFunc("DESTROY TIMER", timer.Destroy)
            eventQueue.Register(timer)
            timer.Start()
        }

        /* //////////
        ** Sound Test
        */ //////////
        piano, err := loadSound(PIANO_WAV)
        if err != nil { panic(err) }
        piano.Play()

        /* ///////// // /////////////////////
        ** Main Loop // Damn Huge Switch Loop
        */ ///////// // /////////////////////
        log("LOOPING...")
        quit := func() { log("QUITING..."); running = false }
        for running { switch e := eventQueue.WaitForEvent(&event); e.(type) {

            // Display Events
            case allegro.DisplayCloseEvent:          log("DisplayCloseEvent"         ); quit()
            case allegro.DisplayExposeEvent:         log("DisplayExposeEvent"        ); allegro.FlipDisplay()
            case allegro.DisplayFoundEvent:          log("DisplayFoundEvent"         );
            case allegro.DisplayLostEvent:           log("DisplayLostEvent"          );
            case allegro.DisplayOrientationEvent:    log("DisplayOrientationEvent"   );
            case allegro.DisplayResizeEvent:         log("DisplayResizeEvent"        ); display.AcknowledgeResize()
            case allegro.DisplaySwitchInEvent:       log("DisplaySwitchInEvent"      );
            case allegro.DisplaySwitchOutEvent:      log("DisplaySwitchOutEvent"     );

            // Joystick Events
            case allegro.JoystickAxisEvent:          log("JoystickAxisEvent"         ); joyState.Get()
            case allegro.JoystickButtonDownEvent:    log("JoystickButtonDownEvent"   ); joyState.Get()
            case allegro.JoystickButtonUpEvent:      log("JoystickButtonUpEvent"     ); joyState.Get()
            case allegro.JoystickConfigurationEvent: log("JoystickConfigurationEvent"); checkJoysticks()

            // Keyboard Events
            case allegro.KeyCharEvent:               log("KeyCharEvent"              ); keyState.Get()
                switch {
                case keyState.IsDown(allegro.KEY_ESCAPE): quit()
                case keyState.IsDown(allegro.KEY_Q     ): quit()
                case keyState.IsDown(allegro.KEY_UP    ): log(" = UP"   );
                case keyState.IsDown(allegro.KEY_DOWN  ): log(" = DOWN" );
                case keyState.IsDown(allegro.KEY_LEFT  ): log(" = LEFT" );
                case keyState.IsDown(allegro.KEY_RIGHT ): log(" = RIGHT");
                }
            case allegro.KeyDownEvent:               log("KeyDownEvent"              ); allegro.FlipDisplay(); allegro.FlipDisplay()
            case allegro.KeyUpEvent:                 log("KeyUpEvent"                ); // Flip-flop background color for testing

            // Mouse Events
            case allegro.MouseAxesEvent:             log("MouseAxesEvent"            );
            case allegro.MouseButtonDownEvent:       log("MouseButtonDownEvent"      );
            case allegro.MouseButtonUpEvent:         log("MouseButtonUpEvent"        );
            case allegro.MouseEnterDisplayEvent:     log("MouseEnterDisplayEvent"    ); timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     log("MouseLeaveDisplayEvent"    ); timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           log("MouseWarpedEvent"          );

            // Timer Events
            case allegro.TimerEvent:                 log("TimerEvent:", timer.Count())
                // Cycle background color for testing
                func(color byte){
                    allegro.ClearToColor(allegro.MapRGB(1*color&0xFF,2*color&0xFF,3*color&0xFF))
                    allegro.FlipDisplay()
                }(byte(timer.Count()))

            // User Events
            case allegro.UserEvent:                  log("UserEvent"                 );

            default:                                 log("UnknownEvent:", e          );
            }
        }
    }) // allegro.Run(func(){...})
}

