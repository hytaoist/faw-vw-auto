package domain

type Usecaser interface {
	Versions(product string) ([]string, error)
}

type Databaser interface {
	Versions(product string) ([]string, error)
}
