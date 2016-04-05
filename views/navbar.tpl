{{define "navbar"}}
<a class = "navbar-brand" href="/"><b>Bee blog</b></a>
<div>
  <ul class="nav navbar-nav">
    <li {{if .IsHome}} class="active" {{end}}><a href="/">Homepage</a> </li>
    <li {{if .IsCategory}} class="active" {{end}}><a href="/category">Category</a> </li>
    <li {{if .IsTopic}} class="active" {{end}}><a href="/topic">Topic</a> </li>
   </ul>
 </div>
  <div class="pull-right">
  	<ul class="nav navbar-nav">
  		{{if .IsLogin}}
                            <li><a href="/login?exit=true">Exit</a></li>
  		{{else}}
                            <li><a href="/login">Login</a></li>
  		{{end}}
  	</ul>
  </div>  
{{end}}