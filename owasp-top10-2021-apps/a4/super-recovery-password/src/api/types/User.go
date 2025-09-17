package types

type User struct {
	ID    string `json:"id" form:"id" query:"id"`
	Login string `json:"login" form:"login" query:"login"`
}

type UserLogin struct {
	ID       string `json:"id" form:"id" query:"id"`
	Login    string `json:"login" form:"login" query:"login"`
	Password string `json:"password" form:"password" query:"password"`
}

type UserRegister struct {
	Login          string `json:"login" form:"login" query:"login"`
	Password       string `json:"password" form:"password" query:"password"`
	RepeatPassword string `json:"repeatPassword" form:"repeatPassword" query:"repeatPassword"`
	Email		   string `json:"email" form:"email" query:"email"`
}

type RecoveryEmail struct {
	Login          string `json:"login" form:"login" query:"login"`
	Email		   string `json:"email" form:"email" query:"email"`
}

type ChangePassword struct {
	Password       string `json:"password" form:"password" query:"password"`
	RepeatPassword string `json:"repeatPassword" form:"repeatPassword" query:"repeatPassword"`
}
