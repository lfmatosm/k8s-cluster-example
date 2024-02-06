package controllers

import (
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/models"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/services"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/utils"
	"net/http"
)

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{service}
}

func (p *PostController) Save(headers http.Header, bytes []byte) *utils.HttpResponse {
	if len(bytes) == 0 {
		return utils.BadRequest(map[string]interface{}{"error": "Byte array must be non-empty"})
	}

	var dto *models.PostDTO = &models.PostDTO{
		Mime:  headers.Get("Content-Type"),
		Image: bytes,
	}

	_, err := p.service.Save(dto)
	if err != nil {
		return utils.InternalServerError(map[string]interface{}{"error": err.Error()})
	}

	return utils.NoContent()
}

func (p *PostController) List() *utils.HttpResponse {
	posts, err := p.service.List()
	if err != nil {
		return utils.InternalServerError(map[string]interface{}{"error": err.Error()})
	}

	return utils.Ok(posts)
}
