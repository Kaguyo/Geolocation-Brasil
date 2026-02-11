package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kaguyo/Geolocation-Brasil/internal/bootstrap"
)

const (
	DefaultMongoURI   = "mongodb://localhost:27017"
	DefaultDBName     = "geolocalizacao_br"
	DefaultCollection = "localizacoes"
	DefaultPort       = "8080"
)

func main() {
	importFlag := flag.Bool("import", false, "Importar dados de exemplo")
	importFileFlag := flag.String("file", "", "Arquivo CSV para importar (formato GeoNames)")
	serveFlag := flag.Bool("serve", false, "Iniciar servidor API")
	portFlag := flag.String("port", DefaultPort, "Porta do servidor")
	mongoURIFlag := flag.String("mongo-uri", DefaultMongoURI, "URI de conexão do MongoDB")

	flag.Parse()

	// 🔹 Bootstrap da aplicação
	app, err := bootstrap.Build(*mongoURIFlag, DefaultDBName, DefaultCollection)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		app.DB.Close(ctx)
	}()

	ctx := context.Background()

	if *importFlag {
		log.Println("🔄 Iniciando importação de dados...")

		if *importFileFlag != "" {
			log.Printf("📂 Importando arquivo: %s", *importFileFlag)
			if err := app.Service.ImportData(ctx, *importFileFlag); err != nil {
				log.Fatalf("❌ Erro ao importar arquivo: %v", err)
			}
		} else {
			log.Println("📂 Importando dados de exemplo (30 principais cidades)")
			if err := app.Service.ImportBrazilianCitiesExampleTest(ctx); err != nil {
				log.Fatalf("❌ Erro ao importar dados: %v", err)
			}
		}

		log.Println("🔧 Criando índices...")
		app.Service.CreateGeoIndex(ctx)
		app.Service.CreateTextIndex(ctx)

		log.Println("✅ Importação concluída com sucesso!")

		if !*serveFlag {
			return
		}
	}

	if *serveFlag {

		app.Service.CreateGeoIndex(ctx)
		app.Service.CreateTextIndex(ctx)

		srv := &http.Server{
			Addr:         ":" + *portFlag,
			Handler:      app.Router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			log.Printf("🚀 Servidor iniciado na porta %s", *portFlag)
			log.Printf("📍 Endpoints disponíveis:")
			log.Printf("   GET /health")
			log.Printf("   GET /location/{municipio}?estado=XX")
			log.Printf("   GET /nearby?lat=XX&lon=YY&distance=50")
			log.Println()

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
			}
		}()

		<-done
		log.Println("🛑 Encerrando servidor...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("❌ Erro ao encerrar servidor: %v", err)
		}

		log.Println("✅ Servidor encerrado com sucesso!")
		return
	}

	if !*importFlag && !*serveFlag {
		log.Println("🌎 API de Geolocalização - Brasil")
		log.Println()
		flag.PrintDefaults()
	}
}
