<p align="center"><img width="200px" src="/ressources/img/logo.png" alt="podcast-cli"/></p>


___

Top-like interface for listening podcasts
`podcast-cli` lets you play your podcasts from the terminal:
<p align="center"><img src="/ressources/img/demo.gif" alt="podcast-cli"/></p>

`podcast-cli` is built entirely with Go, you can run it on `Linux`, `Mac OS` and `Windows`.

## Install
Fetch the [latest release](https://github.com/goulinkh/podcast-cli/releases)

#### Linux

```bash
sudo wget https://github.com/goulinkh/podcast-cli/releases/download/v0.1.1/podcast-cli-0.1.1-linux-amd64 -O /usr/local/bin/podcast-cli
sudo chmod +x /usr/local/bin/podcast-cli
```

#### OS X

```bash
sudo curl -Lo /usr/local/bin/podcast-cli https://github.com/goulinkh/podcast-cli/releases/download/v0.1.1/podcast-cli-0.1.1-darwin-amd64
sudo chmod +x /usr/local/bin/podcast-cli
```

## Usage
`podcast-cli` requires no arguments and uses your default internet settings to access the internet.

### Options

| Options                  | Description                                 |
| ------------------------ | ------------------------------------------- |
| `-h or  --help`          | Print help information                      |
| `-s or --search <query>` | List podcasts that matches the search query |

### Keybindings

| Key        | Action |
| ---------- | ------ |
| `Enter`    | Select |
| `p, Space` | Pause  |
| `Esc`      | Back   |
| `Right`    | +10s   |
| `Left`     | -10s   |
| `q`        | Exit   |
| `s`        | SEARCH |


## Issues

* Unable to get audio length of a remote content, I have to download the audio file before playing it
* Bad Podcasts API, I have to make a free access API or find a better one
