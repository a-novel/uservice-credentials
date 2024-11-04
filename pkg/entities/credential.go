package entities

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/database"
	"github.com/a-novel/golib/grpc"
)

var (
	ErrUnknownRole         = errors.New("unknown value for credentials_role")
	ErrUnsupportedRoleType = errors.New("unsupported type for credentials_role")
)

type Credential struct {
	bun.BaseModel `bun:"table:credentials,alias:credentials"`

	ID uuid.UUID `bun:"id,pk,type:uuid"`

	Email string `bun:"email"`
	Role  Role   `bun:"role,type:credentials_role"`

	EmailValidationTokenID        string `bun:"email_validation_token_id"`
	PendingEmailValidationTokenID string `bun:"pending_email_validation_token_id"`
	PasswordTokenID               string `bun:"password_token_id"`
	ResetPasswordTokenID          string `bun:"reset_password_token_id"`

	CreatedAt time.Time  `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
}

type Role string

const (
	RoleNone               Role = ""
	RoleEarlyAccessProgram Role = "early-access-program"
	RoleAdmin              Role = "admin"
	RoleCore               Role = "core"
)

var RolesSorted = []Role{
	RoleNone,
	RoleEarlyAccessProgram,
	RoleAdmin,
	RoleCore,
}

var (
	_ sql.Scanner   = (*Role)(nil)
	_ driver.Valuer = (*Role)(nil)
)

func (role *Role) String() string {
	if *role == RoleNone {
		return "none"
	}

	return string(*role)
}

func (role *Role) FromString(value string) error {
	if value == "none" {
		*role = RoleNone
		return nil
	}

	for _, roleValue := range RolesSorted {
		if string(roleValue) == value {
			*role = roleValue
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknownRole, value)
}

func (role *Role) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case string:
		return role.FromString(src)
	case []byte:
		return role.FromString(string(src))
	case nil:
		*role = RoleNone
		return nil
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedRoleType, src)
	}
}

func (role Role) Value() (driver.Value, error) {
	return role.String(), nil
}

func RegisterRole(customValidator *validator.Validate) {
	database.MustRegisterValidation(
		customValidator, "role",
		database.ValidateEnum(
			RoleNone,
			RoleEarlyAccessProgram,
			RoleAdmin,
			RoleCore,
		),
	)
}

var RoleConverter = grpc.NewProtoConverter(
	grpc.ProtoMapper[commonv1.UserRole, Role]{
		commonv1.UserRole_USER_ROLE_EARLY_ACCESS_PROGRAM: RoleEarlyAccessProgram,
		commonv1.UserRole_USER_ROLE_ADMIN:                RoleAdmin,
		commonv1.UserRole_USER_ROLE_CORE:                 RoleCore,
	},
	commonv1.UserRole_USER_ROLE_UNSPECIFIED,
	RoleNone,
)

type SortCredentials string

const (
	SortCredentialsNone      SortCredentials = ""
	SortCredentialsEmail     SortCredentials = "email"
	SortCredentialsRole      SortCredentials = "role"
	SortCredentialsCreatedAt SortCredentials = "created_at"
	SortCredentialsUpdatedAt SortCredentials = "updated_at"
)

func RegisterSortCredentials(customValidator *validator.Validate) {
	database.MustRegisterValidation(
		customValidator, "sort_credentials",
		database.ValidateEnum(
			SortCredentialsNone,
			SortCredentialsEmail,
			SortCredentialsRole,
			SortCredentialsCreatedAt,
			SortCredentialsUpdatedAt,
		),
	)
}

var SortCredentialsConverter = grpc.NewProtoConverter(
	grpc.ProtoMapper[credentialsv1.Sort, SortCredentials]{
		credentialsv1.Sort_SORT_BY_EMAIL:      SortCredentialsEmail,
		credentialsv1.Sort_SORT_BY_ROLE:       SortCredentialsRole,
		credentialsv1.Sort_SORT_BY_CREATED_AT: SortCredentialsCreatedAt,
		credentialsv1.Sort_SORT_BY_UPDATED_AT: SortCredentialsUpdatedAt,
	},
	credentialsv1.Sort_SORT_UNSPECIFIED,
	SortCredentialsNone,
)
