# mammut

mammut is an ugly little mastodon TUI client.

[![asciicast](https://asciinema.org/a/391389.svg)](https://asciinema.org/a/391389)

## Install

``` sh
go get github.com/treethought/mammut
```

## Config

Copy mammut.example.yaml and replace it with your own acocunt info. The client ID and client secret can be found fromthe Mastodon developer preferences.



## Keybindings

| key   | action                                 |
|:------|:---------------------------------------|
| TAB   | Switch focue between timeline and menu |
| j     | Move selection up                      |
| k     | Move selection down                    |
| l     | Like toot                              |
| g     | Go to top of timeline                  |
| G     | Go to bottom of timeline               |
| r     | Refresh timeline                       |
| b     | Boost toot                             |
| d     | Delete a toot you published            |
| t     | Compose a new toot                     |
| f     | Follow account                         |
| u     | Unfollow account                       |
| i     | Focus reply input                      |
| c     | Focus reply input                      |
| Enter | Select item (view toot context)        |





