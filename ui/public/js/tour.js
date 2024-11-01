function logHasSeenTour() {
    localStorage.setItem('hasSeenTour', 'true');
}

function hasSeenTour() {
    return localStorage.getItem('hasSeenTour') === 'true';
}

function runTour(force = false) {
    const tour = new Shepherd.Tour({
        useModalOverlay: true,
        defaultStepOptions: {
            scrollTo: true,
            cancelIcon: {
                enabled: true
            },
            classes: 'shadow-md bg-purple-dark',
        }
    });

    tour.addStep({
        text: 'Welcome to the prompt defender, an application to help you protect your LLM Application from attacks',
        buttons: [
            {
                text: 'Next',
                action: tour.next
            }
        ]
    })

    tour.addStep({
        id: 'code-editor-step',
        text: 'This is the prompt editor. Take a prompt you\'re using in your application and paste it here. We\'ll help to score the security of your prompt, and to automatically suggest improvements.',
        attachTo: {
            element: '#editor-container',
            on: 'right'
        },
        buttons: [
            {
                text: 'Next',
                action: tour.next
            }
        ]
    });

    tour.addStep({
        id: 'example-step',
        text: 'If you just want to test it out without your existing prompt, click here to load example prompt.',
        attachTo: {
            element: '#example',
            on: 'right'
        },
        buttons: [
            {
                text: 'Next',
                action: tour.next
            }
        ]
    });

    tour.addStep({
        id: 'score-step',
        text: 'Once you\'ve entered your prompt, click here to score it. - This will give you a breakdown of the different defences that exist and the ones that don\'t as well as provide you with an overall score.',
        attachTo: {
            element: '#score',
            on: 'right'
        },
        buttons: [
            {
                text: 'Next',
                action: tour.next
            }
        ]
    });

    tour.addStep({
        id: 'improve-step',
        text: 'If you want to add security instructions to your prompt, paste your prompt and click here - we\'ll try and add all of the prompt defences we can and show it back to you.',
        attachTo: {
            element: '#improve',
            on: 'right'
        },
        buttons: [
            {
                text: 'Finish',
                action: tour.complete
            }
        ]
    });

    if (!hasSeenTour() || force) {
        tour.start();
        logHasSeenTour();
    }
}