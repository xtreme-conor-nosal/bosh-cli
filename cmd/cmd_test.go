package cmd_test

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

var _ = Describe("Cmd", func() {
	var (
		ui     *fakeui.FakeUI
		confUI *boshui.ConfUI
		fs     *fakesys.FakeFileSystem
		cmd    Cmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		confUI = boshui.NewWrappingConfUI(ui, logger)

		fs = fakesys.NewFakeFileSystem()

		deps := NewBasicDeps(confUI, logger)
		deps.FS = fs

		cmd = NewCmd(boshopts.BoshOpts{}, nil, deps)
	})

	Describe("Execute", func() {
		It("succeeds executing at least one command", func() {
			cmd.Opts = &boshopts.InterpolateOpts{}

			err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			Expect(ui.Blocks).To(Equal([]string{"null\n"}))
		})

		It("prints message if specified", func() {
			cmd.Opts = &boshopts.MessageOpts{Message: "output"}

			err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			Expect(ui.Blocks).To(Equal([]string{"output"}))
		})

		It("allows to enable json output", func() {
			cmd.BoshOpts = boshopts.BoshOpts{JSONOpt: true}
			cmd.Opts = &boshopts.InterpolateOpts{}

			err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			confUI.Flush()

			Expect(ui.Blocks[0]).To(ContainSubstring(`Blocks": [`))
		})

		Describe("color", func() {
			executeCmdAndPrintTable := func() {
				err := cmd.Execute()
				Expect(err).ToNot(HaveOccurred())

				// Tables have emboldened header values
				confUI.PrintTable(boshtbl.Table{Header: []boshtbl.Header{boshtbl.NewHeader("State")}})
			}

			It("has color in the output enabled by default", func() {
				cmd.BoshOpts = boshopts.BoshOpts{}
				cmd.Opts = &boshopts.InterpolateOpts{}

				executeCmdAndPrintTable()

				// Expect that header values are bold
				Expect(ui.Tables[0].HeaderFormatFunc).ToNot(BeNil())
			})

			It("allows to disable color in the output", func() {
				cmd.BoshOpts = boshopts.BoshOpts{NoColorOpt: true}
				cmd.Opts = &boshopts.InterpolateOpts{}

				executeCmdAndPrintTable()

				// Expect that header values are empty because they were not emboldened
				Expect(ui.Tables[0].HeaderFormatFunc).To(BeNil())
			})
		})

		It("returns error if changing tmp root fails", func() {
			fs.ChangeTempRootErr = errors.New("fake-err")

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("fake-err"))
		})

		It("returns error for unknown commands", func() {
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unhandled command: <nil>"))
		})
	})
})
