const chooseFileButton = $('#choose-file-button');
const fileChosen = $('#file-chosen');

chooseFileButton.on('click', function(){
    let files = $.get('/choose-files', function(data){
        console.log(data)
    }).fail(function() {
        console.error("could not get files")
    })
    // fileChosen.textContent = this.files.length + (this.files.length > 1 ? " files chosen" : " file chosen")
});

