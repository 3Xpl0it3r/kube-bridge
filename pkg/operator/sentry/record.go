package sentry

import (
	"context"
	"fmt"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/sentry/proto"
	"sync"
)

type recordCycleServer struct {
	mu *sync.Mutex
	cache map[string]string
	workQueue controller.EventsHook
}

func newEventCycleServer(workQueue controller.EventsHook)*recordCycleServer{
	return &recordCycleServer{
		mu:    &sync.Mutex{},
		cache: nil,
		workQueue: workQueue,
	}
}

func(s *recordCycleServer)OnAdd(ctx context.Context,req *proto.RecordRequest)(resp *proto.RecordResponse,err error){
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache != nil {
		s.cache[req.Name] = "192.168.1.1"
		s.workQueue.OnAdd(map[string]string{req.Name: "192.168.1.1"})
		return &proto.RecordResponse{
			Ok:                   true,
		}, nil
	} else {
		return &proto.RecordResponse{
			Ok:                   false,
		}, fmt.Errorf("cache is not initialize")
	}
}

func(s *recordCycleServer)OnUpdate(ctx context.Context,req *proto.RecordRequest)(resp *proto.RecordResponse,err error){
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache != nil {
		s.cache[req.Name] = "192.168.1.1"
		s.workQueue.OnUpdate(map[string]string{req.Name: "192.168.1.1"})
		return &proto.RecordResponse{
			Ok:                   true,
		}, nil
	} else {
		return &proto.RecordResponse{
			Ok:                   false,
		}, fmt.Errorf("cache is not initialize")
	}
}

func(s *recordCycleServer)OnDelete(ctx context.Context,req *proto.RecordRequest)(resp *proto.RecordResponse,err error){
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache == nil {
		s.cache[req.Name] = "192.168.1.1"
		return &proto.RecordResponse{
			Ok:                   true,
		}, nil
	}
	delete(s.cache, req.Name)
	s.workQueue.OnDelete(map[string]string{req.Name: "192.168.1.1"})
	return &proto.RecordResponse{Ok: true}, nil
}
