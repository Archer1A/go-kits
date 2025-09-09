package mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func NewProducer() rocketmq.Producer {
	producer, err := rocketmq.NewProducer(producer.WithNameServer([]string{"127.0.0.1:9876"}))
	if err != nil {
		panic("生成 producer 失败")
	}
	err = producer.Start()
	if err != nil {
		panic("启动 producer 失败")
	}
	return producer

}

func Send(producer rocketmq.Producer, topic string, body []byte) {
	sync, err := producer.SendSync(context.TODO(), &primitive.Message{Topic: topic, Body: body})
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return
	}

	fmt.Println(sync.Status)

}

func Consumer() {
	pushConsumer, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"127.0.0.1:9876"}), consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("GID_XXXXXX"))
	if err != nil {
		panic("生成 consumer 失败")
	}
	err = pushConsumer.Subscribe("rocketMq", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		if string(ext[0].Body) == "failed" {
			fmt.Println("fail rollback")
			return consumer.Rollback, nil
		}
		fmt.Println(string(ext[0].Body))
		return consumer.ConsumeSuccess, nil
	})
	err = pushConsumer.Start()
	if err != nil {
		panic("启动 consumer 失败:" + err.Error())
	}
	if err != nil {
		fmt.Println(err)
	}
	defer pushConsumer.Shutdown()
	select {}
}
