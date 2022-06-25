package fedora

import (
	"fmt"
	"math/rand"
	"path"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	pipeline "github.com/osbuild/osbuild-composer/internal/distro/pipelines"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild2"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

func qcow2Pipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	partitionTable, err := t.getPartitionTable(customizations.GetFilesystems(), options, rng)
	if err != nil {
		return nil, err
	}

	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, partitionTable)
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, treePipeline.Serialize())

	diskfile := "disk.img"
	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(packageSetSpecs[osPkgsKey], customizations.GetKernel().Name)
	imagePipeline := liveImagePipeline(&buildPipeline, &treePipeline, diskfile, partitionTable, t.arch, kernelVer)
	pipelines = append(pipelines, imagePipeline.Serialize())

	qemuPipeline := qemuPipeline(&buildPipeline, &imagePipeline, diskfile, t.filename, osbuild.QEMUFormatQCOW2, osbuild.QCOW2Options{Compat: "1.1"})
	pipelines = append(pipelines, qemuPipeline.Serialize())

	return pipelines, nil
}

func vhdPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	partitionTable, err := t.getPartitionTable(customizations.GetFilesystems(), options, rng)
	if err != nil {
		return nil, err
	}

	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, partitionTable)
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, treePipeline.Serialize())

	diskfile := "disk.img"
	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(packageSetSpecs[osPkgsKey], customizations.GetKernel().Name)
	imagePipeline := liveImagePipeline(&buildPipeline, &treePipeline, diskfile, partitionTable, t.arch, kernelVer)
	pipelines = append(pipelines, imagePipeline.Serialize())

	qemuPipeline := qemuPipeline(&buildPipeline, &imagePipeline, diskfile, t.filename, osbuild.QEMUFormatVPC, nil)
	pipelines = append(pipelines, qemuPipeline.Serialize())
	return pipelines, nil
}

func vmdkPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	partitionTable, err := t.getPartitionTable(customizations.GetFilesystems(), options, rng)
	if err != nil {
		return nil, err
	}

	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, partitionTable)
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, treePipeline.Serialize())

	diskfile := "disk.img"
	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(packageSetSpecs[osPkgsKey], customizations.GetKernel().Name)
	imagePipeline := liveImagePipeline(&buildPipeline, &treePipeline, diskfile, partitionTable, t.arch, kernelVer)
	pipelines = append(pipelines, imagePipeline.Serialize())

	qemuPipeline := qemuPipeline(&buildPipeline, &imagePipeline, diskfile, t.filename, osbuild.QEMUFormatVMDK, osbuild.VMDKOptions{Subformat: osbuild.VMDKSubformatStreamOptimized})
	pipelines = append(pipelines, qemuPipeline.Serialize())
	return pipelines, nil
}

func openstackPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	partitionTable, err := t.getPartitionTable(customizations.GetFilesystems(), options, rng)
	if err != nil {
		return nil, err
	}

	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, partitionTable)
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, treePipeline.Serialize())

	diskfile := "disk.img"
	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(packageSetSpecs[osPkgsKey], customizations.GetKernel().Name)
	imagePipeline := liveImagePipeline(&buildPipeline, &treePipeline, diskfile, partitionTable, t.arch, kernelVer)
	pipelines = append(pipelines, imagePipeline.Serialize())

	qemuPipeline := qemuPipeline(&buildPipeline, &imagePipeline, diskfile, t.filename, osbuild.QEMUFormatQCOW2, nil)
	pipelines = append(pipelines, qemuPipeline.Serialize())
	return pipelines, nil
}

func ec2CommonPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions,
	repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec,
	rng *rand.Rand, diskfile string) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	partitionTable, err := t.getPartitionTable(customizations.GetFilesystems(), options, rng)
	if err != nil {
		return nil, err
	}

	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, partitionTable)
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, treePipeline.Serialize())

	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(packageSetSpecs[osPkgsKey], customizations.GetKernel().Name)
	imagePipeline := liveImagePipeline(&buildPipeline, &treePipeline, diskfile, partitionTable, t.arch, kernelVer)
	pipelines = append(pipelines, imagePipeline.Serialize())
	return pipelines, nil
}

// ec2Pipelines returns pipelines which produce uncompressed EC2 images which are expected to use RHSM for content
func ec2Pipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	return ec2CommonPipelines(t, customizations, options, repos, packageSetSpecs, rng, t.Filename())
}

//makeISORootPath return a path that can be used to address files and folders in
//the root of the iso
func makeISORootPath(p string) string {
	fullpath := path.Join("/run/install/repo", p)
	return fmt.Sprintf("file://%s", fullpath)
}

func iotInstallerPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	pipelines = append(pipelines, buildPipeline.Serialize())

	installerPackages := packageSetSpecs[installerPkgsKey]
	d := t.arch.distro
	archName := t.Arch().Name()
	kernelVer := rpmmd.GetVerStrFromPackageSpecListPanic(installerPackages, "kernel")
	ostreeRepoPath := "/ostree/repo"
	payloadStages := ostreePayloadStages(options, ostreeRepoPath)
	kickstartOptions, err := osbuild.NewKickstartStageOptions(kspath, "", customizations.GetUsers(), customizations.GetGroups(), makeISORootPath(ostreeRepoPath), options.OSTree.Ref, "fedora")
	if err != nil {
		return nil, err
	}
	ksUsers := len(customizations.GetUsers())+len(customizations.GetGroups()) > 0
	pipelines = append(pipelines, *anacondaTreePipeline(repos, installerPackages, kernelVer, archName, d.product, d.osVersion, "IoT", ksUsers))

	isolabel := fmt.Sprintf(d.isolabelTmpl, archName)
	pipelines = append(pipelines, *bootISOTreePipeline(kernelVer, archName, d.vendor, d.product, d.osVersion, isolabel, kickstartOptions, payloadStages))

	pipelines = append(pipelines, *bootISOPipeline(t.Filename(), d.isolabelTmpl, archName, false))

	return pipelines, nil
}

func iotCorePipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec) (*pipeline.BuildPipeline, *pipeline.OSPipeline, *pipeline.OSTreeCommitPipeline, error) {
	buildPipeline := pipeline.NewBuildPipeline(t.arch.distro.runner)
	buildPipeline.Repos = repos
	buildPipeline.PackageSpecs = packageSetSpecs[buildPkgsKey]
	treePipeline, err := osPipeline(&buildPipeline, t, repos, packageSetSpecs[osPkgsKey], customizations, options, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	commitPipeline := ostreeCommitPipeline(&buildPipeline, &treePipeline, options, t.arch.distro.osVersion)

	return &buildPipeline, &treePipeline, &commitPipeline, nil
}

func iotCommitPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline, treePipeline, commitPipeline, err := iotCorePipelines(t, customizations, options, repos, packageSetSpecs)
	if err != nil {
		return nil, err
	}
	tarPipeline := pipeline.NewTarPipeline(buildPipeline, &commitPipeline.Pipeline, "commit-archive")
	tarPipeline.Filename = t.Filename()
	pipelines = append(pipelines, buildPipeline.Serialize(), treePipeline.Serialize(), commitPipeline.Serialize(), tarPipeline.Serialize())
	return pipelines, nil
}

func iotContainerPipelines(t *imageType, customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0)

	buildPipeline, treePipeline, commitPipeline, err := iotCorePipelines(t, customizations, options, repos, packageSetSpecs)
	if err != nil {
		return nil, err
	}

	nginxConfigPath := "/etc/nginx.conf"
	httpPort := "8080"
	containerTreePipeline := containerTreePipeline(buildPipeline, commitPipeline, repos, packageSetSpecs[containerPkgsKey], options, customizations, nginxConfigPath, httpPort)
	containerPipeline := containerPipeline(buildPipeline, &containerTreePipeline.Pipeline, t, nginxConfigPath, httpPort)

	pipelines = append(pipelines, buildPipeline.Serialize(), treePipeline.Serialize(), commitPipeline.Serialize(), containerTreePipeline.Serialize(), containerPipeline.Serialize())
	return pipelines, nil
}

