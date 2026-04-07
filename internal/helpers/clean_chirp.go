package helper

import (
	"strings"
)

func Clean(unclean string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	uncleanWords := strings.Split(unclean, " ")

	for _, badWord := range badWords {
		for uncleanInd, uncleanWord := range uncleanWords {
			if strings.ToLower(uncleanWord) == badWord {
				uncleanWords[uncleanInd] = "****"
			}
		}
	}
	return strings.Join(uncleanWords, " ")
}
