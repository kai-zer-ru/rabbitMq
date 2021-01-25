package rabbitMqGolang

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	messagesData := map[string][]string{}
	//messagesResponseData := map[string][]string{} // todo check result
	for i := 0; i < 10; i++ {
		messagesData[fmt.Sprintf("channel%v", i)] = []string{}
		for j := 0; j < 10; j++ {
			messagesData[fmt.Sprintf("channel%v", i)] = append(messagesData[fmt.Sprintf("channel%v", i)], fmt.Sprintf("text%v", j))
		}
	}
	r := RabbitMQ{}
	err := r.Connect()
	if err != nil {
		panic(err)
	}
	go func() {
		for channel, texts := range messagesData {
			for _, text := range texts {
				err = r.Send(channel, []byte(text))
				if err != nil {
					t.Errorf(err.Error())
				}
				t.Logf("Message '%s' is sended to channel '%s'", text, channel)
				r := rand.Intn(10)
				t.Logf("Sleep %v sec", r)
				time.Sleep(time.Duration(r) * time.Millisecond)
			}
		}
	}()
	channels := []string{"channel1", "channel2"}
	data, err := r.Listen(channels...)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case m := <-data["channel1"]:
				t.Log("Channel = channel1")
				//_, ok := messagesResponseData["channel1"]
				//if !ok {
				//	messagesResponseData["channel1"] = []string{}
				//}
				exist, _ := inArray(string(m.Body), messagesData["channel1"])
				if !exist {
					t.Errorf("Invalid text: " + string(m.Body))
				}
				//messagesResponseData["channel1"] = append(messagesResponseData["channel1"])
				t.Logf("Message = %s", m.Body)
			case m := <-data["channel2"]:
				t.Log("Channel = channel2")
				//_, ok := messagesResponseData["channel2"]
				//if !ok {
				//	messagesResponseData["channel2"] = []string{}
				//}
				exist, _ := inArray(string(m.Body), messagesData["channel2"])
				if !exist {
					t.Errorf("Invalid text: " + string(m.Body))
				}
				//messagesResponseData["channel2"] = append(messagesResponseData["channel2"])
				t.Logf("Message = %s", m.Body)
			}
		}
	}()
	<-time.After(30 * time.Second)
}

func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
