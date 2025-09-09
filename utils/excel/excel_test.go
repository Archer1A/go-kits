package excel

import (
	"encoding/json"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create excel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Create()
		})
	}
}

func TestGetTag(t *testing.T) {
	type User struct {
		Username string `json:"username" excel:"header:用户名"`
	}
	getTag(User{})
}

func Test_Generator(t *testing.T) {
	excel := &Client[User]{}
	err := excel.Sheet("Sheet1").Header(User{}).Rows([]User{{Username: "archer1a", Phone: "15088888888", Sex: "男", Email: "1231", Age: 18}}).Save("test.xlsx")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("save success")

}

func Test_MapValueToHeader(t *testing.T) {
	user := &User{
		Username: "archer1a",
	}

	headerMap := MapValueToHeader(user)
	marshal, err := json.Marshal(headerMap)
	if err != nil {
		return
	}
	t.Log(string(marshal))
	t.Log(headerMap)
}
