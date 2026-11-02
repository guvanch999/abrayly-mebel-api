package service

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/dto"
	"github.com/guvanch999/abrayly-mebel-backend/internal/internal/banner/repository"
	"github.com/guvanch999/abrayly-mebel-backend/internal/models"
	"github.com/guvanch999/abrayly-mebel-backend/utils/util"
)

type BannerService struct {
	bannerRepository *repository.BannerRepository
}

func NewBannerService(bannerRepository *repository.BannerRepository) *BannerService {
	return &BannerService{bannerRepository: bannerRepository}
}

func (s *BannerService) Create(refLink *string, image *multipart.FileHeader) (*dto.BannerResponse, error) {
	banner := &models.Banner{
		ID:   uuid.New(),
		Link: refLink,
	}

	imagePath, err := util.SaveUploadedFile(image, banner.ID.String(), "uploads/banners")
	if err != nil {
		return nil, err
	}
	banner.Image = imagePath
	createdBanner, err := s.bannerRepository.Create(banner)
	if err != nil {
		_ = util.DeleteFile(imagePath)
		return nil, err
	}
	return dto.BannerToResponse(createdBanner), nil
}

func (s *BannerService) GetAll() ([]*dto.BannerResponse, error) {
	banners, err := s.bannerRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return dto.BannerListToResponse(banners), nil
}

func (s *BannerService) GetByID(id uuid.UUID) (*dto.BannerResponse, error) {
	banner, err := s.bannerRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return dto.BannerToResponse(banner), nil
}

func (s *BannerService) Update(id uuid.UUID, refLink *string) (*dto.BannerResponse, error) {
	banner, err := s.bannerRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	banner.Link = refLink
	updatedBanner, err := s.bannerRepository.Update(banner)
	if err != nil {
		return nil, err
	}
	return dto.BannerToResponse(updatedBanner), nil
}

func (s *BannerService) Delete(id uuid.UUID) error {
	banner, err := s.bannerRepository.GetByID(id)
	if err != nil {
		return err
	}
	err = s.bannerRepository.Delete(id)
	if err != nil {
		return err
	}
	_ = util.DeleteFile(banner.Image)
	return nil
}
