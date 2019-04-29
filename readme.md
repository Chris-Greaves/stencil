# stencil

Stencil is a tool for taking pre-created project templates, and creating the project using the user's configuration.

## Create a project using a template

```bash
stencil github.com/chris-greaves/stencil-template-test
```

Enter the overrides you want, or just keep pressing 'Enter' till it gets to the building of the project.

## How to get it

You can get pre-compiled binaries from the [Release section on GitHub](https://github.com/Chris-Greaves/stencil/releases).

However, if you want to build for a specific OS / Arch, or just want to build the source code on your own machine. Run the following commands:

```bash
go get github.com/chris-greaves/stencil
cd $GOPATH/src/github.com/chris-greaves/stencil
go install
```

## Licence

Stencil is released under the Apache 2.0 license. See [LICENSE](LICENSE)