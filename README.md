[![wercker status](https://app.wercker.com/status/65a46fd0cf1a038c760bea708e5259b3/s/master "wercker status")](https://app.wercker.com/project/byKey/65a46fd0cf1a038c760bea708e5259b3)


# Pair

Pair aims to simplify the pair programming experience


## Installation

The Go way...

``` shell
go get github.com/revolvingcow/pair
```

Once we get to a stable release we will include compiled programs for a
variety of supported platforms. These include GNU/Linux, FreeBSB, and OS X.
Unfortunately Windows doesn't look too good at the moment but you never know.


## Usage

### Key Management

Maintaining the authorized keys can be a pain. After a long pairing session
you may forget to clean out the unused keys. Or maybe you pair with someone
often who trades out keys like some people do shoes. Either way we have
made this much easier to work with.

Each command may take multiple usernames. These usernames are then used to
pull their public keys from both [github.com](https://github.com) and
[gitlab.com](https://gitlab.com) and manage them in the SSH authorized keys
file.

 - [add](doc/pair_add.md) (`pair add bmallred`)
 - [list](doc/pair_list.md) (`pair list`)
 - [remove](doc/pair_remove.md) (`pair remove bmallred`)
 - [sync](doc/pair_sync.md) (`pair sync`)


### Sessions

Alright, so key management is good. Now to another pain point: connecting.
First, one brave soul must declare themselves [as](doc/pair_ar.md) the host
for the pairing session.

``` shell
pair as github.com/bmallred
```

Then one or more participants may join the session [with](doc/pair_with.md)
the host

``` shell
pair with github.com/bmallred
```

Happy hacking!


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


## License

Pair is released under the GNU General Public License (Version 3).
See [LICENSE](https://github.com/revolvingcow/pair/blob/master/LICENSE).

