// Package models defines core domain entities for the sale‑watches application.

// This file declares Review, representing a user’s feedback submitted via the API. It will be fully implemented in the future.
package models

// Review represents a user’s review of a product or service.
// It captures who wrote it, what they said, and the score they assigned.

// Fields:
//   - Content:  the textual body of the review.
//   - Rating:   the numeric score given by the user (e.g., 1–5).
// Comments should begin with the name of the thing being described and end in a period. :contentReference[oaicite:0]{index=0}
type Review struct {
	Content  string `json:"Content"` // the textual body of the review :contentReference[oaicite:2]{index=2}

	Rating   int    `json:"Rating"` // the numeric score assigned by the user :contentReference[oaicite:3]{index=3}
}
