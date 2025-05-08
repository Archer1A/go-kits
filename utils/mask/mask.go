package mask

import (
	"encoding/json"
	"fmt"
	"strings"
)

const CountryCodeEndIndex = 3

type UsernameMask string

type PhoneMask string

type EmailMask string

func (mask UsernameMask) MarshalJSON() ([]byte, error) {
	masked := Username(string(mask))
	return json.Marshal(masked)
}

// Username 张三 --> 张*
// 王一一 -> 王*一
func Username(username string) string {
	runes := []rune(username)
	if length := len(runes); length >= 2 {
		i := make([]rune, 0, length)
		if length == 2 {
			i = append(i, runes[0], '*')
		} else {
			i = append(i, runes[0])
			for index := 0; index < length-2; index++ {
				i = append(i, '*')
			}
			i = append(i, runes[length-1])
		}
		return string(i)
	}
	return username
}

// Phone 18888888881 -> 188****88888
// +861888888888 -> +86 188****88888
func Phone(phone string) string {
	if len(phone) <= 4 {
		return phone
	}

	countryCode := ""
	phoneNumber := phone
	hasCountryCode := false
	if hasCountryCode = strings.HasPrefix(phone, "+"); hasCountryCode {
		countryCode = phone[:CountryCodeEndIndex]
		phoneNumber = phone[CountryCodeEndIndex:]
	}
	masked := ""
	if hasCountryCode {
		masked = fmt.Sprintf("%s %s****%s", countryCode, phoneNumber[:min(3, len(phoneNumber))], phoneNumber[max(0, len(phoneNumber)-4):])
	} else {
		masked = fmt.Sprintf("%s****%s", phoneNumber[:min(3, len(phoneNumber))], phoneNumber[max(0, len(phoneNumber)-4):])
	}
	return masked
}

func (p PhoneMask) MarshalJSON() ([]byte, error) {
	phone := Phone(string(p))
	return json.Marshal(phone)
}

func (e EmailMask) MarshalJSON() ([]byte, error) {
	emailMasked := Email(string(e))
	return json.Marshal(emailMasked)
}

func Email(email string) string {
	index := strings.LastIndex(email, "@")
	if index == -1 {
		return email
	}
	id := email[0:index]
	domain := email[index:]
	switch len(id) {
	case 1:
		return id
	case 2:
		id = id[:1] + "*"
	case 3:
		id = id[:1] + "*" + id[2:]
	default:
		// abcd -> a**d
		masked := strings.Repeat("*", len(id)-2)
		id = id[:1] + masked + id[len(id)-1:]
	}
	return id + domain
}
