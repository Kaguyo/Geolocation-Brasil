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
	"github.com/Kaguyo/Geolocation-Brasil/internal/utils"
)

const (
	DefaultMongoURI   = "mongodb://localhost:27017"
	DefaultDBName     = "geolocalizacao_br"
	DefaultCollection = "localizacoes"
	DefaultPort       = "8080"
	GeoNamesURL       = "http://download.geonames.org/export/dump/BR.zip"
)

func main() {
	importFlag := flag.Bool("import", false, "Importar dados de exemplo (30 principais cidades)")
	importFileFlag := flag.String("file", "", "Arquivo CSV para importar (formato GeoNames)")
	importAllFlag := flag.Bool("importall", false, "Baixar BR.zip do GeoNames, descompactar e importar todos os dados (~5570 municípios)")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if *importAllFlag {
		log.Println("🧹 Limpando coleção antes da importação completa...")
		if err := app.Service.ResetCollection(ctx, DefaultCollection); err != nil {
			log.Fatalf("Erro ao limpar coleção: %v", err)
		}

		log.Println("🔄 Iniciando importação completa do GeoNames...")

		zipFile := "BR.zip"
		extractedFile := "BR.txt"

		// Download
		log.Printf("📥 Baixando %s do GeoNames...", zipFile)
		if err := utils.DownloadFile(GeoNamesURL, zipFile); err != nil {
			log.Fatalf("❌ Erro ao baixar arquivo: %v", err)
		}
		log.Println("✅ Download concluído!")

		// Extract
		log.Println("📦 Descompactando arquivo...")
		if err := utils.UnzipFile(zipFile); err != nil {
			log.Fatalf("❌ Erro ao descompactar: %v", err)
		}
		log.Println("✅ Descompactação concluída!")

		// Import
		log.Printf("📂 Importando dados de %s (~5570 municípios) aguarde...", extractedFile)
		if err := app.Service.ImportData(ctx, extractedFile); err != nil {
			log.Fatalf("❌ Erro ao importar dados: %v", err)
		}

		// Create indices
		log.Println("🔧 Criando índices...")
		app.Service.CreateGeoIndex(ctx)
		app.Service.CreateTextIndex(ctx)

		log.Println("✅ Importação completa concluída com sucesso!")

		if !*serveFlag {
			return
		}
	}

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

		// Only create indices if they don't exist (server should run with pre-imported data)
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

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("❌ Erro ao encerrar servidor: %v", err)
		}

		log.Println("✅ Servidor encerrado com sucesso!")
		return
	}

	if !*importFlag && !*importAllFlag && !*serveFlag {
		log.Println("🌎 API de Geolocalização - Brasil")
		log.Println("🟡 Inicializado em modo de teste. Use as flags para importar dados ou iniciar o servidor.")
		flag.PrintDefaults()
	}
}
