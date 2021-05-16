const chooseFileButton = $('#choose-file-button');
const fileChosen = $('#file-chosen');
const filesTable = $("#dataTable").DataTable()
var files

chooseFileButton.on('click', function(){
    $.get('/choose-files', function(data){
        files = JSON.parse(data);

        fileChosen.textContent = files.length + (files.length > 1 ? " files chosen" : " file chosen")
        files.forEach(file => addFileToTable(file))
    }).fail(function() {
        console.error("could not get files")
    })
    // fileChosen.textContent = this.files.length + (this.files.length > 1 ? " files chosen" : " file chosen")
});

function addFileToTable(file){
    filesTable.clear().draw();
    filesTable.row.add([file.Path, file.Name, "<input class=\"form-control\" type=\"text\">"]).draw( false );
}
