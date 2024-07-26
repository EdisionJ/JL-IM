package websocketSereve

import (
	"IM/globle"
	"IM/service/enum"
	"IM/utils"
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Method != "GET" {
			return false
		}
		return true
	},
}

var UserQ = globle.Db.User
var GroupMemberQ = globle.Db.GroupMember

func Connect(w http.ResponseWriter, r *http.Request) {
	strToken := r.Header.Get("Token")
	Token, err := utils.ParseToken(strToken)
	if err != nil {
		globle.Logger.Warnln("Token解析失败： ", err)
		return
	}
	uid := Token.UID
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		globle.Logger.Warnln("wesocket升级失败: ", err)
		return
	}

	//连接成功
	user := User{
		ID:      uid,
		Rooms:   make([]int64, 0),
		Channel: conn,
	}
	ctx := context.Background()
	rooms, err := GroupMemberQ.WithContext(ctx).
		Where(GroupMemberQ.ID.Eq(uid), GroupMemberQ.Type.Eq(enum.RoomTypeGroup), GroupMemberQ.State.NotLike(enum.Exited)).
		Select(GroupMemberQ.GroupID).Find()
	if err != nil {
		globle.Logger.Warnln("群聊查询失败: ", err)
		return
	}
	//添加聊天室
	for _, r := range rooms {
		JoinRoom(r.GroupID, uid, conn)
		user.Rooms = append(user.Rooms, r.GroupID)
	}
	//加入用户列表
	UserMapSet(&user)
	var isDone chan struct{}
	defer close(isDone)
	go ping(conn, user.ID, isDone)
	go pong(conn, user.ID, isDone)
	<-isDone

}

func SendMessage(target int, id int64, msg []byte) error {
	switch target {
	case PrivateMsg:
		user, ok := UserMapGet(id)
		if !ok {
			//不在线
			return nil
		}
		err := user.Channel.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			globle.Logger.Errorf("消息发送失败： ", err)
			return err
		}
	case GroupMsg:

	}
	return nil
}

func ping(conn *websocket.Conn, uid int64, isDone chan struct{}) {
	ticker := time.NewTicker(HeartRate)
	defer ticker.Stop()

	for {
		select {
		case <-isDone:
			return
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, []byte("PING"))
			if err != nil {
				disConnect(uid)
				globle.Logger.Println("发送心跳失败： ", err)
				isDone <- struct{}{}
				return
			}
		}
	}
}

func pong(conn *websocket.Conn, uid int64, isDone chan struct{}) {
	for {
		select {
		case <-isDone:
			return
		default:
			msgType, _, err := conn.ReadMessage()
			if err != nil {
				disConnect(uid)
				log.Println("heartbeat response error:", err)
				isDone <- struct{}{}
				return
			}
			if msgType == websocket.PongMessage {
				// 收到Pong消息后重设连接的读取截止时间
				conn.SetReadDeadline(time.Now().Add(DeadLine)) // 假设HeartRate是心跳间隔
			}
		}
	}
}

func disConnect(uid int64) {
	user, _ := UserMapGet(uid)
	// 关闭websocket连接
	user.Channel.Close()
	//从用户列表中移除
	UserMapDel(uid)
	//退出群聊用户列表
	for _, rId := range user.Rooms {
		ExitFromRoom(rId, uid)
	}
}
