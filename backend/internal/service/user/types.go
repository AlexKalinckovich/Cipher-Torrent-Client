package user

type CreateUserInput struct {
	Email     string
	PublicKey string
	Nickname  string
	Role      string
}

type UpdateUserInput struct {
	Email     string
	PublicKey string
	Nickname  string
	Role      string
}