func osPipeline(buildPipeline *pipeline.BuildPipeline,
	t *imageType,
	repos []rpmmd.RepoConfig,
	packages []rpmmd.PackageSpec,
	c *blueprint.Customizations,
	options distro.ImageOptions,
	pt *disk.PartitionTable) (pipeline.OSPipeline, error) {

	imageConfig := t.getDefaultImageConfig()

	pl := pipeline.NewOSPipeline(buildPipeline, t.rpmOstree)

	pl.PartitionTable = pt

	if t.Arch().Name() == distro.S390xArchName {
		pl.BootLoader = pipeline.BOOTLOADER_ZIPL
	} else {
		pl.BootLoader = pipeline.BOOTLOADER_GRUB
	}

	pl.UEFI = t.supportsUEFI()
	pl.GRUBLegacy = t.arch.legacy
	pl.Vendor = t.arch.distro.vendor

	var kernelOptions []string
	if t.kernelOptions != "" {
		kernelOptions = append(kernelOptions, t.kernelOptions)
	}
	if bpKernel := c.GetKernel(); bpKernel.Append != "" {
		kernelOptions = append(kernelOptions, bpKernel.Append)
	}
	pl.KernelOptionsAppend = kernelOptions
	pl.KernelName = c.GetKernel().Name

	pl.OSTreeParent = options.OSTree.Parent
	pl.OSTreeURL = options.OSTree.URL

	pl.Repos = repos
	pl.PackageSpecs = packages
	pl.GPGKeyFiles = imageConfig.GPGKeyFiles

	if !t.bootISO {
		// don't put users and groups in the payload of an installer
		// add them via kickstart instead
		pl.Groups = c.GetGroups()
		pl.Users = c.GetUsers()
	}

	services := &blueprint.ServicesCustomization{
		Enabled:  imageConfig.EnabledServices,
		Disabled: imageConfig.DisabledServices,
	}
	if extraServices := c.GetServices(); extraServices != nil {
		services.Enabled = append(services.Enabled, extraServices.Enabled...)
		services.Disabled = append(services.Disabled, extraServices.Disabled...)
	}
	pl.EnabledServices = services.Enabled
	pl.DisabledServices = services.Disabled
	pl.DefaultTarget = imageConfig.DefaultTarget

	pl.Firewall = c.GetFirewall()

	language, keyboard := c.GetPrimaryLocale()
	if language != nil {
		pl.Language = *language
	} else {
		pl.Language = imageConfig.Locale
	}
	if keyboard != nil {
		pl.Keyboard = keyboard
	} else if imageConfig.Keyboard != nil {
		pl.Keyboard = &imageConfig.Keyboard.Keymap
	}

	if hostname := c.GetHostname(); hostname != nil {
		pl.Hostname = *hostname
	}

	timezone, ntpServers := c.GetTimezoneSettings()
	if timezone != nil {
		pl.Timezone = *timezone
	} else {
		pl.Timezone = imageConfig.Timezone
	}

	if len(ntpServers) > 0 {
		pl.NTPServers = ntpServers
	} else if imageConfig.TimeSynchronization != nil {
		pl.NTPServers = imageConfig.TimeSynchronization.Timeservers
	}

	pl.Grub2Config = imageConfig.Grub2Config
	pl.Sysconfig = imageConfig.Sysconfig
	pl.SystemdLogind = imageConfig.SystemdLogind
	pl.CloudInit = imageConfig.CloudInit
	pl.Modprobe = imageConfig.Modprobe
	pl.DracutConf = imageConfig.DracutConf
	pl.SystemdUnit = imageConfig.SystemdUnit
	pl.Authselect = imageConfig.Authselect
	pl.SELinuxConfig = imageConfig.SELinuxConfig
	pl.Tuned = imageConfig.Tuned
	pl.Tmpfilesd = imageConfig.Tmpfilesd
	pl.PamLimitsConf = imageConfig.PamLimitsConf
	pl.Sysctld = imageConfig.Sysctld
	pl.DNFConfig = imageConfig.DNFConfig
	pl.SshdConfig = imageConfig.SshdConfig
	pl.AuthConfig = imageConfig.Authconfig
	pl.PwQuality = imageConfig.PwQuality
	pl.WAAgentConfig = imageConfig.WAAgentConfig

	return pl, nil
}

func ostreeCommitPipeline(buildPipeline *pipeline.BuildPipeline, treePipeline *pipeline.OSPipeline, options distro.ImageOptions, osVersion string) pipeline.OSTreeCommitPipeline {
	p := pipeline.NewOSTreeCommitPipeline(buildPipeline, treePipeline)
	p.Ref = options.OSTree.Ref
	p.OSVersion = osVersion
	p.Parent = options.OSTree.Parent
	return p
}

func containerTreePipeline(buildPipeline *pipeline.BuildPipeline, commitPipeline *pipeline.OSTreeCommitPipeline, repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, c *blueprint.Customizations, nginxConfigPath, listenPort string) pipeline.OSTreeCommitServerTreePipeline {
	p := pipeline.NewOSTreeCommitServerTreePipeline(buildPipeline, commitPipeline)
	p.Repos = repos
	p.PackageSpecs = packages
	p.NginxConfigPath = nginxConfigPath
	p.ListenPort = listenPort
	language, _ := c.GetPrimaryLocale()
	if language != nil {
		p.Language = *language
	}
	return p
}

func containerPipeline(buildPipeline *pipeline.BuildPipeline, treePipeline *pipeline.Pipeline, t *imageType, nginxConfigPath, listenPort string) pipeline.OCIContainerPipeline {
	p := pipeline.NewOCIContainerPipeline(buildPipeline, treePipeline)
	p.Architecture = t.Arch().Name()
	p.Filename = t.Filename()
	p.Cmd = []string{"nginx", "-c", nginxConfigPath}
	p.ExposedPorts = []string{listenPort}
	return p
}

