Keep your git statistics with post-commit hook.

---

#### Usage

```shell
USAGE:
   git-mirror [global options] command [command options] [arguments...]

COMMANDS:
   install  install post-commit hook for adding stats automatically
   add      add stats of latest commit
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --repo value       git repo directory path of the mirror repo (default: "$HOME/.git-mirror")
   --whitelist value  comma seperated file extensions to create stats. eg: go,rs,sh,Makefile (default: all extensions)
   --help, -h         show help (default: false)
```

#### Install
- Grab a [binary](https://github.com/mesuutt/git-mirror/releases) or install via `go install github.com/mesuutt/git-mirror@latest`
- Create a git repository at `~/.git-mirror`(Do not forget `git init`)
- Copy [example config file](https://github.com/mesuutt/git-mirror/blob/main/example.config.toml) to your `~/.git-mirror` repo as `config.toml`. 

#### Usage

git-mirror works as git post-commit hook. 
So after install it you need to add related hook to project repositories by `git-mirror install` in your repos.  

When you add a new commit to your project, stats of the commit will be added to the mirror repository.


----

#### Example generated stats repo content
```shell
$ tree
.
└── 2023
    └── 01
        ├── 01
        │ └── log.go
        └── 02
            ├── log.go
            ├── log.java
            ├── log.md
            ├── log.sql
            └── log.yaml

# Example generated file content after some commits
$ cat 2023/01/02/log.go
fmt.Println("02:42 24 insertion(s), 7 deletion(s)")
fmt.Println("03:39 83 insertion(s), 11 deletion(s)")
```

#### TODO
- [x] file type overwriting (`jsx=js`)
- [x] content template with config
- [ ] filename template with config
- [ ] excluding file extensions
- [ ] separating stat dirs like `personal/`, `work/`
- [ ] custom vars in params and templates(`-vars="project=foo-api-gateway"`)

----

License: MIT