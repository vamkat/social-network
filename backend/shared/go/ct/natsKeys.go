package ct

import "fmt"

func PrivateMessageKey(receiverId any) string {
	return fmt.Sprintf("dm.%v", receiverId)
}

func GroupMessageKey(groupId any) string {
	return fmt.Sprintf("grm.%v", groupId)
}

func NotificationKey(receiverId any) string {
	return fmt.Sprintf("ntf.%v", receiverId)
}
