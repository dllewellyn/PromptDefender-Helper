const tips = [
    "Tip: You can use keyboard shortcuts to navigate faster!",
    "Tip: Double-click on a word to select it.",
    "Tip: Use Ctrl + F to find text quickly.",
    "Tip: Press Ctrl + S to save your work frequently.",
    "Tip: Use the search bar to find specific settings.",
    "Tip: Customize your theme for better visibility."
];

function randomizeTip() {
    const randomIndex = Math.floor(Math.random() * tips.length);
    document.getElementById('tipText').innerText = tips[randomIndex];
}

document.addEventListener("turbo:before-fetch-request", () => {
    console.log("turbo:before-fetch-request");
    showLoader();
    randomizeTip();
});

// Listen for turbo errors
document.addEventListener("turbo:fetch-request-error", () => {
    console.log("turbo:fetch-request-error ");
    window.location.href = "/error";
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

document.addEventListener("turbo:submit-end", (ev) => {
    hideLoader();

    // Check if event is error: 
    if (ev.success === false) {
        console.log("turbo:submit-end");
        window.location.href = "/error";
    }
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
