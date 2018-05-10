package engine

import (
    "github.com/dradtke/go-allegro/allegro"
)

const (
    STICK_THRESHOLD = 0.5
)

const (
    JOY_A = iota
    JOY_B
    JOY_X
    JOY_Y
    JOY_LB
    JOY_RB
    JOY_BACK
    JOY_START
    JOY_XBOX
    JOY_LS
    JOY_RS
)

const (
    KEY_A allegro.KeyCode = allegro.KEY_A
    KEY_B                 = allegro.KEY_B
    KEY_C                 = allegro.KEY_C
    KEY_D                 = allegro.KEY_D
    KEY_E                 = allegro.KEY_E
    KEY_F                 = allegro.KEY_F
    KEY_G                 = allegro.KEY_G
    KEY_H                 = allegro.KEY_H
    KEY_I                 = allegro.KEY_I
    KEY_J                 = allegro.KEY_J
    KEY_K                 = allegro.KEY_K
    KEY_L                 = allegro.KEY_L
    KEY_M                 = allegro.KEY_M
    KEY_N                 = allegro.KEY_N
    KEY_O                 = allegro.KEY_O
    KEY_P                 = allegro.KEY_P
    KEY_Q                 = allegro.KEY_Q
    KEY_R                 = allegro.KEY_R
    KEY_S                 = allegro.KEY_S
    KEY_T                 = allegro.KEY_T
    KEY_U                 = allegro.KEY_U
    KEY_V                 = allegro.KEY_V
    KEY_W                 = allegro.KEY_W
    KEY_X                 = allegro.KEY_X
    KEY_Y                 = allegro.KEY_Y
    KEY_Z                 = allegro.KEY_Z
    KEY_0                 = allegro.KEY_0
    KEY_1                 = allegro.KEY_1
    KEY_2                 = allegro.KEY_2
    KEY_3                 = allegro.KEY_3
    KEY_4                 = allegro.KEY_4
    KEY_5                 = allegro.KEY_5
    KEY_6                 = allegro.KEY_6
    KEY_7                 = allegro.KEY_7
    KEY_8                 = allegro.KEY_8
    KEY_9                 = allegro.KEY_9
    KEY_PAD_0             = allegro.KEY_PAD_0
    KEY_PAD_1             = allegro.KEY_PAD_1
    KEY_PAD_2             = allegro.KEY_PAD_2
    KEY_PAD_3             = allegro.KEY_PAD_3
    KEY_PAD_4             = allegro.KEY_PAD_4
    KEY_PAD_5             = allegro.KEY_PAD_5
    KEY_PAD_6             = allegro.KEY_PAD_6
    KEY_PAD_7             = allegro.KEY_PAD_7
    KEY_PAD_8             = allegro.KEY_PAD_8
    KEY_PAD_9             = allegro.KEY_PAD_9
    KEY_F1                = allegro.KEY_F1
    KEY_F2                = allegro.KEY_F2
    KEY_F3                = allegro.KEY_F3
    KEY_F4                = allegro.KEY_F4
    KEY_F5                = allegro.KEY_F5
    KEY_F6                = allegro.KEY_F6
    KEY_F7                = allegro.KEY_F7
    KEY_F8                = allegro.KEY_F8
    KEY_F9                = allegro.KEY_F9
    KEY_F10               = allegro.KEY_F10
    KEY_F11               = allegro.KEY_F11
    KEY_F12               = allegro.KEY_F12
    KEY_ESCAPE            = allegro.KEY_ESCAPE
    KEY_TILDE             = allegro.KEY_TILDE
    KEY_MINUS             = allegro.KEY_MINUS
    KEY_EQUALS            = allegro.KEY_EQUALS
    KEY_BACKSPACE         = allegro.KEY_BACKSPACE
    KEY_TAB               = allegro.KEY_TAB
    KEY_OPENBRACE         = allegro.KEY_OPENBRACE
    KEY_CLOSEBRACE        = allegro.KEY_CLOSEBRACE
    KEY_ENTER             = allegro.KEY_ENTER
    KEY_SEMICOLON         = allegro.KEY_SEMICOLON
    KEY_QUOTE             = allegro.KEY_QUOTE
    KEY_BACKSLASH         = allegro.KEY_BACKSLASH
    KEY_BACKSLASH2        = allegro.KEY_BACKSLASH2
    KEY_COMMA             = allegro.KEY_COMMA
    KEY_FULLSTOP          = allegro.KEY_FULLSTOP
    KEY_SLASH             = allegro.KEY_SLASH
    KEY_SPACE             = allegro.KEY_SPACE
    KEY_INSERT            = allegro.KEY_INSERT
    KEY_DELETE            = allegro.KEY_DELETE
    KEY_HOME              = allegro.KEY_HOME
    KEY_END               = allegro.KEY_END
    KEY_PGUP              = allegro.KEY_PGUP
    KEY_PGDN              = allegro.KEY_PGDN
    KEY_LEFT              = allegro.KEY_LEFT
    KEY_RIGHT             = allegro.KEY_RIGHT
    KEY_UP                = allegro.KEY_UP
    KEY_DOWN              = allegro.KEY_DOWN
    KEY_PAD_SLASH         = allegro.KEY_PAD_SLASH
    KEY_PAD_ASTERISK      = allegro.KEY_PAD_ASTERISK
    KEY_PAD_MINUS         = allegro.KEY_PAD_MINUS
    KEY_PAD_PLUS          = allegro.KEY_PAD_PLUS
    KEY_PAD_DELETE        = allegro.KEY_PAD_DELETE
    KEY_PAD_ENTER         = allegro.KEY_PAD_ENTER
    KEY_PRINTSCREEN       = allegro.KEY_PRINTSCREEN
    KEY_PAUSE             = allegro.KEY_PAUSE
    KEY_ABNT_C1           = allegro.KEY_ABNT_C1
    KEY_YEN               = allegro.KEY_YEN
    KEY_KANA              = allegro.KEY_KANA
    KEY_CONVERT           = allegro.KEY_CONVERT
    KEY_NOCONVERT         = allegro.KEY_NOCONVERT
    KEY_AT                = allegro.KEY_AT
    KEY_CIRCUMFLEX        = allegro.KEY_CIRCUMFLEX
    KEY_COLON2            = allegro.KEY_COLON2
    KEY_KANJI             = allegro.KEY_KANJI
    KEY_LSHIFT            = allegro.KEY_LSHIFT
    KEY_RSHIFT            = allegro.KEY_RSHIFT
    KEY_LCTRL             = allegro.KEY_LCTRL
    KEY_RCTRL             = allegro.KEY_RCTRL
    KEY_ALT               = allegro.KEY_ALT
    KEY_ALTGR             = allegro.KEY_ALTGR
    KEY_LWIN              = allegro.KEY_LWIN
    KEY_RWIN              = allegro.KEY_RWIN
    KEY_MENU              = allegro.KEY_MENU
    KEY_SCROLLLOCK        = allegro.KEY_SCROLLLOCK
    KEY_NUMLOCK           = allegro.KEY_NUMLOCK
    KEY_CAPSLOCK          = allegro.KEY_CAPSLOCK
    KEY_PAD_EQUALS        = allegro.KEY_PAD_EQUALS
    KEY_BACKQUOTE         = allegro.KEY_BACKQUOTE
    KEY_SEMICOLON2        = allegro.KEY_SEMICOLON2
    KEY_COMMAND           = allegro.KEY_COMMAND
)

