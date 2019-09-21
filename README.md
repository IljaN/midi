# midi
[![Coverage Status](https://coveralls.io/repos/github/IljaN/midi/badge.svg?branch=HEAD)](https://coveralls.io/github/IljaN/midi?branch=HEAD)

Currently only provides a master-clock which can be used to write to /dev/midi.

## Virtual midi device
A virtual midi device can be created to receive the master clock.
```
sudo modprobe snd-virmidi snd_index=1
aconnect -io
aconnect 20:0 21:0
``
