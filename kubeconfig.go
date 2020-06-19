package fetchkube

import (
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Encode(conf *api.Config) ([]byte, error) {
	return clientcmd.Write(*conf)
}

func Decode(conf []byte) (*api.Config, error) {
	return clientcmd.Load(conf)
}

func StartingConfig() (*api.Config, error) {
	conf, err := clientcmd.DefaultClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func Contexts(conf *api.Config) ([]string, error) {
	ctxs := make([]string, 0, len(conf.Contexts))
	for ctx := range conf.Contexts {
		ctxs = append(ctxs, ctx)
	}
	return ctxs, nil
}

func FetchOnlyDefault(conf *api.Config) (*api.Config, error) {
	newConf := api.NewConfig()
	if conf.CurrentContext == "" {
		return nil, fmt.Errorf("context cannot be empty")
	}
	ctx, ok := conf.Contexts[conf.CurrentContext]
	if !ok {
		return nil, fmt.Errorf("not found context %q", conf.CurrentContext)
	}
	newConf.Contexts[conf.CurrentContext] = ctx
	cluster, ok := conf.Clusters[ctx.Cluster]
	if !ok {
		return nil, fmt.Errorf("not found cluster %q", ctx.Cluster)
	}
	newConf.Clusters[ctx.Cluster] = cluster
	auth, ok := conf.AuthInfos[ctx.AuthInfo]
	if !ok {
		return nil, fmt.Errorf("not found users %q", ctx.AuthInfo)
	}
	newConf.AuthInfos[ctx.AuthInfo] = auth
	return newConf, nil
}

func ResolveLocalPaths(conf *api.Config) (*api.Config, error) {
	err := clientcmd.ResolveLocalPaths(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func InlineData(conf *api.Config) (*api.Config, error) {
	for _, cluster := range conf.Clusters {
		err := loadFile(&cluster.CertificateAuthorityData, &cluster.CertificateAuthority)
		if err != nil {
			return nil, err
		}
	}

	for _, auth := range conf.AuthInfos {
		err := loadFile(&auth.ClientCertificateData, &auth.ClientCertificate)
		if err != nil {
			return nil, err
		}
		err = loadFile(&auth.ClientKeyData, &auth.ClientKey)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func loadFile(data *[]byte, filepath *string) error {
	if *filepath == "" {
		return nil
	}
	defer func() {
		*filepath = ""
	}()
	var err error
	if len(*data) != 0 {
		return nil
	}
	*data, err = ioutil.ReadFile(*filepath)
	if err != nil {
		return err
	}
	return nil
}
