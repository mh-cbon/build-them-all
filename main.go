package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mh-cbon/build-them-all/utils"
	"github.com/urfave/cli"
)

var VERSION = "0.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "build-them-all"
	app.Version = VERSION
	app.Usage = "Command line to build go programs to multiple targets"
	app.UsageText = "build-them-all <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:      "clean",
			Usage:     "Clean up directories",
			UsageText: "build-them-all clean <options>",
			Action:    clean,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "dir, d",
					Value: &cli.StringSlice{"build/"},
					Usage: "Directoires to clean",
				},
			},
		},
		{
			Name:      "build",
			Usage:     "Build the binaries",
			UsageText: "build-them-all build <options> <packages>",
			Action:    build,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "os",
					Value: "major",
					Usage: "OS selector",
				},
				cli.StringFlag{
					Name:  "arch",
					Value: "major",
					Usage: "Arch selector",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: "Output pattern",
				},
				cli.StringFlag{
					Name:  "gobin",
					Value: "go",
					Usage: "Go bin path",
				},
				cli.StringSliceFlag{
					Name:  "dir, d",
					Value: &cli.StringSlice{"build/"},
					Usage: "Directoires to clean",
				},
				cli.StringFlag{
					Name:  "wd",
					Value: "",
					Usage: "Working directory",
				},
				cli.BoolFlag{
					Name:  "i",
					Usage: "Installs the packages that are dependencies of the target",
				},
				cli.BoolFlag{
					Name:  "a",
					Usage: "Force rebuilding of packages that are already up-to-date",
				},
				cli.BoolFlag{
					Name:  "n",
					Usage: "Print the commands but do not run them",
				},
				cli.StringFlag{
					Name:  "p",
					Value: "",
					Usage: "The number of programs that can be run in parallel",
				},
				cli.BoolFlag{
					Name:  "v",
					Usage: "Print the names of packages as they are compiled",
				},
				cli.BoolFlag{
					Name:  "x",
					Usage: "Print the commands",
				},
				cli.BoolFlag{
					Name:  "work",
					Usage: "Print the name of the temporary work directory and do not delete it when exiting.",
				},
				cli.BoolFlag{
					Name:  "race",
					Usage: "Enable data race detection",
				},
				cli.BoolFlag{
					Name:  "msan",
					Usage: "Enable interoperation with memory sanitizer",
				},
				cli.StringFlag{
					Name:  "asmflags",
					Value: "",
					Usage: "Arguments to pass on each go tool asm invocation",
				},
				cli.StringFlag{
					Name:  "buildmode",
					Value: "",
					Usage: "Build mode to use. See 'go help buildmode' for more",
				},
				cli.StringFlag{
					Name:  "compiler",
					Value: "",
					Usage: "Name of compiler to use, as in runtime.Compiler (gccgo or gc)",
				},
				cli.StringFlag{
					Name:  "gccgoflags",
					Value: "",
					Usage: "Arguments to pass on each gccgo compiler/linker invocation",
				},
				cli.StringFlag{
					Name:  "gcflags",
					Value: "",
					Usage: "Arguments to pass on each go tool compile invocation",
				},
				cli.StringFlag{
					Name:  "installsuffix",
					Value: "",
					Usage: "A suffix to use in the name of the package installation directory",
				},
				cli.StringFlag{
					Name:  "ldflags",
					Value: "",
					Usage: "Arguments to pass on each go tool link invocation",
				},
				cli.StringFlag{
					Name:  "linkshared",
					Value: "",
					Usage: "Link against shared libraries previously created with -buildmode=shared",
				},
				cli.StringFlag{
					Name:  "pkgdir",
					Value: "",
					Usage: "Install and load all packages from dir instead of the usual locations",
				},
				cli.StringFlag{
					Name:  "tags",
					Value: "",
					Usage: "A list of build tags to consider satisfied during the build",
				},
				cli.StringFlag{
					Name:  "toolexec",
					Value: "",
					Usage: "A program to use to invoke toolchain programs like vet and asm",
				},
			},
		},
	}

	app.Run(os.Args)
}

func clean(c *cli.Context) error {
	errs := utils.CleanDirectories(c.StringSlice("dir"))
	for _, err := range errs {
		fmt.Println(err)
	}
	if len(errs) > 0 {
		return cli.NewExitError("Some files were not properly deleted", 1)
	}
	fmt.Println("Cleaned directories")
	for _, d := range c.StringSlice("dir") {
		fmt.Println(string(d))
	}
	return nil
}

func build(c *cli.Context) error {
	wantOs := c.String("os")
	wantArch := c.String("arch")
	output := c.String("output")
	gobin := c.String("gobin")
	wd := c.String("wd")
	dirs := c.StringSlice("dir")

	packages := make([]string, 0)
	for _, arg := range c.Args() {
		packages = append(packages, arg)
	}
	if len(packages) == 0 {
		return cli.NewExitError("Packages list is required", 1)
	}

	gobin, err := exec.LookPath(gobin)
	if err != nil {
		return cli.NewExitError("Path to go bin '"+gobin+"' was not found", 1)
	}

	if len(wd) == 0 {
		wd, err = os.Getwd()
		if err != nil {
			return cli.NewExitError("Working directory could not be determined", 1)
		}
	}

	errs := utils.CleanDirectories(dirs)
	for _, err := range errs {
		fmt.Println(err)
	}
	if len(errs) > 0 {
		return cli.NewExitError("Some files were not properly deleted", 1)
	}

	fmt.Println("wd=" + wd)
	gotErr := false
	selectedOsArch := utils.SelectOsArch(wantOs, wantArch)
	for _, input := range packages {
		pkg := utils.DetermineOutputName(input)
		if len(pkg) == 0 {
			fmt.Println("Could not determine output name for '" + input + "'")
			gotErr = true
			break
		}
		for cOs, archs := range selectedOsArch {
			for _, cArch := range archs {

				bCmd := utils.BuildCommand{
					Install:       c.Bool("i"),
					RebuildAll:    c.Bool("a"),
					N:             c.Bool("n"),
					P:             c.String("p"),
					Race:          c.Bool("race"),
					Msan:          c.Bool("msan"),
					V:             c.Bool("v"),
					Work:          c.Bool("work"),
					X:             c.Bool("x"),
					Asmflags:      c.String("asmflags"),
					Buildmode:     c.String("buildmode"),
					Compiler:      c.String("compiler"),
					Gccgoflags:    c.String("gccgoflags"),
					Gcflags:       c.String("gcflags"),
					InstallSuffix: c.String("installsuffix"),
					Ldflags:       c.String("ldflags"),
					Linkshared:    c.String("linkshared"),
					Pkgdir:        c.String("pkgdir"),
					Tags:          c.String("tTags"),
					Toolexec:      c.String("toolexec"),
					Output:        utils.GenerateOutputName(output, pkg, cOs, cArch),
					Wd:            wd,
					Gobin:         gobin,
				}

				// build [-o output] [-i] [build flags] [packages]
				oCmd := utils.GenerateBuildCommand(cOs, cArch, bCmd)
				fmt.Println("> GOOS=" + cOs + " " + "GOARCH=" + cArch + " " + strings.Join(oCmd.Args, " "))

				out, err := oCmd.CombinedOutput()
				sOut := string(out)
				if len(sOut) > 0 {
					fmt.Println(sOut)
				}
				if err == nil {
					fmt.Println("Success!")
				} else {
					fmt.Println("Failure!")
					fmt.Println(err)
					gotErr = true
				}
				fmt.Println("")
			}
		}
	}

	if gotErr {
		return cli.NewExitError("There were errors during the build", 1)
	}

	return nil
}
