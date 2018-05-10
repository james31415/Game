package main

import (
    "github.com/james31415/Game/engine"
    "github.com/dradtke/go-allegro/allegro"
    //"github.com/yuin/gopher-lua"
)

const (
    BACKGROUND = "background"
    PLAYER     = "player"
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

// /////////////////////////////////////////////////////////////////////////
func main() {
    engine.LogLvl(engine.LOG_GENERAL, "STARTING..."); defer engine.LogLvl(engine.LOG_GENERAL, "...FINISHED")
    allegro.Run(func() {
        var event allegro.Event

        game := engine.NewGameState(); defer game.Destroy()

        /* ///////////
        ** Sprite Test
        */ ///////////
        game.LoadSprite(BACKGROUND, engine.SPRITE_CENTER|engine.SPRITE_SPAWN)
        game.LoadSprite(PLAYER,     engine.SPRITE_CENTER|engine.SPRITE_SPAWN)
        game.Update()

        /* ///////////
        ** Define Keys
        */ ///////////
        quit := func() { engine.LogLvl(engine.LOG_GENERAL, "QUITING..."); game.Running = false }

        game.KeyMap[allegro.KEY_ESCAPE] = quit
        game.KeyMap[allegro.KEY_Q]      = quit
        game.KeyMap[allegro.KEY_UP]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = UP"   )}
        game.KeyMap[allegro.KEY_DOWN]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = DOWN" )}
        game.KeyMap[allegro.KEY_LEFT]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = LEFT" )}
        game.KeyMap[allegro.KEY_RIGHT]  = func(){engine.LogLvl(engine.LOG_GENERAL, " = RIGHT")}

        game.JoyMap[engine.JOY_A]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_A"    )}
        game.JoyMap[engine.JOY_B]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_B"    )}
        game.JoyMap[engine.JOY_X]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_X"    )}
        game.JoyMap[engine.JOY_Y]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_Y"    )}
        game.JoyMap[engine.JOY_LB]    = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LB"   )}
        game.JoyMap[engine.JOY_RB]    = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RB"   )}
        game.JoyMap[engine.JOY_BACK]  = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_BACK" )}
        game.JoyMap[engine.JOY_START] = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_START")}
        game.JoyMap[engine.JOY_XBOX]  = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_XBOX" )}
        game.JoyMap[engine.JOY_LS]    = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LS"   )}
        game.JoyMap[engine.JOY_RS]    = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RS"   )}

        /* ///////// // ///////////
        ** Main Loop // Switch Loop
        */ ///////// // ///////////
        engine.LogLvl(engine.LOG_GENERAL, "LOOPING...")
        game.Running = engine.LOOPING
        for game.Running { switch e := game.Events.WaitForEvent(&event); e.(type) {

            // Display Events
            case allegro.DisplayCloseEvent:          engine.LogLvl(engine.LOG_EVENTS, "DisplayCloseEvent"              ); quit()
            case allegro.DisplayExposeEvent:         engine.LogLvl(engine.LOG_EVENTS, "DisplayExposeEvent"             ); allegro.FlipDisplay()
            case allegro.DisplayFoundEvent:          engine.LogLvl(engine.LOG_EVENTS, "DisplayFoundEvent"              );
            case allegro.DisplayLostEvent:           engine.LogLvl(engine.LOG_EVENTS, "DisplayLostEvent"               );
            case allegro.DisplayOrientationEvent:    engine.LogLvl(engine.LOG_EVENTS, "DisplayOrientationEvent"        );
            case allegro.DisplayResizeEvent:         engine.LogLvl(engine.LOG_EVENTS, "DisplayResizeEvent"             ); game.Display.AcknowledgeResize()
            case allegro.DisplaySwitchInEvent:       engine.LogLvl(engine.LOG_EVENTS, "DisplaySwitchInEvent"           );
            case allegro.DisplaySwitchOutEvent:      engine.LogLvl(engine.LOG_EVENTS, "DisplaySwitchOutEvent"          );

            // Joystick Events
            case allegro.JoystickConfigurationEvent: engine.LogLvl(engine.LOG_EVENTS, "JoystickConfigurationEvent"     ); game.JoyState = engine.ConfigureJoysticks()

            // Mouse Events
            case allegro.MouseAxesEvent:             engine.LogLvl(engine.LOG_EVENTS, "MouseAxesEvent"                 );
            case allegro.MouseButtonDownEvent:       engine.LogLvl(engine.LOG_EVENTS, "MouseButtonDownEvent"           );
            case allegro.MouseButtonUpEvent:         engine.LogLvl(engine.LOG_EVENTS, "MouseButtonUpEvent"             );
            case allegro.MouseEnterDisplayEvent:     engine.LogLvl(engine.LOG_EVENTS, "MouseEnterDisplayEvent"         ); game.Timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     engine.LogLvl(engine.LOG_EVENTS, "MouseLeaveDisplayEvent"         ); game.Timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           engine.LogLvl(engine.LOG_EVENTS, "MouseWarpedEvent"               );

            // Timer Events
            case allegro.TimerEvent:                 engine.LogLvl(engine.LOG_TIMER,   "TimerEvent:", game.Timer.Count()); game.Update()

            default:
            }
        }
    }) // allegro.Run(func(){...})
}

