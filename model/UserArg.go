package model

type UserArg struct {
	PageArg
	Ttype string `form:"ttype" json:"ttype"`
}
