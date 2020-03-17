package kube_resource

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

func PodRemoteCommandExec(clientSet kubernetes.Interface, restConf *rest.Config,pod *corev1.Pod)error{

	command := []string{"/bin/sh", "-c", "echo '114.114.114.114' >>  /etc/resolv.conf"}
	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme);err != nil{
		panic(err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)

	req.VersionedParams(&corev1.PodExecOptions{
		Command: command,
		Container: "",
		Stdin:false,
		Stdout: true,
		Stderr:true,
		TTY:false,
	}, parameterCodec)

	exec,err := remotecommand.NewSPDYExecutor(restConf, "POST", req.URL())
	fmt.Println(err)
	if err!= nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             nil,
		Stdout:           	os.Stdout,
		Stderr:            os.Stderr,
		Tty:               false,
	})

	fmt.Println("----------->", err)
	fmt.Println(err)
	return err
}
