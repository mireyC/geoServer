package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DockerImage     DockerImage
	GeoServerConfig GeoServerConfig
}

type DockerImage struct {
	ImageName      string `json:"imageName"`
	ImageTag       string `json:"imageTag"`
	DockerHostPort string `json:"dockerHostPort"`
	Source         string `json:"source"`
	Target         string `json:"target"`
}

type GeoServerConfig struct {
	GeoServerURL string `json:"geoServerURL"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Workspace    string `json:"workspace"`
	StoreType    string `json:"storeType"`
	//FileURL      string `json:"fileURL"`
	//ZipFilePath  string `json:"zipFilePath"`
	//UploadZIPURL string `json:"uploadZIPURL"`
}
