package environments

type EnvironmentsLinks struct {
	First    EnvironmentsLink `json:"first"`
	Next     EnvironmentsLink `json:"next"`
	Previous EnvironmentsLink `json:"previous"`
	Last     EnvironmentsLink `json:"last"`
}
