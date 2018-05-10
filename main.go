package main

import (
    "github.com/james31415/Game/engine"
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
func main() {engine.Start(func(){
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

    game.KeyMap[engine.KEY_ESCAPE] = func(){game.Quit()}
    game.KeyMap[engine.KEY_Q]      = func(){game.Quit()}
    game.KeyMap[engine.KEY_UP]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = UP"       )}
    game.KeyMap[engine.KEY_DOWN]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = DOWN"     )}
    game.KeyMap[engine.KEY_LEFT]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = LEFT"     )}
    game.KeyMap[engine.KEY_RIGHT]  = func(){engine.LogLvl(engine.LOG_GENERAL, " = RIGHT"    )}

    game.JoyMap[engine.JOY_A]      = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_A"    )}
    game.JoyMap[engine.JOY_B]      = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_B"    )}
    game.JoyMap[engine.JOY_X]      = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_X"    )}
    game.JoyMap[engine.JOY_Y]      = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_Y"    )}
    game.JoyMap[engine.JOY_LB]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LB"   )}
    game.JoyMap[engine.JOY_RB]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RB"   )}
    game.JoyMap[engine.JOY_BACK]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_BACK" )}
    game.JoyMap[engine.JOY_START]  = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_START")}
    game.JoyMap[engine.JOY_XBOX]   = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_XBOX" )}
    game.JoyMap[engine.JOY_LS]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LS"   )}
    game.JoyMap[engine.JOY_RS]     = func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RS"   )}

    game.Loop()
})}

