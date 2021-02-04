package bootstrap

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
)

// verifyBootstrapInPlace validate the number of control plane replica is one and that installation disk is set
func verifyBootstrapInPlace(installConfig *types.InstallConfig) error {
	if *installConfig.ControlPlane.Replicas != 1 {
		return errors.Errorf("bootstrap in place require single control plane replica, current value: %d", *installConfig.ControlPlane.Replicas)
	}
	if installConfig.BootstrapInPlace == nil {
		return errors.Errorf("Missing bootstrap in place configuration")
	}
	if installConfig.BootstrapInPlace.InstallationDisk == "" {
		return errors.Errorf("bootstrap in place requires installation disk configuration")
	}
	logrus.Warnf("Creating single node bootstrap in place configuration")
	return nil
}

