package rdbkey

import "fmt"

var MailLogin = func(to string) string {
	return fmt.Sprintf("mail:%s:login", to)
}

var MailLoginLock = func() string {
	return "mail:login:locker"
}
