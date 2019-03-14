package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshreldir "github.com/cloudfoundry/bosh-cli/releasedir"
)

type GeneratePackageCmd struct {
	releaseDir boshreldir.ReleaseDir
}

func NewGeneratePackageCmd(releaseDir boshreldir.ReleaseDir) GeneratePackageCmd {
	return GeneratePackageCmd{releaseDir: releaseDir}
}

func (c GeneratePackageCmd) Run(opts boshopts.GeneratePackageOpts) error {
	return c.releaseDir.GeneratePackage(opts.Args.Name)
}
