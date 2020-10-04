package cash

import (
	"time"
	"users/tables"
)

var stor = make(map[int]tables.User)

func Find(id int) (tables.User, bool) {
	var user tables.User
	result, ok := stor[id]
	if ok {
		user = result
	}
	return user, ok
}

func InsertAndDEl(user tables.User) {
	stor[user.Id] = user
	time.Sleep(5000 * time.Millisecond)
	delete(stor, user.Id)
}
