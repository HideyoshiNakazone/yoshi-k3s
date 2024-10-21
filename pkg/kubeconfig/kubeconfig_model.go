package kubeconfig

import "gopkg.in/yaml.v3"

type KubeconfigModel struct {
	ApiVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
	Kind           string   `yaml:"kind"`
	Preferences    struct{} `yaml:"preferences"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func NewKubeconfigModel(kubeconfigData *[]byte) *KubeconfigModel {
	var kubeconfig = new(KubeconfigModel)
	err := yaml.Unmarshal(*kubeconfigData, kubeconfig)
	if err != nil {
		return nil
	}
	return kubeconfig
}
