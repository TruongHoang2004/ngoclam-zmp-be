package domain

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type Product struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Price       int64             `json:"price,omitempty"`
	Variants    *[]ProductVariant `json:"variants,omitempty"`
	Images      *[]ProductImage   `json:"images,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ProductImage struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Product   *Product  `json:"product,omitempty"`
	ImageID   uint      `json:"image_id"`
	Image     *Image    `json:"image,omitempty"`
	Order     int       `json:"order"`
	IsMain    bool      `json:"is_main"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	Stock     int64     `json:"stock,omitempty"`
	Order     int       `json:"order,omitempty"`
	Price     int64     `json:"price,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProductFromModel(m *model.Product) *Product {
	return &Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		// Variants will be set separately
	}
}

func (p *Product) SetVariants(variants []*model.ProductVariant) {
	var variantList []ProductVariant
	for _, v := range variants {
		variantList = append(variantList, ProductVariant{
			ID:        v.ID,
			ProductID: v.ProductID,
			Name:      v.Name,
			Stock:     v.Stock,
			Price:     v.Price,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	p.Variants = &variantList
}

func (p *Product) SetImages(images []*model.ProductImage) {
	if len(images) == 0 {
		empty := make([]ProductImage, 0)
		p.Images = &empty
		return
	}

	var imageList []ProductImage
	for _, img := range images {
		if img == nil {
			continue
		}
		imageList = append(imageList, *NewProductImageFromModel(img))
	}
	p.Images = &imageList
}

func (p *Product) ToModel() *model.Product {
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (pv *ProductVariant) ToModel(productID uint) *model.ProductVariant {
	return &model.ProductVariant{
		ID:        pv.ID,
		ProductID: productID,
		Name:      pv.Name,
		Stock:     pv.Stock,
		Price:     pv.Price,
		CreatedAt: pv.CreatedAt,
		UpdatedAt: pv.UpdatedAt,
	}
}

func NewProductImageFromModel(m *model.ProductImage) *ProductImage {
	if m == nil {
		return nil
	}

	pi := &ProductImage{
		ID:        m.ID,
		ProductID: m.ProductID,
		ImageID:   m.ImageID,
		Order:     m.Order,
		IsMain:    m.IsMain,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return pi
}

func (pi *ProductImage) ToModel() *model.ProductImage {
	return &model.ProductImage{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		ImageID:   pi.ImageID,
		Order:     pi.Order,
		IsMain:    pi.IsMain,
	}
}
