package application

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
	"time"

	ct "social-network/shared/go/customtypes"
)

var processingVariants atomic.Bool

// StartVariantWorker starts a background worker that periodically processes pending file variants
func (m *MediaService) StartVariantWorker(ctx context.Context, interval time.Duration) {
	log.Printf("Initiating variant worker. Interval %s\n", interval.String())
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !processingVariants.CompareAndSwap(false, true) {
					continue
				}
				if err := m.processPendingVariants(ctx); err != nil {
					log.Printf("Error processing pending variants: %v", err)
				}
			case <-ctx.Done():
				log.Println("Variant worker stopped")
				return
			}
		}
	}()
}

// processPendingVariants queries for file_variants with status 'pending' and calls GenerateVariant for each
func (m *MediaService) processPendingVariants(ctx context.Context) error {
	defer processingVariants.Store(false)

	variants, err := m.Queries.GetPendingVariants(ctx)
	if err != nil {
		return err
	}
	log.Printf("Running variant worker for num of variants: %v", len(variants))

	for _, v := range variants {
		// Compute source bucket and object key from the original file
		sourceBucket := m.Cfgs.FileService.Buckets.Originals
		sourceObjectKey := strings.TrimSuffix(v.ObjectKey, "/"+v.Variant.String())

		// Call GenerateVariant
		size, err := m.Clients.GenerateVariant(ctx, sourceBucket, sourceObjectKey, v.Variant)
		if err != nil {
			log.Printf("Failed to generate variant for file %d variant %s: %v", v.Id, v.Variant, err)
			// Update status to failed
			if updateErr := m.Queries.UpdateVariantStatusAndSize(ctx, v.Id, ct.Failed, size); updateErr != nil {
				log.Printf("Failed to update status to failed: %v", updateErr)
			}
		} else {
			log.Printf("Successfully generated variant for file %d variant %s", v.Id, v.Variant)
			// Update status to complete
			if updateErr := m.Queries.UpdateVariantStatusAndSize(ctx, v.Id, ct.Complete, size); updateErr != nil {
				log.Printf("Failed to update status to complete: %v", updateErr)
			}
		}
	}

	return nil
}
