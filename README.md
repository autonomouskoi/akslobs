# akslobs

AKSLOBS is hacked-together OBS Overlay for current tracks on Denon DJ Prime 4. In the future, this
functionality will be incorporated into the AutonomousKoi bot.

It's known to work with Prime 4 and Prime 4+. It may work on Prime Go, SC Live,
SC5000/6000 with X1800/1850 or separate mixer, but these haven't been tested yet.

It's known to work on Macs running Apple Silicon. It hasn't been tested on on Intel-based Macs,
Windows, or Linux.

It's only expected to work with Denon devices running in Standalone mode, not Computer mode
connected to DJ software like Serato.

The Denon devices _must_ be connected to the network via ethernet cable in the _Link_ port. The
Denon devices may _also_ be connected to WiFi, but the StagelinQ protocol used by AKSLOBS doesn't
seem to work over WiFi. The computer running AKSLOBS can be connected via WiFi, but the WiFi
network and wired network must be the same.

## How It Works

StagelinQ is a network communications format created by Denon for doing visualization stuff, I
think. With StagelinQ, another device can talk to the Denon device and get information about what
the Denon device is doing, such as currently playing tracks.

AKSLOBS talks to your Denon device over the network via StagelinQ and instructs the Denon device to
tell it when a track is loaded into a deck and when a deck starts and stops playing. AKSLOBS
receives this information and keeps track of what's going on with each deck.

OBS has support for _Browser_ Sources. A Browser source is just a web page that OBS will display
and update. AKSLOBS makes a web page available to OBS that will show the currently playing tracks.

When a new track is loaded, AKSLOBS will make note of the track but won't change what it's
displaying. When you start playing a track AKSLOBS will update what it displays. Stopping the track
does not remove it from the display. The display updates every second so when a new track starts
playing it will take up to a second to be displayed.

## How To Use It

Download the `.zip` file for your system from the Releases page.

Open the `.zip` file. Place the `akslobs.html` file in your home folder on your PC. The program
can be saved anywhere.

Your Denon device must be plugged into the network with an ethernet cable to the device's Link
port. Having the Denon device on WiFi isn't sufficient; this is kind of lame.

Launch the program. It should open a terminal window and start logging what it's doing. It will
attempt to connect to your Denon for 5 seconds and if it doesn't succeed it will keep trying. It
will also open a new browser tab/window with further instructions and some controls. Closing this
browser tab/window will cause the program to exit.

The controls allow you to toggle whether or not specific decks are displayed. This is useful when
you're only doing two deck mixes, or some channels are playing from external sources like vinyl.

## Customizing

You can customize the look of the overlay by editing CSS in the `akslobs.html` file. Two
customizations are relatively easy: per-deck colors and deck ordering.

To set per-deck colors, find this section:

```
.deck {
    background-color: aqua;
    border: solid black 1px;
    width: 208px;
    font-family: "Titillium Web", sans-serif;
}
```

Change it to be:

```
.deck {
    background-color: aqua;
    border: solid black 1px;
    width: 208px;
    font-family: "Titillium Web", sans-serif;
}
#deck1 {
    background-color: green;
}
#deck2 {
    background-color: blue;
}
#deck3 {
    background-color: red;
}
#deck4 {
    background-color: gray;
}
```

Search for an _HTML Color Picker_ and use it to find a color you like. The color picker should
present the chosen color in the form `#f3b516`. Use that value in place of the color word. For
example, you could replace the `red` in `#deck3` with `#be6a58`.

When you change one of these values, you can select the Browser source you've added in OBS and
click the _Refresh_ button to see these changes.

You can also change the ordering of the decks. Find the line that says:

```
let deckIDs = ["deck1", "deck2", "deck3", "deck4"];
```

and change the order of the decks. For example:

```
let deckIDs = ["deck3", "deck1", "deck2", "deck4"];
```

Similarly to changing colors, to see changes to deck ordering you have to _Refresh_ the Browser
source.

Much more advanced display is possible by editing the HTML, but that's beyond the scope of this
document.