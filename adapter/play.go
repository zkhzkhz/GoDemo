package main

import . "./player"

func main() {
	gameSound := GameSoundPlayer{Src: "music.mp3"}
	gameSoundAdapter := GameSoundAdapter{SoundPlayer: gameSound}
	play(gameSoundAdapter)
}

func play(player Player) {
	player.PlayMusic()
}
