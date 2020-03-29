package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"l0calh0st.cn/k8s-bridge/configure"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/controller/dns"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
	"l0calh0st.cn/k8s-bridge/pkg/controller/sentry"
	"os"
	"os/signal"
	"syscall"
)

var (
	masterUrl = flag.String("masterUrl", "","")
	kubeConfig = flag.String("kubeConfig", "", "")
)

var (
	globalConfig *configure.Config = configure.NewConfig()
	DEBUG bool = true
)


func main() {
	flag.Parse()

	logrus.Warn(globalConfig.Address)


	if DEBUG {logrus.SetLevel(logrus.DebugLevel)}

	logrus.Infoln("Initialize KubeConfig and kubeClientSet")
	restConfig, kubeClientSet, err := initializeRestConfigAndClientSet()
	if err!=nil{
		logrus.WithError(err).Fatal("Unable to build kubeClient and kubeConfig")
	}
	ctx,cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	kubeBridgeDispatchController := controller.NewDispatcher()

	logrus.Infoln("Starting the kubeResourceServiceController...")
	kubeResourceServiceController := kube_resource.NewKubeResourceServiceController(kubeClientSet, restConfig, kubeBridgeDispatchController)
	kubeBridgeDispatchController.RegisterController(kubeResourceServiceController)
	go runController(ctx, kubeResourceServiceController)


	logrus.Infoln("Starting the kubeResourcePodController...")
	kubeResourcePodController := kube_resource.NewKubeResourcePodController(kubeClientSet, restConfig, kubeBridgeDispatchController)
	kubeBridgeDispatchController.RegisterController(kubeResourcePodController)
	go runController(ctx, kubeResourcePodController)

	logrus.Infof("Start dns controller ......")
	dnsController := dns.NewKubeBridgeDnsController(kubeBridgeDispatchController)
	kubeBridgeDispatchController.RegisterController(dnsController)
	go runController(ctx, dnsController)

	logrus.Infof("Staer sentry controller......")
	syncController := sentry.NewKubeBridgeSentryController(kubeBridgeDispatchController)
	kubeBridgeDispatchController.RegisterController(syncController)
	go runController(ctx, syncController)


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
	if masterUrl != nil || kubeConfig != nil {
		restConfig, err = clientcmd.BuildConfigFromFlags(*masterUrl, *kubeConfig)
	} else {
		restConfig, err = rest.InClusterConfig()
	}

	if err!= nil {
		return
	}
	kubeClientSet,err = kubernetes.NewForConfig(restConfig)
	return
}