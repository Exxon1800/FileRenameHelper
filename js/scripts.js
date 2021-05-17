const chooseFileButton = $('#choose-file-button');
const filesTable = $("#dataTable").DataTable()
const fileChosen = $('#file-chosen')[0];
var files

chooseFileButton.on('click', function(){
    $.get('/choose-files', function(data){
        filesTable.clear().draw();
        files = JSON.parse(data);
        fileChosen.textContent = files.length + (files.length > 1 ? " files chosen" : " file chosen")
        files.forEach(file => addFileToTable(file))
    }).fail(function() {
        console.error("could not get files")
    })
    // fileChosen.textContent = this.files.length + (this.files.length > 1 ? " files chosen" : " file chosen")
});

function addFileToTable(file){
    filesTable.row.add([file.TruncatedPath, file.Name, "<input class=\"form-control\" type=\"text\">"]).draw( false );
}
