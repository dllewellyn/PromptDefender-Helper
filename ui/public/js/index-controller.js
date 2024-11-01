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
        const action = event.currentTarget.dataset.actionValue;
        this.formTarget.action = action;
        this.promptTarget.value = this.editor.getValue();
    }

    takeTour() {
        runTour(true);
    }
})