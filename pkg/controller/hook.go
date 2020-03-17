package controller


type Hook interface {
	OnAdd(object interface{})
	OnUpdate(object interface{})
	OnDelete(object interface{})
}


type HookManager struct {
	hooks []Hook
}

func NewHookManager()HookManager{
	return HookManager{hooks: []Hook{}}
}

func(h *HookManager)GetHooks()[]Hook{
	return h.hooks
}

func(h *HookManager)AddHook(hook Hook)error{
	return nil
}

func(h *HookManager)RemoveHook(hook Hook)error{
	return nil
}
