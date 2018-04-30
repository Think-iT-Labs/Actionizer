package notifier

type Notifier interface {
	SendNotification(s string)
}
