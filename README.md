<p align="center"><img width="200px" src="/resources/img/logo.png" alt="podcast-cli"/></p>


___

Top-like interface for listening to podcasts
`podcast-cli` lets you play your favourite podcasts from the terminal:
<p align="center"><img src="/resources/img/demo.gif" alt="podcast-cli"/></p>

`podcast-cli` is entirely built with Go, you can run it on `Linux`, `Mac OS` and `Windows`.

## Install
Fetch the [latest release](https://github.com/goulinkh/podcast-cli/releases)

#### Linux

```bash
sudo wget https://github.com/goulinkh/podcast-cli/releases/download/1.3.0/podcast-cli-1.3.0-linux-amd64 -O /usr/local/bin/podcast-cli
sudo chmod +x /usr/local/bin/podcast-cli
```

#### OS X

```bash
sudo curl -Lo /usr/local/bin/podcast-cli https://github.com/goulinkh/podcast-cli/releases/download/1.3.0/podcast-cli-1.3.0-darwin-amd64
sudo chmod +x /usr/local/bin/podcast-cli
```

## Usage
`podcast-cli` requires no arguments and uses your default internet settings to access the internet.

### Options

| Options                  | Description                                 |
| ------------------------ | ------------------------------------------- |
| `-h or  --help`          | Print help information                      |
| `-s or --search <query>` | List podcasts that matches the search query |
| `-r or --rss <url>`    | Custom podcast rss url source               |
| `-o or --offset <episode number starting with 0>` | Play episode number                         |

### Keybindings

| Key        | Action   |
| ---------- | -------- |
| `Enter`    | Select   |
| `p, Space` | Pause    |
| `Esc`      | Back     |
| `Right`    | +10s     |
| `Left`     | -10s     |
| `u`        | Slowdown |
| `d`        | Speedup  |
| `q`        | Exit     |


## Issues

* Unable to get audio length of a remote content, I have to download the audio file before playing it

