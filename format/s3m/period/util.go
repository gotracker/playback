package period

import (
	"github.com/gotracker/playback/format/s3m/system"
	"github.com/gotracker/playback/frequency"
)

var DefaultC4SampleRate = system.DefaultC4SampleRate
var semitonePeriodTable = [...]float32{27392, 25856, 24384, 23040, 21696, 20480, 19328, 18240, 17216, 16256, 15360, 14496}

// CalcFinetuneC4SampleRate calculates a new frequency after a finetune adjustment
func CalcFinetuneC4SampleRate(finetune uint8) frequency.Frequency {
	switch finetune {
	case 0x0:
		return 7895
	case 0x1:
		return 7941
	case 0x2:
		return 7985
	case 0x3:
		return 8046
	case 0x4:
		return 8107
	case 0x5:
		return 8169
	case 0x6:
		return 8232
	case 0x7:
		return 8280
	case 0x8:
		return 8363
	case 0x9:
		return 8413
	case 0xA:
		return 8463
	case 0xB:
		return 8529
	case 0xC:
		return 8581
	case 0xD:
		return 8651
	case 0xE:
		return 8723
	case 0xF:
		return 8757
	default:
		panic("unhandled")
	}
}
