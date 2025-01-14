package entities

type Identifier struct {
	ID int64 `json:"id"`
}

func (i Identifier) IsNew() bool {
	return i.ID == 0
}
