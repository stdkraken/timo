function Fish(raw) {
    x = 0
    while(x < 200) {
        sleep(10)
        console.log(x)
        x++
    }
}

function Neo(raw) {
    execute("neofetch", 2)
}

function Sylt(raw) {
    execute("xdg-open", 4, "https://www.youtube.com/watch?v=ZryxBDQYOkY")
}

addCustomCommand("fish", Fish)
addCustomCommand("neo", Neo)
addCustomCommand("sylt", Sylt)