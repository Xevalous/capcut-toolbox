# capcut-toolbox

A terminal UI for patching, version-locking, and downloading legacy CapCut builds on Windows.

## Features

- **Patch**: Get many features by patching `VECreator.dll`
- **Lock**: Lock CapCut version
- **Download**: Browse and open legacy CapCut versions (1.0.0-5.4.0 Beta6) in browser

## Build

```
make build
```

Requires Go. Cross-compiles to Windows from any OS (`GOOS=windows GOARCH=amd64`).

## Usage

Run `capcut-toolbox.exe` on a Windows machine with CapCut installed. The TUI walks through prechecks (OS, installation, running state) then presents the main menu.

```
capcut-toolbox.exe --version
```

## License

[MIT](LICENSE)

## Credits

- [iosdevice](https://iosgods.com/profile/6266834-iosdevice/): patches
