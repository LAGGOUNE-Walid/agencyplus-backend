package controllers

import (
	"logispro/internal/web/controllers/contact"
	"logispro/internal/web/controllers/user"
)

type Controller struct {
	UserController    *user.UserController
	ContactController *contact.ContactController
}
