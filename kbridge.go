package main


import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
	"os"
	"os/signal"
	"syscall"
)

var (
	masterUrl = flag.String("masterUrl", "","")
	kubeConfig = flag.String("kubeConfig", "", "")
)

var (
	kubeClientSet *kubernetes.Clientset
)

func main() {
	flag.Parse()
	if err := initializeClientSets();err != nil{
		return
	}
	ctx,cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	kubeResourceServiceController := kube_resource.NewKubeResourceServiceController(kubeClientSet)
	go runController(ctx, kubeResourceServiceController)

	kubeResourcePodController := kube_resource.NewKubeResourcePodController(kubeClientSet)
	go runController(ctx, kubeResourcePodController)

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <- stopCh:
		os.Exit(-1)
	}
}

func runController(ctx context.Context, controller controller.Controller){
	if err := controller.Run(ctx);err !=nil {
		fmt.Println("controller run failed ", err.Error())
	} else {
		fmt.Printf("controller running \n")
	}
}

func initializeClientSets()error{
	var err error
	cfg,err := clientcmd.BuildConfigFromFlags(*masterUrl, *kubeConfig)
	if err != nil{
		return err
	}
	kubeClientSet,err = kubernetes.NewForConfig(cfg)
	if err != nil{
		return err
	}
	return nil

}