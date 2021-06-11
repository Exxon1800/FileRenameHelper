const chooseFileButton = $('#choose-file-button');
const filesTable = $("#dataTable").DataTable({
    "order": [[1, "desc"]],
    "columnDefs": [{"orderable": false, "targets": 0}],
});
const fileChosen = $('#file-chosen')[0];
var files

chooseFileButton.on('click', function () {
    $.get('/choose-files', function (data) {
        filesTable.clear().draw();
        files = JSON.parse(data);
        if (files) {
            fileChosen.textContent = files.length + (files.length > 1 ? " files chosen" : " file chosen")
            files.forEach(file => addFileToTable(file))
        }
    }).fail(function () {
        console.error("could not get files")
        fileChosen.textContent = "No file chosen"
    })
});

function addFileToTable(file) {
    let selectedForRenameId = file.Name + "selected"
    let newFileNameId = file.Name + "newName"
    filesTable.row.add([
        `<input type="checkbox" class="select-item checkbox big-checkbox" name="select-item" value="${selectedForRenameId}" />`,
        file.TruncatedPath,
        file.Name,
        `<input id="${newFileNameId}" class="form-control newFileNameInput" type="text" value="${file.Name}">`
    ]).draw(false);
}

// checkboxes
$(function () {
    //button select all or cancel
    $("#select-all").on("click", function () {
        var all = $("input.select-all")[0];
        all.checked = !all.checked
        var checked = all.checked;
        $("input.select-item").each(function (index, item) {
            item.checked = checked;
        });
    });
    //column checkbox select all or cancel
    $("input.select-all").click(function () {
        var checked = this.checked;
        $("input.select-item").each(function (index, item) {
            item.checked = checked;
        });
    });
    //check selected items
    $("input.select-item").click(function () {
        var checked = this.checked;
        var all = $("input.select-all")[0];
        var total = $("input.select-item").length;
        var len = $("input.select-item:checked:checked").length;
        all.checked = len === total;
    });
});
