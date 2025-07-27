package sms_service

import (
	"context"
	"encoding/json"
	"logispro/internal/db"
	"logispro/internal/utils"
	"logispro/internal/web/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateSmsService struct {
	Queries      *db.Queries
	RabbitMqConn *amqp.Connection
}

func (s *CreateSmsService) Create(req requests.CreateSmsRequest, userId int64, agencyUsers []int64, ctx context.Context) (sms db.SmsQueue, err error) {
	user, err := s.Queries.GetUserById(ctx, userId)
	if err != nil {
		return sms, err
	}

	contacts, err := s.Queries.GetContactsById(ctx, db.GetContactsByIdParams{
		UsersID: agencyUsers,
		Ids:     req.Contacts,
	})

	if len(contacts) == 0 {
		return sms, err
	}

	sms, err = s.Queries.CreateSMSQueue(ctx, db.CreateSMSQueueParams{
		UserID:          userId,
		FromNumber:      user.Phone,
		Content:         req.Content,
		TotalRecipients: int64(len(req.Contacts)),
	})
	if err != nil {
		return sms, err
	}

	// delete sms if err

	for _, contact := range contacts {
		_, err := s.Queries.AddSMSQueueContact(ctx, db.AddSMSQueueContactParams{
			SmsQueueID:  sms.ID,
			PhoneNumber: contact.Phone.String,
		})
		if err != nil {
			return sms, err
		}
	}

	// push to the queue
	ch, err := s.RabbitMqConn.Channel()
	if err != nil {
		return sms, err
	}
	rmq := &utils.RabbitMQ{Conn: s.RabbitMqConn, Channel: ch}
	defer rmq.Channel.Close()
	rmq.DeclareQueue("sms_prepare")
	msgMap := map[string]interface{}{
		"sms_id": sms.ID,
	}
	msg, _ := json.Marshal(msgMap)
	err = rmq.Publish("sms_prepare", msg, amqp.Table{})
	if err != nil {
		return sms, err
	}
	return sms, err
}