const (
    KEYMOD_SHIFT allegro.KeyModifier = allegro.KEYMOD_SHIFT
    KEYMOD_CTRL                      = allegro.KEYMOD_CTRL
    KEYMOD_ALT                       = allegro.KEYMOD_ALT
    KEYMOD_LWIN                      = allegro.KEYMOD_LWIN
    KEYMOD_RWIN                      = allegro.KEYMOD_RWIN
    KEYMOD_MENU                      = allegro.KEYMOD_MENU
    KEYMOD_ALTGR                     = allegro.KEYMOD_ALTGR
    KEYMOD_COMMAND                   = allegro.KEYMOD_COMMAND
    KEYMOD_SCROLLLOCK                = allegro.KEYMOD_SCROLLLOCK
    KEYMOD_NUMLOCK                   = allegro.KEYMOD_NUMLOCK
    KEYMOD_CAPSLOCK                  = allegro.KEYMOD_CAPSLOCK
    KEYMOD_INALTSEQ                  = allegro.KEYMOD_INALTSEQ
    KEYMOD_ACCENT1                   = allegro.KEYMOD_ACCENT1
    KEYMOD_ACCENT2                   = allegro.KEYMOD_ACCENT2
    KEYMOD_ACCENT3                   = allegro.KEYMOD_ACCENT3
    KEYMOD_ACCENT4                   = allegro.KEYMOD_ACCENT4
)

