# replace

Utility for searching and replacing text.

## INSTALLATION

``` shell
> go install github.com/go-curses/coreutils-replace/cmd/rpl@latest
```

## DOCUMENTATION

``` shell
> rpl --help
NAME:
   rpl - search and replace

USAGE:
   rpl [options] <search> <replace> <path> [path...]

VERSION:
   0.2.2

DESCRIPTION:
   command line search and replace

GLOBAL OPTIONS:
   --all, -a                                     include files and directories that start with a "." (default: false)
   --backup, -b                                  make backups before replacing content (default: false)
   --backup-extension value, -B value            specify the backup file suffix to use (default: ".bak")
   --dot-match-nl, -s                            set the dot-match-nl (?s) regexp flag (implies -P) (default: false)
   --dry-run, -n                                 report what would have otherwise been done (default: false)
   --help, -h, --usage                           display command-line usage information (default: false)
   --ignore-case, -i                             perform a case-insensitive search (default: false)
   --interactive, -I                             selectively apply replacement edits (default: false)
   --multi-line, -m                              set the multi-line (?m) regexp flag (implies -P) (default: false)
   --multi-line-dot-match-nl, --ms, -p           convenience flag to set -m and -s (implies -P) (default: false)
   --multi-line-dot-match-nl-insensitive, --msi  convenience flag to set -m, -s and -i (implies -P) (default: false)
   --quiet, -q                                   run silently, ignored if --dry-run is also used (default: false)
   --recurse, -R                                 recurse into sub-directories (default: false)
   --regex, -P                                   search and replace arguments are regular expressions (default: false)
   --show-diff, -D                               include unified diffs of all changes in the output (default: false)
   --verbose, -v                                 run loudly, ignored if --quiet is also used (default: false)
   --version                                     display the version (default: false)
```


## LICENSE

```
Copyright 2023  The Go-Curses Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
