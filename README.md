# go-netcatnat - netcat accross nats!

This is a tiny netcat util that works across nats. It's a toy to demonstrate
netty, a nat traversal tool, and its Go libs:

- [getlantern/natty](https://github.com/getlantern/natty) - extracted from chromium
- [getlantern/waddell](https://github.com/getlantern/waddell) - ICE signaling server in Go (tcp)
- [getlantern/go-natty](https://github.com/getlantern/go-natty) - natty bindings in Go
- [getlantern/nattywad](https://github.com/getlantern/nattywad) - natty and waddell client libs

Note that this is over raw UDP, so you might see packets reordered or dropped.
Someone can add SCTP.

![](http://jbenet.static.s3.amazonaws.com/eba94f4/natcat)

## Install

You'll need to install [a great many things](https://www.youtube.com/watch?v=Nembr1ZeRc8):

1. Go (tested on 1.3)
2. Install `netcatnat`:

    ```
    go get github.com/jbenet/go-natnetcat
    go install github.com/jbenet/go-natnetcat/natnetcat
    ```

3. Install `waddell` signaling server:

    ```
    go get github.com/getlantern/waddell
    go install github.com/getlantern/waddell/waddell
    ```

4. You'll need these packages. If you've Chrome/Chromium, should already have
    them. Yes, two of these are A/V stuff; natty punches nat holes with
    ultrasound. No, just kidding. natty got pulled out of chrome src, so for
    now the code still depends on that.

    ```sh
    # debian package names. your pkg mgr's might differ
    libnss3
    libasound2
    libpulse-dev
    ```

## Usage

First, launch `waddell` in an accessible machine A (no nat traversal required)
This is the signaling service that will coordinate for the clients.
Note the IP of this machine, as we'll need it.

```sh
> waddell
waddell: Starting waddell at :62443
```

Ok, waddell is now running at: `<ip-of-machine-A>:62443`.
Next, in machine B launch `netcatnat`:

```sh
> netcatnat --waddell <ip-of-machine-A>:62443
using waddell server: <ip-of-machine-A>:62443
nattywad: Waddell address changed
nattywad: Connected to Waddell!! Id is: f67b6696-5d71-45e7-ada0-ba514c33fc91
```

Note the Id you get (here `f67b6696-5d71-45e7-ada0-ba514c33fc91`). Call it `<machine-B-ID>`.
Finally, in machine C launch `netcatnat` again, with `<machine-B-ID>`

```sh
> netcatnat --waddell <ip-of-machine-A>:62443 --id <machine-B-ID>
using waddell server: :62443
attempting to connect to: f67b6696-5d71-45e7-ada0-ba514c33fc91
nattywad: Configuring nat traversal client with 1 server peers
nattywad: Attempting traversal to f67b6696-5d71-45e7-ada0-ba514c33fc91
nattywad: Connected to Waddell!! Id is: d241153a-3f39-4b3e-a752-53980320195c
connected 192.168.1.56:64284 to 192.168.1.56:59149

>> connected! enter text below to transmit <<

```

Now, if all went well, you can enter text in either terminal and see it appear
in the other. netcat. across nats.

![](http://jbenet.static.s3.amazonaws.com/81644e8/teddy-deal-with-it)
