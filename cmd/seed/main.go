package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/beavercli/beaver_api/common/config"
	"github.com/beavercli/beaver_api/common/database"
	"github.com/beavercli/beaver_api/internal/storage"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	numContributors           = 100
	numTags                   = 50
	numUsers                  = 10
	numSnippets               = 500
	maxTagsPerSnippet         = 5
	maxContributorsPerSnippet = 5
)

var languages = []string{"C", "Go", "Python"}

var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Christopher", "Lisa", "Daniel", "Nancy",
	"Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra", "Donald", "Ashley",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
}

var tagPrefixes = []string{
	"algorithm", "data-structure", "network", "security", "database", "api",
	"testing", "performance", "memory", "concurrency", "async", "sync",
	"cache", "queue", "stack", "tree", "graph", "sort", "search", "hash",
	"encryption", "compression", "parsing", "serialization", "validation",
}

var tagSuffixes = []string{
	"basic", "advanced", "intro", "deep-dive", "practical", "example",
	"pattern", "tutorial", "guide", "tip", "trick", "hack", "best-practice",
}

var codeSnippets = map[string][]string{
	"C": {
		`#include <stdio.h>

int main() {
    printf("Hello, World!\\n");
    return 0;
}`,
		`#include <stdlib.h>

int* create_array(int size) {
    return (int*)malloc(size * sizeof(int));
}`,
		`#include <string.h>

void reverse_string(char* str) {
    int len = strlen(str);
    for (int i = 0; i < len / 2; i++) {
        char temp = str[i];
        str[i] = str[len - 1 - i];
        str[len - 1 - i] = temp;
    }
}`,
	},
	"Go": {
		`package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`,
		`func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}`,
		`func quickSort(arr []int) []int {
    if len(arr) < 2 {
        return arr
    }
    pivot := arr[0]
    var left, right []int
    for _, v := range arr[1:] {
        if v <= pivot {
            left = append(left, v)
        } else {
            right = append(right, v)
        }
    }
    return append(append(quickSort(left), pivot), quickSort(right)...)
}`,
	},
	"Python": {
		`def hello_world():
    print("Hello, World!")

if __name__ == "__main__":
    hello_world()`,
		`def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)`,
		`def quicksort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[0]
    left = [x for x in arr[1:] if x <= pivot]
    right = [x for x in arr[1:] if x > pivot]
    return quicksort(left) + [pivot] + quicksort(right)`,
	},
}

func main() {
	ctx := context.Background()
	cfg := config.New()

	pool, err := database.New(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := storage.New(pool)

	fmt.Println("Seeding database...")

	seedLanguages(ctx, queries)
	seedTags(ctx, queries)
	seedContributors(ctx, queries)
	seedUsers(ctx, queries)
	seedSnippets(ctx, queries)

	fmt.Println("Done!")
}

func seedLanguages(ctx context.Context, q *storage.Queries) {
	fmt.Printf("Inserting %d languages...\n", len(languages))
	for _, lang := range languages {
		if err := q.UpsertLanguage(ctx, pgtype.Text{String: lang, Valid: true}); err != nil {
			panic(err)
		}
	}
}

func seedTags(ctx context.Context, q *storage.Queries) {
	fmt.Printf("Inserting %d tags...\n", numTags)
	for i := 0; i < numTags; i++ {
		prefix := tagPrefixes[rand.Intn(len(tagPrefixes))]
		suffix := tagSuffixes[rand.Intn(len(tagSuffixes))]
		tag := fmt.Sprintf("%s-%s", prefix, suffix)

		if err := q.UpsertTag(ctx, pgtype.Text{String: tag, Valid: true}); err != nil {
			panic(err)
		}
	}
}

func seedContributors(ctx context.Context, q *storage.Queries) {
	fmt.Printf("Inserting %d contributors...\n", numContributors)
	for i := 0; i < numContributors; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s.%d@example.com", firstName, lastName, i)

		if err := q.UpsertContributor(ctx, storage.UpsertContributorParams{
			FirstName: pgtype.Text{String: firstName, Valid: true},
			LastName:  pgtype.Text{String: lastName, Valid: true},
			Email:     pgtype.Text{String: email, Valid: true},
		}); err != nil {
			panic(err)
		}
	}
}

