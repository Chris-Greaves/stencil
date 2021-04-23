# stencil

Stencil is a tool for taking pre-created project templates, and creating the project using the user's configuration.

## Create a project using a template

Example:

```bash
stencil github.com/Chris-Greaves/stencil-template-test
```

Enter the overrides you want, or just keep pressing 'Enter' till it gets to the building of the project.

## How to get it

You can get pre-compiled binaries from the [Release section on GitHub](https://github.com/Chris-Greaves/stencil/releases).

However, if you want to build for a specific OS / Arch, or just want to build the source code on your own machine. Run the following commands:

```bash
go get github.com/Chris-Greaves/stencil
cd $GOPATH/src/github.com/Chris-Greaves/stencil
go install
```

## Download and run source for contribution

First, pull down the repo:

```bash
git clone github.com/Chris-Greaves/stencil
```

Then use go to download the dependencies and build the cli tool

```bash
go get
go build
```

Test it by using the resulting executable:

```bash
# Example based on Windows
stencil.exe example-templates\basic-example
```

## Licence

Stencil is released under the Apache 2.0 license. See [LICENSE](LICENSE)