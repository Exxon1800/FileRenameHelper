const chooseFileButton = $('#choose-file-button');
const fileChosen = $('#file-chosen');
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
    $('#dataTable tr:last').after(`    
    <tr>
        <td>${file.Path}</td>
        <td>${file.Name}</td>
        <td><input className="form-control" type="text"></td>
    </tr>
    `);

}
