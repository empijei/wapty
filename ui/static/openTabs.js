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
				a.addEventListener("onClick", "openTab(event, 'Name')");
				var name = document.createTextNode(tabTitle);
				a.appendChild(name);
				li.appendChild(a);
				document.getElementById("tabsTitle").appendChild(li);

				var div = document.createElement("DIV");
    div.classList.add("tabcontent");
    div.id = tabTitle;
		  document.getElementById("tabs").appendChild(div)
}
