package user

type User struct {
	Uid      int
	Username string `json:"username"`
	Password string `json:"password"` // fixed typo
}
