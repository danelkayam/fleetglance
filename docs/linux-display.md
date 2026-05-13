# Running Fleetglance on a Linux machine with an attached display

Fleetglance is designed to work well on a Linux machine connected to a display.

This can be a Raspberry Pi, mini PC, old laptop, or any small Linux machine.

> Tested it only on Raspberry Pi 4 with cheap "7 800x480 display. Higher resultion displays will result better visuality.

## Recommended setup

- Linux machine with attached display
- `fleetglance-console` installed
- `fleetglance.yaml` configured
- Network access from the console machine to each agent
- X11, Openbox, and Alacritty for a simple full-screen kiosk setup

## Install display packages

On a headless distribution such as Raspberry Pi OS Lite, install a minimal graphical environment:

```sh
sudo apt update
sudo apt install --no-install-recommends \
  xserver-xorg \
  xinit \
  openbox \
  x11-xserver-utils \
  alacritty \
  fontconfig \
  wget \
  unclutter \
  unzip
```

## Install JetBrainsMono Nerd Font

Fleetglance uses terminal glyphs/icons. A Nerd Font is recommended.

```sh
mkdir -p ~/.local/share/fonts
cd /tmp
wget https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.zip
unzip JetBrainsMono.zip -d ~/.local/share/fonts/JetBrainsMonoNerd
fc-cache -fv
```

Verify the font:

```sh
fc-match "JetBrainsMono Nerd Font"
```

## Configure Alacritty

Create the config directory:

```sh
mkdir -p ~/.config/alacritty
```

Edit the config file:

```sh
vim ~/.config/alacritty/alacritty.toml
```

Recommended starting config for small attached displays:

```toml
[window]
decorations = "None"
startup_mode = "Fullscreen"
dynamic_title = false
padding = { x = 0, y = 0 }

[font]
normal = { family = "JetBrainsMono Nerd Font", style = "Regular" }
bold = { family = "JetBrainsMono Nerd Font", style = "Bold" }
italic = { family = "JetBrainsMono Nerd Font", style = "Italic" }
bold_italic = { family = "JetBrainsMono Nerd Font", style = "Bold Italic" }
size = 10.0

[cursor]
style = { shape = "Block", blinking = "Never" }
unfocused_hollow = false

[colors.primary]
background = "#0B0F17"
foreground = "#D7DEE9"
```

Removing padding gives Fleetglance more usable terminal space.

## Configure Openbox autostart

Create the Openbox config directory:

```sh
mkdir -p ~/.config/openbox
```

Edit the autostart file:

```sh
vim ~/.config/openbox/autostart
```

Add:

```sh
# Disable screen blanking and power management.
xset s off
xset s noblank
xset -dpms

# Hide mouse cursor after 1 second if unclutter is installed.
unclutter -idle 1 &

# Start Fleetglance in fullscreen Alacritty.
exec alacritty \
  -e /usr/local/bin/fleetglance-console \
  -f "$HOME/.config/fleetglance/fleetglance.yaml"
```

## Configure xinit

Edit:

```sh
vim ~/.xinitrc
```

Add:

```sh
exec openbox-session
```

## Start the kiosk only on local tty1

Edit:

```sh
vim ~/.bash_profile
```

Add:

```sh
# Start Fleetglance graphical kiosk only on the local console.
# SSH sessions should remain normal.
if [ -z "$SSH_CONNECTION" ] && [ -z "$DISPLAY" ] && [ "$(tty)" = "/dev/tty1" ]; then
  # hacky waiting system to finish booting
  sleep 5

  startx
fi
```

This starts the graphical Fleetglance display only on the local console. SSH sessions remain normal.

## Run manually for testing

Before enabling autostart, test the console manually:

```sh
startx
```

Or run Alacritty directly from an existing graphical session:

```sh
alacritty -e /usr/local/bin/fleetglance-console -f "$HOME/.config/fleetglance/fleetglance.yaml"
```

## Check terminal size

Inside the Fleetglance terminal, run:

```sh
stty size
```

If the display is too small, ship names may be trimmed or some visual elements may not fit.

Try:

- reducing terminal padding
- reducing font size
- using full-screen mode
- reducing the number of configured ships

## Restart Fleetglance without rebooting

If Fleetglance is started from `tty1`, restart the local console session from SSH:

```sh
sudo systemctl restart getty@tty1.service
```

## Disable automatic start

Remove the Fleetglance `startx` block from:

```sh
~/.bash_profile
```

Then restart the local console session:

```sh
sudo systemctl restart getty@tty1.service
```

## Troubleshooting

### SSH opens Fleetglance unexpectedly

Check that the `~/.bash_profile` block includes this condition:

```sh
[ -z "$SSH_CONNECTION" ]
```

### Fleetglance does not start on the display

Check that `startx` works manually:

```sh
startx
```

Check that the console binary exists:

```sh
ls -l /usr/local/bin/fleetglance-console
```

Check that the config exists:

```sh
ls -l ~/.config/fleetglance/fleetglance.yaml
```

### Icons look wrong

Verify the font:

```sh
fc-match "JetBrainsMono Nerd Font"
```

If the font is not found, rerun `fc-cache -fv`.

### Screen blanks after a while

Check that the Openbox autostart file contains:

```sh
xset s off
xset s noblank
xset -dpms
```
