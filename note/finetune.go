package note

import (
	"fmt"
)

// Finetune is a 1/64th of a Semitone
type Finetune int16

func (f Finetune) String() string {
	return fmt.Sprintf("%s(%d)", Normal(f/64), f%64)
}
