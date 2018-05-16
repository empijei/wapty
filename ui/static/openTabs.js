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
