PiCarGo
=======

This app is designed to run on a Raspberry Pi and enhance your vehicle's
capabilities. Most of these features are already being provided by the newest
vehicles on the roads, but for those of us who can't afford those, or at the
least aren't driving them, this provides a relatively-low-cost option for adding
these features to vehicles that don't have them. Plus, there are some advantages
to collecting some of these features into the same place, rather than spreading
them across multiple devices, even ones that can talk to each other.

Since it's written in Go, it should run equally well on other systems, too, so
if you aren't using a Raspberry Pi for this, that's probably not a problem.

I'm building this mostly for my own use (I've assembled all the hardware to make
it work, so now it's time for the software), but I'm releasing it here so others
can benefit and/or submit improvements. I haven't really found anything quite
like it out there, so hopefully I'm not reinventing the same wheel yet again. :D

Planned features:
1.  GPS/AIS location/tracking
2.  Navigation maps
3.  OBD II diagnostics/control
4.  Cellular hotspot (with wired/wireless roaming support)
5.  Hands-free via Bluetooth (uses connected mic, if available)
  - A2DP
  - AVRCP
  - DIP
  - HFP
  - HSP
  - PAN
  - PBA(P)
  - PXP
  - SDAP
  - SAP/SIM/rSAP
6.  Audio-only media center (supports internal and external storage)
7.  ~~Voice commands (I hope!)~~ _Too resource-intensive for a Pi, sorry..._ 

Currently supported:
- Nothing!

Watch this space for more information.
