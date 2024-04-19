package event

import (
	"github.com/satori-protocol-go/satori-model-go/pkg/channel"
	"github.com/satori-protocol-go/satori-model-go/pkg/guild"
	"github.com/satori-protocol-go/satori-model-go/pkg/guildmember"
	"github.com/satori-protocol-go/satori-model-go/pkg/guildrole"
	"github.com/satori-protocol-go/satori-model-go/pkg/interaction"
	"github.com/satori-protocol-go/satori-model-go/pkg/login"
	"github.com/satori-protocol-go/satori-model-go/pkg/message"
	"github.com/satori-protocol-go/satori-model-go/pkg/user"
)

type EventType string

const (
	// 群组事件

	// 加入群组时触发。必需资源：guild。
	EventTypeGuildAdded EventType = "guild-added"
	// 群组被修改时触发。必需资源：guild。
	EventTypeGuildUpdated EventType = "guild-updated"
	// 退出群组时触发。必需资源：guild。
	EventTypeGuildRemoved EventType = "guild-removed"
	// 接收到新的入群邀请时触发。必需资源：guild。
	EventTypeGuildRequest EventType = "guild-request"

	// 群组成员事件

	// 群组成员增加时触发。必需资源：guild，member，user。
	EventTypeGuildMemberAdded EventType = "guild-member-added"
	// 群组成员信息更新时触发。必需资源：guild，member，user。
	EventTypeGuildMemberUpdated EventType = "guild-member-updated"
	// 群组成员移除时触发。必需资源：guild，member，user。
	EventTypeGuildMemberRemoved EventType = "guild-member-removed"
	// 接收到新的加群请求时触发。必需资源：guild，member，user。
	EventTypeGuildMemberRequest EventType = "guild-member-request"

	// 群组角色事件

	// 群组角色被创建时触发。必需资源：guild，role。
	EventTypeGuildRoleCreated EventType = "guild-role-created"
	// 群组角色被修改时触发。必需资源：guild，role。
	EventTypeGuildRoleUpdated EventType = "guild-role-updated"
	// 群组角色被删除时触发。必需资源：guild，role。
	EventTypeGuildRoleDeleted EventType = "guild-role-deleted"

	// 交互事件

	// 类型为 action 的按钮被点击时触发。必需资源：button。
	EventTypeInteractionButton EventType = "interaction/button"
	// 调用斜线指令时触发。资源 argv 或 message 中至少包含其一。
	EventTypeInteractionCommand EventType = "interaction/command"

	// 登录信息事件

	// 登录被创建时触发。必需资源：login。
	EventTypeLoginAdded EventType = "login-added"
	// 登录被删除时触发。必需资源：login。
	EventTypeLoginRemoved EventType = "login-removed"
	// 登录信息更新时触发。必需资源：login。
	EventTypeLoginUpdated EventType = "login-updated"

	// 消息事件

	// 当消息被创建时触发。必需资源：channel，message，user。
	EventTypeMessageCreated EventType = "message-created"
	// 当消息被编辑时触发。必需资源：channel，message，user。
	EventTypeMessageUpdated EventType = "message-updated"
	// 当消息被删除时触发。必需资源：channel，message，user。
	EventTypeMessageDeleted EventType = "message-deleted"

	// 表态事件

	// 当表态被添加时触发。
	EventTypeReactionAdded EventType = "reaction-added"
	// 当表态被移除时触发。
	EventTypeReactionRemoved EventType = "reaction-removed"

	// 用户事件

	// 接收到新的好友申请时触发。必需资源：user。
	EventTypeFriendRequest EventType = "friend-request"
)

type Event struct {
	Id        int64                    `json:"id"`                 // 事件 ID
	Type      EventType                `json:"type"`               // 事件类型
	Platform  string                   `json:"platform"`           // 接收者的平台名称
	SelfId    string                   `json:"self_id"`            // 接收者的平台账号
	Timestamp int64                    `json:"timestamp"`          // 事件的时间戳
	Argv      *interaction.Argv        `json:"argv,omitempty"`     // 交互指令
	Button    *interaction.Button      `json:"button,omitempty"`   // 交互按钮
	Channel   *channel.Channel         `json:"channel,omitempty"`  // 事件所属的频道
	Guild     *guild.Guild             `json:"guild,omitempty"`    // 事件所属的群组
	Login     *login.Login             `json:"login,omitempty"`    // 事件的登录信息
	Member    *guildmember.GuildMember `json:"member,omitempty"`   // 事件的目标成员
	Message   *message.Message         `json:"message,omitempty"`  // 事件的消息
	Operator  *user.User               `json:"operator,omitempty"` // 事件的操作者
	Role      *guildrole.GuildRole     `json:"role,omitempty"`     // 事件的目标角色
	User      *user.User               `json:"user,omitempty"`     // 事件的目标用户
}
