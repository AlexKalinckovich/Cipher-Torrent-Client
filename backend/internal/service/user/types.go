package user

type CreateUserInput struct {
	Email        string
	Nickname     string
	PasswordHash string
	Role         string
}

type UpdateUserInput struct {
	Email        string
	Nickname     string
	PasswordHash string
	Role         string
}
