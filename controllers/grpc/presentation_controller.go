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

func (t *Server) GetArtist(ctx context.Context, request *presentationGrpcs.ArtistRequest) (*presentationGrpcs.Artist, error) {
	artistService := t.getArtistService()
	artist, err := artistService.GetOneWithReleases(int(request.Id))
	if err != nil {
		return nil, err
	}

	return artistConverter.ToArtistMessage(artist), nil
}

func (t *Server) GetRelease(ctx context.Context, request *presentationGrpcs.ReleaseRequest) (*presentationGrpcs.Release, error) {
	releaseService := t.getReleaseService()
	release, err := releaseService.GetOne(int(request.Id))
	if err != nil {
		return nil, err
	}

	return artistConverter.ToReleaseMessage(release), nil
}

func (t *Server) getArtistService() *artistServices.ArtistService {
	if _artistService == nil {
		_artistService = do.MustInvoke[*artistServices.ArtistService](t.Injector)
	}

	return _artistService
}

func (t *Server) getReleaseService() *artistServices.ReleaseService {
	if _artistService == nil {
		_releaseService = do.MustInvoke[*artistServices.ReleaseService](t.Injector)
	}

	return _releaseService
}

var _artistService *artistServices.ArtistService
var _releaseService *artistServices.ReleaseService
