package activity

type activityType string

const (
	MOUSE_CURSOR_MOVEMENT activityType = "cursor-move"
	MOUSE_LEFT_CLICK      activityType = "left-mouse-click"
	SCREEN_CHANGE         activityType = "screen-change"
)

type Type struct {
	ActivityType activityType
}
