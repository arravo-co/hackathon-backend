package dtos

type BasicLoginDTO struct {
	Identifier string ` validate:"required" json:"identifier"`
	Password   string ` validate:"required" json:"password"`
}
