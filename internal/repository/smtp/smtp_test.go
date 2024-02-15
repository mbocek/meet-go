package smtp

import (
	"github.com/mbocek/meet-go/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMail(t *testing.T) {
	smtp := NewSMTP(config.SMTP{
		Host:     "localhost",
		Port:     "9025",
		From:     "test@test.com",
		User:     "myuser",
		Password: "mysecretpassword",
	})

	err := smtp.Send("recipient@test.com", `<!DOCTYPE html>
			<html>
			<body>
			<h3>Name:</h3><span>{{.Name}}</span><br/><br/>
			<h3>Email:</h3><span>{{.Message}}</span><br/>
			</body>
			</html>`, struct {
		Name    string
		Message string
	}{
		Name:    "Puneet Singh",
		Message: "This is a test message in a HTML template",
	})
	assert.Nil(t, err)
}
