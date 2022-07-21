inputFG = 0x5B // light red
setInputFg(inputFG)

while(true) {
    for (var color=30; color++;) {
        sleep(1000)
        setInputFg(color)
        if (color == 37) {
            color = 29
        }
    }
}