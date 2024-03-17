package common

// ActionRoute defines configuration options for actions handlers which should be as well http routes handlers.
// Invocations to such handlers, will contain RoutingContext from FastHttp instead of the message with payload.
type ActionRoute struct {
	Methods []string
	Path    string
}
