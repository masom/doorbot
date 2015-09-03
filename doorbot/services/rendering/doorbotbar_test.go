package rendering

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoorBotBar(t *testing.T) {
	assert.Implements(t, (*Renderer)(nil), DoorbotBar())
}

func TestDoorBotBarRender(t *testing.T) {
	dbb := DoorbotBar()

	assert.Equal(
		t,
		"This is a test",
		dbb.Render("This is a test", map[string]string{}),
	)

	assert.Equal(
		t,
		"Hello bob",
		dbb.Render("Hello {{name}}", map[string]string{"name": "bob"}),
	)

	assert.Equal(
		t,
		"Hello bob bobbob{{}}",
		dbb.Render("Hello {{name}} {{name}}{{name}}{{}}", map[string]string{"name": "bob"}),
	)

	assert.Equal(
		t,
		"Hello John Rambo",
		dbb.Render("Hello {{first}} {{last}}", map[string]string{"first": "John", "last": "Rambo"}),
	)
}
