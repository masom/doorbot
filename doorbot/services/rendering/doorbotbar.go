package rendering

import (
	"strings"
)

type Renderer interface {
	Render(input string, parameters map[string]string) string
}

//
func DoorbotBar() Renderer {
	return &doorbotBar{}
}

type doorbotBar struct {
}

// Render the given input, replacing the given parameters keys with their set values.
// Parameters in the template are detected when wrapped inside {{ }} without spaces.
// Ex:
//     input := "Hello {{name}}"
//     db := doorbotbar.DoorBar()
//     output := db.Render(input, map[string]string{"name": "John"})
//
func (bb *doorbotBar) Render(input string, parameters map[string]string) string {
	for k, v := range parameters {
		input = strings.Replace(input, "{{"+k+"}}", v, -1)
	}

	return input
}
