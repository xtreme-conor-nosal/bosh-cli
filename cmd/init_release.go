package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshreldir "github.com/cloudfoundry/bosh-cli/releasedir"
)

type InitReleaseCmd struct {
	releaseDir boshreldir.ReleaseDir
}

func NewInitReleaseCmd(releaseDir boshreldir.ReleaseDir) InitReleaseCmd {
	return InitReleaseCmd{releaseDir: releaseDir}
}

func (c InitReleaseCmd) Run(opts boshopts.InitReleaseOpts) error {
	return c.releaseDir.Init(opts.Git)
}
