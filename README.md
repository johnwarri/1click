# crossrun

A skeleton for shipping cross-platform, zero-install command-line logic behind
a single download link. You drop your code into one file; a build script turns
it into standalone binaries for Windows, macOS, and Linux; and a small web page
hands each visitor the right one automatically.

## Layout

    app.go                       <- YOUR logic goes here
    platform.go                  <- OS detection + command runner (don't need to touch)
    main.go                      <- entry point
    build.sh                     <- builds a binary for every OS/arch into ./dist
    web/index.html               <- OS-detecting download page (one link in, right file out)
    .github/workflows/release.yml<- builds + publishes binaries on every version tag

## 1. Add your logic

Open `app.go` and replace the EXAMPLE block inside `Run()`. You get a
`Platform` value with three helpers for OS-level work:

    p.Run("some command")          // run in the native shell, stream output
    out, _ := p.Capture("cmd")     // run and get the output back as a string
    p.Exec("program", "arg1")      // run a program directly, no shell

Because OS-level commands differ per system, branch with `p.IsWindows()`:

    cmd := "ls -la"
    if p.IsWindows() { cmd = "dir" }
    p.Run(cmd)

## 2. Build

    ./build.sh

Produces standalone binaries in `dist/` — nothing needs to be installed on the
end user's machine. You can run this on any one OS; Go cross-compiles the rest.
Run it on a Mac and it also produces a universal Mac binary (Intel + Apple
Silicon). Rename the app by editing `APP=` at the top of `build.sh` (match it
in `web/index.html`).

## 3. Publish + the single link

Push a version tag and the GitHub Action builds everything and attaches the
binaries to a Release:

    git tag v0.1.0
    git push origin v0.1.0

Then set two lines in `web/index.html` (`APP` and your repo's
`releases/latest/download/` URL) and host that page on GitHub Pages,
Cloudflare Pages, or any static host. Visitors hit one link, the page detects
their OS, and the correct binary downloads on its own. GitHub serves release
assets as forced downloads, so the download starts without a second click.

## The one caveat

The plumbing is pure Go, so cross-compiling works. If your logic later pulls in
a library that needs C (CGO), cross-compilation breaks and you'd need to build
each OS on its own machine — switch the Action from a single `macos-latest` job
to a matrix across `windows-latest`, `macos-latest`, and `ubuntu-latest`.
