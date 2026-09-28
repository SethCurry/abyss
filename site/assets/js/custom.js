// Put your custom JS code here
import * as AsciinemaPlayer from 'asciinema-player';

document.querySelectorAll('.asciinema-video').forEach((el) => {
    var playerDiv = document.createElement("div");
    var playerDivID = el.id + "-player";
    playerDiv.id = playerDivID;
    el.appendChild(playerDiv);
    var asciinemaURL = el.getAttribute("src");
    AsciinemaPlayer.create(asciinemaURL, playerDiv);
});
