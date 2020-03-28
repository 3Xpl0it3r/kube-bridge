package sentry

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/sentry/proto"
	"net"
	"time"
)

type Operator interface {
	Run(ctx context.Context)error
	OnAdd(object interface{}) error
	OnUpdate(object interface{})error
	OnDelete(object interface{})error
}

type realSentryOperator struct {
	server *recordCycleServer
	workQueue controller.EventsHook
}

func NewRealSyncOperator(workQueue controller.EventsHook)*realSentryOperator{
	return &realSentryOperator{server: newEventCycleServer(workQueue), workQueue: workQueue}
}

func(s *realSentryOperator)Run(ctx context.Context)error{
	ctx,cancel  := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := s.runServer(ctx);err != nil {
			logrus.WithError(err)
		}
	}()

	<- ctx.Done()
	return ctx.Err()
}

func(s *realSentryOperator)runServer(ctx context.Context)error{
	chErr := make(chan error)
	defer close(chErr)
	lis,err := net.Listen("tcp", ":9999")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	proto.RegisterRecordCycleServiceServer(server, s.server)
	select {
	case chErr<- server.Serve(lis):
		return fmt.Errorf("%s",chErr)
	case <-ctx.Done():
		return ctx.Err()
	}
}



func(s *realSentryOperator)OnAdd(object interface{}) error{
	conn,err := grpc.Dial("", grpc.WithInsecure())
	defer conn.Close()
	c := proto.NewRecordCycleServiceClient(conn)
	err = onAddRecordRequest(c, &proto.RecordRequest{
		Name:                 object.(string),
	})
	return err
}

func(s *realSentryOperator)OnDelete(object interface{})error{
	conn,err := grpc.Dial("", grpc.WithInsecure())
	defer conn.Close()
	c := proto.NewRecordCycleServiceClient(conn)
	err = onDeleteRecordRequest(c, &proto.RecordRequest{Name: object.(string)})
	return err

}
func(s *realSentryOperator)OnUpdate(object interface{})error{

	conn,err := grpc.Dial("", grpc.WithInsecure())
	defer conn.Close()
	c := proto.NewRecordCycleServiceClient(conn)
	err = onUpdateRecordRequest(c, &proto.RecordRequest{Name: object.(string)})
	return err
}



func onAddRecordRequest(client proto.RecordCycleServiceClient, req *proto.RecordRequest)error{
	ctx,cancel := context.WithTimeout(context.Background(), time.Duration(5 * time.Second))
	defer cancel()
	_,err := client.OnAdd(ctx, req)
	return err
}
func onUpdateRecordRequest(client proto.RecordCycleServiceClient, req *proto.RecordRequest)error{
	ctx,cancel := context.WithTimeout(context.Background(), time.Duration(5 * time.Second))
	defer cancel()
	_,err := client.OnUpdate(ctx, req)
	return err
}
func onDeleteRecordRequest(client proto.RecordCycleServiceClient, req *proto.RecordRequest)error{
	ctx,cancel := context.WithTimeout(context.Background(), time.Duration(5 * time.Second))
	defer cancel()
	_,err := client.OnDelete(ctx, req)
	return err
}

