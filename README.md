GXUI - A Go cross platform UI library.
=======

[![Join the chat at https://gitter.im/google/gxui](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/google/gxui?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


Disclaimer
---
All code in this package **is experimental and will have frequent breaking
changes**. Please feel free to play, but please don't be upset when the API has significant reworkings.

The code is currently undocumented, and is certainly **not idiomatic Go**. It will be heavily refactored over the coming months.

This is not an official Google product (experimental or otherwise), it is just code that happens to be owned by Google.

Dependencies
---

### Linux:

In order to build GXUI on linux, you will need the following packages installed:

    sudo apt-get install libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev mesa-common-dev libgl1-mesa-dev

### Common:

After setting up ```GOPATH``` (see [Go documentation](https://golang.org/doc/code.html)), you will first need to fetch the required dependencies:

    go get code.google.com/p/freetype-go/freetype
    go get github.com/go-gl/gl/v2.1/gl
    go get github.com/go-gl/glfw/v3.1/glfw


Once these have been fetched, you can then fetch the GXUI library:

    go get github.com/google/gxui

Samples
---
Samples can be found in [`gxui/samples`](https://github.com/google/gxui/tree/master/samples). 

To build all samples run:

    go install github.com/google/gxui/samples/...

And they will be built into ```GOPATH/bin```.

If you add ```GOPATH/bin``` to your PATH, you can simply type the name of a sample to run it. For example: ```image_viewer```. 

Fonts
---
Many of the samples require a font to render text. The dark theme (and currently the only theme) uses `Roboto`.
This is built into the gxfont package.

Make sure to mention this font in any notices file distributed with your application.

Contributing
---
GXUI was written by a couple of Googlers as an experiment, but with help of the open-source community GXUI could mature into something far more interesting.

Contributions, however small are extremely welcome but will require the author to have signed the [Google Individual Contributor License Agreement](https://developers.google.com/open-source/cla/individual?csw=1).

The CLA is necessary mainly because you own the copyright to your changes, even after your contribution becomes part of our codebase, so we need your permission to use and distribute your code. We also need to be sure of various other things—for instance that you'll tell us if you know that your code infringes on other people's patents. You don't have to sign the CLA until after you've submitted your code for review and a member has approved it, but you must do it before we can put your code into our codebase. Before you start working on a larger contribution, you should get in touch with us first through the issue tracker with your idea so that we can help out and possibly guide you. Coordinating up front makes it much easier to avoid frustration later on.
