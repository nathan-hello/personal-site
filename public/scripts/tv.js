async function check() {
    const { ok } = await fetch("https://tv.reluekiss.com/mystream/index.m3u8");
    if (ok) {
        document.getElementById("video-player").src = "https://tv.reluekiss.com/mystream";
        document.getElementById("video-parent").classList.remove("is-offline");
        document.getElementById("video-parent").classList.add("is-online");
    } else {
        document.getElementById("video-player").src = "";
        document.getElementById("video-parent").classList.remove("is-online");
        document.getElementById("video-parent").classList.add("is-offline");
    }
}
await check();
setInterval(check, 5000);
