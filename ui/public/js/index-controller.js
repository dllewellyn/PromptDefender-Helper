const application = Stimulus.Application.start();

application.register("code-editor", class extends Stimulus.Controller {

    static targets = ["tour", "example", "score", "improve", "editor", "form", "prompt"];

    connect() {
        var editorTarget = this.editorTarget.value;
        this.editor = CodeMirror.fromTextArea(this.editorTarget, {
            lineNumbers: true,
            lineWrapping: true,
            theme: "default"
        });

        if (editorTarget !== '') {
            console.log('Setting editor value ' + editorTarget);
            this.editor.setValue(editorTarget);
            this.editor.refresh();
        }
    }

    refresh() {
        this.editor.refresh();
    }

    loadExample() {
        this.editor.setValue(`Your job is to translate users input from English into French:\n\n<user_input>{user_input}</user_input>`);
    }

    setFormAction(event) {
        let prompt = this.editor.getValue();

        if (prompt === '') {
            alert('Please enter a prompt');
            event.preventDefault();
            return;
        }

        const action = event.currentTarget.dataset.actionValue;
        this.formTarget.action = action;
        this.promptTarget.value = this.editor.getValue();

    }

    takeTour() {
        runTour(true);
    }
})