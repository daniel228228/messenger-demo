package dto

type GetUser struct {
	ID string `json:"id"`
}

type User struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type CreateUser struct {
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type UserID struct {
	ID string `json:"id"`
}
