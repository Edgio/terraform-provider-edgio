package models

type PropertiesLinksModel struct {
	First    PropertiesLinkModel `tfsdk:"first"`
	Next     PropertiesLinkModel `tfsdk:"next"`
	Previous PropertiesLinkModel `tfsdk:"previous"`
	Last     PropertiesLinkModel `tfsdk:"last"`
}
