function addMember() {
    var membersList = document.getElementById("members-list");
    var index = membersList.children.length;

    var newMemberDiv = document.createElement("div");
    newMemberDiv.classList.add("member", "mb-3");

    newMemberDiv.innerHTML = `
        <label for="member_id">User ID</label>
        <input type="text" name="Members[${index}][ID]" required>
        <label for="member_privilege">Privilege</label>
        <select name="Members[${index}][Privilege]">
            <option value="Admin">Admin</option>
            <option value="Member">Member</option>
        </select>
    `;
    membersList.appendChild(newMemberDiv);
}

function addColumn() {
    var columnsList = document.getElementById("column-list");
    var index = columnsList.children.length;

    var newColumnDiv = document.createElement("div");
    newColumnDiv.classList.add("member", "mb-3");
    newColumnDiv.id = `column${index}`
    
    newColumnDiv.innerHTML = `
    <label for="column${index}_name">Name</label>
    <input type="text" name="Columns[${index}][Name]" id="column${index}_name">
    <label for="column${index}_type">Type</label>
    <select name="Columns[${index}][Type]" id="column${index}_type" oninput="setType(${index})">
        <option value="text">Text</option>
        <option value="select">Dropdown</option>
        <option value="member">Member</option>
    </select>
    <fieldset id="column${index}_options">

    </fieldset>
    `;
    columnsList.appendChild(newColumnDiv);
}
function setType(index) {
    var type = document.getElementById(`column${index}_type`).value;
    switch (type) {
        case 'select' :
            addOption(index)
            break;

        default:
            removeOptions(index, 'all');
            document.getElementById(`column${index}_adoption`).remove();
            break;
    }

}

function addOption(index) {
    var options = document.getElementById(`column${index}_options`)
    var newOption = document.createElement('input') //`<input name="" type="text"></input>`
    var div = document.createElement('div');
    newOption.type = 'text';
    newOption.name = `Columns[${index}][Options[]]`
    div.appendChild(newOption);
    div.appendChild(document.createElement('br'));
    options.appendChild(div);
    if (!document.getElementById(`column${index}_adoption`)) {
        var column = document.getElementById(`column${index}`)
        var newBtn = document.createElement('btn');
        newBtn.id = `column${index}_adoption`;
        newBtn.onclick = function(){addOption(index)};
        newBtn.classList.add('btn', 'btn-secondary');
        newBtn.textContent = 'Add Option';
        column.appendChild(newBtn);
    }
}
function removeOptions(index, which) {
    var options = document.getElementById(`column${index}_options`);

    if (which == 'all') {
        options.innerHTML = '';
    } else {
        options.removeChild(options.querySelector(`#option${which}`));
    }
}

function addRow() {
    var table = document.getElementById('table');
    var newRow = document.createElement('tr');
    newRow.innerHTML = document.getElementById('row-template').innerHTML;
    table.appendChild(newRow);
}