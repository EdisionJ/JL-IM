package test

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"testing"
)

type Room struct {
	RoomID int64
	Users  map[int64]string
}

func TestCmap(t *testing.T) {
	var roomMap = cmap.New[*Room]()
	roomMap.Set("data_1", &Room{RoomID: 1, Users: map[int64]string{1: "value"}})
	data, _ := roomMap.Get("data_1")
	fmt.Println(data)
	data.Users[1] = "new_value"
	new_data, _ := roomMap.Get("data_1")
	fmt.Println(new_data)

}
