package domain

type UserRepository interface {
	Save(user User) error
	GetByEmail(email string) (User, error)
}
