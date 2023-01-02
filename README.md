Create git activities from commits by commits.

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
   --whitelist value  comma seperated file extensions to create stats. eg: go,rs,sh,Makefile
   --help, -h         show help (default: false)
```

#### Install
- Install `git-miror` by `go install github.com/mesuutt/git-mirror`
- Create a git repository at `~/.git-mirror`
- Go to project directory and add git hook by `git-mirror install`

After that when you add a new commit to your project, stats of the commit will be added to the mirror repository.

You need to push the new commit of mirror repository to see contributions on your Github profile.

----

License: MIT