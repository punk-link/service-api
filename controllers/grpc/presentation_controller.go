package grpc

import (
	"context"

	artistServices "main/services/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type Server struct {
	Injector *do.Injector
	presentationContracts.UnimplementedPresentationServer
}

func (t *Server) GetArtist(ctx context.Context, request *presentationContracts.ArtistRequest) (*presentationContracts.Artist, error) {
	grpcArtistService := t.getGrpcArtistService()
	return grpcArtistService.GetOne(request)
}

func (t *Server) GetRelease(ctx context.Context, request *presentationContracts.ReleaseRequest) (*presentationContracts.Release, error) {
	grpcReleaseService := t.getGrpcReleaseService()
	return grpcReleaseService.GetOne(request)
}

func (t *Server) getGrpcArtistService() *artistServices.GrpcArtistService {
	if _artistService == nil {
		_artistService = do.MustInvoke[*artistServices.GrpcArtistService](t.Injector)
	}

	return _artistService
}

func (t *Server) getGrpcReleaseService() *artistServices.GrpcReleaseService {
	if _artistService == nil {
		_releaseService = do.MustInvoke[*artistServices.GrpcReleaseService](t.Injector)
	}

	return _releaseService
}

var _artistService *artistServices.GrpcArtistService
var _releaseService *artistServices.GrpcReleaseService
