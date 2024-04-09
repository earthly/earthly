package earthfile2llb

type cloneable[T any] interface {
	Clone() T
}
