document.addEventListener("turbo:before-fetch-request", () => {
    console.log("turbo:before-fetch-request");
    Pace.restart();
});

document.addEventListener("turbo:before-fetch-response", () => {
    console.log("turbo:before-fetch-response");
    Pace.stop();
});