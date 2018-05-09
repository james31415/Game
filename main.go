package main

import (
    "fmt" // for log
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/image"
    "github.com/dradtke/go-allegro/allegro/audio"
    "github.com/dradtke/go-allegro/allegro/acodec"
    //"github.com/yuin/gopher-lua"
    "os" // for loadSound check
)

const (
    // DEBUG
    LOGGING = false
    LOOPING = true
    // JOYSTICK
    STICK_THRESHOLD = 0.5
    // SOUND
    PIANO_WAV = "dat/snd/piano.wav"
    // DISPLAY
    TEST_ICON = "dat/img/icecat.png"
    WINX = 800
    WINY = 600
    // MAIN
    FPS  = 60
)

type sprite struct {
    Bitmap         *allegro.Bitmap
    DrawFlags       allegro.DrawFlags
    BitmapFilename  string
    OffsetX         float32
    OffsetY         float32
    ScaleX          float32
    ScaleY          float32
    Height          float32
    Width           float32
    Angle           float32
    X               float32
    Y               float32
}

func (s *sprite) Draw() { s.Bitmap.DrawScaledRotated(s.OffsetX,s.OffsetY,s.X,s.Y,s.ScaleX,s.ScaleY,s.Angle,s.DrawFlags) }
//func (s *sprite) Draw() { s.Bitmap.Draw(s.X+s.OffsetX,s.Y+s.OffsetY,s.DrawFlags) }

func (s *sprite) Load() (err error) {
    if bitmap, err := allegro.LoadBitmap(s.BitmapFilename); err != nil {
        return err
    } else {
        s.Bitmap  = bitmap
        s.Width   = float32(bitmap.Width())
        s.Height  = float32(bitmap.Height())
        s.OffsetY = s.Height/2
        s.OffsetX = s.Width/2
        s.ScaleX  = 1.0
        s.ScaleY  = 1.0
    }
    return nil
}

func (s *sprite) Unload() {
    s.Bitmap.Destroy()
    s.OffsetX = 0
    s.OffsetY = 0
    s.Height  = 0
    s.Width   = 0
}

