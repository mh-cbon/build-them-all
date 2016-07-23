# build-them-all

Bin util to build go programs to multiple targets

## Install

Pick an msi package [here](https://github.com/mh-cbon/build-them-all/releases)!

__deb/rpm__

```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/build-them-all sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/build-them-all sh -xe
```

__go__

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
git clone https://github.com/mh-cbon/build-them-all.git
cd build-them-all
glide install
go install
```

# Usage

```
NAME:
   build-them-all - Command line to build go programs to multiple targets

USAGE:
   build-them-all <cmd> <options>

VERSION:
   0.0.1

COMMANDS:
     clean	Clean up directories
     build	Build the binaries

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

#### clean

```
NAME:
   build-them-all clean - Clean up directories

USAGE:
   build-them-all clean [command options] [arguments...]

OPTIONS:
   --dir value, -d value	Directoires to clean (default: "build/")

EXAMPLES:
  build-them-all clean -d build/ -d some/
```

#### build

```
NAME:
   build-them-all build - Build the binaries

USAGE:
   build-them-all build [command options] [arguments...]

OPTIONS:
   --os value                   OS selector (default: "major")
   --arch value                 Arch selector (default: "major")
   --output value, -o value     Output pattern
   --gobin value                Go bin path (default: "go")
   --dir value, -d value        Directoires to clean (default: "build/")
   --wd value                   Working directory
   -i                           Installs the packages that are dependencies of the target
   -a                           Force rebuilding of packages that are already up-to-date
   -n                           Print the commands but do not run them
   -p value                     The number of programs that can be run in parallel
   -v                           Print the names of packages as they are compiled
   -x                           Print the commands
   --work                       Print the name of the temporary work directory and do not delete it when exiting.
   --race                       Enable data race detection
   --msan                       Enable interoperation with memory sanitizer
   --asmflags value             Arguments to pass on each go tool asm invocation
   --buildmode value            Build mode to use. See 'go help buildmode' for more
   --compiler value             Name of compiler to use, as in runtime.Compiler (gccgo or gc)
   --gccgoflags value           Arguments to pass on each gccgo compiler/linker invocation
   --gcflags value              Arguments to pass on each go tool compile invocation
   --installsuffix value        A suffix to use in the name of the package installation directory
   --ldflags value              Arguments to pass on each go tool link invocation
   --linkshared value           Link against shared libraries previously created with -buildmode=shared
   --pkgdir value               Install and load all packages from dir instead of the usual locations
   --tags value                 A list of build tags to consider satisfied during the build
   --toolexec value             A program to use to invoke toolchain programs like vet and asm

EXAMPLES:
  build-them-all build main.go -o "build/&pkg-&os-&arch" --ldflags "-X main.VERSION=0.0.1"
  build-them-all build main.go -o "build/&pkg-&os-&arch" --ldflags "-X main.VERSION=0.0.1" --arch "amd64"
  build-them-all build main.go -o "build/&pkg-&os-&arch" --ldflags "-X main.VERSION=0.0.1" --os "windows"
```
