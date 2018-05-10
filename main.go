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
    game.LoadSprite(BACKGROUND,0)
    game.LoadSprite(PLAYER,0)
    game.Update()

    /* ///////////
    ** Define Keys
    */ ///////////

    game.KeyMap = engine.KeyboardMap {
        engine.KEY_ESCAPE: func(){game.Quit()},
        engine.KEY_Q:      func(){game.Quit()},
        engine.KEY_P:      func(){game.GetSprite(PLAYER).Play(engine.SPAWN)},
        engine.KEY_X:      func(){game.GetSprite(PLAYER).DrawFlags ^= engine.FLIP_X},
        engine.KEY_Y:      func(){game.GetSprite(PLAYER).DrawFlags ^= engine.FLIP_Y},
        engine.KEY_UP:     func(){engine.LogLvl(engine.LOG_GENERAL, " = UP"       )},
        engine.KEY_DOWN:   func(){engine.LogLvl(engine.LOG_GENERAL, " = DOWN"     )},
        engine.KEY_LEFT:   func(){engine.LogLvl(engine.LOG_GENERAL, " = LEFT"     )},
        engine.KEY_RIGHT:  func(){engine.LogLvl(engine.LOG_GENERAL, " = RIGHT"    )},
    }

    game.JoyMap = engine.JoystickMap {
        engine.JOY_A:      func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_A"    )},
        engine.JOY_B:      func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_B"    )},
        engine.JOY_X:      func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_X"    )},
        engine.JOY_Y:      func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_Y"    )},
        engine.JOY_LB:     func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LB"   )},
        engine.JOY_RB:     func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RB"   )},
        engine.JOY_BACK:   func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_BACK" )},
        engine.JOY_START:  func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_START")},
        engine.JOY_XBOX:   func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_XBOX" )},
        engine.JOY_LS:     func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_LS"   )},
        engine.JOY_RS:     func(){engine.LogLvl(engine.LOG_GENERAL, " = JOY_RS"   )},
    }

    game.Loop()

})}

