package player

import "fmt"

type Player interface {
	PlayMusic()
}

type MusicPlayer struct {
	Src string
}
type GameSoundPlayer struct {
	Src string
}

type GameSoundAdapter struct {
	SoundPlayer GameSoundPlayer
}

func (g GameSoundPlayer) PlaySound() {
	fmt.Println("play sound :" + g.Src)
}

func (a GameSoundAdapter) PlayMusic() {
	a.SoundPlayer.PlaySound()
}
