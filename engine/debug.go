package engine

import (
    "fmt"
)

type DebugLevel int

const (
    LOG_GENERAL DebugLevel = 1 << iota
    LOG_SPRITES
    LOG_BITMAPS
    LOG_SOUNDS
    LOG_EVENTS
    LOG_TIMER
)

const (
    LOGLVL  = LOG_GENERAL
    LOGGING = true
    LOOPING = true
)

func Log(                   msg ...interface{}) { if LOGGING          { fmt.Println(msg...)          } }
func LogLvl(lvl DebugLevel, msg ...interface{}) { if lvl & LOGLVL > 0 { Log(msg...)                  } }
func LogFunc(                   msg string, function func())          { Log(msg);         function() }
func LogLvlFunc(lvl DebugLevel, msg string, function func())          { LogLvl(lvl, msg); function() }

