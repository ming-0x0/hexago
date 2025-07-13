package banner

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/hexago/internal/shared/errors"
	"github.com/oklog/ulid/v2"
)

type ID string

const (
	maxNameLength = 255
)

//go:generate sh -c "$(go list -m -f '{{.Dir}}')/bin/accessor -type=Banner"
type Banner struct {
	id   ID
	name string
}

func New(
	name string,
) (*Banner, error) {
	banner := &Banner{
		id:   ID(ulid.Make().String()),
		name: name,
	}
	if err := banner.validate(); err != nil {
		return nil, err
	}

	return banner, nil
}

func FromRepository(
	id ID,
	name string,
) (*Banner, error) {
	banner := &Banner{
		id:   id,
		name: name,
	}
	if err := banner.validate(); err != nil {
		return nil, err
	}

	return banner, nil
}

func (b *Banner) validate() error {
	err := validation.ValidateStruct(b,
		validation.Field(
			&b.name,
			validation.Required,
			validation.Length(1, maxNameLength),
		),
	)
	if err != nil {
		return errors.NewDomainError(errors.Validation, err.Error())
	}
	return nil
}
