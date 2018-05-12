//Loading tabs titles
var oReq = new XMLHttpRequest();
oReq.open("GET", "/modules/enabled_modules.json", true);
oReq.onload = function () { 
				console.log(this.responseText); 
				var enabledModules = JSON.parse(this.responseText);

				var tabs = document.getElementById('tab_modules');
				console.log(tabs)
				tabs.innerHTML = '';

				for (var module in enabledModules) {
								var li = document.createElement('li');
								li.classList.add("tablinks");
								tabs.appendChild(li);

								var anchor = document.createElement('a');
								anchor.addEventListener('click', handler, false);
								var text = document.createTextNode(enabledModules[module]);
								anchor.appendChild(text);
								tabs.appendChild(anchor);
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
				document.getElementById(tabTitle).style.display = "block";
				evt.path[1].className += " is-active";
}

function createTab(evt, tabTitle) {
				var li = document.createElement("LI");
				li.classList.add("tablinks");
				var a = document.createElement("A");
				a.addEventListener("click", "openTab(event, 'Name')");
				var name = document.createTextNode(tabTitle);
				a.appendChild(name);
				li.appendChild(a);
				document.getElementById("tabsTitle").appendChild(li);

				var div = document.createElement("DIV");
				div.classList.add("tabcontent");
				div.id = tabTitle;
				document.getElementById("tabs").appendChild(div)
}
