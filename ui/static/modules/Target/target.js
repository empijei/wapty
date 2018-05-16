function newTarget() {
	//get target content
	var content = document.getElementById("targetfield").value;

	//insert new row in table
	var checkbox = '<td width="5%"><input type="checkbox" name="active" checked></i></td>'
	var targetText = '<td>'+content+'</td>'
	var deleteBtn = '<td><a class="button is-small is-danger" onclick="deleteRow(event)">Delete</a></td>'

	var newRow = document.getElementById("targetbody");
	newRow.insertAdjacentHTML('beforeend', '<tr>'+checkbox+targetText+deleteBtn+'</tr>');
}

function deleteRow(evt) {
	evt.path[2].remove();
}
