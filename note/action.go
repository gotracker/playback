package note

import "fmt"

// Action is the action to take on a note
type Action uint8

const (
	// ActionCut will cut the old note/instrument playback immediately
	// (this is the default for most tracker formats)
	ActionCut = Action(iota)
	// ActionContinue will continue the old note/instrument playback indefinitely
	ActionContinue
	// ActionRelease will perform a release (key-off) on the old note/instrument playback
	ActionRelease
	// ActionFadeout will fade out the old note/instrument playback
	// (if the instrument's fadeout volume is 0, then this effectively becomes a NewNoteActionContinue)
	ActionFadeout
	// ActionRetrigger will perform a key-on for the note/instrument playback immediately
	ActionRetrigger
)

func (a Action) String() string {
	switch a {
	case ActionCut:
		return "ActionCut"
	case ActionContinue:
		return "ActionContinue"
	case ActionRelease:
		return "ActionRelease"
	case ActionFadeout:
		return "ActionFadeout"
	case ActionRetrigger:
		return "ActionRetrigger"
	default:
		return fmt.Sprintf("Unknown[%d]", int(a))
	}
}
