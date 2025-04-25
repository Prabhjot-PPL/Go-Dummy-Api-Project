package dto

type LoginResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Image        string `json:"image"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
