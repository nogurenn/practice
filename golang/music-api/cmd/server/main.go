package main

import (
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kelseyhightower/envconfig"
	"github.com/nogurenn/practice/golang/music-api/internal/artist"
	"github.com/nogurenn/practice/golang/music-api/internal/genre"
	"github.com/nogurenn/practice/golang/music-api/internal/server"
	"github.com/nogurenn/practice/golang/music-api/internal/spotify"
	"github.com/nogurenn/practice/golang/music-api/internal/transaction"
)

func main() {
	cfg, err := newServerConfig()
	if err != nil {
		panic(err)
	}

	spotifyClient := spotify.NewClient(cfg.SpotifyConfig)

	db, err := transaction.NewPostgresConnection(cfg.PostgresConfig)
	if err != nil {
		// I think we should panic here, because our endpoints depend on the database connection
		// in relation to Spotify integration costs and rate limits.
		// Otherwise, I am okay with logging the error and continuing the execution,
		// and returning an empty list of genres until the issue is resolved.
		// Good use case for feature flags at the response and/or dependency levels.
		panic(err)
	}

	genreRepo := genre.NewPostgresRepository(db)
	genreService := genre.NewService(genreRepo)
	genreHandler := server.NewGenreHandler(genreService)

	artistRepo := artist.NewPostgresRepository(db)
	artistService := artist.NewService(artistRepo, genreRepo, spotifyClient)
	artistHandler := server.NewArtistHandler(artistService)

	// I think I could do better design-wise here.
	router := server.NewConvenientRouter()
	router.Handle("/artists/", artistHandler)
	router.Handle("/genres/", genreHandler)

	fmt.Printf("Server is running on port %s\n", cfg.HTTPConfig.Port)
	fmt.Printf("Access the API at localhost:%s\n", cfg.HTTPConfig.Port)
	fmt.Println("Press CTRL+C to stop the server")
	http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPConfig.Port), router)
}

type serverConfig struct {
	PostgresConfig *transaction.PostgresConfig
	SpotifyConfig  *spotify.Config
	HTTPConfig     *httpConfig
}

type httpConfig struct {
	Port string `envconfig:"HTTP_SERVER_PORT" default:"3000"`
}

func newServerConfig() (*serverConfig, error) {
	postgresConfig := &transaction.PostgresConfig{}
	err := envconfig.Process("POSTGRES", postgresConfig)
	if err != nil {
		return nil, err
	}

	spotifyConfig := &spotify.Config{}
	err = envconfig.Process("SPOTIFY", spotifyConfig)
	if err != nil {
		return nil, err
	}

	httpConfig := &httpConfig{}
	err = envconfig.Process("HTTP", httpConfig)
	if err != nil {
		return nil, err
	}

	return &serverConfig{
		PostgresConfig: postgresConfig,
		SpotifyConfig:  spotifyConfig,
		HTTPConfig:     httpConfig,
	}, nil
}
