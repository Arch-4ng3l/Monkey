package object

type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{
		store: s,
		outer: nil,
	}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func extendFunctionEnv(fn *Function, args []Object) *Env {
	env := NewEnclosedEnv(fn.Env)

	for i, p := range fn.Params {
		env.Set(p.Value, args[i])
	}

	return env
}
