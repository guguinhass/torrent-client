# torrent-client

A small BitTorrent client written in Go. Give it a `.torrent` file and a destination path, and it downloads the file from other peers in the swarm.

This project is based on [veggiedefender/torrent-client](https://github.com/veggiedefender/torrent-client) and its accompanying blog post, ["Building a BitTorrent client from scratch in Go"](https://blog.jse.li/posts/torrent/). It is distributed under the same license, the **GNU GPL v3.0** (see [License](#license)).

## Usage

Build:

```sh
go build -o torrent-client .
```

Run:

```sh
./torrent-client <file.torrent> <output-path>
```

Try it with Debian, which has many seeders and uses a tracker. Download a `.torrent` from <https://cdimage.debian.org/debian-cd/current/amd64/bt-cd/> and run:

```sh
./torrent-client debian-netinst.iso.torrent debian.iso
```

The `.torrent` file can be anywhere; paths are relative to the directory you run the command from.

## Limitations

- Only supports `.torrent` files (no magnet links, no DHT)
- Only supports HTTP trackers (no UDP trackers)
- Only supports single-file torrents
- Leech only: it downloads but does not upload pieces to other peers

## How it works

A `.torrent` file is a small [bencoded](https://en.wikipedia.org/wiki/Bencode) description of a download. It holds the tracker URL, the file name and length, the piece length, and the SHA-1 hash of every piece. The file itself comes from other peers. The download goes through these stages:

1. **Parse the `.torrent`.** Decode the bencode and compute the *info hash*: the SHA-1 of the bencoded `info` dictionary. It identifies the torrent to the tracker and to peers.
2. **Ask the tracker for peers.** Send an HTTP GET to the announce URL with the info hash, a random 20-byte peer ID, and how many bytes are left. The tracker replies with a compact list of peers (6 bytes each: 4 for the IP, 2 for the port).
3. **Connect and handshake.** Open a TCP connection to each peer and exchange a 68-byte handshake containing the protocol string, the info hash and the peer ID. Peers whose info hash differs are dropped.
4. **Exchange messages.** Peers speak a length-prefixed binary protocol (`choke`, `unchoke`, `interested`, `have`, `bitfield`, `request`, `piece`, ...). The peer's `bitfield` tells us which pieces it has.
5. **Download pieces concurrently.** One goroutine per peer pulls pieces from a shared work queue. Each piece is requested in 16 KiB blocks, with several requests in flight at once (pipelining) to keep the connection busy.
6. **Verify and assemble.** Each finished piece is hashed with SHA-1 and compared with the hash from the `.torrent`. Corrupt pieces are put back in the queue for another peer. Verified pieces are copied into place and the result is written to disk.

## Project layout

| Package | Responsibility |
|---|---|
| `main.go` | Command-line entry point: opens the torrent and starts the download |
| `torrentFile` | Parses the `.torrent`, computes the info hash, talks to the tracker, and drives the download (Go package name: `torrentfile`) |
| `peers` | Parses the tracker's compact peer list into IP and port pairs |
| `handshake` | Builds, sends and reads the BitTorrent handshake |
| `message` | Serializes and parses protocol messages and their payloads |
| `bitfield` | Compact bit set recording which pieces a peer has |
| `client` | A connection to one peer: dial, handshake, read and send messages, track choke state |
| `p2p` | The download engine: work queue, per-peer workers, pipelined block requests, integrity checks |

## Requirements

- Go (see `go.mod` for the minimum version)
- Dependency: [`jackpal/bencode-go`](https://github.com/jackpal/bencode-go) for bencode parsing

## Credits

- [veggiedefender/torrent-client](https://github.com/veggiedefender/torrent-client) by Jesse Li, the project this one is based on, and the blog post that walks through it.
- [`jackpal/bencode-go`](https://github.com/jackpal/bencode-go) for bencode.

## License

This project is a derivative work of [veggiedefender/torrent-client](https://github.com/veggiedefender/torrent-client), which is licensed under the GNU General Public License v3.0. In accordance with that license, this project is also distributed under the **GNU GPL v3.0**. See the [LICENSE](LICENSE) file for the full text.

Please only download content you have the right to share.
