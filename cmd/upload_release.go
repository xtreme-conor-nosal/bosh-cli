package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	semver "github.com/cppforlife/go-semi-semantic/version"

	boshopts "github.com/cloudfoundry/bosh-cli/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshrel "github.com/cloudfoundry/bosh-cli/release"
	boshreldir "github.com/cloudfoundry/bosh-cli/releasedir"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type UploadReleaseCmd struct {
	releaseDirFactory    func(boshopts.DirOrCWDArg) (boshrel.Reader, boshreldir.ReleaseDir)
	releaseArchiveWriter boshrel.Writer

	director              boshdir.Director
	releaseArchiveFactory func(string) boshdir.ReleaseArchive

	cmdRunner boshsys.CmdRunner
	fs        boshsys.FileSystem
	ui        boshui.UI
}

func NewUploadReleaseCmd(
	releaseDirFactory func(boshopts.DirOrCWDArg) (boshrel.Reader, boshreldir.ReleaseDir),
	releaseArchiveWriter boshrel.Writer,
	director boshdir.Director,
	releaseArchiveFactory func(string) boshdir.ReleaseArchive,
	cmdRunner boshsys.CmdRunner,
	fs boshsys.FileSystem,
	ui boshui.UI,
) UploadReleaseCmd {
	return UploadReleaseCmd{
		releaseDirFactory:    releaseDirFactory,
		releaseArchiveWriter: releaseArchiveWriter,

		director:              director,
		releaseArchiveFactory: releaseArchiveFactory,

		cmdRunner: cmdRunner,
		fs:        fs,
		ui:        ui,
	}
}

func (c UploadReleaseCmd) Run(opts boshopts.UploadReleaseOpts) error {
	switch {
	case opts.Release != nil:
		return c.uploadRelease(opts.Release, opts)
	case opts.Args.URL.IsRemote():
		return c.uploadIfNecessary(opts, c.uploadRemote)
	case opts.Args.URL.IsGit():
		return c.uploadIfNecessary(opts, c.uploadGit)
	default:
		return c.uploadFile(opts)
	}
}

func (c UploadReleaseCmd) uploadRemote(opts boshopts.UploadReleaseOpts) error {
	return c.director.UploadReleaseURL(string(opts.Args.URL), opts.SHA1, opts.Rebase, opts.Fix)
}

func (c UploadReleaseCmd) uploadGit(opts boshopts.UploadReleaseOpts) error {
	repoPath, err := c.fs.TempDir("bosh-upload-release-git-clone")
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating tmp dir for git cloning")
	}

	defer c.fs.RemoveAll(repoPath)

	_, _, _, err = c.cmdRunner.RunCommand("git", "clone", opts.Args.URL.GitRepo(), "--depth", "1", repoPath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Cloning git repo")
	}

	newOpts := boshopts.UploadReleaseOpts{
		Directory: boshopts.DirOrCWDArg{Path: repoPath},
		Name:      opts.Name,
		Version:   opts.Version,
		Fix:       opts.Fix,
	}

	return c.uploadFile(newOpts)
}

func (c UploadReleaseCmd) uploadFile(opts boshopts.UploadReleaseOpts) error {
	if c.releaseDirFactory == nil {
		return bosherr.Errorf("Cannot upload non-remote release")
	}

	releaseReader, releaseDir := c.releaseDirFactory(opts.Directory)

	var release boshrel.Release
	var err error

	path := opts.Args.URL.FilePath()

	if len(path) > 0 {
		release, err = releaseReader.Read(path)
		if err != nil {
			return err
		}
	} else {
		release, err = releaseDir.FindRelease(opts.Name, semver.Version(opts.Version))
		if err != nil {
			return err
		}
	}

	defer release.CleanUp()

	return c.uploadRelease(release, opts)
}

func (c UploadReleaseCmd) uploadRelease(release boshrel.Release, opts boshopts.UploadReleaseOpts) error {
	var pkgFpsToSkip []string
	var err error

	if !opts.Fix {
		pkgFpsToSkip, err = c.director.MatchPackages(release.Manifest(), release.IsCompiled())
		if err != nil {
			return err
		}
	}

	path, err := c.releaseArchiveWriter.Write(release, pkgFpsToSkip)
	if err != nil {
		return err
	}

	defer c.fs.RemoveAll(path)

	file, err := c.releaseArchiveFactory(path).File()
	if err != nil {
		return bosherr.WrapErrorf(err, "Opening release")
	}

	return c.director.UploadReleaseFile(file, opts.Rebase, opts.Fix)
}

func (c UploadReleaseCmd) uploadIfNecessary(opts boshopts.UploadReleaseOpts, uploadFunc func(boshopts.UploadReleaseOpts) error) error {
	necessary, err := c.needToUpload(opts)
	if err != nil || !necessary {
		return err
	}

	return uploadFunc(opts)
}

func (c UploadReleaseCmd) needToUpload(opts boshopts.UploadReleaseOpts) (bool, error) {
	if opts.Fix {
		return true, nil
	}

	version := semver.Version(opts.Version).AsString()

	found, err := c.director.HasRelease(opts.Name, version, opts.Stemcell)
	if err != nil {
		return true, err
	}

	if found {
		if opts.Stemcell.IsProvided() {
			c.ui.PrintLinef("Release '%s/%s' for stemcell '%s' already exists.", opts.Name, version, opts.Stemcell)
		} else {
			c.ui.PrintLinef("Release '%s/%s' already exists.", opts.Name, version)
		}

		return false, nil
	}

	return true, nil
}