type keyboardMap map[allegro.KeyCode]func()

func (keyMap keyboardMap) Check() {
    var keyState allegro.KeyboardState
    keyState.Get()
    for k, f := range keyMap {
        if keyState.IsDown(k) { f() }
    }
}

type joystickMap map[int]func()

func (joyMap joystickMap) Check(joyState *allegro.JoystickState) {
    if joyState == nil { return }
    joyState.Get()
    for b, f := range joyMap {
        if joyState.Button[b] > 0 { f() }
    }

    // TODO
    Axis_LX := joyState.Stick[0].Axis[0]
    Axis_LY := joyState.Stick[0].Axis[1]
    Axis_TL := joyState.Stick[1].Axis[0]
    Axis_RX := joyState.Stick[1].Axis[1]
    Axis_RY := joyState.Stick[2].Axis[0]
    Axis_TR := joyState.Stick[2].Axis[1]
    Axis_DX := joyState.Stick[3].Axis[0]
    Axis_DY := joyState.Stick[3].Axis[1]

    if Axis_LY < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = L_STICK_UP"   ) }
    if Axis_LY >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = L_STICK_DOWN" ) }
    if Axis_LX < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = L_STICK_LEFT" ) }
    if Axis_LX >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = L_STICK_RIGHT") }

    if Axis_RY < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = R_STICK_UP"   ) }
    if Axis_RY >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = R_STICK_DOWN" ) }
    if Axis_RX < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = R_STICK_LEFT" ) }
    if Axis_RX >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = R_STICK_RIGHT") }

    if Axis_DY < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = DIR_PAD_UP"   ) }
    if Axis_DY >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = DIR_PAD_DOWN" ) }
    if Axis_DX < -STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = DIR_PAD_LEFT" ) }
    if Axis_DX >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = DIR_PAD_RIGHT") }

    if Axis_TL >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = TRIGGER_LEFT" ) }
    if Axis_TR >  STICK_THRESHOLD { LogLvl(LOG_GENERAL, " = TRIGGER_RIGHT") }

}

func ConfigureJoysticks() (joyState *allegro.JoystickState) {
    LogLvl(LOG_GENERAL, "CONFIGURING JOYSTICKS")
    if allegro.ReconfigureJoysticks() { LogLvl(LOG_GENERAL, " = RECONFIGURED") }
    if joys := allegro.NumJoysticks(); joys > 0 {
        for joy := 0; joy < joys; joy++ {
            if joystick, err := allegro.GetJoystick(joy); err != nil {
                panic(err)
            } else {
                LogLvl(LOG_GENERAL, " = JOYSTICK:", joy, "=", joystick.Name())
                joyState = joystick.State()
            }
        }
    }
    return
}

