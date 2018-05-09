package main

import (
    "fmt" // for log
    "io/ioutil"
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/image"
    "github.com/dradtke/go-allegro/allegro/audio"
    "github.com/dradtke/go-allegro/allegro/acodec"
    //"github.com/yuin/gopher-lua"
)

type DebugLevel int

const (
    GENERAL DebugLevel = 1 << iota
)

const ( // DEBUG
    LOGLVL  = 0
    LOGGING = true
    LOOPING = true
)

const (
    // JOYSTICK
    STICK_THRESHOLD = 0.5
    // SOUND
    PIANO_WAV = "dat/snd/piano.wav"
    // DISPLAY
    WINX = 800
    WINY = 600
    // MAIN
    FPS  = 60
    // SPRITES
    SPRITEDATA = "dat/sprites/"
    BACKGROUND = "background"
    PLAYER     = "player"
    // SOUND
    SPAWN = "spawn.wav"
)

type sprite struct {
    Name       string
    Folder     string
    Sound      map[string]*audio.Sample
    Bitmap    *allegro.Bitmap
    DrawFlags  allegro.DrawFlags
    OffsetX    float32
    OffsetY    float32
    ScaleX     float32
    ScaleY     float32
    Height     float32
    Width      float32
    Angle      float32
    X          float32
    Y          float32
    Draw       func()
    Spawn      func()
}

func (s *sprite) DrawNormal() {
    dx := s.X-s.OffsetX
    dy := s.Y-s.OffsetY
    df := s.DrawFlags
    s.Bitmap.Draw(dx,dy,df)
}

func (s *sprite) DrawScaled() {
    sx := s.X-s.OffsetX
    sy := s.Y-s.OffsetY
    sw := s.Width
    sh := s.Height
    dx := sx
    dy := sy
    dw := sw
    dh := sh
    df := s.DrawFlags
    s.Bitmap.DrawScaled(sx,sy,sw,sh,dx,dy,dw,dh,df)
}

func (s *sprite) DrawRotated() {
    cx := s.X-s.OffsetX
    cy := s.Y-s.OffsetY
    dx := cx
    dy := cy
    da := s.Angle
    df := s.DrawFlags
    s.Bitmap.DrawRotated(cx,cy,dx,dy,da,df)
}

//func (s *sprite) Draw() { s.Bitmap.DrawScaledRotated(s.OffsetX,s.OffsetY,s.X,s.Y,s.ScaleX,s.ScaleY,s.Angle,s.DrawFlags) }

func (s *sprite) Load(name string) {
    s.Name = name
    s.Folder = SPRITEDATA + name

    loadSound := func(sound string) {
        if sample, err := audio.LoadSample(s.Folder+"/snd/"+sound); err != nil {
            panic(err)
        } else {
            s.Sound[sound] = sample
        }
    }

    s.Sound = make(map[string]*audio.Sample)
    if sounds, err := ioutil.ReadDir(s.Folder+"/snd/"); err != nil {
        log(err)
    } else {
        for _, sound := range sounds {
            log(sound.Name())
            loadSound(sound.Name())
        }
    }

    if bitmap, err := allegro.LoadBitmap(s.Folder+"/bitmap"); err != nil {
        panic(err)
    } else {
        s.Bitmap  = bitmap
        s.Width   = float32(bitmap.Width())
        s.Height  = float32(bitmap.Height())
        s.OffsetY = s.Height/2
        s.OffsetX = s.Width/2
        s.ScaleX  = 10.0
        s.ScaleY  = 10.0
        s.Draw    = s.DrawNormal
        s.Spawn   = func(){
            s.Play(SPAWN)
            s.Draw()
        }
    }
}

func (s *sprite) Play(sound string) {
    instance := audio.CreateSampleInstance(s.Sound[sound])
    if err := instance.AttachToMixer(audio.DefaultMixer()); err != nil {
        panic(err)
    }
    if err := instance.Play(); err != nil {
        panic(err)
    }
}

func (s *sprite) Unload() {
    s.Bitmap.Destroy()
    s.OffsetX = 0
    s.OffsetY = 0
    s.Height  = 0
    s.Width   = 0
}

func (s *sprite) Center(display *allegro.Display) {
    s.Y = float32(display.Height()/2)
    s.X = float32(display.Width()/2)
}

