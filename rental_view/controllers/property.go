package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	web.Controller
}

func (c *PropertyController) Index() {
	c.TplName = "index.tpl"
}

func (c *PropertyController) List() {
	// Mock data - replace with actual API call
	properties := []map[string]interface{}{
		{
			"id":           1,
			"title":        "Mid Superior Villa At Palm Paradise Luxury Villa",
			"price":        98791,
			"bedrooms":     4,
			"amenities":    []string{"Pool", "Beach View"},
			"rating":       4.9,
			"reviewCount":  24,
			"image":        "https://images.unsplash.com/photo-1582719478250-c89cae4dc85b",
			"badgeText":    "New",
		},
		{
			"id":           2,
			"title":        "Beaches Holiday Homes - Capital Bay",
			"price":        52000,
			"bedrooms":     2,
			"amenities":    []string{"City View", "Gym"},
			"rating":       4.8,
			"reviewCount":  16,
			"image":        "https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd",
			"badgeText":    "Featured",
		},
		{
			"id":           3,
			"title":        "Luxury Pool Access Apartment",
			"price":        3278,
			"bedrooms":     3,
			"amenities":    []string{"Pool", "Spa"},
			"rating":       4.7,
			"reviewCount":  32,
			"image":        "https://images.unsplash.com/photo-1594560913095-8cf34baf3a39",
			"badgeText":    "Popular",
		},
	}

	c.Data["json"] = properties
	c.ServeJSON()
}