package presentation

import (
	"context"

	artistServices "main/services/artists"

	"github.com/samber/do"
)

type Server struct {
	Injector *do.Injector
	UnimplementedPresentationServer
}

func (t *Server) GetArtist(ctx context.Context, request *ArtistRequest) (*ArtistResponse, error) {
	service := t.getArtistService()
	album, _ := service.GetOne(int(request.Id))

	return &ArtistResponse{
		Id:      int32(album.Id),
		LabelId: int32(album.LabelId),
		Name:    album.Name,
		ImageDetails: &ImageDetails{
			AltText: album.ImageDetails.AltText,
			Height:  int32(album.ImageDetails.Width),
			Url:     album.ImageDetails.Url,
			Width:   int32(album.ImageDetails.Width),
		},
	}, nil
}

func (t *Server) getArtistService() *artistServices.ArtistService {
	if _artistService == nil {
		_artistService = do.MustInvoke[*artistServices.ArtistService](t.Injector)
	}

	return _artistService
}

var _artistService *artistServices.ArtistService
