package main

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

