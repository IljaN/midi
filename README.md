# midi

Currently only provides a master-clock which can be used to write to /dev/midi.

## Virtual midi device
A virtual midi device can be created to receive the master clock.
```
sudo modprobe snd-virmidi snd_index=1
aconnect -io
aconnect 20:0 21:0
``
