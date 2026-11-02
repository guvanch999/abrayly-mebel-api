package dto

import (
	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
)

// CategoryResponseDto is a node of the category tree.
type CategoryResponseDto struct {
	ID            uuid.UUID             `json:"id"`
	ImageUrl      string                `json:"image_url"`
	NameTm        string                `json:"name_tm"`
	NameRu        string                `json:"name_ru"`
	NameEn        string                `json:"name_en"`
	ParentId      *uuid.UUID            `json:"parent_id"`
	Subcategories []CategoryResponseDto `json:"subcategories"`
}

func FromCategoryToResponseDto(category *models.Category) *CategoryResponseDto {
	return &CategoryResponseDto{
		ID:            category.ID,
		ImageUrl:      category.ImageUrl,
		NameTm:        category.NameTm,
		NameRu:        category.NameRu,
		NameEn:        category.NameEn,
		ParentId:      category.ParentId,
		Subcategories: []CategoryResponseDto{},
	}
}

// BuildCategoryTree turns the flat, self-referencing category list into a
// tree of arbitrary depth. Roots are categories without a parent; the input
// order is preserved among siblings. Categories whose parent is missing from
// the list (or that form a cycle) are not reachable from a root and are left
// out, so the function always terminates.
func BuildCategoryTree(categories []*models.Category) []*CategoryResponseDto {
	children := make(map[uuid.UUID][]*models.Category, len(categories))
	roots := make([]*models.Category, 0)
	for _, category := range categories {
		if category == nil {
			continue
		}
		if category.ParentId == nil {
			roots = append(roots, category)
			continue
		}
		children[*category.ParentId] = append(children[*category.ParentId], category)
	}

	visited := make(map[uuid.UUID]bool, len(categories))
	var build func(category *models.Category) CategoryResponseDto
	build = func(category *models.Category) CategoryResponseDto {
		visited[category.ID] = true
		node := FromCategoryToResponseDto(category)
		for _, child := range children[category.ID] {
			if visited[child.ID] {
				continue
			}
			node.Subcategories = append(node.Subcategories, build(child))
		}
		return *node
	}

	tree := make([]*CategoryResponseDto, 0, len(roots))
	for _, root := range roots {
		if visited[root.ID] {
			continue
		}
		node := build(root)
		tree = append(tree, &node)
	}
	return tree
}
