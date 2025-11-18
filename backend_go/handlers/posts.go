package handlers

import (
	"math"

	"backend_go/db"
	"backend_go/models"

	"github.com/gofiber/fiber/v2"
)

// ListPosts returns paginated posts.
func ListPosts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 12)
	if limit < 1 {
		limit = 12
	}

	var posts []models.Post
	var total int64

	db.DB.Model(&models.Post{}).Count(&total)
	db.DB.Order("created_at desc").Limit(limit).Offset((page - 1) * limit).Find(&posts)

	pages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"posts": posts,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"count": total,
			"pages": pages,
			"prev_page": func() interface{} {
				if page > 1 {
					return page - 1
				}
				return nil
			}(),
			"next_page": func() interface{} {
				if page < pages {
					return page + 1
				}
				return nil
			}(),
		},
	})
}
