package events

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"com.pikpo.events/app/database"
	"com.pikpo.events/app/utils"
)

type EventHandler struct{}

func (event *EventItem) validate() (err error) {
	if len(event.Name) < 3 {
		return errors.New("Field Name required min 3 characters")
	}
	if len(event.Date) < 0 {
		return errors.New("Field 'date' must be not empty")
	}
	if len(event.StartTime) < 0 {
		return errors.New("Field 'From date' must be not empty")
	}
	if len(event.EndTime) < 0 {
		return errors.New("Field 'To date' must be not empty")
	}
	if len(event.TimeZone) < 0 {
		return errors.New("Field time_zone is required")
	}
	// Execute the prepared query with value of each parameter for the target columns
	_, err = utils.ToTime(event.StartTime)
	if err != nil {
		return err
	}
	_, err = utils.ToTime(event.EndTime)
	if err != nil {
		return err
	}

	date, err := utils.ToDateTime(event.Date)
	if err != nil {
		return err
	}

	if int(date.Weekday()) == 0 || int(date.Weekday()) == 6 {
		return errors.New("Just allow weekdays event")
	}

	return nil
}

func (e *EventItem) validateConflict() (*sql.Rows, error) {
	log.Println("Validate event ", e)
	switch e.EventType {
	case EVENT_TYPE_DAYLY:
		return validateDaylyEvent(e)
	case EVENT_TYPE_WEEKLY:
		return validateWeeklyEvent(e)
	case EVENT_TYPE_MONTHLY:
		return validateMonthlyEvent(e)
	case EVENT_TYPE_ONE_DAY:
		return validateOneDayEvent(e)
	}
	return nil, errors.New("Unrecognized event type")
}

func validateOneDayEvent(event *EventItem) (*sql.Rows, error) {
	startTime, _ := utils.ToTime(event.StartTime)
	endTime, _ := utils.ToTime(event.EndTime)
	date, _ := utils.ToDateTime(event.Date)
	return database.Database.DB.Query(`
		SELECT 
			count(id)
		FROM events 
		WHERE
		(	
			(TO_TIMESTAMP(start_time, 'HH24:MM:SS') >= $1
		AND
			TO_TIMESTAMP(end_time, 'HH24:MM:SS') <= $2)
		OR
			allday = true
		)
		AND 
			date = $3
		`,
		startTime,
		endTime,
		date,
	)
}

func validateDaylyEvent(event *EventItem) (*sql.Rows, error) {
	startTime, _ := utils.ToTime(event.StartTime)
	endTime, _ := utils.ToTime(event.EndTime)

	return database.Database.DB.Query(`
		SELECT 
			count(id)
		FROM events 
		WHERE
		(	
			(TO_TIMESTAMP(start_time, 'HH24:MM:SS') >= $1
		AND
			TO_TIMESTAMP(end_time, 'HH24:MM:SS') <= $2)
		OR
			allday = true
		)
		AND 
			event_type = $3
	`,
		startTime,
		endTime,
		event.EventType,
	)
}

func validateWeeklyEvent(event *EventItem) (*sql.Rows, error) {
	startTime, _ := utils.ToTime(event.StartTime)
	endTime, _ := utils.ToTime(event.EndTime)
	date, _ := utils.ToDateTime(event.Date)

	return database.Database.DB.Query(`
		SELECT 
			count(id)
		FROM events 
		WHERE
		(	
			(TO_TIMESTAMP(start_time, 'HH24:MM:SS') >= $1
		AND
			TO_TIMESTAMP(end_time, 'HH24:MM:SS') <= $2)
		OR
			allday = true
		)
		AND 
			event_type = $3
		AND
			extract(dow from events.date) = $4
	`,
		startTime,
		endTime,
		event.EventType,
		int(date.Weekday()),
	)
}

func validateMonthlyEvent(event *EventItem) (*sql.Rows, error) {
	startTime, _ := utils.ToTime(event.StartTime)
	endTime, _ := utils.ToTime(event.EndTime)
	date, _ := utils.ToDateTime(event.Date)

	return database.Database.DB.Query(`
		SELECT 
			count(id)
		FROM events 
		WHERE
		(	
			(TO_TIMESTAMP(start_time, 'HH24:MM:SS') >= $1
		AND
			TO_TIMESTAMP(end_time, 'HH24:MM:SS') <= $2)
		OR
			allday = true
		)
		AND 
			event_type = $3
		AND
			extract(day from events.date) = $4
	`,
		startTime,
		endTime,
		event.EventType,
		int(date.Day()),
	)
}

func validateConflictEvent(event *EventItem) (err error) {
	var eventCount int32
	var rows *sql.Rows
	rows, err = event.validateConflict()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&eventCount,
		)
		if err != nil {
			return err
		}
	}
	if eventCount > 0 {
		return errors.New("conflict-event")
	}
	return nil
}

