package dbmodels

import (
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string `json:"name" gorm:"not null"`
	Email                    string `json:"email" gorm:"unique;not null"`
	Handle                   string `json:"handle" gorm:"unique;not null"`
	Password                 string `json:"password" gorm:"not null"`
	Bio                      string `json:"bio"`
	Avatar                   string `json:"avatar"`
	IsVerified               bool   `json:"is_verified" gorm:"not null;default:false"`

	Orgs []Organization `gorm:"many2many:user_organizations;"` // many to many
}

type UserOrganizations struct {
	UserUUID         uuid.UUID `json:"user_uuid" gorm:"type:uuid;primaryKey"`
	OrganizationUUID uuid.UUID `json:"organization_uuid" gorm:"type:uuid;primaryKey"`
	Role             string    `json:"role" gorm:"not null;default:member"`
}
