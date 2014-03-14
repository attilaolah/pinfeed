# PinFeed Pinterest Feed Parser ![Public Domain](https://pypip.in/license/intperm/badge.png)

A simple Heroku application that acts as a proxy to Pinterest feeds. It parses
each item in the feed and updates the embedded image by changing the thumbinail
to the original size image.

## Why is this useful?

This is very useful together with [IFTTT][1] when generating content (e.g.
tweets) based on a Pinterest feed.

[1]: https://ifttt.com/

## The API

Just send requests to `/{username}`. That's it.
