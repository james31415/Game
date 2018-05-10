package main

import (
    "fmt"
)

type DebugLevel int

const (
    GENERAL DebugLevel = 1 << iota
    SPRITES
    BITMAPS
    SOUNDS
    TIMER
)

const ( // DEBUG
    LOGLVL  = SPRITES
    LOGGING = true
    LOOPING = true
)

/* ///////
** Logging
*/ ///////
func log(                   msg ...interface{}) { if LOGGING          { fmt.Println(msg...)          } }
func logLvl(lvl DebugLevel, msg ...interface{}) { if lvl & LOGLVL > 0 { log(msg...)                  } }
func logFunc(                   msg string, function func())          { log(msg);         function() }
func logLvlFunc(lvl DebugLevel, msg string, function func())          { logLvl(lvl, msg); function() }

