package events

import (
	"context"
	"log"
	"testing"
	"time"

	"com.pikpo.events/app/database"
	"com.pikpo.events/app/utils"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCRUD_MonthlyEvent(t *testing.T) {

	utils.BaseTest(t)

	// Start the rpc client
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	eventClient := NewEventClient(conn)

	// // Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test create event
	event := &EventItem{
		Name:      "Band Show 1",
		Date:      "2022-04-05",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_MONTHLY,
		TimeZone:  "UTC",
	}
	result, err := eventClient.CreateEvent(ctx, event)
	assert.Equal(t, nil, err, "Should be not error")
	assert.NotEqual(t, 0, len(result.Id), "Should be not error")
	log.Println(result.Id)

	// Conflict event
	event = &EventItem{
		Name:      "Band Show 2",
		Date:      "2022-05-05",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_MONTHLY,
		TimeZone:  "UTC",
	}

	_, err = eventClient.CreateEvent(ctx, event)
	assert.NotNil(t, err, "Should be error")

	// Test update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 2",
		Date:      "2022-05-05",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_MONTHLY,
		TimeZone:  "UTC",
	}
	result, err = eventClient.UpdateEvent(ctx, event)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Band Show 2", result.Name)

	// Test conflict update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 3",
		Date:      "2022-06-05",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_MONTHLY,
		TimeZone:  "UTC",
	}
	_, err = eventClient.UpdateEvent(ctx, event)
	assert.NotEqual(t, nil, err)

	// Test Get events
	res, err := eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 1, len(res.Data))

	// Test Delete event
	_, err = eventClient.DeleteEvent(ctx, &EventByID{Id: result.Id})
	assert.Equal(t, nil, err)

	res, err = eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 0, len(res.Data))

	// Stop UT
	utils.BaseTest(t)
	database.Database.DB.Close()
}

func TestCRUD_WeeklyEvent(t *testing.T) {

	utils.BaseTest(t)

	// Start the rpc client
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	eventClient := NewEventClient(conn)

	// // Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test create event
	event := &EventItem{
		Name:      "Band Show 1",
		Date:      "2022-04-05",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_WEEKLY,
		TimeZone:  "UTC",
	}
	result, err := eventClient.CreateEvent(ctx, event)
	assert.Equal(t, nil, err, "Should be not error")
	assert.NotEqual(t, 0, len(result.Id), "Should be not error")
	log.Println(result.Id)

	// Conflict event
	event = &EventItem{
		Name:      "Band Show 2",
		Date:      "2022-04-12",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_WEEKLY,
		TimeZone:  "UTC",
	}

	_, err = eventClient.CreateEvent(ctx, event)
	assert.NotNil(t, err, "Should be error")

	// Test update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 2",
		Date:      "2022-04-12",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_WEEKLY,
		TimeZone:  "UTC",
	}
	result, err = eventClient.UpdateEvent(ctx, event)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Band Show 2", result.Name)

	// Test conflict update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 3",
		Date:      "2022-04-12",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_WEEKLY,
		TimeZone:  "UTC",
	}
	_, err = eventClient.UpdateEvent(ctx, event)
	assert.NotEqual(t, nil, err)

	// Test Get events
	res, err := eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 1, len(res.Data))

	// Test Delete event
	_, err = eventClient.DeleteEvent(ctx, &EventByID{Id: result.Id})
	assert.Equal(t, nil, err)

	res, err = eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 0, len(res.Data))

	// Stop UT
	utils.BaseTest(t)
	database.Database.DB.Close()
}

func TestCRUD_DaylyEvent(t *testing.T) {

	utils.BaseTest(t)

	// Start the rpc client
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	eventClient := NewEventClient(conn)

	// // Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test create event
	event := &EventItem{
		Name:      "Band Show 1",
		Date:      "2022-03-09",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_DAYLY,
		TimeZone:  "UTC",
	}
	result, err := eventClient.CreateEvent(ctx, event)
	assert.Equal(t, nil, err, "Should be not error")
	assert.NotEqual(t, 0, len(result.Id), "Should be not error")
	log.Println(result.Id)

	// Conflict event
	event = &EventItem{
		Name:      "Band Show 2",
		Date:      "2022-03-16",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_DAYLY,
		TimeZone:  "UTC",
	}

	_, err = eventClient.CreateEvent(ctx, event)
	assert.NotNil(t, err, "Should be error")

	// Test update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 2",
		Date:      "2022-03-17",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_DAYLY,
		TimeZone:  "UTC",
	}
	result, err = eventClient.UpdateEvent(ctx, event)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Band Show 2", result.Name)

	// Test conflict update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 3",
		Date:      "2022-03-24",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_DAYLY,
		TimeZone:  "UTC",
	}
	_, err = eventClient.UpdateEvent(ctx, event)
	assert.NotEqual(t, nil, err)

	// Test Get events
	res, err := eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 1, len(res.Data))

	// Test Delete event
	_, err = eventClient.DeleteEvent(ctx, &EventByID{Id: result.Id})
	assert.Equal(t, nil, err)

	res, err = eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 0, len(res.Data))

	// Stop UT
	utils.BaseTest(t)
	database.Database.DB.Close()
}

func TestCRUD_OneDayEvent(t *testing.T) {

	utils.BaseTest(t)

	// Start the rpc client
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	eventClient := NewEventClient(conn)

	// // Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test create event
	event := &EventItem{
		Name:      "Band Show 1",
		Date:      "2022-04-05",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_ONE_DAY,
		TimeZone:  "UTC",
	}
	result, err := eventClient.CreateEvent(ctx, event)
	assert.Equal(t, nil, err, "Should be not error")
	assert.NotEqual(t, 0, len(result.Id), "Should be not error")
	log.Println(result.Id)

	// Conflict event
	event = &EventItem{
		Name:      "Band Show 2",
		Date:      "2022-04-05",
		StartTime: "15:00:00",
		EndTime:   "18:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_ONE_DAY,
		TimeZone:  "UTC",
	}

	_, err = eventClient.CreateEvent(ctx, event)
	assert.NotNil(t, err, "Should be error")

	// Test update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 2",
		Date:      "2022-04-05",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_ONE_DAY,
		TimeZone:  "UTC",
	}
	result, err = eventClient.UpdateEvent(ctx, event)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Band Show 2", result.Name)

	// Test conflict update event
	event = &EventItem{
		Id:        result.Id,
		Name:      "Band Show 3",
		Date:      "2022-04-05",
		StartTime: "19:00:00",
		EndTime:   "20:00:00",
		Allday:    false,
		EventType: EVENT_TYPE_ONE_DAY,
		TimeZone:  "UTC",
	}
	_, err = eventClient.UpdateEvent(ctx, event)
	assert.NotEqual(t, nil, err)

	// Test Get events
	res, err := eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 1, len(res.Data))

	// Test Delete event
	_, err = eventClient.DeleteEvent(ctx, &EventByID{Id: result.Id})
	assert.Equal(t, nil, err)

	res, err = eventClient.GetEvents(ctx, &GetEventRequest{})
	assert.Equal(t, nil, err, "Should be not error")
	assert.Equal(t, 0, len(res.Data))

	// Stop UT
	utils.BaseTest(t)
	database.Database.DB.Close()
}
