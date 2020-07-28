package aio

import (
	"encoding/json"
	"fmt"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

const (
	aioIgnFilename = "aio.ign"
)

// aioTemplateData is the data to use to replace values in all-in-one
// template files.
type aioTemplateData struct {
	ReleaseImage string
}

// AIO is an asset that generates the ignition config for an all-in-one node.
type AIO struct {
	Config *igntypes.Config
	File   *asset.File
}

var _ asset.WritableAsset = (*AIO)(nil)

// Dependencies returns the assets on which the AIO asset depends.
func (a *AIO) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.EtcdCABundle{},
		&tls.EtcdMetricCABundle{},
		&tls.EtcdMetricSignerCertKey{},
		&tls.EtcdMetricSignerClientCertKey{},
		&tls.EtcdSignerCertKey{},
		&tls.EtcdSignerClientCertKey{},
		&releaseimage.Image{},
	}
}

// Generate generates the ignition config for the all-in-one asset.
func (a *AIO) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	releaseImage := &releaseimage.Image{}
	dependencies.Get(installConfig, releaseImage)

	templateData, err := a.getTemplateData(installConfig.Config, releaseImage.PullSpec)

	if err != nil {
		return errors.Wrap(err, "failed to get bootstrap templates")
	}

	a.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	err = ignition.AddStorageFiles(a.Config, "/", "aio/files", templateData)
	if err != nil {
		return err
	}

	enabled := map[string]struct{}{
		"kubelet.service": {},
		"aiokube.service": {},
	}

	err = ignition.AddSystemdUnits(a.Config, "aio/systemd/units", templateData, enabled)
	if err != nil {
		return err
	}

	a.Config.Passwd.Users = append(
		a.Config.Passwd.Users,
		igntypes.PasswdUser{Name: "core", SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
			igntypes.SSHAuthorizedKey(installConfig.Config.SSHKey),
		}},
	)

	data, err := json.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	a.File = &asset.File{
		Filename: aioIgnFilename,
		Data:     data,
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *AIO) Name() string {
	return "All-in-one Ignition Config"
}

// Files returns the files generated by the asset.
func (a *AIO) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// getTemplateData returns the data to use to execute all-in-one templates.
func (a *AIO) getTemplateData(installConfig *types.InstallConfig, releaseImage string) (*aioTemplateData, error) {
	if *installConfig.ControlPlane.Replicas != 1 {
		return nil, fmt.Errorf("All-in-one configurations must use a single control plane replica")
	}

	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil && *worker.Replicas != 0 {
			return nil, fmt.Errorf("All-in-one configurations do not support compute replicas")
		}
	}

	return &aioTemplateData{
		ReleaseImage: releaseImage,
	}, nil
}

// Load returns the all-in-one ignition from disk.
func (a *AIO) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(aioIgnFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", aioIgnFilename)
	}

	a.File, a.Config = file, config
	return true, nil
}
