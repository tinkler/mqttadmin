package acl

import (
	"context"
	"strings"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"github.com/tinkler/mqttadmin/pkg/rabbitmq"
)

type aclmMode string

const (
	AM_SERVER aclmMode = "server"
	AM_CLIENT aclmMode = "client"
)

type ACLM struct {
	ctx      context.Context
	closeCtx context.CancelFunc
	mode     aclmMode
	ch       *amqp.Channel
	wg       sync.WaitGroup
}

func newACLM(mode aclmMode) *ACLM {
	// init rabbitmq
	conn := rabbitmq.AmqpConn()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a := &ACLM{
		ctx:      ctx,
		closeCtx: cancel,
		ch:       ch,
		mode:     mode,
	}
	switch mode {
	case AM_SERVER:
		// extract roles from role message
		a.listenRoleAdd()
		a.listenRoleRemove()
		a.listenTokenCheck()
	case AM_CLIENT:
	}

	return a
}

func (a *ACLM) Close() error {
	a.ch.Close()
	a.closeCtx()
	a.wg.Wait()
	return a.ch.Close()
}

func (a *ACLM) Mode() aclmMode {
	return a.mode
}

func (a *ACLM) listenRoleAdd() {
	roleQueue, err := a.ch.QueueDeclare(
		QueueRoleAdd,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.ch.Consume(
		roleQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for d := range msgs {
			if a.ctx.Err() != nil {
				return
			}
			parts := strings.SplitN(string(d.Body), ":", 2)
			userID := strings.TrimSpace(parts[0])
			roleList := strings.TrimSpace(parts[1])
			if len(roleList) == 0 {
				continue
			}
			roleNames := strings.Split(parts[1], ",")
			var roleIDs []string
			if err := db.DB().Raw("SELECT id FROM authv1.role WHERE name IN ('" + strings.Join(roleNames, "','") + "')").Scan(&roleIDs).Error; err != nil {
				logger.Error("Failed to get role id: %s", err)
				continue
			}
			if len(roleIDs) == 0 {
				logger.Warn("Failed to get role id: %v", roleNames)
			} else {
				if err := db.DB().Exec("INSERT INTO authv1.user_role (user_id, role_id) VALUES ('" + userID + "', '" + strings.Join(roleIDs, "'), ('"+userID+"', '") + "') ON CONFLICT (user_id,role_id) DO UPDATE SET delete_at = null").Error; err != nil {
					logger.Error("Failed to add user role: %s", err)
					continue
				}
			}
			SetAllDeviceRemoveFlag(userID)

			if err := d.Ack(false); err != nil {
				logger.Error("Failed to ack role add message: %s", err)
			}
		}
	}()

}

func (a *ACLM) listenRoleRemove() {
	roleQueue, err := a.ch.QueueDeclare(
		QueueRoleRemove,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.ch.Consume(
		roleQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for d := range msgs {
			if a.ctx.Err() != nil {
				return
			}
			parts := strings.SplitN(string(d.Body), ":", 2)
			userID := strings.TrimSpace(parts[0])
			roleList := strings.TrimSpace(parts[1])
			if len(roleList) == 0 {
				continue
			}
			roleNames := strings.Split(parts[1], ",")

			se := db.DB().Exec("UPDATE authv1.user_role SET delete_at = CURRENT_TIMESTAMP WHERE user_id = ? AND role_id IN (SELECT id FROM authv1.role WHERE name IN ('"+strings.Join(roleNames, "','")+"'))", userID)
			if se.Error != nil {
				logger.Error("Failed to remove user role: %s", se.Error)
				continue
			}
			SetAllDeviceRemoveFlag(userID)

			err := d.Ack(false)

			if err != nil {
				logger.Error("Failed to ack role add message: %s", err)
			}
		}
	}()
}

func (a *ACLM) listenTokenCheck() {
	roleQueue, err := a.ch.QueueDeclare(
		QueueTokenCheck,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.ch.Consume(
		roleQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for d := range msgs {
			if a.ctx.Err() != nil {
				return
			}
			parts := strings.SplitN(strings.TrimSpace(string(d.Body)), ":", 2)
			if len(parts) != 2 {
				logger.Error("Invalid token check message: %s", string(d.Body))
				if err := a.ch.Ack(d.DeliveryTag, false); err != nil {
					logger.Error("Failed to ack token check message: %s", err)
				}
				continue
			}
			userID := parts[0]
			token := parts[1]

			cacheDeviceID := getDeviceID(userID, token)

			err = a.ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(cacheDeviceID),
				MessageId:     d.MessageId,
			})
			if err != nil {
				logger.Error("Failed to ack token check message: %s", err)
				continue
			}

			if err := a.ch.Ack(d.DeliveryTag, false); err != nil {
				logger.Error("Failed to ack token check message: %s", err)
			}
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Error("%s: %s", msg, err)
	}
}
