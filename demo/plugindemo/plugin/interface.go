package plugin

type Plugin interface {
	Worker(map[string]interface{})
}
