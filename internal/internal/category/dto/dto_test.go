package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/validator"
)

func category(name string, parent *models.Category) *models.Category {
	c := &models.Category{ID: uuid.New(), NameEn: name, NameRu: name + "-ru", NameTm: name + "-tm"}
	if parent != nil {
		c.ParentId = &parent.ID
	}
	return c
}

// names flattens a tree into "Parent>Child" paths for easy comparison.
func names(nodes []CategoryResponseDto, prefix string) []string {
	var out []string
	for _, n := range nodes {
		path := prefix + n.NameEn
		out = append(out, path)
		out = append(out, names(n.Subcategories, path+">")...)
	}
	return out
}

func treeNames(tree []*CategoryResponseDto) []string {
	flat := make([]CategoryResponseDto, 0, len(tree))
	for _, n := range tree {
		flat = append(flat, *n)
	}
	return names(flat, "")
}

func TestBuildCategoryTreeNestsAnyDepth(t *testing.T) {
	bedroom := category("Bedroom", nil)
	beds := category("Beds", bedroom)
	doubleBeds := category("Double beds", beds)
	wardrobes := category("Wardrobes", bedroom)
	kitchen := category("Kitchen", nil)

	// Children listed before their parents must still be attached.
	tree := BuildCategoryTree([]*models.Category{doubleBeds, bedroom, beds, kitchen, wardrobes})

	got := strings.Join(treeNames(tree), " | ")
	want := "Bedroom | Bedroom>Beds | Bedroom>Beds>Double beds | Bedroom>Wardrobes | Kitchen"
	if got != want {
		t.Errorf("tree = %s\nwant   %s", got, want)
	}

	if len(tree) != 2 {
		t.Fatalf("roots = %d, want 2", len(tree))
	}
	if tree[0].ParentId != nil {
		t.Error("root must have a nil parent_id")
	}
	if child := tree[0].Subcategories[0]; child.ParentId == nil || *child.ParentId != bedroom.ID {
		t.Errorf("child parent_id = %v, want %s", child.ParentId, bedroom.ID)
	}
	if tree[0].NameTm != "Bedroom-tm" || tree[0].NameRu != "Bedroom-ru" {
		t.Error("trilingual names were not copied")
	}
}

func TestBuildCategoryTreeLeavesHaveEmptySubcategories(t *testing.T) {
	tree := BuildCategoryTree([]*models.Category{category("Kitchen", nil)})

	body, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(body), `"subcategories":[]`) {
		t.Errorf("leaf should serialize subcategories as [] (got %s)", body)
	}
}

func TestBuildCategoryTreeIgnoresOrphansAndCycles(t *testing.T) {
	root := category("Root", nil)
	orphan := category("Orphan", &models.Category{ID: uuid.New()})

	a := category("A", nil)
	b := category("B", a)
	a.ParentId = &b.ID // A <-> B cycle, unreachable from any root

	tree := BuildCategoryTree([]*models.Category{root, orphan, a, b, nil})

	if got := strings.Join(treeNames(tree), " | "); got != "Root" {
		t.Errorf("tree = %q, want only the reachable root", got)
	}
}

func TestBuildCategoryTreeEmpty(t *testing.T) {
	tree := BuildCategoryTree(nil)
	if tree == nil || len(tree) != 0 {
		t.Errorf("BuildCategoryTree(nil) = %#v, want an empty non-nil slice", tree)
	}
}

func TestCategoryDtoValidation(t *testing.T) {
	v := validator.NewValidator()
	parent := uuid.New()

	tests := []struct {
		name    string
		dto     CategoryDto
		wantErr bool
	}{
		{name: "root category", dto: CategoryDto{NameTm: "tm", NameRu: "ru", NameEn: "en"}},
		{name: "subcategory", dto: CategoryDto{NameTm: "tm", NameRu: "ru", NameEn: "en", ParentId: &parent}},
		{name: "missing Turkmen name", dto: CategoryDto{NameRu: "ru", NameEn: "en"}, wantErr: true},
		{name: "missing Russian name", dto: CategoryDto{NameTm: "tm", NameEn: "en"}, wantErr: true},
		{name: "missing English name", dto: CategoryDto{NameTm: "tm", NameRu: "ru"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(&tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
