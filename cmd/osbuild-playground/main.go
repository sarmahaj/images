package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/distroregistry"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

var ImageTypes = make(map[string]ImageType)

type ImageType interface {
	Name() string
	InstantiateManifest(m *manifest.Manifest, repos []rpmmd.RepoConfig, runner string) error
	GetExports() []string
	GetCheckpoints() []string
}

func AddImageType(img ImageType) {
	ImageTypes[img.Name()] = img
}

// osbuild-playground is a utility command and is often run from within the
// source tree.  Find the dnf-json binary in case the osbuild-composer package
// isn't installed.  This prioritises the local source version over the system
// version if run from within the source tree.
func findDnfJsonBin() string {
	locations := []string{"./dnf-json", "/usr/libexec/osbuild-composer/dnf-json", "/usr/lib/osbuild-composer/dnf-json"}
	for _, djPath := range locations {
		_, err := os.Stat(djPath)
		if !os.IsNotExist(err) {
			return djPath
		}
	}

	// can't run: panic
	panic(fmt.Sprintf("could not find 'dnf-json' in any of the known paths: %+v", locations))
}

func main() {
	var distroArg string
	flag.StringVar(&distroArg, "distro", "", "distro to build from")
	var archArg string
	flag.StringVar(&archArg, "arch", common.CurrentArch(), "architecture to build for")
	var imageTypeArg string
	flag.StringVar(&imageTypeArg, "type", "my-image", "image type to build")
	flag.Parse()

	// Path to options or '-' for stdin
	optionsArg := flag.Arg(0)

	img := ImageTypes[imageTypeArg]
	if optionsArg != "" {
		var reader io.Reader
		if optionsArg == "-" {
			reader = os.Stdin
		} else {
			var err error
			reader, err = os.Open(optionsArg)
			if err != nil {
				panic("Could not open path to image options: " + err.Error())
			}
		}
		file, err := ioutil.ReadAll(reader)
		if err != nil {
			panic("Could not read image options: " + err.Error())
		}
		err = json.Unmarshal(file, img)
		if err != nil {
			panic("Could not parse image options: " + err.Error())
		}
	}

	distros := distroregistry.NewDefault()
	var d distro.Distro
	if distroArg != "" {
		d = distros.GetDistro(distroArg)
		if d == nil {
			panic(fmt.Sprintf("distro '%s' not supported\n", distroArg))
		}
	} else {
		d = distros.FromHost()
		if d == nil {
			panic("host distro not supported")
		}
	}

	arch, err := d.GetArch(archArg)
	if err != nil {
		panic(fmt.Sprintf("arch '%s' not supported\n", archArg))
	}

	repos, err := rpmmd.LoadRepositories([]string{"./"}, d.Name())
	if err != nil {
		panic("could not load repositories for distro " + d.Name())
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("os.UserHomeDir(): " + err.Error())
	}

	state_dir := path.Join(home, ".local/share/osbuild-playground/")

	RunPlayground(img, d, arch, repos, state_dir)
}
