package pagehelp

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

type User string

func Test_panic(t *testing.T) {
	var m map[string]int
	if _, ok := m["a"]; ok {

	}
	m["a"] = 1
	mutex := sync.Mutex{}
	mutex.Lock()
}

type response struct {
	total int
	rows  []User
}

func (r *response) GetTotal() int {
	return r.total
}

func (r *response) GetRows() []User {
	return r.rows
}

func listAll(ctx context.Context, page int, pageSize int) (IResponse[User], error) {
	fmt.Printf("page %d , pageSize %d \n", page, pageSize)
	return &response{
		total: 460,
		rows:  []User{"1"},
	}, nil
}

func TestClient_SyncAll(t *testing.T) {
	client := NewClient[User]()
	all, err := client.SyncAll(context.TODO(), listAll)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(all))

}
