package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type DeleteSnapshotCmd struct {
	ui         boshui.UI
	deployment boshdir.Deployment
}

func NewDeleteSnapshotCmd(ui boshui.UI, deployment boshdir.Deployment) DeleteSnapshotCmd {
	return DeleteSnapshotCmd{ui: ui, deployment: deployment}
}

func (c DeleteSnapshotCmd) Run(opts boshopts.DeleteSnapshotOpts) error {
	err := c.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	return c.deployment.DeleteSnapshot(opts.Args.CID)
}
