document.addEventListener("turbo:submit-start", (ev) => {
    if (ev != undefined && ev !== null) {
        console.log("turbo:before-fetch-request");
        showLoader();
    }
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