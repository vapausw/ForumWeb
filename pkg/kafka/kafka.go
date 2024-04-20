package kafka

import (
	"ForumWeb/model"
	"ForumWeb/pkg/smtp"
	"ForumWeb/setting"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// NewKafkaProducer 创建一个新的生产者
func NewKafkaProducer(brokers []string, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return w
}

// NewKafkaConsumer 创建一个新的消费者
func NewKafkaConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func StartEmailProducer(msgBytes []byte) {
	go func() {
		kafkaProducer := NewKafkaProducer(setting.Conf.MyKafkaConfig.Brokers, setting.Conf.MyKafkaConfig.EmailTopic)
		// 发送消息，注意不使用 `c`
		if err := kafkaProducer.WriteMessages(context.Background(), kafka.Message{
			Value: msgBytes,
		}); err != nil {
			zap.L().Error("failed to send Kafka message", zap.Error(err))
		}
		kafkaProducer.Close()
	}()
}

func StartEmailConsumer(ctx context.Context) {
	reader := NewKafkaConsumer(setting.Conf.MyKafkaConfig.Brokers, setting.Conf.MyKafkaConfig.EmailTopic, setting.Conf.MyKafkaConfig.GroupID)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			zap.L().Error("failed to read message from Kafka", zap.Error(err))
			continue
		}

		var msg model.WelcomeEmailMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			zap.L().Error("failed to unmarshal Kafka message", zap.Error(err))
			continue
		}

		// 使用你的邮件发送逻辑
		err = smtp.SendEmail(msg.Email, []byte(msg.Message))
		if err != nil {
			zap.L().Error("failed to send email", zap.Error(err))
		}
	}
}
