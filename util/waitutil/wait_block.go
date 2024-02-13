package waitutil

import "context"

// WaitBlock stores items within a WAIT / END block
type WaitBlock interface {
	Wait(ctx context.Context, push, save bool) error
	AddItem(item WaitItem)
	SetDoSaves()
	SetDoPushes()
}
