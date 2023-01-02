Create git activities from your commits.

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

#### Usage

git-mirror works as git post-commit hook. 
So after install it you need to add related hook to project repositories by `git-mirror install` in your repos.  

When you add a new commit to your project, stats of the commit will be added to the mirror repository.
You need to push the new commit of mirror repository in order to see contributions at your GitHub profile and contribution chart.

----

#### Example generated mirror repo content
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
$ cat 2023/01/02/log.java
7 insertion(s), 2 deletion(s)
59 insertion(s), 6 deletion(s)
9 insertion(s), 1 deletion(s)
```

#### TODO

- [ ] Write better readme
- [ ] File stat aliases: `yaml=yml,java=gradle,xml,properties`
- [ ] commit, filename templates with config
```shell
-vars="company=mycorp,project=api-gateway"
filename_format="{company}.{ext}"
log_format="{company}->{project}: {insert_count} insertion(s), {delete_count} deletion(s)"
```
- [ ] Funny commit messages
- [ ] Create os packages for homebrew,apt etc.
----

License: MIT