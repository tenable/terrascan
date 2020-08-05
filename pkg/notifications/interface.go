package notifications

// Notifier defines the interface which every type of notification provider
// needs to implement to claim support in terrascan
type Notifier interface {
	Init(interface{}) error
	SendNotification(interface{}) error
}