func configureJoysticks() (joyState *allegro.JoystickState) {
    log("CONFIGURING JOYSTICKS")
    if allegro.ReconfigureJoysticks() { log(" = RECONFIGURED") }
    if joys := allegro.NumJoysticks(); joys > 0 {
        for joy := 0; joy < joys; joy++ {
            if joystick, err := allegro.GetJoystick(joy); err != nil {
                log(err)
            } else {
                log(" = JOYSTICK:", joy, "=", joystick.Name())
                joyState = joystick.State()
            }
        }
    }
    return
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
            joyState   *allegro.JoystickState
            keyState    allegro.KeyboardState
            timer      *allegro.Timer
            running     bool = LOOPING
            err         error
        )

        /* /////////////////
        ** Configure Display
        */ /////////////////
        log("CONFIG DISPLAY")
        var flags allegro.DisplayFlags = 0

        // Manually set DisplayFlags for now
        // planning on setting with a config

        // flags |= allegro.WINDOWED
        // flags |= allegro.FULLSCREEN
        // flags |= allegro.OPENGL
        // flags |= allegro.RESIZABLE
        // flags |= allegro.FRAMELESS
        // flags |= allegro.NOFRAME
        // flags |= allegro.GENERATE_EXPOSE_EVENTS
        // flags |= allegro.OPENGL_3_0
        // flags |= allegro.OPENGL_FORWARD_COMPATIBLE
           flags |= allegro.FULLSCREEN_WINDOW

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

        /* ///////////////////
        ** Install Image Addon
        */ ///////////////////
        log("INSTALL IMAGE")
        if err = image.Install(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL IMAGE", image.Uninstall)
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
            //eventQueue.RegisterEventSource(allegro.JoystickEventSource())
            joyState = configureJoysticks()
        }
        
        /* /////////////////
        ** Register Keyboard
        */ /////////////////
        log("INSTALL KEYBOARD")
        if err = allegro.InstallKeyboard(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL KEYBOARD", allegro.UninstallKeyboard)
            //if keyboard, err := allegro.KeyboardEventSource(); err != nil {
            //    panic(err)
            //} else {
            //    eventQueue.RegisterEventSource(keyboard)
            //}
        }
        
        /* //////////////
        ** Register Mouse
        */ //////////////
        log("INSTALL MOUSE")
        if err = allegro.InstallMouse(); err != nil {
            panic(err)
        } else {
            defer logFunc("UNINSTALL MOUSE", allegro.UninstallMouse)
            //if mouse, err := allegro.MouseEventSource(); err != nil {
            //    panic(err)
            //} else {
            //    eventQueue.RegisterEventSource(mouse)
            //}
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


        /* ///////////
        ** Sprite Test
        */ ///////////
        var icon sprite
        icon.BitmapFilename = TEST_ICON
        if err = icon.Load(); err != nil { log(err) }
        icon.X = float32(display.Width()/2)
        icon.Y = float32(display.Height()/2)

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
            case allegro.JoystickConfigurationEvent: log("JoystickConfigurationEvent"); joyState = configureJoysticks()

            // Mouse Events
            case allegro.MouseAxesEvent:             log("MouseAxesEvent"            );
            case allegro.MouseButtonDownEvent:       log("MouseButtonDownEvent"      );
            case allegro.MouseButtonUpEvent:         log("MouseButtonUpEvent"        );
            case allegro.MouseEnterDisplayEvent:     log("MouseEnterDisplayEvent"    ); timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     log("MouseLeaveDisplayEvent"    ); timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           log("MouseWarpedEvent"          );

            // Timer Events
            case allegro.TimerEvent: {            // log("TimerEvent:", timer.Count())

                // Keyboard State
                keyState.Get()
                if keyState.IsDown(allegro.KEY_ESCAPE) { quit() }
                if keyState.IsDown(allegro.KEY_Q     ) { quit() }
                if keyState.IsDown(allegro.KEY_UP    ) { log(" = UP"   ) }
                if keyState.IsDown(allegro.KEY_DOWN  ) { log(" = DOWN" ) }
                if keyState.IsDown(allegro.KEY_LEFT  ) { log(" = LEFT" ) }
                if keyState.IsDown(allegro.KEY_RIGHT ) { log(" = RIGHT") }
             
                // Joystick State
                if joyState != nil {
                    joyState.Get()

                    Axis_LX := joyState.Stick[0].Axis[0]
                    Axis_LY := joyState.Stick[0].Axis[1]
                    Axis_TL := joyState.Stick[1].Axis[0]
                    Axis_RX := joyState.Stick[1].Axis[1]
                    Axis_RY := joyState.Stick[2].Axis[0]
                    Axis_TR := joyState.Stick[2].Axis[1]
                    Axis_DX := joyState.Stick[3].Axis[0]
                    Axis_DY := joyState.Stick[3].Axis[1]

                    Button_A     := joyState.Button[0x0]
                    Button_B     := joyState.Button[0x1]
                    Button_X     := joyState.Button[0x2]
                    Button_Y     := joyState.Button[0x3]
                    Button_LB    := joyState.Button[0x4]
                    Button_RB    := joyState.Button[0x5]
                    Button_BACK  := joyState.Button[0x6]
                    Button_START := joyState.Button[0x7]
                    Button_XBOX  := joyState.Button[0x8]
                    Button_LS    := joyState.Button[0x9]
                    Button_RS    := joyState.Button[0xA]

                    if Axis_LY < -STICK_THRESHOLD { log(" = L_STICK_UP"   ) }
                    if Axis_LY >  STICK_THRESHOLD { log(" = L_STICK_DOWN" ) }
                    if Axis_LX < -STICK_THRESHOLD { log(" = L_STICK_LEFT" ) }
                    if Axis_LX >  STICK_THRESHOLD { log(" = L_STICK_RIGHT") }

                    if Axis_RY < -STICK_THRESHOLD { log(" = R_STICK_UP"   ) }
                    if Axis_RY >  STICK_THRESHOLD { log(" = R_STICK_DOWN" ) }
                    if Axis_RX < -STICK_THRESHOLD { log(" = R_STICK_LEFT" ) }
                    if Axis_RX >  STICK_THRESHOLD { log(" = R_STICK_RIGHT") }

                    if Axis_DY < -STICK_THRESHOLD { log(" = DIR_PAD_UP"   ) }
                    if Axis_DY >  STICK_THRESHOLD { log(" = DIR_PAD_DOWN" ) }
                    if Axis_DX < -STICK_THRESHOLD { log(" = DIR_PAD_LEFT" ) }
                    if Axis_DX >  STICK_THRESHOLD { log(" = DIR_PAD_RIGHT") }

                    if Axis_TL >  STICK_THRESHOLD { log(" = TRIGGER_LEFT" ) }
                    if Axis_TR >  STICK_THRESHOLD { log(" = TRIGGER_RIGHT") }

                    if Button_A     > 0 { log(" = JOY_A"    ) }
                    if Button_B     > 0 { log(" = JOY_B"    ) }
                    if Button_X     > 0 { log(" = JOY_X"    ) }
                    if Button_Y     > 0 { log(" = JOY_Y"    ) }
                    if Button_LB    > 0 { log(" = JOY_LB"   ) }
                    if Button_RB    > 0 { log(" = JOY_RB"   ) }
                    if Button_BACK  > 0 { log(" = JOY_BACK" ) }
                    if Button_START > 0 { log(" = JOY_START") }
                    if Button_XBOX  > 0 { log(" = JOY_XBOX" ) }
                    if Button_LS    > 0 { log(" = JOY_LS"   ) }
                    if Button_RS    > 0 { log(" = JOY_RS"   ) }

                    icon.X      += Axis_LX * 8
                    icon.Y      += Axis_LY * 8
                    icon.ScaleX += Axis_RX
                    icon.ScaleY -= Axis_RY
                    icon.Angle  -= Axis_TL / 8
                    icon.Angle  += Axis_TR / 8
                    icon.Draw()

                    // Cycle background color
                    // c :=      -Axis_RY * 63 + 63
                    // r := byte( Axis_LX * c + c) & 0xFF
                    // g := byte( Axis_RX * c + c) & 0xFF
                    // b := byte(-Axis_LY * c + c) & 0xFF
                    // allegro.ClearToColor(allegro.MapRGB(r,g,b))

                } // Joystick State

                // Display State
                allegro.FlipDisplay()

            } // Timer Events

            default:
            }
        }
    }) // allegro.Run(func(){...})
}

