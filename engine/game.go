package engine

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

var Start = allegro.Run

// Game State
type gameState struct {
    Timer    *allegro.Timer
    Display  *allegro.Display
    Events   *allegro.EventQueue
    JoyState *allegro.JoystickState
    JoyMap    JoystickMap
    KeyMap    KeyboardMap
    Sprite []*sprite
    Running   bool
}

func (game *gameState) LoadSprite(name string, flags int) {
    var s sprite
    s.Load(name)
    if flags & SPRITE_NOCENTER == 0 { s.Center(game.Display) }
    if flags & SPRITE_NOSPAWN  == 0 { s.Spawn() }
    game.Sprite = append(game.Sprite, &s)
}

func (game *gameState) GetSprite(name string) (*sprite) {
    for _, sprite := range game.Sprite {
        if sprite.Name == name {
            return sprite
        }
    }
    return nil
}

func (game *gameState) Update() {
    game.KeyMap.Check()
    game.JoyMap.Check(game.JoyState)
    for _, sprite := range game.Sprite {
        sprite.Update()
        sprite.Draw()
    }
    allegro.FlipDisplay()
}

func (game *gameState) Destroy() {
    LogLvlFunc(LOG_GENERAL, "DESTROY TIMER",                   game.Timer.Destroy)
    LogLvlFunc(LOG_GENERAL, "UNINSTALL MOUSE",                 allegro.UninstallMouse)
    LogLvlFunc(LOG_GENERAL, "UNINSTALL KEYBOARD",              allegro.UninstallKeyboard)
    LogLvlFunc(LOG_GENERAL, "UNINSTALL JOYSTICK",              allegro.UninstallJoystick)
    LogLvlFunc(LOG_GENERAL, "DESTROY DISPLAY",                 game.Display.Destroy)
    LogLvlFunc(LOG_GENERAL, "DESTROY EVENT QUEUE",             game.Events.Destroy)
    LogLvlFunc(LOG_GENERAL, "UNINSTALL AUDIO AND AUDIO CODEC", audio.Uninstall)
    LogLvlFunc(LOG_GENERAL, "UNINSTALL IMAGE",                 image.Uninstall)
    LogLvl(    LOG_GENERAL, "...FINISHED")
}

func NewGameState() (game gameState) {
    LogLvl(LOG_GENERAL, "STARTING...")
    var err error
    game.KeyMap = make(KeyboardMap)
    game.JoyMap = make(JoystickMap)

    LogLvl(LOG_GENERAL, "CONFIG DISPLAY")
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

    if flags == 0 {            LogLvl(LOG_GENERAL, " = DEFAULT"                  ) } else {
        if flags & 0x001 > 0 { LogLvl(LOG_GENERAL, " = WINDOWED"                 ) }
        if flags & 0x002 > 0 { LogLvl(LOG_GENERAL, " = FULLSCREEN"               ) }
        if flags & 0x004 > 0 { LogLvl(LOG_GENERAL, " = OPENGL"                   ) }
        if flags & 0x010 > 0 { LogLvl(LOG_GENERAL, " = RESIZABLE"                ) }
        if flags & 0x020 > 0 { LogLvl(LOG_GENERAL, " = FRAMELESS"                ) } // FRAMELESS == NOFRAME
        if flags & 0x020 > 0 { LogLvl(LOG_GENERAL, " = NOFRAME"                  ) } // NOFRAME == FRAMELESS
        if flags & 0x040 > 0 { LogLvl(LOG_GENERAL, " = GENERATE_EXPOSE_EVENTS"   ) }
        if flags & 0x080 > 0 { LogLvl(LOG_GENERAL, " = OPENGL_3_0"               ) }
        if flags & 0x100 > 0 { LogLvl(LOG_GENERAL, " = OPENGL_FORWARD_COMPATIBLE") }
        if flags & 0x200 > 0 { LogLvl(LOG_GENERAL, " = FULLSCREEN_WINDOW"        ) }
    }

    LogLvl(LOG_GENERAL, "INSTALL IMAGE")
    if err = image.Install(); err != nil {
        panic(err)
    }

    LogLvl(LOG_GENERAL, "INSTALL AUDIO")
    if err = audio.Install(); err != nil {
        panic(err)
    } else {
        audio.ReserveSamples(SAMPLEMAX)
        LogLvl(LOG_GENERAL, "INSTALL AUDIO CODEC")
        if err = acodec.Install(); err != nil {
            panic(err)
        }
    }

    LogLvl(LOG_GENERAL, "CREATE EVENT QUEUE")
    if game.Events, err = allegro.CreateEventQueue(); err != nil {
        panic(err)
    }
    
    LogLvl(LOG_GENERAL, "CREATE DISPLAY")
    allegro.SetNewDisplayFlags(flags)
    if game.Display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
        panic(err)
    } else {
        game.Display.SetWindowTitle("Game")
        game.Events.Register(game.Display)
        game.Events.RegisterEventSource(game.Display.EventSource()) // Redundant?
    }

    LogLvl(LOG_GENERAL, "INSTALL JOYSTICK")
    if err = allegro.InstallJoystick(); err != nil {
        panic(err)
    } else {
        //game.Events.RegisterEventSource(allegro.JoystickEventSource())
        game.JoyState = ConfigureJoysticks()
    }
    
    LogLvl(LOG_GENERAL, "INSTALL KEYBOARD")
    if err = allegro.InstallKeyboard(); err != nil {
        panic(err)
    }
    
    LogLvl(LOG_GENERAL, "INSTALL MOUSE")
    if err = allegro.InstallMouse(); err != nil {
        panic(err)
    }
    
    LogLvl(LOG_GENERAL, "CREATE TIMER")
    if game.Timer, err = allegro.CreateTimer(1.0 / FPS); err != nil {
        panic(err)
    } else {
        game.Events.Register(game.Timer)
        game.Timer.Start()
    }

    return
}

