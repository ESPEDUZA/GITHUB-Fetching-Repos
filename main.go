package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ESPEDUZA/CC-GO/pkg"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	user := os.Getenv("GITHUB_USER")
	token := os.Getenv("GITHUB_TOKEN")

	repos, err := pkg.FetchRepositories(user, token)
	if err != nil {
		log.Fatal(err)
	}

	destDir := fmt.Sprintf("repos-%s", user)

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("repositories.csv")
	if err != nil {
		log.Fatal("Erreur lors de la création du fichier CSV:", err)
	}
	defer file.Close()

	// Créer un writer CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'en-tête du fichier CSV
	err = writer.Write([]string{"Name", "Clone URL", "Description", "Last Updated"})
	if err != nil {
		log.Fatal("Erreur lors de l'écriture de l'en-tête du fichier CSV:", err)
	}

	// Écrire chaque dépôt dans le fichier CSV
	for _, repo := range repos {
		err = writeRepoToCSV(file, writer, repo)
		if err != nil {
			log.Println("Erreur lors de l'écriture du dépôt dans le fichier CSV:", err)
		}
	}

	fmt.Println("Les informations des dépôts ont été écrites dans repositories.csv")

	const numWorkers = 10 // Nombre de goroutines s'exécutant simultanément

	// Créer un canal pour les travaux
	jobs := make(chan pkg.Repository, len(repos))

	// Ajouter les travaux au canal
	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs) // Fermer le canal pour indiquer qu'il n'y a plus de travaux à ajouter

	// Créer un canal pour les erreurs
	errors := make(chan error, len(repos))

	// Lancer les travailleurs
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				err := cloneRepo(repo, destDir, token)
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	// Attendre que tous les travailleurs aient terminé
	wg.Wait()
	close(errors) // Fermer le canal d'erreurs

	// Traiter les erreurs
	for err := range errors {
		log.Println("Error cloning the repository:", err)
	}

	zipFileName := "repos-archives-" + os.Getenv("GITHUB_USER") + ".zip"

	cmd := exec.Command("zip", "-r", zipFileName, destDir)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Erreur lors de la création de l'archive ZIP: %v", err)
	}

	fmt.Println("Tous les dépôts ont été clonés.")
	fmt.Println("Les dépôts ont été archivés avec succès dans", zipFileName)

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		username := os.Getenv("GITHUB_USER")
		filepath := "repos-archives-" + username + ".zip"
		http.ServeFile(w, r, filepath)
	})

	// Démarrer le serveur HTTP
	port := "8080"
	log.Printf("Serveur démarré sur :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func writeRepoToCSV(file *os.File, writer *csv.Writer, repo pkg.Repository) error {
	err := writer.Write([]string{repo.Name, repo.CloneURL, repo.Description, repo.UpdatedAt})
	if err != nil {
		return err
	}
	return nil
}

func cloneRepo(repo pkg.Repository, destDir string, token string) error {
	repoDir := filepath.Join(destDir, repo.Name)

	cloneURL := repo.CloneURL
	if token != "" {
		cloneURL = strings.Replace(cloneURL, "https://", fmt.Sprintf("https://%s:x-oauth-basic@", token), 1)
	}

	cmd := exec.Command("git", "clone", cloneURL, repoDir)
	err := cmd.Run()
	if err != nil {
		return err
	}

	err = pkg.ExecuteGitCommands(repoDir)
	if err != nil {
		return err
	}

	return nil
}
