package webroot

func init() {
	webFiles[""] = indexHtml
}

const indexHtml = `
<HTML>
	<HEAD>
		<!--	<SCRIPT src="/angular.min.js"></SCRIPT> -->
		<SCRIPT src="/index.js"></SCRIPT>
		<SCRIPT src="/angular.min.js"></SCRIPT>
	</HEAD>
	<BODY>
		<div ng-app="app" ng-controller="ProxyCtrl">
			<button id="forwardOriginal" type="button" onclick="clickhandler();">Forward Original</button>
			<button id="forwardModified" type="button" onclick="clickhandler();">Forward Modified</button>
			<button id="drop" type="button" onclick="clickhandler();">Drop</button>
			<button id="provideResponse" type="button" onclickngular.min.js="clickhandler();">Provide Response</button>
			<label><input type="checkbox" id="interceptToggle" value="Intercept" onclick="clickhandler();">Intercept</label>


			<button id="menu" type="button">Action</button>
			<br>
			<br>

			<textarea id="proxybuffer" name="proxybuffer" cols="100" rows="40"></textarea>
		</div>

		<div ng-app="app" ng-controller="ProxyCtrl">
			<table border="1" id="historyTable">
				<tr>
					<th>ID</th>
					<th>Host</th>
					<th>Path</th>
				</tr>
			</table>
		</div>

	</BODY>
</HTML>
`
