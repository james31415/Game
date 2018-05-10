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

// Game State
type gameState struct {
    timer    *allegro.Timer
    display  *allegro.Display
    events   *allegro.EventQueue
    joyState *allegro.JoystickState
    joyMap    joystickMap
    keyMap    keyboardMap
    sprite  []sprite
    running   bool
}

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

func (game *gameState) Destroy() {
    logLvlFunc(GENERAL, "DESTROY TIMER", game.timer.Destroy)
    logLvlFunc(GENERAL, "UNINSTALL MOUSE", allegro.UninstallMouse)
    logLvlFunc(GENERAL, "UNINSTALL KEYBOARD", allegro.UninstallKeyboard)
    logLvlFunc(GENERAL, "UNINSTALL JOYSTICK", allegro.UninstallJoystick)
    logLvlFunc(GENERAL, "DESTROY DISPLAY", game.display.Destroy)
    logLvlFunc(GENERAL, "DESTROY EVENT QUEUE", game.events.Destroy)
    logLvlFunc(GENERAL, "UNINSTALL AUDIO AND AUDIO CODEC", audio.Uninstall)
    logLvlFunc(GENERAL, "UNINSTALL IMAGE", image.Uninstall)
}

func newGameState() (game gameState) {
    var err error
    game.keyMap = make(keyboardMap)
    game.joyMap = make(joystickMap)

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

    logLvl(GENERAL, "INSTALL IMAGE")
    if err = image.Install(); err != nil {
        panic(err)
    }

    logLvl(GENERAL, "INSTALL AUDIO")
    if err = audio.Install(); err != nil {
        panic(err)
    } else {
        audio.ReserveSamples(SAMPLEMAX)
        logLvl(GENERAL, "INSTALL AUDIO CODEC")
        if err = acodec.Install(); err != nil {
            panic(err)
        }
    }

    logLvl(GENERAL, "CREATE EVENT QUEUE")
    if game.events, err = allegro.CreateEventQueue(); err != nil {
        panic(err)
    }
    
    logLvl(GENERAL, "CREATE DISPLAY")
    allegro.SetNewDisplayFlags(flags)
    if game.display, err = allegro.CreateDisplay(WINX,WINY); err != nil {
        panic(err)
    } else {
        game.display.SetWindowTitle("Game")
        game.events.Register(game.display)
        game.events.RegisterEventSource(game.display.EventSource()) // Redundant?
    }

    logLvl(GENERAL, "INSTALL JOYSTICK")
    if err = allegro.InstallJoystick(); err != nil {
        panic(err)
    } else {
        //game.events.RegisterEventSource(allegro.JoystickEventSource())
        game.joyState = configureJoysticks()
    }
    
    logLvl(GENERAL, "INSTALL KEYBOARD")
    if err = allegro.InstallKeyboard(); err != nil {
        panic(err)
    }
    
    logLvl(GENERAL, "INSTALL MOUSE")
    if err = allegro.InstallMouse(); err != nil {
        panic(err)
    }
    
    logLvl(GENERAL, "CREATE TIMER")
    if game.timer, err = allegro.CreateTimer(1.0 / FPS); err != nil {
        panic(err)
    } else {
        game.events.Register(game.timer)
        game.timer.Start()
    }

    return
}

