application.register("upload-file", class extends Stimulus.Controller {

    static targets = ["form", "fileUpload", "prompt", "fileContent", "message", "okmessage"];

    connect() {
        const pond = FilePond.create(this.fileUploadTarget, {
            server: false,
            allowMultiple: false,
            instantUpload: true,
            acceptedFileTypes: ['text/*'],
            labelIdle: 'Drag & Drop your prompt file or <span class="filepond--label-action">Browse</span>',
            labelFileTypeNotAllowed: 'Please select a text file',
            onaddfile: this.handleFileUpload.bind(this),
            onremovefile: this.handleFileRemove.bind(this)
        });
    }

    setFormAction(event) {
        if (!this.promptTarget.value || this.promptTarget.value === '') {
            alert('Please upload a file');
            event.preventDefault();
            return;
        }

        const action = event.currentTarget.dataset.actionValue;
        this.formTarget.action = action;
    }

    handleFileUpload(error, file) {
        if (error) {
            console.log('Error adding file:', error);
            return;
        }

        const reader = new FileReader();
        const promptTarget = this.promptTarget;
        const okMessage = this.okmessageTarget;
        const messageTarget = this.messageTarget;

        reader.onload = function (e) {
            try {
                const content = e.target.result;
                promptTarget.value = content;
                
                okMessage.textContent = 'File uploaded successfully';

            } catch (error) {
                console.error('Error processing file:', error);
                messageTarget.textContent = 'Error processing file: ' + error;
            }
        };

        reader.onerror = (error) => {
            console.error('Error reading file:', error);
            messageTarget.innerHTML = '<p>Error reading file: ' + error + '</p>';
        };

        reader.readAsText(file.file);
    }

    handleFileRemove(error, file) {
        if (error) {
            console.log('Error removing file:', error);
            return;
        }
    }
});
