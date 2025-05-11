package public

type RegisterReq struct {
	Login    *string `json:"login" binding:"required,gte=1,lte=120"`
	Password *string `json:"password" binding:"required,gte=8,lte=60" validate:"c-password"`
}
