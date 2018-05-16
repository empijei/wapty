function openTab(evt, tabTitle) {
	var i, tabcontent, tablinks;
	tabcontent = document.getElementsByClassName("tabcontent");
	for (i = 0; i < tabcontent.length; i++) {
		tabcontent[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tablinks");
	for (i = 0; i < tablinks.length; i++) {
		tablinks[i].className = tablinks[i].className.replace(" is-active", "");
	};
	//console.log(document.getElementById(tabTitle));
	document.getElementById(tabTitle).style.display = "block";
	document.getElementById("li"+tabTitle).className += " is-active";
}

var tabNum = 1;

function createTab(evt, tabTitle) {
	//Adding a tab in the navbar
	tabNum = tabNum + 1;
	var tabName = 'Tab' + tabNum;
	var li = '<li class="tablinks" id="li'+tabName+'">';
	var a = '<a onclick="openTab(event, \'' + tabName + '\')">';
	var button = '<button class="delete is-small" onclick="closeTab(\''+tabName+'\')"></button>';
	var newtabTitle = document.getElementById('titles');
	newtabTitle.insertAdjacentHTML('beforeend', li + a + tabName + button + '</a></li>');

	//Adding the content of the tab
	var div = '<div class="tabcontent" id="'+tabName+'">';
	var textarea = '<textarea class="textarea" placeholder="New Tab" rows="30"></textarea>';
	var newtabContent = document.getElementById('contents');
	newtabContent.insertAdjacentHTML('beforeend', div + textarea + '</div>');
	openTab(evt, tabName)
}

function closeTab(tabName) {
	document.getElementById('li'+tabName).remove();
	document.getElementById(tabName).remove();
}
