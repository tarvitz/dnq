Files
=====

- bad-file.txt is just a plain file, it does not do anything to ogg
- opus-headers-only.ogg -- it contains only the first two pages of OGG file
  without a payload (a real data).
- opus-broken-segment-sizes.ogg -- it contains ogg page, but it's not complete.
  Metadata.PageSegments is set to 64, however, the real amount of page segments
  (i.e. sizes) is 3. Thus far page reader will meet EOF soon enough.
- opus-broken-segments.ogg -- it contains ogg page with wrong amount of segments.
- vorbis-headers-only.ogg -- contains the same decoded audio with vorbis codec.
  https://core.telegram.org/bots/api#sendvoice. Despite vorbis works too, opus
  encoded messages looks in telegram much better ;).

Notes
-----
The original source for ogg was:
`Original source <https://yadi.sk/d/mXaTbBD478w75g>`_

All binary data (ogg files preparations) has been performed by the help of
`Sweetscape 010 Editor <https://www.sweetscape.com/010editor/>`_
