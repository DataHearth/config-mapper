<a name="{{ (index .Versions 0).Tag.Name }}"></a>
## {{ if (index .Versions 0).Tag.Previous }}[{{ (index .Versions 0).Tag.Name }}]{{ else }}{{ (index .Versions 0).Tag.Name }}{{ end }} - {{ datetime "2006-01-02" (index .Versions 0).Tag.Date }}
{{ range (index .Versions 0).CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if (index .Versions 0).RevertCommits -}}
### Reverts
{{ range (index .Versions 0).RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if (index .Versions 0).MergeCommits -}}
### Pull Requests
{{ range (index .Versions 0).MergeCommits -}}
- {{ .Header }}
{{ end }}
{{ end -}}

{{- if (index .Versions 0).NoteGroups -}}
{{ range (index .Versions 0).NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}