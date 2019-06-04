{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "entity.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "entity.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "entity.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Make image.
GOLANG VERSION must be >= 1.11
!!! The template expect next params in the its scope:
  .repository
  .name
  .tag
*/}}
{{- define "entity.image" -}}
  {{- printf "%s/%s:%s" .repository .name .tag -}}
{{- end -}}

{{/*
Common annotations
*/}}
{{- define "entity.annotations" }}
app.kubernetes.io/name: {{ include "entity.name" . | quote }}
helm.sh/chart: {{ include "entity.chart" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end -}}


{{/*
common labels
*/}}
{{- define "entity.common.labels" }}
app: {{ .Release.Name | quote }}
{{- end -}}


{{/*
labels
repeat entity.common.labels
!!! .module prm must be set in the template scope
*/}}
{{- define "entity.labels" }}
app: {{ .releaseName | quote }}
module: {{ .module | quote }}
{{- end -}}
