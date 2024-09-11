package properties

type Links struct {
	First    Link `json:"first"`
	Next     Link `json:"next"`
	Previous Link `json:"previous"`
	Last     Link `json:"last"`
}
