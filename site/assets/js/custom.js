// Put your custom JS code here
import * as AsciinemaPlayer from 'asciinema-player';
import mermaid from 'mermaid';
import Panzoom from 'svg-pan-zoom';
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

document.querySelectorAll(".mermaid svg").forEach(x => {
    const bbox = x.parentNode.getBoundingClientRect();
    const height = Math.ceil(bbox.bottom - bbox.top).toString() + "px";
    const width = Math.ceil(bbox.right - bbox.left).toString() + "px";
    console.log("height" + height.toString());
    console.log("wdith" + width.toString())

    x.style.height = "100%";

    const parent = x.parentNode;

    // Initialize Panzoom
    const panzoomInstance = Panzoom(x, {
        maxScale: 5,
        minScale: 0.5,
        step: 0.1,
        enableMouseWheelZoom: true,
        controlIconsEnabled: false,
        fit: true,
    });

    parent.style.height = height;
    parent.style.width = width;
});

document.querySelectorAll('.asciinema-video').forEach((el) => {
    var playerDiv = document.createElement("div");
    var playerDivID = el.id + "-player";
    playerDiv.id = playerDivID;
    el.appendChild(playerDiv);
    var asciinemaURL = el.getAttribute("src");
    AsciinemaPlayer.create(asciinemaURL, playerDiv);
});
