package cmd

import (
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type DeleteVMCmd struct {
	ui         boshui.UI
	deployment boshdir.Deployment
}

func NewDeleteVMCmd(ui boshui.UI, deployment boshdir.Deployment) DeleteVMCmd {
	return DeleteVMCmd{ui: ui, deployment: deployment}
}

func (c DeleteVMCmd) Run(opts boshopts.DeleteVMOpts) error {
	err := c.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	return c.deployment.DeleteVM(opts.Args.CID)
}
