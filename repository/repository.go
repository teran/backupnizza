package repository

type Repository interface {
	Label() string
	password() string
}

type repository struct {
	label  string
	passwd string
}

func NewRepository(address, passwd string) (Repository, error) {
	return &repository{
		label:  address,
		passwd: passwd,
	}, nil
}

func (r *repository) Label() string {
	return r.label
}

func (r *repository) password() string {
	return r.passwd
}
