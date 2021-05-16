const actualBtn = document.getElementById('actual-btn');

const fileChosen = document.getElementById('file-chosen');

actualBtn.addEventListener('change', function(){
    fileChosen.textContent = this.files.length + (this.files.length > 1 ? " files chosen" : " file chosen")
    let files = [];
    for (let i = 0; i < this.files.length; i++) {
        let file = this.files[i]
        files.push({
            'lastModified': file.lastModified,
            'lastModifiedDate': file.lastModifiedDate,
            'name': file.name,
            'size': file.size,
            'type': file.type
        })
    }
    $.post('/submit-files', JSON.stringify(files))
})

