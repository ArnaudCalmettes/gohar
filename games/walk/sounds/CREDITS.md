# Where the sounds come from

## snap.wav

- **Source**: "finger-snap-stereo-11.wav", by newagesoup, on Freesound
  (sound #364732):
  <https://freesound.org/people/newagesoup/sounds/364732/>
- **License**: Creative Commons 0 (public domain).
- **Edits**: cut 3 ms before the attack, with a 1 ms fade in; 150 ms of
  tail ending in a 20 ms fade out; mono, 48 kHz, 16 bits.

## GeneralUser GS

The double bass, the piano and the drum kit come from GeneralUser GS
v2.0.3, by S. Christian Collins
(<https://github.com/mrbumpy409/GeneralUser-GS>), under the GeneralUser
GS License v2.0: "Please feel free to use it in your software projects,
and to modify the SoundFont bank or its packaging to suit your needs."
The same license warns that the origin of some of its samples is
uncertain, which "may concern you if you intend to use GeneralUser GS
in a commercial software product".

## walk.sf2

GeneralUser GS slimmed down to what Walk with me plays: the double bass
(0:32), the piano (0:0), and four keys of the Jazz kit (128:32: the
pedal hi-hat, ride 1, the high and low wood blocks), 2.5 MB rather than
32. `make slim` makes it from GeneralUser GS, which `make sounds`
downloads into the user's cache directory and checks its hash; every
note kept sounds as in the full soundfont (synth/sf2, TestGeneralUser).
