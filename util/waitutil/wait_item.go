package waitutil

// WaitItem is an item to wait for
type WaitItem interface {
	SetDoPush()
	SetDoSave()
}
