package mask

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUsername2(t *testing.T) {
	type User struct {
		Username UsernameMask `json:"username"`
	}

	marshal, err := json.Marshal(User{Username: "张三"})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(string(marshal))
}

func TestUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "张三",
			args: args{
				username: "张三",
			},
			want: "张*",
		},
		{
			name: "empty",
			args: args{
				username: "",
			},
			want: "",
		},
		{
			name: "张三哦",
			args: args{
				username: "张三哦",
			},
			want: "张*哦",
		},
		{
			name: "上官三哦",
			args: args{
				username: "上官三哦",
			},
			want: "上**哦",
		},
		{
			name: "reload",
			args: args{
				username: "reload",
			},
			want: "r****d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Username(tt.args.username); got != tt.want {
				t.Errorf("Username() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "length 4",
			args: args{phone: "1888"},
			want: "1888",
		},
		{
			name: "empty",
			args: args{phone: ""},
			want: "",
		},
		{
			name: "1888888888",
			args: args{phone: "1888888888"},
			want: "188****8888",
		},
		{
			name: "+861888888888",
			args: args{phone: "+861888888888"},
			want: "+86 188****8888",
		},
		{
			name: "177777",
			args: args{phone: "17777"},
			want: "177****7777",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Phone(tt.args.phone); got != tt.want {
				t.Errorf("Phone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhoneMask_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       PhoneMask
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsername1(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Username(tt.args.username); got != tt.want {
				t.Errorf("Username() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsernameMask_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		mask    UsernameMask
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mask.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "abc@gmail.com",
			args: args{email: "abc@gmail.com"},
			want: "a*c@gmail.com",
		},
		{
			name: "abcd@gmail.com",
			args: args{email: "abcd@gmail.com"},
			want: "a**d@gmail.com",
		},
		{
			name: "err",
			args: args{email: "abcdgmail.com"},
			want: "abcdgmail.com",
		},
		{
			name: "err",
			args: args{email: "abc@"},
			want: "a*c@",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Email(tt.args.email); got != tt.want {
				t.Errorf("Email() = %v, want %v", got, tt.want)
			}
		})
	}
}
