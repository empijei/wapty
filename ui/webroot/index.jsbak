var intercept = {
			EDITORCHANNEL: "proxy/intercept/editor",
			SETTINGSCHANNEL: "proxy/intercept/options",
			HISTORYCHANNEL: "proxy/httpHistory",
}
var waptyServer= new WebSocket("ws://localhost:8081/ws");
waptyServer.onopen = function(event){
			console.log("WebSocket connected");
			var msg = {
						Action: "intercept",
						Channel: intercept.SETTINGSCHANNEL
			}
			waptyServer.send(JSON.stringify(msg));

}

//scope this
var tmpHistory = null;
var historyTbody = document.getElementById("historyTbody");
var interceptOn = false;
var historyReqBuffer = document.getElementById("historyReqBuffer");
var historyResBuffer = document.getElementById("historyResBuffer");

waptyServer.onmessage = function(event){
			//	console.log(event.data);
			msg = JSON.parse(event.data);
			//console.log(msg);
			switch (msg.Channel){
						case intercept.EDITORCHANNEL:
									//if ('Payload' in msg){
									document.getElementById("proxybuffer").value=atob(msg.Payload);
									document.getElementById("endpointIndicator").innerText=(msg.Args[0][2]==="q"?"Request for: ":"Response from: ")  + msg.Args[1];
									controls = true;
									//}
									break;
						case intercept.SETTINGSCHANNEL:
									switch (msg.Action){
												case "intercept":
															console.log(msg);
															btn = document.getElementById("interceptToggle");
															if (msg.Args[0] === "true"){
																		btn.className="btn btn-success";
																		btn.innerText="Intercept is on ";
																		interceptOn = true;
															}else{
																		btn.className="btn btn-danger";
																		btn.innerText="Intercept is off";
																		interceptOn = false;
															}
									}
									break;
						case intercept.HISTORYCHANNEL:
									switch (msg.Action){
												case "metaData":
															var metaData = JSON.parse(msg.Args[0])
															if (tmpHistory===null){
																		var h = document.getElementById("historyHeader")
																		for (var key in metaData) {
																					if (metaData.hasOwnProperty(key)) {
																								var cell = h.insertCell(-1)
																								cell.innerText = key
																					}
																		}
																		tmpHistory = {}
															}
															var stringID=""+metaData.Id;
															console.log("Got metaData for id " + stringID);
															if (stringID in tmpHistory){
																		tmp = tmpHistory[stringID]
																		for (var key in metaData) {
																					if (metaData.hasOwnProperty(key)) {
																								tmp[key].innerText=metaData[key]
																					}
																		}
																		//FIXME this is commented because it looks like the page
																		//receives the same metadata multiple times.
																		//delete tmpHistory[stringID]
															}else{
																		var row=historyTbody.insertRow(-1);
																		var tmp={}
																		for (var key in metaData) {
																					if (metaData.hasOwnProperty(key)) {
																								var cell = row.insertCell(-1)
																								cell.innerText = metaData[key]
																								tmp[key] = cell
																					}
																		}
																		tmpHistory[stringID]=tmp
															}
															/*case "metaData":*/
															//var problem = false;
															//var metaData = JSON.parse(msg.Args[0])
															//console.log("Metadata for request " + metaData.Id + " received:");
															//console.log(metaData)
															//if (""+metaData.Id in debugHistory){
															//if (debugHistory[""+metaData.Id]==1){
															//debugHistory[""+metaData.Id]=2
															//}else{
															//console.log("Problem with request " + metaData.Id);
															//problem=true;
															//}
															//}else{
															//debugHistory[""+metaData.Id]=1
															//}
															//document.getElementById("historyTable").innerHTML=document.getElementById("historyTable").innerHTML + "<tr"+
															//(problem?" style='color:red;' ":"")+
															//"><td>"+metaData.Id+"</td>"+
															//"<td>"+metaData.Host+"</td>"+
															//"<td>"+metaData.Path+"</td>"+
															//"</tr>";
															break;
												case "fetch":
															var pl = JSON.parse(atob(msg.Payload))
															console.log(atob(pl.RawReq))
															historyReqBuffer.innerText = atob(pl.RawReq)
															console.log(atob(pl.RawRes))
															historyResBuffer.innerText = atob(pl.RawRes)
															console.log(atob(pl.RawEditedReq))
															console.log(atob(pl.RawEditedRes))
															break;
									}
									break;
			}
}
waptyServer.onclose=function(event){
			var value = ("Server connection lost, would you like to try to reconnect?")
			if (value){
						location.reload()
			}
}

var controls = false;

function clickhandler(){
			if (!controls){
						return;
			}
			switch (event.target.id){
						case "forwardOriginal":
									var msg = {
												Action: "forward",
												Channel: intercept.EDITORCHANNEL
									}
									controls = false;
									document.getElementById("proxybuffer").value="";
									document.getElementById("endpointIndicator").innerText="";
									waptyServer.send(JSON.stringify(msg));
									break;
						case "forwardModified":
									var payload = btoa(document.getElementById("proxybuffer").value);
									var msg = {
												Action: "edit",
												Channel: intercept.EDITORCHANNEL,
												Payload: payload
									}
									controls = false;
									document.getElementById("proxybuffer").value="";
									document.getElementById("endpointIndicator").innerText="";
									waptyServer.send(JSON.stringify(msg));
									//var xhr = new XMLHttpRequest();
									//xhr.open("POST", "/edit", true);
									//xhr.setRequestHeader('Content-Type', 'application/json');
									//xhr.send(JSON.stringify(msg));
									break;
						case "drop":
									var msg = {
												Action: "drop",
												Channel: intercept.EDITORCHANNEL,
									}
									controls = false;
									document.getElementById("proxybuffer").value="";
									document.getElementById("endpointIndicator").innerText="";
									waptyServer.send(JSON.stringify(msg));
									break;
						case "provideResponse":
									var payload = btoa(document.getElementById("proxybuffer").value);

									var msg = {
												Action: 	"provideResp",
												Channel: intercept.EDITORCHANNEL,
												Payload: payload
									}
									controls = false;
									document.getElementById("proxybuffer").value="";
									document.getElementById("endpointIndicator").innerText="";
									waptyServer.send(JSON.stringify(msg));
									//var xhr = new XMLHttpRequest();
									//xhr.open("POST", "/edit", true);
									//xhr.setRequestHeader('Content-Type', 'application/json');
									//xhr.send(JSON.stringify(msg));
									break;
						default:
									console.log("unknown event")
			}
}
function toggler(){
			var msg = {
						Action: "intercept",
						Channel: intercept.SETTINGSCHANNEL,
						Args: [""+!interceptOn]
			}
			waptyServer.send(JSON.stringify(msg));
}
function fetchHistory(id){
			var msg = {
						Action: "fetch",
						Channel: intercept.HISTORYCHANNEL,
						Args: [""+id]
			}
			waptyServer.send(JSON.stringify(msg));
}
function historyTclick(){
			fetchHistory(event.target.parentNode.children[0].innerText)
}
