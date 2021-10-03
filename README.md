# chip8.go

## Overview

A [chip-8 implementation]((http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx0A)) written in Go. Core implementation exposes a simple interface for plugging in a display responsible for IO (keyboard + rendering). `pixeldisplay` contains a sample implemention of the interface using the [pixel](https://github.com/faiface/pixel) library for rendering. In the future a GopherJS implementation of the frontend.