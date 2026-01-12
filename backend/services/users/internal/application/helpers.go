package application

import (
	"context"
	"errors"
	"fmt"
	notifpb "social-network/shared/gen-go/notifications"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// HashPassword hashes a password using bcrypt.
// func hashPassword(password string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hash), err
// }

func checkPassword(storedPassword, newHashedPassword string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	// return err == nil
	tele.Debug(context.Background(), "Comparing passwords. @1 @2", "stored", storedPassword, "new", newHashedPassword)
	return storedPassword == newHashedPassword
}

func (s *Application) removeFailedImages(ctx context.Context, imagesToDelete []int64) {
	//==============pinpoint not found images============

	if len(imagesToDelete) > 0 {
		tele.Info(ctx, "removing avatar ids @1 from users", "failedImageIds", imagesToDelete)
		err := s.RemoveImages(context.WithoutCancel(ctx), imagesToDelete)
		if err != nil {
			tele.Warn(ctx, "failed  to delete failed images @1 from users: @2", "failedImageIds", imagesToDelete, "error", err.Error())
		}
	}

}

func (s *Application) removeFailedImagesAsync(ctx context.Context, imagesToDelete []int64) {
	go s.removeFailedImages(ctx, imagesToDelete)
}

func (s *Application) removeFailedImage(ctx context.Context, err error, imageId int64) {
	var commonError *ce.Error
	if errors.As(err, &commonError) {
		if err.(*ce.Error).IsClass(ce.ErrNotFound) {

			tele.Info(ctx, "removing avatar id @1 from users", "failedImageId", imageId)
			err := s.RemoveImages(context.WithoutCancel(ctx), []int64{imageId})
			if err != nil {
				tele.Warn(ctx, "failed  to delete failed image @1 from users: @2", "failedImageId", imageId, "error", err.Error())
			}
		}

	}
}

func (s *Application) removeFailedImageAsync(ctx context.Context, err error, imageId int64) {
	go s.removeFailedImage(ctx, err, imageId)
}

// Helper to create and send a notification event
func (s *Application) createAndSendNotificationEvent(ctx context.Context, event *notifpb.NotificationEvent) error {
	// Extract metadata
	requestId, ok := ctx.Value(ct.ReqID).(string)
	if !ok {
		tele.Error(ctx, "could not get request id")
		requestId = "unknown"
	}
	traceId, ok := ctx.Value(ct.TraceId).(string)
	if !ok {
		tele.Error(ctx, "could not get trace id")
		traceId = "unknown"
	}

	metadata := map[string]string{
		"source":     "users",
		"request_id": requestId,
		"trace_id":   traceId,
	}

	// Populate common fields
	event.EventId = uuid.NewString()
	event.OccurredAt = timestamppb.Now()
	event.Metadata = metadata

	// Serialize
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal notification event: %w", err)
	}

	// Send to Kafka
	err = s.eventProducer.Send(ctx, ct.NotificationTopic, eventBytes)
	if err != nil {
		return fmt.Errorf("failed to send notification event: %w", err)
	}

	tele.Info(ctx, "Notification event sent: @1", "eventType", event.EventType)
	return nil
}
