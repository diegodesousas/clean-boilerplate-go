package errors

type Conflict struct {
	Errors ConflictErrors `json:"errors"`
}

type ConflictErrors struct {
	ReferenceId ConflictErrorMessages `json:"reference_id"`
}

type ConflictErrorMessages struct {
	Messages Unique `json:"messages"`
}

type Unique struct {
	Unique string `json:"unique"`
}

func (e Conflict) Error() string {
	return e.Errors.ReferenceId.Messages.Unique
}

func NewConflict(message string) *Conflict {
	return &Conflict{Errors: ConflictErrors{ReferenceId: ConflictErrorMessages{Messages: Unique{Unique: message}}}}
}
