package models

/**
    @date: 2022/10/19
**/

type UserRoom struct {
	UserIdentity string `json:"user_identity"`
	RoomIdentity string `json:"room_identity"`
	RoomType     int    `json:"room_type"` //房间类型[1.单独房间 2.群聊房间]
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (UserRoom) CollectionNameRoom() string {
	return "user_room"
}
