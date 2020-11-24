Podbook: Youtube-dl -> Podcast manager
======================================

Podbook is a tool that retrieves media using youtube-dl and stores audio files in a simple database. This makes sure that the same item is never downloaded twice even if the files are renamed. In addition podbook can generate rss-files that can be read by any podcast player.

Building
--------

Install the latest version of Go (Golang) and a few other tools (this procedure will soon be handled by Docker). Then run

$ make

Quick start
-----------

Start by initializing podbook with the url where you will store your media and rss-files.

$ ./out/podbook init https://myserver.com/podcasts

To download a few Audiobooks from Vibravox and generate a feed.rss, run

$ ./out/podbook get -rss feed.rss https://www.youtube.com/watch?v=EEmq9Lk31lc https://www.youtube.com/watch?v=Dy3IPNjflEo https://www.youtube.com/watch?v=0PWnawMMGbY

All books will be downloaded in parallel and feed.rss will be updated every time a download has finished.
