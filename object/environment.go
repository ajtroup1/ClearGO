package object

// Instantiates & returns a new instance of Environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Our environment struct contains the entire environment 'tool'
// Environment is just a fancy way to associate strings with objects
// For now, we can just use a hashmap to associate these
type Environment struct {
	store map[string]Object
}

// Simple getters and setters for manipulating environment vars
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
