package client

import (
	yaml "gopkg.in/yaml.v2"

	"github.com/dchest/uniuri"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
)

const (
	TILLER_HOST = "tiller-deploy.kube-system.svc.cluster.local:44134"
	CHART_PATH  = "/mariadb-0.6.1.tgz"
)

func Install(releaseName, namespace string) error {
	vals, err := yaml.Marshal(map[string]interface{}{
		"mariadbRootPassword": uniuri.New(),
		"mariadbDatabase":     "dbname",
	})
	if err != nil {
		return err
	}
	helmClient := helm.NewClient(helm.Host(TILLER_HOST))
	_, err = helmClient.InstallRelease(CHART_PATH, namespace, helm.ReleaseName(releaseName), helm.ValueOverrides(vals))
	if err != nil {
		return err
	}
	return nil
}

func Delete(releaseName string) error {
	helmClient := helm.NewClient(helm.Host(TILLER_HOST))
	if _, err := helmClient.DeleteRelease(releaseName); err != nil {
		return err
	}
	return nil
}

func GetPassword(releaseName, namespace string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	secret, err := clientset.Core().Secrets(namespace).Get(releaseName + "-mariadb")
	if err != nil {
		return "", err
	}
	return string(secret.Data["mariadb-root-password"]), nil
}