func (game *gameState) Quit() { LogLvl(LOG_GENERAL, "QUITING..."); game.Running = false }

func (game *gameState) Loop() {
    LogLvl(LOG_GENERAL, "LOOPING...")
    game.Running = LOOPING
    var event allegro.Event
    for game.Running { switch e := game.Events.WaitForEvent(&event); e.(type) {

        // Display Events
        case allegro.DisplayCloseEvent:          LogLvl(LOG_EVENTS, "DisplayCloseEvent"              ); game.Quit()
        case allegro.DisplayExposeEvent:         LogLvl(LOG_EVENTS, "DisplayExposeEvent"             ); allegro.FlipDisplay()
        case allegro.DisplayFoundEvent:          LogLvl(LOG_EVENTS, "DisplayFoundEvent"              );
        case allegro.DisplayLostEvent:           LogLvl(LOG_EVENTS, "DisplayLostEvent"               );
        case allegro.DisplayOrientationEvent:    LogLvl(LOG_EVENTS, "DisplayOrientationEvent"        );
        case allegro.DisplayResizeEvent:         LogLvl(LOG_EVENTS, "DisplayResizeEvent"             ); game.Display.AcknowledgeResize()
        case allegro.DisplaySwitchInEvent:       LogLvl(LOG_EVENTS, "DisplaySwitchInEvent"           );
        case allegro.DisplaySwitchOutEvent:      LogLvl(LOG_EVENTS, "DisplaySwitchOutEvent"          );

        // Joystick Events
        case allegro.JoystickConfigurationEvent: LogLvl(LOG_EVENTS, "JoystickConfigurationEvent"     ); game.JoyState = ConfigureJoysticks()

        // Mouse Events
        case allegro.MouseAxesEvent:             LogLvl(LOG_EVENTS, "MouseAxesEvent"                 );
        case allegro.MouseButtonDownEvent:       LogLvl(LOG_EVENTS, "MouseButtonDownEvent"           );
        case allegro.MouseButtonUpEvent:         LogLvl(LOG_EVENTS, "MouseButtonUpEvent"             );
        case allegro.MouseEnterDisplayEvent:     LogLvl(LOG_EVENTS, "MouseEnterDisplayEvent"         ); game.Timer.Start() // Pauses when mouse
        case allegro.MouseLeaveDisplayEvent:     LogLvl(LOG_EVENTS, "MouseLeaveDisplayEvent"         ); game.Timer.Stop()  // leaves the window
        case allegro.MouseWarpedEvent:           LogLvl(LOG_EVENTS, "MouseWarpedEvent"               );

        // Timer Events
        case allegro.TimerEvent:                 LogLvl(LOG_TIMER,   "TimerEvent:", game.Timer.Count()); game.Update()

        default:
        }
    }
}

