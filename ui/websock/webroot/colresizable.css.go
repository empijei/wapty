package webroot

func init() {
	webFiles["colresizable.css"] = colresizableCss
}

const colresizableCss = `
* {
	cursor:default;
	margin:0;
	padding:0}

html,body {
	overflow-y:auto}


a {
	text-decoration:none;
	cursor:pointer;
	border-bottom:1px dotted #888;
	font-style:italic;
	color:#222}

ul {
	margin-left:20px;
	margin-top:10px}

li {
	list-style:circle}

.abs {
	position:absolute;
	top:auto;
	left:auto;
	right:auto;
	bottom:auto}

.align_center {
	text-align:center}

.align_right {
	text-align:right}

.float_left {
	float:left}

img {
	-ms-interpolation-mode:bicubic}

.center {
	text-align:left;
	margin-left:auto;
	margin-right:auto;
	width:910px;
	height:100%}

.content {
	width:910px;
	top:251px;
	bottom:0;
}

.section {
	width:100%;
	background-image:url(../img/inner.png)}

.button {
	line-height:30px;
	text-align:center;
	display:inline-block;
	height:30px;
	cursor:pointer}

.button:hover {
	text-shadow:#000 2px 2px 2px}

.dropItem {
	display:block;
	width:234px;
	text-indent:10px;
	height:30px;
	margin-left:27px}

.dropItem:hover {
	background-image:url(../img/transBlack.png);
	cursor:pointer}

i {
	font-weight:bold}

h2 {
	font-family:Arial, Helvetica, sans-serif;
	font-size:18px;
	padding-top:20px;
	padding-bottom:10px}

h3 {
	display:inline-block}

h4 {
	font-family:Arial, Helvetica, sans-serif;
	font-size:13px;
	font-weight:normal;
	font-style:italic;
	display:inline-block;
	text-indent:8px}

.side {
	height:161px;
	width:196px;
	background-image:url(../img/side.png);
	position:absolute;
	left:40px;
	top:0}

.side .bound {
	position:absolute;
	top:-24px}

.side a {
	display:block;
	margin-left:19px;
	width:110px;
	font-style:normal;
	font-family:'Trebuchet MS', Verdana, Helvetica, Arial, sans-serif;
	font-size:12px;
	line-height:1.4;
	text-decoration:none;
	border-bottom:1px solid #EEE;
	padding:.2em 10px}

.code {
	border:1px solid #aaa;
	border-radius:5px;
	display:block;
	padding:15px;
	overflow:hidden;
	}
	
pre{
	display:block;
}
.overflowContainer{
    width:100%;
    position:relative;
    overflow-x: scroll;
}
.whiteGradient{
    position:relative;
    width: 100%;
}
.whiteGradient:after{
    display:block;
    content:"";
    width:50px;
    height:100%;
    position:absolute;
    top:0px;
    right:0px;
    z-index:5;
    pointer-events: none;
    background: -moz-linear-gradient(left,  rgba(250,250,250,0) 0%, rgba(250,250,250,1) 93%, rgba(250,250,250,1) 100%); /* FF3.6-15 */
    background: -webkit-linear-gradient(left,  rgba(250,250,250,0) 0%,rgba(250,250,250,1) 93%,rgba(250,250,250,1) 100%); /* Chrome10-25,Safari5.1-6 */
    background: linear-gradient(to right,  rgba(250,250,250,0) 0%,rgba(250,250,250,1) 93%,rgba(250,250,250,1) 100%); /* W3C, IE10+, FF16+, Chrome26+, Opera12+, Safari7+ */
    filter: progid:DXImageTransform.Microsoft.gradient( startColorstr='#00fafafa', endColorstr='#fafafa',GradientType=1 ); /* IE6-9 */



}

.sampleOptions {
	border:1px solid #ccc;
	border-radius:2px;
	box-shadow:2px 2px 2px #AAA;
	padding:3px}

.sampleOptions a {
	border:none;
	text-decoration:none;
	text-indent:21px;
	background-repeat:no-repeat;
	display:inline-block;
	margin-right:7px;
	margin-left:5px;
	font-style:normal}

th {
	height:30px;
	background-repeat:no-repeat;
	text-shadow:#012b4d 2px 2px 2px;
	text-align:center}

td {
	text-indent:5px;
	border-bottom:1px solid #bbb;
	border-left:1px solid #bbb}

td.left {
	border-left:1px solid #2e638e}

td.right {
	border-right:1px solid #2e638e}

td.bottom {
	border-bottom:1px solid #2e638e}

.grip {
	width:20px;
	height:30px;
	margin-top:-3px;
	margin-left:-5px;
	position:relative;
	z-index:88;
	cursor:e-resize}

.grip:hover {
	background-position-x:-20px}

.dragging .grip {
	background-position-x:-40px}

.grip2{
	width:20px;
	height:15px;
	margin-top:-3px;
	margin-left:-5px;
	position:relative;
	z-index:88;
	cursor:e-resize;
    background-repeat: no-repeat;
}

.grip2:hover{
	background-position-x:-20px;
}

.dragging .grip2{
	background-position-x:-40px;
}
.JCLRLastGrip .grip2{
    background-position-y:-18px;
    left:-2px;
}	
	
.sampleText {
	position:relative;
	width:100%}

.dotted {
	background-repeat:repeat-y}


label.check {
	margin-left:6px}

.rangeGrip {
	width:10px;
	height:19px;
	cursor:e-resize;
	background-image:url(../img/slider.png);
	z-index:8}

.rangeDrag .rangeGrip,.rangeGrip:hover {
	background-position:right}

a.dwn {
	width:19px;
	height:20px;
	background-image:url(../img/dwn.png);
	outline:none;
	display:block;
	text-decoration:none;
	border:none;
	margin:3px 3px 3px 22px}

.ui-helper-hidden-accessible {
	position:absolute!important;
	clip:rect(1px,1px,1px,1px)}

.ui-helper-reset {
	border:0;
	outline:0;
	line-height:1.3;
	text-decoration:none;
	font-size:100%;
	list-style:none;
	margin:0;
	padding:0}

.ui-helper-clearfix:after {
	content:".";
	display:block;
	height:0;
	clear:both;
	visibility:hidden}

* html .ui-helper-clearfix {
	height:1%}

.ui-helper-zfix {
	width:100%;
	height:100%;
	top:0;
	left:0;
	position:absolute;
	opacity:0;
	filter:Alpha(Opacity=0)}

.ui-icon {
	display:block;
	text-indent:-99999px;
	width:16px;
	height:16px;
	background-image:url(../img/resizer.png);
	background-position:bottom right;
	background-repeat:no-repeat}

.ui-resizable {
	position:relative}

.ui-resizable-handle {
	position:absolute;
	font-size:.1px;
	z-index:99999;
	display:block}

.ui-resizable-n {
	cursor:n-resize;
	height:7px;
	width:100%;
	top:-5px;
	left:0}

.ui-resizable-s {
	cursor:s-resize;
	height:7px;
	width:100%;
	bottom:-5px;
	left:0}

.ui-resizable-e {
	cursor:e-resize;
	width:7px;
	right:-5px;
	top:0;
	height:100%}

.ui-resizable-w {
	cursor:w-resize;
	width:7px;
	left:-5px;
	top:0;
	height:100%}

.ui-resizable-se {
	cursor:se-resize;
	width:12px;
	height:12px;
	right:1px;
	bottom:1px}

.ui-resizable-sw {
	cursor:sw-resize;
	width:9px;
	height:9px;
	left:-5px;
	bottom:-5px}

.ui-resizable-nw {
	cursor:nw-resize;
	width:9px;
	height:9px;
	left:-5px;
	top:-5px}

.ui-resizable-ne {
	cursor:ne-resize;
	width:9px;
	height:9px;
	right:-5px;
	top:-5px}

ol.linenums {
	margin-top:0;
	margin-bottom:0}

li.L0,li.L1,li.L2,li.L3,li.L5,li.L6,li.L7,li.L8 {
	list-style-type:none}

li.L1,li.L3,li.L5,li.L7,li.L9 {
	background:#eee}

.float_right,#sample2Txt {
	float:right}

.sampleCode,#links,.ui-helper-hidden,.ui-resizable-disabled .ui-resizable-handle,.ui-resizable-autohide .ui-resizable-handle {
	display:none}

#features{
	clear:both;
}
.bluebtn {
  margin-top: 20px;
  float: right;
  background: #3498db;
  background-image: -webkit-linear-gradient(top, #3498db, #2980b9);
  background-image: -moz-linear-gradient(top, #3498db, #2980b9);
  background-image: -ms-linear-gradient(top, #3498db, #2980b9);
  background-image: -o-linear-gradient(top, #3498db, #2980b9);
  background-image: linear-gradient(to bottom, #3498db, #2980b9);
  -webkit-border-radius: 28;
  -moz-border-radius: 28;
  border-radius: 28px;
  font-size: 14px;
  padding: 7px 15px 7px 15px;
  text-decoration: none;
  font-style:normal;
}

.bluebtn:hover {
  background: #3cb0fd;
  background-image: -webkit-linear-gradient(top, #3cb0fd, #3498db);
  background-image: -moz-linear-gradient(top, #3cb0fd, #3498db);
  background-image: -ms-linear-gradient(top, #3cb0fd, #3498db);
  background-image: -o-linear-gradient(top, #3cb0fd, #3498db);
  background-image: linear-gradient(to bottom, #3cb0fd, #3498db);
  text-decoration: none;
}
#fiddle{
  margin-left:5px;
 }
 
 .repo, .repo img{
   outline:none !important;
   border:none;
   cursor:pointer !important;
   display:inline-block;
   margin-bottom:20px;
 }
 .git{
	margin-left:8px;
 }
.bower{
    position:absolute;
    top: 41px;
    left: 121px;
}

.crosslinks a{
    text-decoration: none;
    border-bottom: none;
    font-size: 9px;
}



/* Ultimo update */

@-webkit-keyframes swing {
    10% {
        -webkit-transform: rotate3d(0, 0, 1, 15deg);
        transform: rotate3d(0, 0, 1, 15deg)
    }
    30% {
        -webkit-transform: rotate3d(0, 0, 1, -10deg);
        transform: rotate3d(0, 0, 1, -10deg)
    }
    40% {
        -webkit-transform: rotate3d(0, 0, 1, 5deg);
        transform: rotate3d(0, 0, 1, 5deg)
    }
    60% {
        -webkit-transform: rotate3d(0, 0, 1, -5deg);
        transform: rotate3d(0, 0, 1, -5deg)
    }
    80% {
        -webkit-transform: rotate3d(0, 0, 1, 0deg);
        transform: rotate3d(0, 0, 1, 0deg)
    }
    100% {
        -webkit-transform: rotate3d(0, 0, 1, 0deg);
        transform: rotate3d(0, 0, 1, 0deg)
    }
 }

@keyframes swing {
    10% {
        -webkit-transform: rotate3d(0, 0, 1, 15deg);
        transform: rotate3d(0, 0, 1, 15deg)
    }
    30% {
        -webkit-transform: rotate3d(0, 0, 1, -10deg);
        transform: rotate3d(0, 0, 1, -10deg)
    }
    40% {
        -webkit-transform: rotate3d(0, 0, 1, 5deg);
        transform: rotate3d(0, 0, 1, 5deg)
    }
    60% {
        -webkit-transform: rotate3d(0, 0, 1, -5deg);
        transform: rotate3d(0, 0, 1, -5deg)
    }
    80% {
        -webkit-transform: rotate3d(0, 0, 1, 0deg);
        transform: rotate3d(0, 0, 1, 0deg)
    }
    100% {
        -webkit-transform: rotate3d(0, 0, 1, 0deg);
        transform: rotate3d(0, 0, 1, 0deg)
    }
}
.swing {
    -webkit-transform-origin: top center;
    transform-origin: top center;
    -webkit-animation: swing 4s infinite; 
    animation: swing 4s infinite;
    animation-delay: 1s;
}

`
