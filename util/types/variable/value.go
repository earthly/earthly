package variable

import "github.com/earthly/earthly/domain"

type Value struct {
	Str      string
	ComeFrom domain.Target
}
