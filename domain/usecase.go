package domain

type usecase struct {
	db Databaser
}

func NewUsecase(db Databaser) *usecase {
	return &usecase{db}
}

func (u *usecase) Versions(product string) ([]string, error) {
	versions, err := u.db.Versions(product)
	if err != nil {
		return nil, err
	}
	return versions, nil
}
