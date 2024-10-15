package utility

type MockMethod int

const (
	MockCreate MockMethod = iota
	MockGet
	MockUpdate
	MockDelete
	MockUpload
)
