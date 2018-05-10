package main

import (
    "github.com/dradtke/go-allegro/allegro"
    //"github.com/yuin/gopher-lua"
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
    logLvl(GENERAL, "STARTING..."); defer logLvl(GENERAL, "...FINISHED")
    allegro.Run(func() {
        var event allegro.Event

        game := newGameState(); defer game.Destroy()

        /* ///////////
        ** Sprite Test
        */ ///////////
        game.LoadSprite(BACKGROUND, SPRITE_CENTER|SPRITE_SPAWN)
        game.LoadSprite(PLAYER,     SPRITE_CENTER|SPRITE_SPAWN)
        game.Update()

        /* ///////////
        ** Define Keys
        */ ///////////
        quit := func() { logLvl(GENERAL, "QUITING..."); game.running = false }

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
        game.running = LOOPING
        for game.running { switch e := game.events.WaitForEvent(&event); e.(type) {

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
            case allegro.MouseEnterDisplayEvent:     logLvl(GENERAL, "MouseEnterDisplayEvent"    ); game.timer.Start() // Pauses when mouse
            case allegro.MouseLeaveDisplayEvent:     logLvl(GENERAL, "MouseLeaveDisplayEvent"    ); game.timer.Stop()  // leaves the window
            case allegro.MouseWarpedEvent:           logLvl(GENERAL, "MouseWarpedEvent"          );

            // Timer Events
            case allegro.TimerEvent:                 logLvl(TIMER,   "TimerEvent:", game.timer.Count()); game.Update()

            default:
            }
        }
    }) // allegro.Run(func(){...})
}

