package main

import (
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/image"
    "github.com/dradtke/go-allegro/allegro/audio"
    "github.com/dradtke/go-allegro/allegro/acodec"
    //"github.com/yuin/gopher-lua"
)

const (
    // SOUND
    SAMPLEMAX = 8
    // DISPLAY
    WINX = 800
    WINY = 600
    // MAIN
    FPS  = 60
)

/* /////////////////////////////////////////////////////////////////////////////
** Under Construction //                                                     //
** ///////////////////////////////////////////////////////////////////////////
**
** Kludged together and tested displays, events, sounds, keyboards, joysticks,
** and sprites.  Now that the library portions have largely been implemented,
** I've started breaking up the code into more manageable chunks.  Will need
** to start discussing implementation strategies as far as how the engine
** handles it's own data structures.  Once the mechanics have been ironed
** out, we'll need to start adding in external configuration and script
** support.
*/

// Keyboard
type keyboardMap map[allegro.KeyCode]func()

func (keyMap keyboardMap) Check() {
    var keyState allegro.KeyboardState
    keyState.Get()
    for k, f := range keyMap {
        if keyState.IsDown(k) { f() }
    }
}

// Game State
type gameState struct {
    display  *allegro.Display
    joyState *allegro.JoystickState
    joyMap    joystickMap
    keyMap    keyboardMap
    sprite  []sprite
}

const (
    SPRITE_CENTER = 1 << iota
    SPRITE_SPAWN
)

func (game *gameState) LoadSprite(name string, flags int) {
    var s sprite
    s.Load(name)
    if flags & SPRITE_CENTER > 0 { s.Center(game.display) }
    if flags & SPRITE_SPAWN  > 0 { s.Spawn() }
    game.sprite = append(game.sprite, s)
}

func (game *gameState) Update() {
    game.keyMap.Check()
    game.joyMap.Check(game.joyState)
    for sprite := range game.sprite { game.sprite[sprite].Update() }
    for sprite := range game.sprite { game.sprite[sprite].Draw()   }
    allegro.FlipDisplay()
}

func newGameState() (game gameState) {
    game.keyMap = make(keyboardMap)
    game.joyMap = make(joystickMap)
    return
}

// /////////////////////////////////////////////////////////////////////////
func main() {
    logLvl(GENERAL, "STARTING..."); defer logLvl(GENERAL, "...FINISHED")
    allegro.Run(func() {
        var (
            event       allegro.Event
            eventQueue *allegro.EventQueue
            timer      *allegro.Timer
            running     bool = LOOPING
            err         error
        )

        game := newGameState()

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
            audio.ReserveSamples(SAMPLEMAX)
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
        if game.display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
            panic(err)
        } else {
            defer logLvlFunc(GENERAL, "DESTROY DISPLAY", game.display.Destroy)
            game.display.SetWindowTitle("Game")
            eventQueue.Register(game.display)
            eventQueue.RegisterEventSource(game.display.EventSource()) // Redundant?
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
            game.joyState = configureJoysticks()
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

        /* ///////////
        ** Sprite Test
        */ ///////////
        game.LoadSprite(BACKGROUND, SPRITE_CENTER|SPRITE_SPAWN)
        game.LoadSprite(PLAYER,     SPRITE_CENTER|SPRITE_SPAWN)
        game.Update()

        /* ///////////
        ** Define Keys
        */ ///////////
        quit := func() { logLvl(GENERAL, "QUITING..."); running = false }

        game.keyMap[allegro.KEY_ESCAPE] = quit
        game.keyMap[allegro.KEY_Q]      = quit
        game.keyMap[allegro.KEY_UP]     = func(){logLvl(GENERAL, " = UP"   )}
        game.keyMap[allegro.KEY_DOWN]   = func(){logLvl(GENERAL, " = DOWN" )}
        game.keyMap[allegro.KEY_LEFT]   = func(){logLvl(GENERAL, " = LEFT" )}
        game.keyMap[allegro.KEY_RIGHT]  = func(){logLvl(GENERAL, " = RIGHT")}

        game.joyMap[JOY_A]     = func(){logLvl(GENERAL, " = JOY_A"    )}
        game.joyMap[JOY_B]     = func(){logLvl(GENERAL, " = JOY_B"    )}
        game.joyMap[JOY_X]     = func(){logLvl(GENERAL, " = JOY_X"    )}
        game.joyMap[JOY_Y]     = func(){logLvl(GENERAL, " = JOY_Y"    )}
        game.joyMap[JOY_LB]    = func(){logLvl(GENERAL, " = JOY_LB"   )}
        game.joyMap[JOY_RB]    = func(){logLvl(GENERAL, " = JOY_RB"   )}
        game.joyMap[JOY_BACK]  = func(){logLvl(GENERAL, " = JOY_BACK" )}
        game.joyMap[JOY_START] = func(){logLvl(GENERAL, " = JOY_START")}
        game.joyMap[JOY_XBOX]  = func(){logLvl(GENERAL, " = JOY_XBOX" )}
        game.joyMap[JOY_LS]    = func(){logLvl(GENERAL, " = JOY_LS"   )}
        game.joyMap[JOY_RS]    = func(){logLvl(GENERAL, " = JOY_RS"   )}

        /* ///////// // ///////////
        ** Main Loop // Switch Loop
        */ ///////// // ///////////
        logLvl(GENERAL, "LOOPING...")
        for running { switch e := eventQueue.WaitForEvent(&event); e.(type) {

            // Display Events
            case allegro.DisplayCloseEvent:          logLvl(GENERAL, "DisplayCloseEvent"         ); quit()
            case allegro.DisplayExposeEvent:         logLvl(GENERAL, "DisplayExposeEvent"        ); allegro.FlipDisplay()
            case allegro.DisplayFoundEvent:          logLvl(GENERAL, "DisplayFoundEvent"         );
            case allegro.DisplayLostEvent:           logLvl(GENERAL, "DisplayLostEvent"          );
            case allegro.DisplayOrientationEvent:    logLvl(GENERAL, "DisplayOrientationEvent"   );
            case allegro.DisplayResizeEvent:         logLvl(GENERAL, "DisplayResizeEvent"        ); game.display.AcknowledgeResize()
            case allegro.DisplaySwitchInEvent:       logLvl(GENERAL, "DisplaySwitchInEvent"      );
            case allegro.DisplaySwitchOutEvent:      logLvl(GENERAL, "DisplaySwitchOutEvent"     );

            // Joystick Events
            case allegro.JoystickConfigurationEvent: logLvl(GENERAL, "JoystickConfigurationEvent"); game.joyState = configureJoysticks()

            // Mouse Events
            case allegro.MouseAxesEvent:             logLvl(GENERAL, "MouseAxesEvent"            );
            case allegro.MouseButtonDownEvent:       logLvl(GENERAL, "MouseButtonDownEvent"      );
            case allegro.MouseButtonUpEvent:         logLvl(GENERAL, "MouseButtonUpEvent"        );
            case allegro.MouseEnterDisplayEvent:     logLvl(GENERAL, "MouseEnterDisplayEvent"    ); timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     logLvl(GENERAL, "MouseLeaveDisplayEvent"    ); timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           logLvl(GENERAL, "MouseWarpedEvent"          );

            // Timer Events
            case allegro.TimerEvent:                 logLvl(TIMER,   "TimerEvent:", timer.Count()); game.Update()

            default:
            }
        }
    }) // allegro.Run(func(){...})
}

