// Package vegr (hel...vegr - the road to helheim) is a place to store library
// calls that hel mocks need in order to work. We put any non-trivial logic here
// so that it can have tests, and then allow hel mocks to import and use this
// properly-tested logic.
//
// Since this is only intended for use by generated code, it eagerly panics for
// obviously incorrect usage.
package vegr
