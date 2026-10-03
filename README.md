# winbinds

Tiny native Windows exes for window hotkeys (used with Vicinae app hotkeys).

- `maximize-window.exe` — maximize active window
- `close-window.exe` — close active window (WM_CLOSE)

Skips Vicinae's launcher window (allows "Vicinae Settings"). No network, no deps.

## Build

On Windows:

```
go build -trimpath -ldflags "-s -w -H windowsgui" -o maximize-window.exe ./cmd/maximize
go build -trimpath -ldflags "-s -w -H windowsgui" -o close-window.exe ./cmd/close
```

Cross-compile from Linux (fish):

```
env GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -H windowsgui" -o maximize-window.exe ./cmd/maximize
env GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -H windowsgui" -o close-window.exe ./cmd/close
```

## Setup

1. Put exes in a folder, add it to Vicinae → Application directories.
2. Assign hotkeys (e.g. Win+F / Win+Q) to the two apps.
3. Run Vicinae elevated (Task Scheduler, highest privileges) so it works over admin windows.
