package controllers

import (
	"logispro/internal/web/controllers/building"
	"logispro/internal/web/controllers/contact"
	"logispro/internal/web/controllers/sms"
	"logispro/internal/web/controllers/user"
)

type Controller struct {
	UserController     *user.UserController
	ContactController  *contact.ContactController
	BuildingController *building.BuildingController
	SmsController      *sms.SmsController
}
