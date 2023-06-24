package types

import (
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
)

type RegistryType uint8

const (
	Undefined    RegistryType = 0
	Azure        RegistryType = 1
	DigitalOcean RegistryType = 2
	DockerHub    RegistryType = 3
	Google       RegistryType = 4
	GitHub       RegistryType = 5
	Gitlab       RegistryType = 6
	Private      RegistryType = 7
)

// GetRegistryType is a function that determines the underlying registry type
func GetRegistryType(secret koyeb.Secret) RegistryType {
	if secret.AzureContainerRegistry != nil {
		return Azure
	}
	if secret.DigitalOceanRegistry != nil {
		return DigitalOcean
	}
	if secret.DockerHubRegistry != nil {
		return DockerHub
	}
	if secret.GcpContainerRegistry != nil {
		return Google
	}
	if secret.GithubRegistry != nil {
		return GitHub
	}
	if secret.GitlabRegistry != nil {
		return Gitlab
	}
	if secret.PrivateRegistry != nil {
		return Private
	}
	return Undefined
}

func (x RegistryType) String() string {
	switch x {
	case Undefined:
		return ""
	case Azure:
		return "Azure"
	case DigitalOcean:
		return "DigitalOcean"
	case Google:
		return "GCR"
	case GitHub:
		return "GitHub"
	case Gitlab:
		return "Gitlab"
	case Private:
		return "private"
	default:
		return ""
	}
}
