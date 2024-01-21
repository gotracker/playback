package voice

type VoiceFactory[TPeriod Period, TGlobalVolume, TMixingVolume, TVolume Volume, TPanning Panning] interface {
	NewVoice(settings VoiceConfig[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]) RenderVoice[TPeriod, TGlobalVolume, TMixingVolume, TVolume, TPanning]
}
