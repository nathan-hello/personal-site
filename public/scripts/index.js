(function asdf() {
    document.getElementById("comment-form").addEventListener("htmx:afterRequest", function (event) {
        if (event.detail.successful) {
            event.target.reset();
        }
    });
    document.body.addEventListener("click", (e) => {
        if (e.target.id === "draw-btn") {
            Tegaki.open({
                onDone: () => {
                    Tegaki.flatten().toBlob((blob) => {
                        const dt = new DataTransfer();
                        dt.items.add(new File([blob], "drawing.png", { type: "image/png" }));
                        const input = document.getElementById("comment-image-input");
                        input.files = dt.files;
                        document.getElementById("fileName").textContent = "drawing.png";
                        Tegaki.hide();
                    }, "image/png");
                },
                onCancel: () => {},
                width: 480,
                height: 480,
            });
        }
    });
})();
