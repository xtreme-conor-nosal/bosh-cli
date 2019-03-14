package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshreldir "github.com/cloudfoundry/bosh-cli/releasedir"
)

type GenerateJobCmd struct {
	releaseDir boshreldir.ReleaseDir
}

func NewGenerateJobCmd(releaseDir boshreldir.ReleaseDir) GenerateJobCmd {
	return GenerateJobCmd{releaseDir: releaseDir}
}

func (c GenerateJobCmd) Run(opts boshopts.GenerateJobOpts) error {
	return c.releaseDir.GenerateJob(opts.Args.Name)
}
