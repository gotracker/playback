package channel

import s3mfile "github.com/gotracker/goaudiofile/music/tracked/s3m"

func VolumeFactory(mem *Memory, v s3mfile.Volume) EffectS3M {
	return SetVolume(v)
}
