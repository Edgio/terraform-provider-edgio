package models

type EnvironmentsLinksModel struct {
	First    EnvironmentsLinkModel `tfsdk:"first"`
	Next     EnvironmentsLinkModel `tfsdk:"next"`
	Previous EnvironmentsLinkModel `tfsdk:"previous"`
	Last     EnvironmentsLinkModel `tfsdk:"last"`
}
