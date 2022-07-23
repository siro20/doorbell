# Doorbell

A simple application to use the Raspberry PI as doorbell replacement.
Written in golang.

## Application

This application waits for GPIO events using
[github.com/faiface/beep](github.com/faiface/beep) on GPIO17.

On event it plays a sound file using [github.com/faiface/beep](github.com/faiface/beep)
passed as argument. It uses ALSA/pulseaudio as backend to play the file.

To make this work when running on systemd you must take ownership of a sound card in
exclusive mode or configure pulseaudio to run in system mode. 

## Arguments

Usage:

```
./doorbell <mp3 sound file> <volume>
```

Example:
```
./doorbell ./dingdong.mp3 -4
```
This will play dingdong.mp3 with a volume of -4dB.

## Getting the sound working when running under systemd

*First method:*

You can use the environment variable AUDIODEV to use a device in exclusive mode.
```
AUDIODEV="hw:X,y"
```

In this mode no other application can use this device! pulseaudio will not work any more.

*Second method:*

Run pulseaudio in **system mode**.

Modify `/etc/pulse/daemon.conf` and set :

```
daemonize = yes
system-instance = yes
```

Modify `/etc/default/pulseaudio` and set:

```
PULSEAUDIO_SYSTEM_START=1
```

Modify `/etc/pulse/client.conf` and set:

```
default-server = /var/run/pulse/native
autospawn = no
```

## Testing audio

In order to test if audio works without toggling the GPIO you can add an additional
argument:

```
./doorbell <mp3 sound file> <volume> -t
```

This will play the sound file at given volume and then exit.
