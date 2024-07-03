package webhook

import (
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/infinitare/disgo"
)

func New() *Body {
	body := Body{
		Username:  "Copped AIO",
		AvatarUrl: branding.Icon,
	}

	return &body
}

func NewField(name string, value string) *disgo.EmbedField {

	field := disgo.EmbedField{
		Name:  name,
		Value: value,
	}

	return &field

}
