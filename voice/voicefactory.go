package voice

type VoiceFactory[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	NewVoice() RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
}
