{{/*
Generate a full name for the app
*/}}
{{- define "recipe-app.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}
