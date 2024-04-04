# vpc_flow_logs
The golang project template builds in some automation to get off the ground quickly. Here's how to get started.


Let's get the tests running right away. 'make static' will run all of the local checks
```shell
cd vpc_flow_logs
git init .
git add -A
git commit -am 'initial'
go mod init
go mod tidy
make statc
```

Now let's build the exampls executables
```shell
make build
tree build
build
├── 72ff9f1cf1d5c56a9c9261774bd906efce6245c4
│   ├── darwin
│   │   └── amd64
│   │       ├── prog1
│   │       ├── prog2
│   │       └── version.txt
│   └── linux
│       └── amd64
│           ├── prog1
│           ├── prog2
│           └── version.txt
└── current -> /Users/nmarks/go/src/github.com/natemarks/stayback/build/72ff9f1cf1d5c56a9c9261774bd906efce6245c4

6 directories, 6 files
```


Let's say you were happy with the project and you wanted to release it, you'd bump the version. 
```shell
# clean-venv creates a python virtual environment jsut to use the bump2version python package
make clean-venv
# this does a patch level semver bump
make part=patch bump
```

Now create the release contents with naming that's ok for github releases:
```shell
make release
warning: ignoring symlink /Users/nmarks/go/src/github.com/natemarks/stayback/build/current
warning: ignoring symlink /Users/nmarks/go/src/github.com/natemarks/stayback/build/current
f1de3a5aaa990726bec76b23f0368b5d90d5fa86/linux/amd64
f1de3a5aaa990726bec76b23f0368b5d90d5fa86/darwin/amd64
f1de3a5aaa990726bec76b23f0368b5d90d5fa86/linux/amd64
f1de3a5aaa990726bec76b23f0368b5d90d5fa86/darwin/amd64
rm -f build/current
ln -s /Users/nmarks/go/src/github.com/natemarks/stayback/build/f1de3a5aaa990726bec76b23f0368b5d90d5fa86 /Users/nmarks/go/src/github.com/natemarks/stayback/build/current
mkdir -p release/0.0.1
a .
a ./prog2
a ./version.txt
a ./prog1
a .
a ./prog2
a ./version.txt
a ./prog1
❯ tree release
release
└── 0.0.1
    ├── stayback_0.0.1_darwin_amd64.tar.gz
    └── stayback_0.0.1_linux_amd64.tar.gz

1 directory, 2 files
```