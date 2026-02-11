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
)

const (
	DefaultMongoURI   = "mongodb://localhost:27017"
	DefaultDBName     = "geolocalizacao_br"
	DefaultCollection = "localizacoes"
	DefaultPort       = "8080"
)

func main() {
	// Flags de linha de comando
	importFlag := flag.Bool("import", false, "Importar dados de exemplo")
	importFileFlag := flag.String("file", "", "Arquivo CSV para importar (formato GeoNames)")
	serveFlag := flag.Bool("serve", false, "Iniciar servidor API")
	portFlag := flag.String("port", DefaultPort, "Porta do servidor")
	mongoURIFlag := flag.String("mongo-uri", DefaultMongoURI, "URI de conexão do MongoDB")

	flag.Parse()

	// Conectar ao MongoDB
	db, err := ConnectDB(*mongoURIFlag, DefaultDBName, DefaultCollection)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db.Close(ctx)
	}()

	ctx := context.Background()

	// Modo de importação
	if *importFlag {
		log.Println("🔄 Iniciando importação de dados...")

		if *importFileFlag != "" {
			// Importar de arquivo GeoNames
			log.Printf("📂 Importando arquivo: %s", *importFileFlag)
			if err := db.ImportData(ctx, *importFileFlag); err != nil {
				log.Fatalf("❌ Erro ao importar arquivo: %v", err)
			}
		} else {
			// Importar dados de exemplo (capitais)
			log.Println("📂 Importando dados de exemplo (30 principais cidades)")
			if err := db.ImportBrazilianCities(ctx); err != nil {
				log.Fatalf("❌ Erro ao importar dados: %v", err)
			}
		}

		// Criar índices
		log.Println("🔧 Criando índices...")
		if err := db.CreateGeoIndex(ctx); err != nil {
			log.Printf("⚠️  Aviso ao criar índice geo: %v", err)
		}
		if err := db.CreateTextIndex(ctx); err != nil {
			log.Printf("⚠️  Aviso ao criar índice texto: %v", err)
		}

		log.Println("✅ Importação concluída com sucesso!")

		if !*serveFlag {
			return
		}
	}

	// Modo servidor
	if *serveFlag {
		// Criar índices se não existirem
		db.CreateGeoIndex(ctx)
		db.CreateTextIndex(ctx)

		// Configurar API
		api := NewAPI(db)
		router := api.SetupRoutes()

		// Configurar servidor
		srv := &http.Server{
			Addr:         ":" + *portFlag,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		// Canal para capturar sinais de shutdown
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Iniciar servidor em goroutine
		go func() {
			log.Printf("🚀 Servidor iniciado na porta %s", *portFlag)
			log.Printf("📍 Endpoints disponíveis:")
			log.Printf("   GET /health - Verificar status da API")
			log.Printf("   GET /location/{municipio}?estado=XX - Buscar por município")
			log.Printf("   GET /nearby?lat=XX&lon=YY&distance=50 - Buscar próximos")
			log.Println()
			log.Printf("💡 Exemplos:")
			log.Printf("   curl http://localhost:%s/health", *portFlag)
			log.Printf("   curl http://localhost:%s/location/São%%20Paulo", *portFlag)
			log.Printf("   curl http://localhost:%s/location/Campinas?estado=SP", *portFlag)
			log.Printf("   curl \"http://localhost:%s/nearby?lat=-23.5505&lon=-46.6333&distance=100\"", *portFlag)
			log.Println()

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
			}
		}()

		// Aguardar sinal de shutdown
		<-done
		log.Println("🛑 Encerrando servidor...")

		// Graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("❌ Erro ao encerrar servidor: %v", err)
		}

		log.Println("✅ Servidor encerrado com sucesso!")
		return
	}

	// Se nenhum flag foi especificado, mostrar ajuda
	if !*importFlag && !*serveFlag {
		log.Println("🌎 API de Geolocalização - Brasil")
		log.Println()
		log.Println("Uso:")
		log.Println("  Importar dados de exemplo e iniciar servidor:")
		log.Println("    go run . -import -serve")
		log.Println()
		log.Println("  Importar arquivo GeoNames:")
		log.Println("    go run . -import -file=BR.txt")
		log.Println()
		log.Println("  Apenas iniciar servidor:")
		log.Println("    go run . -serve")
		log.Println()
		log.Println("  Especificar porta:")
		log.Println("    go run . -serve -port=3000")
		log.Println()
		log.Println("Flags disponíveis:")
		flag.PrintDefaults()
	}
}
