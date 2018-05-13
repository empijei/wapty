//Loading tabs titles
var oReq = new XMLHttpRequest();
oReq.open("GET", "/modules/enabled_modules.json", true);
oReq.onload = function () { 
	//console.log(this.responseText); 
	var enabledModules = JSON.parse(this.responseText);

	var tabs = document.getElementById('tab_modules');
	tabs.innerHTML = '';

	for (var module in enabledModules) {
		var li = document.createElement('li');
		li.classList.add("tablinks");
		if (enabledModules[module] == "Proxy") {
			li.classList.add("is-active");
		}
		tabs.appendChild(li);

		var anchor = document.createElement('a');
		anchor.addEventListener('click', handler, false);
		var text = document.createTextNode(enabledModules[module]);
		anchor.appendChild(text);
		li.appendChild(anchor);
	}
};
oReq.send(null);

//Loding Proxy tab first
var req = new XMLHttpRequest();
req.open('GET', 'modules/Proxy/index.html', true);
req.onload = function () {
	var workspace = document.getElementById('workspace');
	workspace.innerHTML = this.responseText;
}
req.send(null);

//Loading tab content
function handler(e, data) {
	var module_name = e.toElement.innerHTML;
	//Removing is-active attributes
	var actTab = document.getElementsByClassName("is-active");
	actTab[0].classList.remove("is-active");
	//Adding is-active class to selected tab
	e.path[1].classList.add("is-active");

	var req = new XMLHttpRequest();
	req.open('GET', 'modules/' + module_name + '/index.html', true);
	req.onload = function () {
		var workspace = document.getElementById('workspace');
		workspace.innerHTML = this.responseText;
	}
	req.send(null);
}

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
	console.log(evt);
}
