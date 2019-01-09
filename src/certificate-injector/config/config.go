package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const ImportCertificatePs = `
$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }
$certFile=[System.IO.Path]::GetTempFileName()
$decodedCertData = [Convert]::FromBase64String("%s")
[IO.File]::WriteAllBytes($certFile, $decodedCertData)
Import-Certificate -CertStoreLocation Cert:\\LocalMachine\Root -FilePath $certFile
Remove-Item $certFile
`

type Container struct{}

type ContainerJSON struct {
	Process process `json:"process"`
}

type process struct {
	Args []string `json:"args"`
	Cwd  string   `json:"cwd"`
}

func NewContainer() Container {
	return Container{}
}

// Creates a powershell script to write the certs
// to a file and import the certificate. It appends
// this script as a process to a config.json that will
// be run on the container.
func (c Container) Write(certData []byte) error {
	command := fmt.Sprintf(ImportCertificatePs, string(certData))

	encodedCommand := base64.StdEncoding.EncodeToString([]byte(command))

	containerJSON := ContainerJSON{
		Process: process{
			Args: []string{"powershell.exe", "-EncodedCommand", encodedCommand},
			Cwd:  `C:\`,
		},
	}

	marshalledContainerJSON, err := json.Marshal(containerJSON)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("config.json", marshalledContainerJSON, 0644)
	if err != nil {
		return fmt.Errorf("Write config.json failed: %s", err)
	}

	return nil
}
