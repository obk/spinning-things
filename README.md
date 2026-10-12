# Spinning Cube

A spinning 3D cube and a spinning donut, rendered as ASCII art in the terminal. Written in Go with only the standard library.

```
                        ##$####****!!!!!!!!*!!*****####$$$$$$$$$$
                      #####*****!!*!!!!!=!==!!!!!!!****##$$$$$$$$$#
                    *#####**!!!!!!===;;;;;;=======!!!!***####$$$$$$$#
                   #####***!!!==;;;;:::~~~~~:::;;====!*****###$$$$$$$#!
                 *####****!!==;:::~~--,,,----~~:::;===!!****#####$$$###*
                **##****!!===;:~~-,,.........,,-~~:;;==!!****###########*
               **##*****!!==;:~-,,.           .,--~::;==!!****##########*!
```

## Usage

Requires Go 1.27 or later.

```sh
go run . cube    # default
go run . donut
```

Press **Ctrl+C** to quit.

The terminal must be at least 80×40 characters and support ANSI escape codes (Windows Terminal, and most macOS and Linux terminals, do).

## How it works

Both shapes are sampled as points on their surface. Each point is rotated, projected onto an 80×40 character grid, and kept only if it is nearer than what is already in that cell (a z-buffer). Horizontal positions are doubled because terminal cells are about twice as tall as they are wide.

- **Cube** (`cube.go`): each of the six faces has its own character and color.
- **Donut** (`donut.go`): a torus lit by a single light; each point's brightness picks a character from `.,-~:;=!*#$@`.
- **`main.go`**: picks the shape and runs the redraw loop at about 60 frames per second.

Size, speed and projection settings are constants at the top of each file.

## License

[MIT](LICENSE)
