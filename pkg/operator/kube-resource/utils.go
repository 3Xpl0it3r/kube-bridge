package kube_resource

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"l0calh0st.cn/k8s-bridge/configure"
	"os"
)


var kubeResourceConfig *configure.Config

func init() {
	kubeResourceConfig = configure.NewConfig()
	kubeResourceConfig.ReadFromFile()
}



func PodRemoteCommandExec(clientSet kubernetes.Interface, restConf *rest.Config,pod *corev1.Pod, cmd ...string)error{

	command := cmd
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
	if err!= nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             nil,
		Stdout:           	os.Stdout,
		Stderr:            os.Stderr,
		Tty:               false,
	})

	return err
}