func seedUsers(ctx context.Context, q *storage.Queries) {
	fmt.Printf("Inserting %d users...\n", numUsers)
	for i := 0; i < numUsers; i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		passwordHash := fmt.Sprintf("hash_%d", i)

		if err := q.UpsertUser(ctx, storage.UpsertUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}); err != nil {
			panic(err)
		}
	}
}

func seedSnippets(ctx context.Context, q *storage.Queries) {
	fmt.Printf("Inserting %d snippets...\n", numSnippets)

	// Get all languages
	allLanguages, err := q.ListAllLanguages(ctx)
	if err != nil {
		panic(err)
	}
	langMap := make(map[string]int64)
	for _, l := range allLanguages {
		if l.Name.Valid {
			langMap[l.Name.String] = l.ID
		}
	}

	// Get all tags
	allTags, err := q.ListAllTags(ctx)
	if err != nil {
		panic(err)
	}
	var tagIDs []int64
	for _, t := range allTags {
		tagIDs = append(tagIDs, t.ID)
	}

	// Get all contributors
	allContributors, err := q.ListAllContributors(ctx)
	if err != nil {
		panic(err)
	}
	var contributorIDs []int64
	for _, c := range allContributors {
		contributorIDs = append(contributorIDs, c.ID)
	}

	// Get all users
	allUsers, err := q.ListAllUsers(ctx)
	if err != nil {
		panic(err)
	}
	var userIDs []int64
	for _, u := range allUsers {
		userIDs = append(userIDs, u.ID)
	}

	for i := 0; i < numSnippets; i++ {
		lang := languages[rand.Intn(len(languages))]
		langID := langMap[lang]
		snippets := codeSnippets[lang]
		code := snippets[rand.Intn(len(snippets))]

		title := fmt.Sprintf("%s-snippet-%d-%d", lang, i, time.Now().UnixNano())
		projectURL := fmt.Sprintf("https://github.com/example/project-%d", rand.Intn(100))

		var userID pgtype.Int8
		if rand.Float32() < 0.7 && len(userIDs) > 0 {
			userID = pgtype.Int8{Int64: userIDs[rand.Intn(len(userIDs))], Valid: true}
		}

		snippetID, err := q.UpsertSnippet(ctx, storage.UpsertSnippetParams{
			Title:      pgtype.Text{String: title, Valid: true},
			Code:       pgtype.Text{String: code, Valid: true},
			ProjectUrl: pgtype.Text{String: projectURL, Valid: true},
			LanguageID: pgtype.Int8{Int64: langID, Valid: true},
			UserID:     userID,
			CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			panic(err)
		}

		// Add random tags
		numTagsToAdd := rand.Intn(maxTagsPerSnippet) + 1
		usedTags := make(map[int64]bool)
		for j := 0; j < numTagsToAdd && j < len(tagIDs); j++ {
			tagID := tagIDs[rand.Intn(len(tagIDs))]
			if usedTags[tagID] {
				continue
			}
			usedTags[tagID] = true
			if err := q.LinkSnippetTag(ctx, storage.LinkSnippetTagParams{
				SnippetID: snippetID,
				TagID:     tagID,
			}); err != nil {
				panic(err)
			}
		}

		// Add random contributors
		numContribsToAdd := rand.Intn(maxContributorsPerSnippet) + 1
		usedContribs := make(map[int64]bool)
		for j := 0; j < numContribsToAdd && j < len(contributorIDs); j++ {
			contribID := contributorIDs[rand.Intn(len(contributorIDs))]
			if usedContribs[contribID] {
				continue
			}
			usedContribs[contribID] = true
			if err := q.LinkSnippetContributor(ctx, storage.LinkSnippetContributorParams{
				SnippetID:     snippetID,
				ContributorID: contribID,
			}); err != nil {
				panic(err)
			}
		}

		if (i+1)%100 == 0 {
			fmt.Printf("  Inserted %d snippets...\n", i+1)
		}
	}
}