func configureJoysticks() (joyState *allegro.JoystickState) {
    logLvl(GENERAL, "CONFIGURING JOYSTICKS")
    if allegro.ReconfigureJoysticks() { logLvl(GENERAL, " = RECONFIGURED") }
    if joys := allegro.NumJoysticks(); joys > 0 {
        for joy := 0; joy < joys; joy++ {
            if joystick, err := allegro.GetJoystick(joy); err != nil {
                panic(err)
            } else {
                logLvl(GENERAL, " = JOYSTICK:", joy, "=", joystick.Name())
                joyState = joystick.State()
            }
        }
    }
    return
}

/* ///////
** Logging
*/ ///////
func log(                   msg ...interface{}) { if LOGGING          { fmt.Println(msg...)          } }
func logLvl(lvl DebugLevel, msg ...interface{}) { if lvl & LOGLVL > 0 { log(msg)                     } }
func logFunc(                   msg string, function func())          { log(msg);         function() }
func logLvlFunc(lvl DebugLevel, msg string, function func())          { logLvl(lvl, msg); function() }

////////////////////////////////////////////////////////////////////////////////
func main() {
    logLvl(GENERAL, "STARTING..."); defer logLvl(GENERAL, "...FINISHED")
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
        logLvl(GENERAL, "CONFIG DISPLAY")
        var flags allegro.DisplayFlags = 0

        // Manually set DisplayFlags for now
        // planning on setting with a config

           flags |= allegro.WINDOWED
        // flags |= allegro.FULLSCREEN
        // flags |= allegro.OPENGL
        // flags |= allegro.RESIZABLE
        // flags |= allegro.FRAMELESS
        // flags |= allegro.NOFRAME
        // flags |= allegro.GENERATE_EXPOSE_EVENTS
        // flags |= allegro.OPENGL_3_0
        // flags |= allegro.OPENGL_FORWARD_COMPATIBLE
        // flags |= allegro.FULLSCREEN_WINDOW

        if flags == 0 {            logLvl(GENERAL, " = DEFAULT"                  ) } else {
            if flags & 0x001 > 0 { logLvl(GENERAL, " = WINDOWED"                 ) }
            if flags & 0x002 > 0 { logLvl(GENERAL, " = FULLSCREEN"               ) }
            if flags & 0x004 > 0 { logLvl(GENERAL, " = OPENGL"                   ) }
            if flags & 0x010 > 0 { logLvl(GENERAL, " = RESIZABLE"                ) }
            if flags & 0x020 > 0 { logLvl(GENERAL, " = FRAMELESS"                ) } // FRAMELESS == NOFRAME
            if flags & 0x020 > 0 { logLvl(GENERAL, " = NOFRAME"                  ) } // NOFRAME == FRAMELESS
            if flags & 0x040 > 0 { logLvl(GENERAL, " = GENERATE_EXPOSE_EVENTS"   ) }
            if flags & 0x080 > 0 { logLvl(GENERAL, " = OPENGL_3_0"               ) }
            if flags & 0x100 > 0 { logLvl(GENERAL, " = OPENGL_FORWARD_COMPATIBLE") }
            if flags & 0x200 > 0 { logLvl(GENERAL, " = FULLSCREEN_WINDOW"        ) }
        }

        /* ///////////////////
        ** Install Image Addon
        */ ///////////////////
        logLvl(GENERAL, "INSTALL IMAGE")
        if err = image.Install(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "UNINSTALL IMAGE", image.Uninstall)
        }

        /* ///////////////////////////
        ** Install Audio & Audio Codec
        */ ///////////////////////////
        logLvl(GENERAL, "INSTALL AUDIO")
        if err = audio.Install(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "UNINSTALL AUDIO AND AUDIO CODEC", audio.Uninstall)
            audio.ReserveSamples(8)
            logLvl(GENERAL, "INSTALL AUDIO CODEC")
            if err = acodec.Install(); err != nil {
                panic(err)
            }
        }

        /* //////////////////
        ** Create Event Queue
        */ //////////////////
        logLvl(GENERAL, "CREATE EVENT QUEUE")
        if eventQueue, err = allegro.CreateEventQueue(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "DESTROY EVENT QUEUE", eventQueue.Destroy)
        }
        
        /* ////////////////
        ** Register Display
        */ ////////////////
        logLvl(GENERAL, "CREATE DISPLAY")
        allegro.SetNewDisplayFlags(flags)
        if display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "DESTROY DISPLAY", display.Destroy)
            display.SetWindowTitle("Game")
            eventQueue.Register(display)
            eventQueue.RegisterEventSource(display.EventSource()) // Redundant?
        }

        /* /////////////////
        ** Register Joystick
        */ /////////////////
        logLvl(GENERAL, "INSTALL JOYSTICK")
        if err = allegro.InstallJoystick(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "UNINSTALL JOYSTICK", allegro.UninstallJoystick)
            //eventQueue.RegisterEventSource(allegro.JoystickEventSource())
            joyState = configureJoysticks()
        }
        
        /* /////////////////
        ** Register Keyboard
        */ /////////////////
        logLvl(GENERAL, "INSTALL KEYBOARD")
        if err = allegro.InstallKeyboard(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "UNINSTALL KEYBOARD", allegro.UninstallKeyboard)
            //if keyboard, err := allegro.KeyboardEventSource(); err != nil {
            //    panic(err)
            //} else {
            //    eventQueue.RegisterEventSource(keyboard)
            //}
        }
        
        /* //////////////
        ** Register Mouse
        */ //////////////
        logLvl(GENERAL, "INSTALL MOUSE")
        if err = allegro.InstallMouse(); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "UNINSTALL MOUSE", allegro.UninstallMouse)
            //if mouse, err := allegro.MouseEventSource(); err != nil {
            //    panic(err)
            //} else {
            //    eventQueue.RegisterEventSource(mouse)
            //}
        }
        
        /* //////////////
        ** Register Timer
        */ //////////////
        logLvl(GENERAL, "CREATE TIMER")
        if timer, err = allegro.CreateTimer(1.0 / FPS); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "DESTROY TIMER", timer.Destroy)
            eventQueue.Register(timer)
            timer.Start()
        }

        /* //////////
        ** Sound Test
        */ //////////
        // piano, err := loadSound(PIANO_WAV)
        // if err != nil { panic(err) }
        // piano.Play()


        /* ///////////
        ** Sprite Test
        */ ///////////
        var background sprite
        background.Load(BACKGROUND)
        background.Center(display)
        background.Spawn()

        var player sprite
        player.Load(PLAYER)
        player.Center(display)
        player.Spawn()

        allegro.FlipDisplay()
        background.Draw()
        player.Draw()
        allegro.FlipDisplay()

        /* ///////// // /////////////////////
        ** Main Loop // Damn Huge Switch Loop
        */ ///////// // /////////////////////
        logLvl(GENERAL, "LOOPING...")
        quit := func() { logLvl(GENERAL, "QUITING..."); running = false }
        for running { switch e := eventQueue.WaitForEvent(&event); e.(type) {

            // Display Events
            case allegro.DisplayCloseEvent:          logLvl(GENERAL, "DisplayCloseEvent"         ); quit()
            case allegro.DisplayExposeEvent:         logLvl(GENERAL, "DisplayExposeEvent"        ); allegro.FlipDisplay()
            case allegro.DisplayFoundEvent:          logLvl(GENERAL, "DisplayFoundEvent"         );
            case allegro.DisplayLostEvent:           logLvl(GENERAL, "DisplayLostEvent"          );
            case allegro.DisplayOrientationEvent:    logLvl(GENERAL, "DisplayOrientationEvent"   );
            case allegro.DisplayResizeEvent:         logLvl(GENERAL, "DisplayResizeEvent"        ); display.AcknowledgeResize()
            case allegro.DisplaySwitchInEvent:       logLvl(GENERAL, "DisplaySwitchInEvent"      );
            case allegro.DisplaySwitchOutEvent:      logLvl(GENERAL, "DisplaySwitchOutEvent"     );

            // Joystick Events
            case allegro.JoystickConfigurationEvent: logLvl(GENERAL, "JoystickConfigurationEvent"); joyState = configureJoysticks()

            // Mouse Events
            case allegro.MouseAxesEvent:             logLvl(GENERAL, "MouseAxesEvent"            );
            case allegro.MouseButtonDownEvent:       logLvl(GENERAL, "MouseButtonDownEvent"      );
            case allegro.MouseButtonUpEvent:         logLvl(GENERAL, "MouseButtonUpEvent"        );
            case allegro.MouseEnterDisplayEvent:     logLvl(GENERAL, "MouseEnterDisplayEvent"    ); timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     logLvl(GENERAL, "MouseLeaveDisplayEvent"    ); timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           logLvl(GENERAL, "MouseWarpedEvent"          );

            // Timer Events
            case allegro.TimerEvent: {            // logLvl(GENERAL, "TimerEvent:", timer.Count())

                // Keyboard State
                keyState.Get()
                if keyState.IsDown(allegro.KEY_ESCAPE) { quit() }
                if keyState.IsDown(allegro.KEY_Q     ) { quit() }
                if keyState.IsDown(allegro.KEY_P     ) { player.Play(SPAWN) }
                if keyState.IsDown(allegro.KEY_UP    ) { logLvl(GENERAL, " = UP"   ) }
                if keyState.IsDown(allegro.KEY_DOWN  ) { logLvl(GENERAL, " = DOWN" ) }
                if keyState.IsDown(allegro.KEY_LEFT  ) { logLvl(GENERAL, " = LEFT" ) }
                if keyState.IsDown(allegro.KEY_RIGHT ) { logLvl(GENERAL, " = RIGHT") }
             
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

                    if Axis_LY < -STICK_THRESHOLD { logLvl(GENERAL, " = L_STICK_UP"   ) }
                    if Axis_LY >  STICK_THRESHOLD { logLvl(GENERAL, " = L_STICK_DOWN" ) }
                    if Axis_LX < -STICK_THRESHOLD { logLvl(GENERAL, " = L_STICK_LEFT" ) }
                    if Axis_LX >  STICK_THRESHOLD { logLvl(GENERAL, " = L_STICK_RIGHT") }

                    if Axis_RY < -STICK_THRESHOLD { logLvl(GENERAL, " = R_STICK_UP"   ) }
                    if Axis_RY >  STICK_THRESHOLD { logLvl(GENERAL, " = R_STICK_DOWN" ) }
                    if Axis_RX < -STICK_THRESHOLD { logLvl(GENERAL, " = R_STICK_LEFT" ) }
                    if Axis_RX >  STICK_THRESHOLD { logLvl(GENERAL, " = R_STICK_RIGHT") }

                    if Axis_DY < -STICK_THRESHOLD { logLvl(GENERAL, " = DIR_PAD_UP"   ) }
                    if Axis_DY >  STICK_THRESHOLD { logLvl(GENERAL, " = DIR_PAD_DOWN" ) }
                    if Axis_DX < -STICK_THRESHOLD { logLvl(GENERAL, " = DIR_PAD_LEFT" ) }
                    if Axis_DX >  STICK_THRESHOLD { logLvl(GENERAL, " = DIR_PAD_RIGHT") }

                    if Axis_TL >  STICK_THRESHOLD { logLvl(GENERAL, " = TRIGGER_LEFT" ) }
                    if Axis_TR >  STICK_THRESHOLD { logLvl(GENERAL, " = TRIGGER_RIGHT") }

                    if Button_A     > 0 { logLvl(GENERAL, " = JOY_A"    ) }
                    if Button_B     > 0 { logLvl(GENERAL, " = JOY_B"    ) }
                    if Button_X     > 0 { logLvl(GENERAL, " = JOY_X"    ) }
                    if Button_Y     > 0 { logLvl(GENERAL, " = JOY_Y"    ) }
                    if Button_LB    > 0 { logLvl(GENERAL, " = JOY_LB"   ) }
                    if Button_RB    > 0 { logLvl(GENERAL, " = JOY_RB"   ) }
                    if Button_BACK  > 0 { logLvl(GENERAL, " = JOY_BACK" ) }
                    if Button_START > 0 { logLvl(GENERAL, " = JOY_START") }
                    if Button_XBOX  > 0 { logLvl(GENERAL, " = JOY_XBOX" ) }
                    if Button_LS    > 0 { logLvl(GENERAL, " = JOY_LS"   ) }
                    if Button_RS    > 0 { logLvl(GENERAL, " = JOY_RS"   ) }

                    //player.ScaleX =  1 // Axis_RX
                    //player.ScaleY =  1 // Axis_RY
                    player.Draw()

                } // Joystick State

                // Display State
                allegro.FlipDisplay()

            } // Timer Events

            default:
            }
        }
    }) // allegro.Run(func(){...})
}

