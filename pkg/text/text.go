package text

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func ClassifyText(ctx context.Context, client *language.Client, text string) (*languagepb.ClassifyTextResponse, error) {
	return client.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}
