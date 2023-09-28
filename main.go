package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ESPEDUZA/CC-GO/pkg"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	file, err := os.Create(filepath.Join(destDir, "repositories.csv"))
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Name", "Clone URL", "Description", "Last Updated"})
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo pkg.Repository) {
			defer wg.Done()
			err := writeRepoToCSV(file, writer, repo)
			if err != nil {
				log.Println("Erreur lors de l'écriture du dépôt dans le CSV:", err)
			}
			err = cloneRepo(repo, destDir)
			if err != nil {
				log.Println("Erreur lors du clonage du dépôt:", err)
			}
		}(repo)
	}
	wg.Wait()

	zipFileName := "repos-archive.zip"

	cmd := exec.Command("zip", "-r", zipFileName, destDir)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Erreur lors de la création de l'archive ZIP: %v", err)
	}
	fmt.Println("Les dépôts ont été archivés avec succès dans", zipFileName)

	fmt.Println("Tous les dépôts ont été clonés, mis à jour et archivés.")

}

func writeRepoToCSV(file *os.File, writer *csv.Writer, repo pkg.Repository) error {
	err := writer.Write([]string{repo.Name, repo.CloneURL, repo.Description, repo.UpdatedAt})
	if err != nil {
		return err
	}
	return nil
}

func cloneRepo(repo pkg.Repository, destDir string) error {
	repoDir := filepath.Join(destDir, repo.Name)
	cmd := exec.Command("git", "clone", repo.CloneURL, repoDir)
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
