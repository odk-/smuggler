# Smuggler

Smuggler is very simple software to smuggle file between servers. Sometimes server is hidden behind firewalls and SCP file from it may be problematic. Big enterprise employees will understand. This is a HTTP server that accepts POST with binary data and offers to receive over HTTP. Typical flow:

1. Register "package"

		$ curl https://smuggler.uri/smuggle

2. In response you will get instructions like this

		To send data use:
		curl -X POST --data-binary "@file" https://smuggler.uri/pack?packID=611d17d616fc5a1fbd4429f3c35fc552
		To retrieve them use:
		curl https://smuggler.uri/pack?packID=611d17d616fc5a1fbd4429f3c35fc552

Smuggler doesn't cache file on transfer so curl with POST will block until receiver comes in. If receiver will connect before sender it will also block and wait.