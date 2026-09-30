// Put your custom JS code here
import * as AsciinemaPlayer from 'asciinema-player';
import mermaid from 'mermaid';
/*
theme?: 'default' | 'base' | 'dark' | 'forest' | 'neutral' | 'neo' | 'neo-dark' | 'redux' | 'redux-dark' | 'redux-color' | 'redux-dark-color' | 'null';
look?: 'classic' | 'handDrawn' | 'neo';
*/
mermaid.registerIconPacks([
  {
    name: 'pixel',
    loader: () =>
      fetch('https://unpkg.com/@iconify-json/streamline-pixel@1.2.0/icons.json').then((res) => res.json()),
  },
]);
mermaid.initialize({
    theme: 'redux-dark-color',
    look: 'neo'
});

document.querySelectorAll('.asciinema-video').forEach((el) => {
    var playerDiv = document.createElement("div");
    var playerDivID = el.id + "-player";
    playerDiv.id = playerDivID;
    el.appendChild(playerDiv);
    var asciinemaURL = el.getAttribute("src");
    AsciinemaPlayer.create(asciinemaURL, playerDiv);
});
