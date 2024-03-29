{{- define "pgHdrStart" }}
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no" />
    <link rel="stylesheet" type="text/css" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous" />
    <link rel="stylesheet" type="text/css" href="/static/lhlunch.css" />
{{- end }}

{{- define "pgHdrTitle" }}
    <title>{{.}}</title>
{{- end }}

<!-- Global site tag (gtag.js) - Google Analytics -->
{{- define "pgHdrScript" }}
    <script async src="https://www.googletagmanager.com/gtag/js?id={{.}}"></script>
    <script type="text/javascript" src="/static/lhlunch.js"></script>
    <script>
      gtag('js', new Date());
      gtag('config', '{{.}}');
    </script>
{{- end }}

{{- define "pgHdrEnd" }}
  </head>
{{- end }}

{{- define "bodyStart" }}
  <body class="bodybg">
    <div id="content" class="mt-3 mx-auto rounded">
{{- end }}

{{- define "bodyEnd" }}
    </div> <!-- div content -->
  </body>
</html>
{{- end }}

{{- define "p3" }}
      <div class="w-100 p-3">&nbsp;</div>
{{- end }}

{{- define "pb3" }}
      <div class="w-100 pb-3">&nbsp;</div>
{{- end }}

{{- define "default.html" }}
{{- template "pgHdrStart" -}}
{{- template "pgHdrTitle" "Lunch2Day" -}}
{{- template "pgHdrScript" .Gtag -}}
{{- template "pgHdrEnd" -}}
{{- template "bodyStart" -}}
{{- template "p3" }}
      <h1 class="pghdr h5 text-center">Choose location and format for lunch menu</h1>
{{- template "pb3" -}}
{{- template "pb3" }}
      <h2 class="dlistitem h6"><a href="html/se/gbg/lindholmen">Sweden / Gothenburg / Lindholmen (html)</a></h2>
      <h2 class="dlistitem h6"><a href="text/se/gbg/lindholmen">Sweden / Gothenburg / Lindholmen (text)</a></h2>
      <h2 class="dlistitem h6"><a href="json/se/gbg/lindholmen">Sweden / Gothenburg / Lindholmen (json)</a></h2>
{{- template "pb3" -}}
{{- template "bodyEnd" -}}
{{- end }}

{{- define "lunchlist.html" }}
{{- template "pgHdrStart" -}}
{{- template "pgHdrTitle" "Lunchlist countries" -}}
{{- template "pgHdrScript" .Gtag -}}
{{- template "pgHdrEnd" -}}
{{- template "bodyStart" -}}
{{- template "p3" }}
      <h1 class="pghdr h5 text-center">Lunchlist countries</h1>
{{- template "pb3" -}}
{{- template "pb3" }}
      {{- range .Countries }}
      <h2 class="dlistitem h6"><a href="{{.ID}}">{{.Name}}</a></h2>
      {{- end}}
{{- template "pb3" -}}
{{- template "bodyEnd" -}}
{{- end }}

{{- define "country.html" }}
{{- template "pgHdrStart" }}
    <title>Cities in {{.Name}}</title>
{{- template "pgHdrScript" .Gtag -}}
{{- template "pgHdrEnd" -}}
{{- template "bodyStart" -}}
{{- template "p3" }}
      <h1 class="pghdr h5 text-center">Cities in {{.Name}}</h1>
{{- template "pb3" -}}
{{- template "pb3" }}
      {{- range .Cities }}
      <h2 class="dlistitem h6"><a href="{{.ID}}">{{.Name}}</a></h2>
      {{- end}}
{{- template "pb3" -}}
{{- template "bodyEnd" -}}
{{- end }}

{{- define "city.html" }}
{{- template "pgHdrStart" }}
    <title>Sites in {{.Name}}</title>
{{- template "pgHdrScript" .Gtag -}}
{{- template "pgHdrEnd" -}}
{{- template "bodyStart" -}}
{{- template "p3" }}
      <h1 class="pghdr h5 text-center">Sites in {{.Name}}</h1>
{{- template "pb3" -}}
{{- template "pb3" }}
      {{- range .Sites }}
      <h2 class="dlistitem h6"><a href="{{.ID}}">{{.Name}}</a></h2>
      {{- end}}
{{- template "pb3" -}}
{{- template "bodyEnd" -}}
{{- end }}

{{- define "site.html" }}
{{- template "pgHdrStart" }}
    <title>Lunch @ {{.Name}} today</title>
{{- template "pgHdrScript" .Gtag -}}
{{- template "pgHdrEnd" -}}
{{- template "bodyStart" -}}
{{- template "p3" }}
      <h1 class="pghdr h5 text-center">Lunch @ {{.Name}} today</h1>
      <div class="toggledetails text-center mt-3" onclick="toggledetail();">[ Show / hide all ]</div>
{{- template "pb3" }}
      <div class="parsed ml-1">Updated: {{.ParsedHumanDate}}</div>
      {{- range .Restaurants }}
      <div class="restaurant m-2">
        <details open class="pb-3">
          <summary>
            <h2 class="name h6">
				<a href="{{.URL}}">{{.Name}}</a>
			{{- if .Address }}
				- <a href="{{.MapURL}}" target="_blank">{{.Address}}</a>
			{{- end }}
			</h2>
          </summary>
          <div class="dishes ml-2 p-2 shadow rounded">
            {{- range .Dishes }}
            <div class="dish m-2">
              <h3 class="name h6 d-inline">{{.Name}}</h3>
              <p class="desc d-inline">{{.Desc}}</p>
              <span class="price">{{.Price}}</span>
            </div> <!-- div dish -->
            {{- end}}
          </div> <!-- div dishes -->
        </details>
      </div> <!-- div restaurant -->
      {{- end}}
{{- template "pb3" -}}
{{- template "bodyEnd" -}}
{{- end }}


