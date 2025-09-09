package mq

import (
	"testing"
)

func TestConsumer(t *testing.T) {
	go Consumer()
	//time.Sleep(3 * time.Second)
	//producer := NewProducer()
	//todo := context.TODO()
	//for i := 0; i < 10; i++ {
	//	if i == 5 {
	//		sync, err := producer.SendSync(todo, &primitive.Message{Topic: "rocketMq", Body: []byte("hello world" + strconv.Itoa(i))})
	//		if err != nil {
	//			return
	//		}
	//		fmt.Println("Sync message" + sync.MsgID)
	//	} else {
	//		sync, err := producer.SendSync(todo, &primitive.Message{Topic: "rocketMq", Body: []byte("hello world" + strconv.Itoa(i))})
	//		if err != nil {
	//			return
	//		}
	//		fmt.Println("Sync message" + sync.MsgID)
	//	}
	//
	//}

	select {}

}
