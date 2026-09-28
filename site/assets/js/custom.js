// Put your custom JS code here
import * as AsciinemaPlayer from 'asciinema-player';

document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll('.asciinema-video').forEach((el) => {
        console.log(el);
        var playerDiv = document.createElement("div");
        var playerDivID = el.id + "-player";
        playerDiv.id = playerDivID;
        el.appendChild(playerDiv);
        console.log("creating player");
        var asciinemaURL = el.getAttribute("src");
        console.log(asciinemaURL);
        AsciinemaPlayer.create(asciinemaURL, playerDiv);
    });
});
