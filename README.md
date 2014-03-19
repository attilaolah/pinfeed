# PinFeed: a Pinterest peed proxy ![Public Domain](https://pypip.in/license/intperm/badge.png)

A simple Heroku application that acts as a proxy to Pinterest feeds. It parses
each item in the feed and updates the embedded image by changing the thumbinail
to the original size image.

This is useful together with [IFTTT][1] or other automation tools when
generating content (e.g.  tweets) based on a Pinterest feed. You can auto-post
the original size images from Pinterest as Twitter images.

The API is simple:

* Use `https://pinfeed.herokuapp.com/attilaolah` to get the feed for user
  `attilaolah`. This will use `https://www.pinterest.com/attilaolah/feed.rss`
  as the source.
* Use `https://pinfeed.herokuapp.com/attilaolah/for-the-home` to get the feed
  for the specific pin board. This will use
  https://www.pinterest.com/attilaolah/for-the-home.rss` as the source.

[1]: https://ifttt.com/
