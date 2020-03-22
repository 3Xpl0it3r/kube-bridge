package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/controller/dns"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
	"os"
	"os/signal"
	"syscall"
)

var (
	masterUrl = flag.String("masterUrl", "","")
	kubeConfig = flag.String("kubeConfig", "", "")
)

var (
	DEBUG bool = true
)


func main() {
	flag.Parse()


	if DEBUG {logrus.SetLevel(logrus.DebugLevel)}

	logrus.Infoln("Initialize KubeConfig and kubeClientSet")
	restConfig, kubeClientSet, err := initializeRestConfigAndClientSet()
	if err!=nil{
		logrus.WithError(err).Fatal("Unable to build kubeClient and kubeConfig")
	}
	ctx,cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	kubeBridgeSyncController := controller.NewSynchronize()

	logrus.Infoln("Starting the kubeResourceServiceController...")
	kubeResourceServiceController := kube_resource.NewKubeResourceServiceController(kubeClientSet, restConfig, kubeBridgeSyncController)
	kubeBridgeSyncController.RegisterController(kubeResourceServiceController)
	go runController(ctx, kubeResourceServiceController)


	logrus.Infoln("Starting the kubeResourcePodController...")
	kubeResourcePodController := kube_resource.NewKubeResourcePodController(kubeClientSet, restConfig, kubeBridgeSyncController)
	kubeBridgeSyncController.RegisterController(kubeResourcePodController)
	go runController(ctx, kubeResourcePodController)

	logrus.Infof("Start dns controller ......")
	dnsController := dns.NewKubeBridgeDnsController(kubeBridgeSyncController)
	kubeBridgeSyncController.RegisterController(dnsController)
	go runController(ctx, dnsController)

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <- stopCh:
		os.Exit(-1)
	}
}

func runController(ctx context.Context, controller controller.Controller){
	if err := controller.Run(ctx);err !=nil {
		logrus.WithError(err).Fatalln("Unable to run the controller")
	}
}

func initializeRestConfigAndClientSet()(restConfig *rest.Config,kubeClientSet kubernetes.Interface,err error){

	restConfig, err = clientcmd.BuildConfigFromFlags(*masterUrl, *kubeConfig)
	if err!= nil {
		return
	}
	kubeClientSet,err = kubernetes.NewForConfig(restConfig)
	return
}