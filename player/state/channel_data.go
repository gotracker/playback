package state

import "github.com/gotracker/playback"

type ChannelData[TChannelData any, TChannelState playback.ChannelState] struct {
	txn     ChannelDataTransaction[TChannelData, TChannelState]
	prevTxn ChannelDataTransaction[TChannelData, TChannelState]
}

// GetData returns the interface to the current channel song pattern data
func (cd *ChannelData[TChannelData, TChannelState]) GetData() *TChannelData {
	if cd.txn == nil {
		return nil
	}

	return cd.txn.GetData()
}

func (cd *ChannelData[TChannelData, TChannelState]) SetData(cdata *TChannelData, cs *TChannelState) error {
	if cd.txn == nil {
		return nil
	}

	return cd.txn.SetData(cdata, cs)
}

// AdvanceRow will update the current state to make room for the next row's state data
func (cd *ChannelData[TChannelData, TChannelState]) AdvanceRow(txn ChannelDataTransaction[TChannelData, TChannelState]) {
	cd.prevTxn, cd.txn = cd.txn, txn
}

func (cd *ChannelData[TChannelData, TChannelState]) GetTxn() ChannelDataTransaction[TChannelData, TChannelState] {
	return cd.txn
}

// AddNoteOp adds a note operation to the channel transaction
func (cd *ChannelData[TChannelData, TChannelState]) AddNoteOp(op NoteOp[TChannelState]) {
	if cd.txn != nil {
		cd.txn.AddNoteOp(op)
	}
}
