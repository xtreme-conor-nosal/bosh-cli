package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type AttachDiskCmd struct {
	deployment boshdir.Deployment
}

func NewAttachDiskCmd(deployment boshdir.Deployment) AttachDiskCmd {
	return AttachDiskCmd{
		deployment: deployment,
	}
}

func (c AttachDiskCmd) Run(opts boshopts.AttachDiskOpts) error {
	return c.deployment.AttachDisk(opts.Args.Slug, opts.Args.DiskCID, opts.DiskProperties)
}
