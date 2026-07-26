package main

import (
    "fmt"
    "github.com/locne/analysis-service/internal/usecase"
    "github.com/locne/analysis-service/internal/interface/handler"
    "github.com/gin-gonic/gin"
    "os"
    "encoding/json"
    "github.com/gin-contrib/cors"
    
)

func loadOpeningsFromCache(filename string) ([]usecase.Opening, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var openings []usecase.Opening
    err = json.Unmarshal(data, &openings)
    return openings, err
}


func main() {
    cacheFiles := []string{
        "cmd/combined_openings.json",
        "cmd/openings.json",
        "combined_openings.json",
        "openings.json",
    }
    var openings []usecase.Opening
    var err error

    for _, file := range cacheFiles {
        openings, err = loadOpeningsFromCache(file)
        if err == nil {
            fmt.Printf("Loaded %d openings from %s\n", len(openings), file)
            break
        }
    }

    if err != nil || len(openings) == 0 {
        fmt.Printf("Failed to load opening cache: %v\n", err)
    }
    
    
    trie := usecase.NewOpeningTrie()
    trie.BuildFromOpenings(openings)

    router := gin.Default()

        router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    
    handler.RegisterAnalysisRoutes(router, trie)
    router.Run(":8080")
}