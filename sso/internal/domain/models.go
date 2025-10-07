package domain

type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
}

type App struct {
	ID     int64
	Secret string
}
