document.addEventListener("turbo:before-fetch-request", () => {
    console.log("turbo:before-fetch-request");
    showLoader();
});

function showLoader() {
    console.log("showLoader");
    document.getElementById("container").style.display = "none";
    document.getElementById("loader").style.display = "block";
}


function hideLoader() {
    console.log("hideLoader");
    document.getElementById("container").style.display = "block";
    document.getElementById("loader").style.display = "none";
}

document.addEventListener("turbo:before-frame-render", () => {
    hideLoader();
});

document.addEventListener("turbo:load", () => {
    hideLoader();
})


const turboEvents = [
    "turbo:click",
    "turbo:before-visit",
    "turbo:visit",
    "turbo:before-cache",
    "turbo:before-render",
    "turbo:render",
    "turbo:load",
    "turbo:morph",
    "turbo:before-morph-element",
    "turbo:before-morph-attribute",
    "turbo:morph-element",
    "turbo:submit-start",
    "turbo:submit-end",
    "turbo:before-frame-render",
    "turbo:frame-render",
    "turbo:frame-load",
    "turbo:frame-missing",
    "turbo:before-stream-render",
    "turbo:before-fetch-request",
    "turbo:before-fetch-response",
    "turbo:before-prefetch",
    "turbo:fetch-request-error"
];

turboEvents.forEach(eventType => {
    document.addEventListener(eventType, event => {
        console.log(`Event: ${event.type}`, event.detail);
    });
});
