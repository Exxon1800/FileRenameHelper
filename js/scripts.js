const chooseFileButton = $('#choose-file-button');
const RenameSelectedFilesButton = $('#rename-selected-files-button');
const selectedPrefix = "selected-"
const newNamePrefix = "new-name-"
const filesTable = $("#dataTable").DataTable({
    "order": [[1, "desc"]],
    "columnDefs": [{"orderable": false, "targets": 0}],
});
const fileChosen = $('#file-chosen')[0];

let rawFiles
let files = new Map();

chooseFileButton.on('click', function () {
    $.get('/choose-files', function (data) {
        filesTable.clear().draw();
        rawFiles = JSON.parse(data);
        if (rawFiles) {
            files.clear()
            fileChosen.textContent = rawFiles.length + (rawFiles.length > 1 ? " files chosen" : " file chosen")
            rawFiles.forEach(function (file, i) {
                rawFiles[i].UUID = UUIDv4()
                files.set(rawFiles[i].UUID, file)
                addFileToTable(file)
            })
        }
    }).fail(function () {
        console.error("could not get files")
        fileChosen.textContent = "No file chosen"
    })
});

// addFileToTable adds the file to the table by adding a new row to th filesTable
// in this row it adds a checkbox with a selectedPrefix + file.UUID as id to later identify the checkbox
// in this row it also adds a input field for the new name and newNamePrefix + file.UUID to identify the input field
function addFileToTable(file) {
    filesTable.row.add([
        `<input type="checkbox" id="${selectedPrefix}${file.UUID}" class="select-item checkbox big-checkbox" name="select-item" />`,
        file.TruncatedPath,
        file.Name,
        `<input id="${newNamePrefix}${file.UUID}" class="form-control newFileNameInput" type="text" onchange="handleNewFileNameInput(this)" value="${file.Name}">`
    ]).draw(false);
}

function handleNewFileNameInput(elem) {
    let newName = elem.value
    let UUID = elem.id.replace(newNamePrefix, "")
    files.get(UUID).NewName = newName
}

RenameSelectedFilesButton.on('click', function () {
    let selectedFiles = []
    let selectedItemElements = $('.select-item:checkbox:checked')

    selectedItemElements.each(function(i, elem) {
        let UUID = elem.id.replace(selectedPrefix, "")
        selectedFiles.push(files.get(UUID))
    });

    $.ajax({
        type: "post",
        url: '/rename-selected-files',
        data: JSON.stringify(selectedFiles),
        dataType: 'json',
        contentType: 'application/json',
    });
})

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


function UUIDv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}