func createEvents(event *EventItem) (*EventItem, error) {
	// Prepare the insert query
	queryStr := fmt.Sprintf(`
	INSERT INTO events(
		name,
		date,
		start_time,
		end_time,
		event_type,
		allday,
		timezone
	) VALUES(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	) returning id`)

	log.Println(queryStr)
	query, err := database.Database.DB.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	date, err := utils.ToDateTime(event.Date)
	if err != nil {
		return nil, err
	}

	err = query.QueryRow(
		event.Name, date, event.StartTime, event.EndTime,
		event.EventType, event.Allday, event.TimeZone,
	).Scan(&event.Id)

	log.Println("Created event : ", event.Id)
	return event, nil
}

func getEvents(req *GetEventRequest) (*GetEventResponse, error) {
	events := []*EventItem{}
	rows, err := database.Database.DB.Query(`
		SELECT 
			id,
			name,
			date,
			start_time,
			end_time,
			event_type,
			allday,
			timezone,
			created_at,
			updated_at,
			deleted_at
		FROM events WHERE deleted_at is null
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event *EventItem
		event = &EventItem{}
		var deletedAt sql.NullString
		err = rows.Scan(
			&event.Id,
			&event.Name,
			&event.Date,
			&event.StartTime,
			&event.EndTime,
			&event.EventType,
			&event.Allday,
			&event.TimeZone,
			&event.CreatedAt,
			&event.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if len(deletedAt.String) > 0 {
			event.DeletedAt = deletedAt.String
		}
		events = append(events, event)
	}
	res := GetEventResponse{
		Meta: req,
		Data: events,
	}
	return &res, nil
}

func getEventsByID(id string) (*EventItem, error) {
	event := EventItem{}
	row := database.Database.DB.QueryRow(`
		SELECT 
			id,
			name,
			date,
			start_time,
			end_time,
			event_type,
			allday,
			timezone,
			created_at,
			updated_at,
			deleted_at
		FROM events WHERE id=$1 AND deleted_at is null
	`, id)

	var deletedAt sql.NullString
	err := row.Scan(
		&event.Id,
		&event.Name,
		&event.Date,
		&event.StartTime,
		&event.EndTime,
		&event.EventType,
		&event.Allday,
		&event.TimeZone,
		&event.CreatedAt,
		&event.UpdatedAt,
		&deletedAt,
	)
	if err != nil || len(event.Id) <= 0 {
		return nil, errors.New("Event not found")
	}
	if len(deletedAt.String) > 0 {
		event.DeletedAt = deletedAt.String
	}

	return &event, nil
}

func updateEvent(event *EventItem) (*EventItem, error) {
	if len(event.Id) <= 0 {
		return nil, errors.New("Unrecognized data")
	}
	_, err := getEventsByID(event.Id)
	if err != nil {
		return nil, errors.New("Event not found")
	}

	date, err := utils.ToDateTime(event.Date)
	if err != nil {
		return nil, err
	}

	_, err = database.Database.DB.Exec(
		`UPDATE events SET
			name = $1,
			date = $2,
			start_time = $3,
			end_time = $4,
			event_type = $5,
			allday = $6,
			timezone = $7,
			updated_at= $8
		WHERE id = $9
		`,
		event.Name,
		date,
		event.StartTime,
		event.EndTime,
		event.EventType.Number(),
		event.Allday,
		event.TimeZone,
		time.Now(),
		event.Id,
	)
	if err != nil {
		return nil, err
	}
	updatedEvent, _ := getEventsByID(event.Id)
	log.Println(updatedEvent)
	return updatedEvent, nil
}

func deleteEvent(id string) (err error) {
	_, err = database.Database.DB.Exec(
		`UPDATE events SET
			deleted_at = $1
		WHERE id = $2
		`,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *EventHandler) CreateEventHandler(event *EventItem) (*EventItem, error) {
	// Validate the input
	err := event.validate()
	if err != nil {
		return nil, err
	}
	// Validte data to database (if there is duplicate or conflict data)
	err = validateConflictEvent(event)
	if err != nil {
		return nil, err
	}
	// Do data transaction to database
	event, err = createEvents(event)
	if err != nil {
		return nil, err
	}

	return event, err
}

func (m *EventHandler) GetEventsHandler(req *GetEventRequest) (*GetEventResponse, error) {
	return getEvents(req)
}

func (m *EventHandler) UpdateEventHandler(event *EventItem) (*EventItem, error) {
	err := event.validate()
	if err != nil {
		return nil, err
	}

	err = validateConflictEvent(event)
	if err != nil {
		return nil, err
	}

	event, err = updateEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (m *EventHandler) DeleteEventHandler(req *EventByID) (*DeleteResponse, error) {
	err := deleteEvent(req.Id)
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{Message: "Success"}, nil
}
