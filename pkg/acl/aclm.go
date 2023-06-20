package acl

import (
	"context"
	"strings"
	"sync"
	"time"

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
	ctx          context.Context
	closeCtx     context.CancelFunc
	mode         aclmMode
	roleAddCh    *amqp.Channel
	roleRemoveCh *amqp.Channel
	roleTokenCh  *amqp.Channel
	wg           sync.WaitGroup
}

func newACLM(mode aclmMode) *ACLM {
	// init rabbitmq
	ctx, cancel := context.WithCancel(context.Background())
	a := &ACLM{
		ctx:      ctx,
		closeCtx: cancel,
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
	if a.roleAddCh != nil {
		if err := a.roleAddCh.Close(); err != nil {
			return err
		}
		if err := a.roleRemoveCh.Close(); err != nil {
			return err
		}
		if err := a.roleTokenCh.Close(); err != nil {
			return err
		}
	}
	a.closeCtx()
	a.wg.Wait()
	return nil
}

func (a *ACLM) Mode() aclmMode {
	return a.mode
}

func (a *ACLM) listenRoleAdd() {
	var err error
	a.roleAddCh, err = rabbitmq.AmqpChannel()
	if err != nil {
		panic(err)
	}
	roleQueue, err := a.roleAddCh.QueueDeclare(
		QueueRoleAdd,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.roleAddCh.Consume(
		roleQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	closed := make(chan struct{})
	go func() {
		go func() {
			for {
				ec := a.roleAddCh.NotifyClose(make(chan *amqp.Error))
				if ec == nil {
					time.Sleep(rabbitmq.RetryDuration)
					continue
				}
				e := <-ec
				logger.Error(e)
				close(closed)
				time.Sleep(rabbitmq.RetryDuration)
				a.listenRoleAdd()
				return
			}
		}()
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			select {
			case d, ok := <-msgs:
				if ok {
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

					if err := d.Ack(true); err != nil {
						logger.Error("Failed to ack role add message: %s", err)
					}
				}
			case <-a.ctx.Done():
				return
			case <-closed:
				return
			}
		}
	}()

}

func (a *ACLM) listenRoleRemove() {
	var err error
	a.roleRemoveCh, err = rabbitmq.AmqpChannel()
	if err != nil {
		panic(err)
	}

	roleQueue, err := a.roleRemoveCh.QueueDeclare(
		QueueRoleRemove,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.roleRemoveCh.Consume(
		roleQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	closed := make(chan struct{})
	go func() {
		go func() {
			for {
				ec := a.roleRemoveCh.NotifyClose(make(chan *amqp.Error))
				if ec == nil {
					time.Sleep(rabbitmq.RetryDuration)
					continue
				}
				e := <-ec
				logger.Error(e)
				close(closed)
				time.Sleep(rabbitmq.RetryDuration)
				a.listenRoleRemove()
				return
			}
		}()
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		for {
			select {
			case d, ok := <-msgs:
				if ok {
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

					err := d.Ack(true)

					if err != nil {
						logger.Error("Failed to ack role add message: %s", err)
					}
				}
			case <-a.ctx.Done():
				return
			case <-closed:
				return
			}
		}
	}()
}

func (a *ACLM) listenTokenCheck() {
	var err error
	a.roleTokenCh, err = rabbitmq.AmqpChannel()
	if err != nil {
		panic(err)
	}

	closed := make(chan struct{})
	go func() {
		go func() {
			for {
				ec := a.roleTokenCh.NotifyClose(make(chan *amqp.Error))
				if ec == nil {
					time.Sleep(rabbitmq.RetryDuration)
					continue
				}
				e := <-ec
				logger.Error(e)
				close(closed)
				time.Sleep(rabbitmq.RetryDuration)
				a.listenTokenCheck()
				return
			}
		}()
	}()

	roleQueue, err := a.roleTokenCh.QueueDeclare(
		QueueTokenCheck,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := a.roleTokenCh.Consume(
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
		for {
			select {
			case d, ok := <-msgs:
				if ok {
					parts := strings.SplitN(strings.TrimSpace(string(d.Body)), ":", 2)
					if len(parts) != 2 {
						logger.Error("Invalid token check message: %s", string(d.Body))
						if err := a.roleTokenCh.Ack(d.DeliveryTag, false); err != nil {
							logger.Error("Failed to ack token check message: %s", err)
						}
						continue
					}
					userID := parts[0]
					token := parts[1]

					cacheDeviceID := getDeviceID(userID, token)

					err = a.roleTokenCh.Publish("", d.ReplyTo, false, false, amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(cacheDeviceID),
						MessageId:     d.MessageId,
					})
					if err != nil {
						logger.Error("Failed to ack token check message: %s", err)
						continue
					}

					if err := a.roleTokenCh.Ack(d.DeliveryTag, false); err != nil {
						logger.Error("Failed to ack token check message: %s", err)
					}
				}
			case <-a.ctx.Done():
				return
			case <-closed:
				return
			}
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Error("%s: %s", msg, err)
	}
}
