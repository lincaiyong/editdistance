package editdistance

const (
	OpInsert  = "+"
	OpDelete  = "-"
	OpReplace = "*"
	OpKeep    = "="
)

func NewInsertOp(to string, toIndex int) *EditOp {
	return &EditOp{
		Type:    OpInsert,
		To:      to,
		ToIndex: toIndex,
	}
}

func NewDeleteOp(from string, fromIndex int) *EditOp {
	return &EditOp{
		Type:      OpDelete,
		From:      from,
		FromIndex: fromIndex,
	}
}

func NewReplaceOp(from, to string, fromIndex, toIndex int) *EditOp {
	return &EditOp{
		Type:      OpReplace,
		From:      from,
		To:        to,
		FromIndex: fromIndex,
		ToIndex:   toIndex,
	}
}

func NewKeepOp(from, to string, fromIndex, toIndex int) *EditOp {
	return &EditOp{
		Type:      OpKeep,
		From:      from,
		To:        to,
		FromIndex: fromIndex,
		ToIndex:   toIndex,
	}
}

type EditOp struct {
	Type      string `json:"type"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	FromIndex int    `json:"from_index,omitempty"`
	ToIndex   int    `json:"to_index,omitempty"`
}
