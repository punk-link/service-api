package presentation

import "context"

type Server struct {
	UnimplementedPresentationServer
}

func (t *Server) GetArtist(ctx context.Context, request *ArtistRequest) (*ArtistResponse, error) {
	return &ArtistResponse{
		Id: request.Id,
	}, nil
}
