package grpc

import (
	"context"

	presentationGrpcs "main/grpc/presentations"
	artistServices "main/services/artists"
	artistConverter "main/services/artists/converters"

	"github.com/samber/do"
)

type Server struct {
	Injector *do.Injector
	presentationGrpcs.UnimplementedPresentationServer
}

func (t *Server) GetArtist(ctx context.Context, request *presentationGrpcs.ArtistRequest) (*presentationGrpcs.ArtistResponse, error) {
	artistService := t.getArtistService()
	artist, _ := artistService.GetOneWithReleases(int(request.Id))

	return artistConverter.ToArtistResponseMessage(artist), nil
}

func (t *Server) getArtistService() *artistServices.ArtistService {
	if _artistService == nil {
		_artistService = do.MustInvoke[*artistServices.ArtistService](t.Injector)
	}

	return _artistService
}

var _artistService *artistServices.ArtistService
