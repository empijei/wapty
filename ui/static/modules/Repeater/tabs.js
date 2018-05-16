function openTab(evt, tabTitle) {
	var i, tabcontent, tablinks;
	tabcontent = document.getElementsByClassName("tabcontent");
	for (i = 0; i < tabcontent.length; i++) {
		tabcontent[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tablinks");
	for (i = 0; i < tablinks.length; i++) {
		tablinks[i].className = tablinks[i].className.replace(" is-active", "");
	}
	//console.log(document.getElementById(tabTitle));
	document.getElementById(tabTitle).style.display = "block";
	evt.path[1].className += " is-active";
}

function createTab(evt, tabTitle) {
	//Adding a tab in the navbar
	var newtabTitle = document.getElementById('titles');
	newtabTitle.insertAdjacentHTML('beforeend', '<li class="tablinks"><a onclick="openTab(event, \'Name\')">Name <button class="delete is-small" onclick="closeTab(event)"></button></a></li>');

	//Adding the content of the tab
	//FIXME: all contents are shown at the same time
	var newtabContent = document.getElementById('contents');
	newtabContent.insertAdjacentHTML('beforeend', '<div class="tabcontent" id="Name"><textarea class="textarea" placeholder="New Tab" rows="30"></textarea></div>');
}

function closeTab(evt) {
	evt.path[2].remove();
}

