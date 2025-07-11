package domain

type Usecaser interface {
	Versions(product string) ([]string, error)
	SumScore() (int16, error)
}

type Databaser interface {
	Versions(product string) ([]string, error)
	SumScore() (int16, error)
}
