package config

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
)

type K8sConfig struct {
	kubeconfigPath string
	insecure       bool
}

func NewK8sConfig(path string, insecure bool) *K8sConfig {
	return &K8sConfig{
		kubeconfigPath: path,
		insecure:       insecure,
	}
}

func (kc *K8sConfig) k8sRestConfigDefaultOrDie(insecure bool) *rest.Config {

	// 获取本机默认kubeconfig   Linux： ~   /home/xxx
	//home, err := os.UserHomeDir()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defaultConfigPath := fmt.Sprintf("%s/.kube/config", home)

	config, err := clientcmd.BuildConfigFromFlags("", kc.kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = insecure
	return config
}

// initDynamicClientOrDie 初始化DynamicClient
func (kc *K8sConfig) initDynamicClientOrDie() dynamic.Interface {
	client, err := dynamic.NewForConfig(kc.k8sRestConfigDefaultOrDie(kc.insecure))
	if err != nil {
		klog.Fatal(err)
	}
	return client
}

// initClient 初始化 clientSet
func (kc *K8sConfig) initClientOrDie() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(kc.k8sRestConfigDefaultOrDie(kc.insecure))
	if err != nil {
		klog.Fatal(err)
	}
	return c
}

// NewRestMapperOrDie 获取 api group resource
func (kc *K8sConfig) NewRestMapperOrDie() *meta.RESTMapper {
	gr, err := restmapper.GetAPIGroupResources(kc.initClientOrDie().Discovery())
	if err != nil {
		log.Fatal(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	return &mapper
}

// InitWatchFactoryAndRestConfig 初始化 dynamic client informerFactory, restConfig
func (kc *K8sConfig) InitWatchFactoryAndRestConfig() (dynamicinformer.DynamicSharedInformerFactory, *rest.Config) {
	dynClient := kc.initDynamicClientOrDie()
	return dynamicinformer.NewDynamicSharedInformerFactory(dynClient, 0), kc.k8sRestConfigDefaultOrDie(kc.insecure)
}
