package weeb

type Container struct {
	singletons map[string]interface{}
}

func NewContainer() *Container {
	return &Container{singletons: map[string]interface{}{}}
}

func (c *Container) Get(name string) interface{} {
	return c.singletons[name]
}

func (c *Container) Singleton(name string, value interface{}) {
	c.singletons[name] = value
}
