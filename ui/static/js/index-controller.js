const application = Stimulus.Application.start();

application.register("code-editor", class extends Stimulus.Controller {

    static targets = ["tour", "example", "score", "improve", "editor", "form", "prompt"];

    connect() {
        this.editor = CodeMirror.fromTextArea(this.editorTarget, {
            lineNumbers: true,
            lineWrapping: true,
            theme: "default"
        });
    }

    loadExample() {
        this.editor.setValue(`Your job is to translate users input from English into French:\n\n<user_input>{user_input}</user_input>`);
    }

    setFormAction(event) {
        let prompt = this.editor.getValue();

        if (prompt === '') {
            alert('Please enter a prompt');
            return;
        }

        // Do a http post as javascript to the /api/improve endpoint
        fetch('/api/improve', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({prompt})
        })
            .then(response => response.json())
            .then(data => {
                this.editor.setValue(data.result);
            })
            .catch(error => console.error('Error:', error));

        event.preventDefault()
    }

    takeTour() {
        runTour(true);
    }
})