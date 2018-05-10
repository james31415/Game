package main

import (
    "io/ioutil"
    "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/audio"
)

const (
    // SPRITES
    SPRITEDATA = "dat/sprites/"
    BACKGROUND = "background"
    PLAYER     = "player"
    SPAWN      = "spawn.wav"
)

type sprite struct {
    Name       string
    Folder     string
    Sound      map[string]*audio.Sample
    Bitmap    *allegro.Bitmap
    DrawFlags  allegro.DrawFlags
    OffsetX    float32
    OffsetY    float32
    ScaleX     float32
    ScaleY     float32
    Height     float32
    Width      float32
    Angle      float32
    X          float32
    Y          float32
    Draw       func()
    Spawn      func()
}

func (s *sprite) DrawNormal() {
    dx := s.X-s.OffsetX
    dy := s.Y-s.OffsetY
    df := s.DrawFlags
    s.Bitmap.Draw(dx,dy,df)
}

func (s *sprite) DrawScaled() { // TODO / BROKEN
    sx := s.X-s.OffsetX
    sy := s.Y-s.OffsetY
    sw := s.Width
    sh := s.Height
    dx := sx
    dy := sy
    dw := sw
    dh := sh
    df := s.DrawFlags
    s.Bitmap.DrawScaled(sx,sy,sw,sh,dx,dy,dw,dh,df)
}

func (s *sprite) DrawRotated() { // TODO / BROKEN
    cx := s.X-s.OffsetX
    cy := s.Y-s.OffsetY
    dx := cx
    dy := cy
    da := s.Angle
    df := s.DrawFlags
    s.Bitmap.DrawRotated(cx,cy,dx,dy,da,df)
}

func (s *sprite) Play(sound string) {
    logLvl(SPRITES|SOUNDS, " = PLAY:", s.Name, sound)

    instance := audio.CreateSampleInstance(s.Sound[sound])
    if err := instance.AttachToMixer(audio.DefaultMixer()); err != nil {
        logLvl(SPRITES|SOUNDS, " = FAIL:", s.Name, sound, err)
        panic(err)
    }
    if err := instance.Play(); err != nil {
        logLvl(SPRITES|SOUNDS, " = FAIL:", s.Name, sound, err)
        panic(err)
    }
}

func (s *sprite) Unload() { // TODO
    s.Bitmap.Destroy()
    s.OffsetX = 0
    s.OffsetY = 0
    s.Height  = 0
    s.Width   = 0
}

func (s *sprite) Center(display *allegro.Display) {
    s.Y = float32(display.Height()/2)
    s.X = float32(display.Width()/2)
}

func (s *sprite) Load(name string) {
    logLvl(SPRITES, "LOADING SPRITE:", name)
    s.Name = name
    s.Folder = SPRITEDATA + name

    loadSound := func(sound string) {
        logLvl(SPRITES|SOUNDS, " = SOUND:", sound)
        if sample, err := audio.LoadSample(s.Folder+"/snd/"+sound); err != nil {
            panic(err)
        } else {
            s.Sound[sound] = sample
        }
    }

    s.Sound = make(map[string]*audio.Sample)
    if sounds, err := ioutil.ReadDir(s.Folder+"/snd/"); err != nil {
        logLvl(SPRITES|SOUNDS, " = FAIL:", err)
    } else {
        for _, sound := range sounds {
            loadSound(sound.Name())
        }
    }

    if bitmap, err := allegro.LoadBitmap(s.Folder+"/bitmap"); err != nil {
        logLvl(SPRITES|BITMAPS, " = FAIL:", err)
    } else {
        logLvl(SPRITES|BITMAPS, " = BITMAP LOADED")
        s.Bitmap  = bitmap
        s.Width   = float32(bitmap.Width())
        s.Height  = float32(bitmap.Height())
        s.OffsetY = s.Height/2
        s.OffsetX = s.Width/2
        s.ScaleX  = 10.0
        s.ScaleY  = 10.0
        s.Draw    = s.DrawNormal
        s.Spawn   = func(){
            logLvl(SPRITES, " = SPAWN:", s.Name)
            s.Play(SPAWN)
            s.Draw()
        }
    }
}

func (s *sprite) Update() { }

