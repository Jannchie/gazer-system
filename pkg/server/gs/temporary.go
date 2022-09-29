package gs

type temporary interface {
	Temporary() bool // IsTemporary returns true if err is temporary.
}

func IsTemporary(err error) bool {
	te, ok := err.(temporary)
	return ok && te.Temporary()
}

type TemporaryError struct {
	error
}

func (t *TemporaryError) Temporary() bool {
	return true
}
