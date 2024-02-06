package services

import (
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/models"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	repository repositories.Repository
	collection string
}

func NewPostService(repository repositories.Repository, collection string) *PostService {
	return &PostService{repository, collection}
}

func (p *PostService) Save(dto *models.PostDTO) (bool, error) {
	post := &models.Post{
		Id:   primitive.NewObjectID(),
		Mime: dto.Mime,
		Image: primitive.Binary{
			Data:    dto.Image,
			Subtype: 0x00,
		},
		LastUpdated: time.Now(),
	}
	return p.repository.Save(p.collection, post)
}

func (p *PostService) List() (*[]models.PostDTO, error) {
	var rawResp []models.Post
	if err := p.repository.List(p.collection, &rawResp); err != nil {
		return nil, err
	}

	var resp []models.PostDTO = []models.PostDTO{}
	for i := range rawResp {
		resp = append(resp, models.PostDTO{
			Mime:  rawResp[i].Mime,
			Image: rawResp[i].Image.Data,
		})
	}

	return &resp, nil
}