func ostreePayloadStages(options distro.ImageOptions, ostreeRepoPath string) []*osbuild.Stage {
	stages := make([]*osbuild.Stage, 0)

	// ostree commit payload
	stages = append(stages, osbuild.NewOSTreeInitStage(&osbuild.OSTreeInitStageOptions{Path: ostreeRepoPath}))
	stages = append(stages, osbuild.NewOSTreePullStage(
		&osbuild.OSTreePullStageOptions{Repo: ostreeRepoPath},
		osbuild.NewOstreePullStageInputs("org.osbuild.source", options.OSTree.Parent, options.OSTree.Ref),
	))

	return stages
}

func anacondaTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, kernelVer, arch, product, osVersion, variant string, users bool) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "anaconda-tree"
	p.Build = "name:build"
	p.AddStage(osbuild.NewRPMStage(osbuild.NewRPMStageOptions(repos), osbuild.NewRpmStageSourceFilesInputs(packages)))
	p.AddStage(osbuild.NewBuildstampStage(buildStampStageOptions(arch, product, osVersion, variant)))
	p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US.UTF-8"}))

	rootPassword := ""
	rootUser := osbuild.UsersStageOptionsUser{
		Password: &rootPassword,
	}

	installUID := 0
	installGID := 0
	installHome := "/root"
	installShell := "/usr/libexec/anaconda/run-anaconda"
	installPassword := ""
	installUser := osbuild.UsersStageOptionsUser{
		UID:      &installUID,
		GID:      &installGID,
		Home:     &installHome,
		Shell:    &installShell,
		Password: &installPassword,
	}
	usersStageOptions := &osbuild.UsersStageOptions{
		Users: map[string]osbuild.UsersStageOptionsUser{
			"root":    rootUser,
			"install": installUser,
		},
	}

	p.AddStage(osbuild.NewUsersStage(usersStageOptions))
	p.AddStage(osbuild.NewAnacondaStage(osbuild.NewAnacondaStageOptions(users)))
	p.AddStage(osbuild.NewLoraxScriptStage(loraxScriptStageOptions(arch)))
	p.AddStage(osbuild.NewDracutStage(dracutStageOptions(kernelVer, arch, []string{
		"anaconda",
		"rdma",
		"rngd",
		"multipath",
		"fcoe",
		"fcoe-uefi",
		"iscsi",
		"lunmask",
		"nfs",
	})))
	p.AddStage(osbuild.NewSELinuxConfigStage(&osbuild.SELinuxConfigStageOptions{State: osbuild.SELinuxStatePermissive}))

	return p
}

func bootISOTreePipeline(kernelVer, arch, vendor, product, osVersion, isolabel string, ksOptions *osbuild.KickstartStageOptions, payloadStages []*osbuild.Stage) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "bootiso-tree"
	p.Build = "name:build"

	p.AddStage(osbuild.NewBootISOMonoStage(bootISOMonoStageOptions(kernelVer, arch, vendor, product, osVersion, isolabel), osbuild.NewBootISOMonoStagePipelineTreeInputs("anaconda-tree")))
	p.AddStage(osbuild.NewKickstartStage(ksOptions))
	p.AddStage(osbuild.NewDiscinfoStage(discinfoStageOptions(arch)))

	for _, stage := range payloadStages {
		p.AddStage(stage)
	}

	return p
}
func bootISOPipeline(filename, isolabel, arch string, isolinux bool) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "bootiso"
	p.Build = "name:build"

	p.AddStage(osbuild.NewXorrisofsStage(xorrisofsStageOptions(filename, isolabel, arch, isolinux), osbuild.NewXorrisofsStagePipelineTreeInputs("bootiso-tree")))
	p.AddStage(osbuild.NewImplantisomd5Stage(&osbuild.Implantisomd5StageOptions{Filename: filename}))

	return p
}

func liveImagePipeline(buildPipeline *pipeline.BuildPipeline, treePipeline *pipeline.OSPipeline, outputFilename string, pt *disk.PartitionTable, arch *architecture, kernelVer string) pipeline.LiveImgPipeline {
	p := pipeline.NewLiveImgPipeline(buildPipeline, treePipeline)

	p.Filename = outputFilename

	if arch.name == distro.S390xArchName {
		p.BootLoader = pipeline.BOOTLOADER_ZIPL
	} else {
		p.BootLoader = pipeline.BOOTLOADER_GRUB
		p.GRUBLegacy = arch.legacy
	}

	p.KernelVer = kernelVer
	p.PartitionTable = *pt

	return p
}

func qemuPipeline(buildPipeline *pipeline.BuildPipeline, imagePipeline *pipeline.LiveImgPipeline, inputFilename, outputFilename string, format osbuild.QEMUFormat, formatOptions osbuild.QEMUFormatOptions) pipeline.QemuPipeline {
	p := pipeline.NewQemuPipeline(buildPipeline, imagePipeline, string(format))
	p.InputFilename = inputFilename
	p.OutputFilename = outputFilename
	p.Format = format
	p.FormatOptions = formatOptions

	return p
}
