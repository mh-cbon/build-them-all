package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CleanDirectories(dirs []string) []error {
	errs := make([]error, 0)
	for _, d := range dirs {
		sD := string(d)
		if len(sD) > 0 {
			err := os.RemoveAll(sD)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func DetermineOutputName(input string) string {
	pkg := filepath.Base(input)
	pkg = strings.Replace(pkg, filepath.Ext(pkg), "", -1)
	if pkg == "main" {
		if _, err := os.Stat(input); !os.IsNotExist(err) {
			abs, err := filepath.Abs(input)
			if err != nil {
				return pkg // :-(
			}
			return DetermineOutputName(filepath.Dir(abs))
		}
	}
	return pkg
}

func SelectOsArch(os string, arch string) map[string][]string {
	if os == "all" {
		os = strings.Join(GetKeys(OsArchMap()), ",")
	} else if os == "major" {
		os = strings.Join([]string{"linux", "darwin", "windows"}, ",")
	}
	if arch == "all" {
		arch = "386,amd64,arm,arm64,ppc64,ppc64le,mips64,mips64le"
	} else if arch == "major" {
		arch = "386,amd64"
	}
	selected := make(map[string][]string)
	for cOs, aArch := range OsArchMap() {
		if strings.Index(os, cOs) > -1 {
			selectedArchs := make([]string, 0)
			for _, cArch := range aArch {
				if strings.Index(arch, cArch) > -1 {
					selectedArchs = append(selectedArchs, cArch)
				}
			}
			if len(selectedArchs) > 0 {
				selected[cOs] = selectedArchs
			}
		}
	}
	return selected
}

func OsArchMap() map[string][]string {
	valueMap := make(map[string][]string)

	valueMap["darwin"] = []string{"386", "amd64", "arm", "arm64"}
	valueMap["dragonfly"] = []string{"amd64"}
	valueMap["freebsd"] = []string{"386", "amd64", "arm"}
	valueMap["linux"] = []string{
		"386",
		"amd64",
		"arm",
		"arm64",
		"ppc64",
		"ppc64le",
		"mips64",
		"mips64le",
	}
	valueMap["netbsd"] = []string{"386", "amd64", "arm"}
	valueMap["openbsd"] = []string{"386", "amd64", "arm"}
	valueMap["plan9"] = []string{"386", "amd64"}
	valueMap["solaris"] = []string{"amd64"}
	valueMap["windows"] = []string{"386", "amd64"}

	return valueMap
}

func GetKeys(values map[string][]string) []string {
	keys := make([]string, 0)
	for k := range values {
		keys = append(keys, k)
	}
	return keys
}

func GenerateOutputName(pattern string, wantPkg string, wantOs string, wantArch string) string {
	output := ""
	output = strings.Replace(pattern, "&pkg", wantPkg, -1)
	output = strings.Replace(output, "&os", wantOs, -1)
	output = strings.Replace(output, "&arch", wantArch, -1)
	if wantOs == "windows" {
		output = output + ".exe"
	}
	return output
}

type BuildCommand struct {
	Install       bool
	RebuildAll    bool
	N             bool
	P             string
	Race          bool
	Msan          bool
	V             bool
	Work          bool
	X             bool
	Asmflags      string
	Buildmode     string
	Compiler      string
	Gccgoflags    string
	Gcflags       string
	InstallSuffix string
	Ldflags       string
	Linkshared    string
	Pkgdir        string
	Tags          string
	Toolexec      string
	Output        string
	Wd            string
	Gobin         string
}

func GenerateBuildCommand(wantOs string, wantArch string, bCmd BuildCommand) *exec.Cmd {
	// build [-o output] [-i] [build flags] [packages]
	args := make([]string, 0)
	args = append(args, "build")
	if len(bCmd.Output) > 0 {
		args = append(args, []string{"-o", bCmd.Output}...)
	}
	if bCmd.Install {
		args = append(args, "-i")
	}
	if bCmd.RebuildAll {
		args = append(args, "-a")
	}
	if bCmd.N {
		args = append(args, "-n")
	}
	if len(bCmd.P) > 0 {
		args = append(args, []string{"-p", bCmd.P}...)
	}
	if bCmd.V {
		args = append(args, "-v")
	}
	if bCmd.Race {
		args = append(args, "-race")
	}
	if bCmd.Msan {
		args = append(args, "-msan")
	}
	if bCmd.Work {
		args = append(args, "-work")
	}
	if bCmd.X {
		args = append(args, "-x")
	}
	if len(bCmd.Asmflags) > 0 {
		args = append(args, []string{"-asmflags", bCmd.Asmflags}...)
	}
	if len(bCmd.Buildmode) > 0 {
		args = append(args, []string{"-buildmode", bCmd.Buildmode}...)
	}
	if len(bCmd.Compiler) > 0 {
		args = append(args, []string{"-compiler", bCmd.Compiler}...)
	}
	if len(bCmd.Gccgoflags) > 0 {
		args = append(args, []string{"-gccgoflags", bCmd.Gccgoflags}...)
	}
	if len(bCmd.Gcflags) > 0 {
		args = append(args, []string{"-gcflags", bCmd.Gcflags}...)
	}
	if len(bCmd.InstallSuffix) > 0 {
		args = append(args, []string{"-installsuffix", bCmd.InstallSuffix}...)
	}
	if len(bCmd.Ldflags) > 0 {
		args = append(args, []string{"-ldflags", bCmd.Ldflags}...)
	}
	if len(bCmd.Linkshared) > 0 {
		args = append(args, []string{"-linkshared", bCmd.Linkshared}...)
	}
	if len(bCmd.Pkgdir) > 0 {
		args = append(args, []string{"-pkgdir", bCmd.Pkgdir}...)
	}
	if len(bCmd.Tags) > 0 {
		args = append(args, []string{"-tags", bCmd.Tags}...)
	}
	if len(bCmd.Toolexec) > 0 {
		args = append(args, []string{"-toolexec", bCmd.Toolexec}...)
	}

	oCmd := exec.Command(bCmd.Gobin, args...)
	oCmd.Dir = bCmd.Wd
	oCmd.Env = []string{
		"GOOS=" + wantOs,
		"GOARCH=" + wantArch,
		"GOPATH=" + os.Getenv("GOPATH"),
	}

	return oCmd
}
