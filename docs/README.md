---
home: true
heroImage: /note.png
features:
- title: Simple
  details: Just a tiny command line tool without any config.
- title: Tiny
  details: Tiny binary file powered by golang, support all platform.
- title: Efficient
  details: DB driven use bolt which is stable and efficient. Also support sqlite driven.
footer: MIT Licensed | Copyright © 2018-present zcong1993
---

# Note

### Install

```bash
$ curl -sfL https://git.io/vpCMX | sh
```

### Usage

```bash
# show help
$ note --help
# show version
$ note version
# show all notes
$ note # or note list
# get notes by limit and offset
$ note get -l {limit} -o {offset}
$ note get --limit {limit} --offset {offset}
# add a note
$ note add "test note"
# delete a note
$ note delete {index}
# update a note
$ note update {index} {content}
# delete all
$ note delete-all
# flush db
$ note flush
```
