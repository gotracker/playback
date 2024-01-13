package feature

import "github.com/heucuva/optional"

type StartOrderAndRow struct {
	Order optional.Value[int]
	Row   optional.Value[int]
}
