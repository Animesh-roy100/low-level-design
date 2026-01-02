package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ==================== Interfaces ====================

type NotificationObserver interface {
	SendNotification(description string)
}

type MeetingObserver interface {
	NotifyUsers()
}

// ==================== User ====================
type User struct {
	UserID int
	Name   string
	Email  string
	Phone  string
}

func NewUser(id int, name, email, phone string) *User {
	return &User{
		UserID: id,
		Name:   name,
		Email:  email,
		Phone:  phone,
	}
}

func (u *User) SendNotification(description string) {
	fmt.Printf("Email sent to %s and the description: %s", u.Email, description)
}

// ==================== Meeting ====================
type Meeting struct {
	MeetingID     int
	MeetingRoomID int
	Title         string
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	Participants  []NotificationObserver
	Host          NotificationObserver
	mu            sync.RWMutex // For thread safe access
}

func NewMeeting(meetingId, meetingRoomId int, title, description string, startTime, endTime time.Time, host NotificationObserver) *Meeting {
	return &Meeting{
		MeetingID:     meetingId,
		MeetingRoomID: meetingRoomId,
		Title:         title,
		Description:   description,
		StartTime:     startTime,
		EndTime:       endTime,
		Participants:  make([]NotificationObserver, 0),
		Host:          host,
	}
}

func (m *Meeting) AddParticipants(participant NotificationObserver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Participants = append(m.Participants, participant)
}

func (m *Meeting) NotifyUsers() {
	// for _, attendee := range m.Participants {
	// 	attendee.SendNotification(m.Description)
	// }

	// m.Host.SendNotification(m.Description)

	// --------------------------------------------------

	// Use go routines for concurrent notification sending
	var wg sync.WaitGroup

	for _, attendee := range m.Participants {
		wg.Add(1)
		go func(a NotificationObserver) {
			defer wg.Done()
			a.SendNotification(fmt.Sprintf("Meeting: %s - %s", m.Title, m.Description))
		}(attendee)
	}

	// Notify host
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.Host.SendNotification(fmt.Sprintf("You're hosting: %s - %s", m.Title, m.Description))
	}()

	wg.Wait()
}

// ==================== TimeSlot ====================
type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewTimeSlot(startTime time.Time, duration time.Duration) *TimeSlot {
	return &TimeSlot{
		StartTime: startTime,
		EndTime:   startTime.Add(duration),
	}
}

func (ts *TimeSlot) Overlaps(other *TimeSlot) bool {
	return ts.StartTime.Before(other.EndTime) && other.StartTime.Before(ts.EndTime)
}

// ==================== MeetingRoom ====================
type MeetingRoom struct {
	MeetingRoomID int
	Calendar      *Calendar
	Capacity      int
	Name          string
	mu            sync.RWMutex
}

func NewMeetingRoom(meetingRoomId int, name string, capacity int) *MeetingRoom {
	return &MeetingRoom{
		MeetingRoomID: meetingRoomId,
		Capacity:      capacity,
		Name:          name,
		Calendar:      NewCalendar(meetingRoomId),
	}
}

func (mr *MeetingRoom) IsAvailable(slot *TimeSlot) bool {
	return mr.Calendar.IsSlotAvailable(slot)
}

func (mr *MeetingRoom) BookMeetingRoom(meeting *Meeting) error {
	return mr.Calendar.ScheduleMeeting(meeting)
}

func (mr *MeetingRoom) GetRoomID() int {
	return mr.MeetingRoomID
}

// ==================== Calendar ====================
type Calendar struct {
	ScheduledMeetings map[int]*Meeting
	MeetingRoomID     int
	mu                sync.Mutex
}

func NewCalendar(meetingRoomId int) *Calendar {
	return &Calendar{
		MeetingRoomID:     meetingRoomId,
		ScheduledMeetings: make(map[int]*Meeting),
	}
}

func (c *Calendar) IsSlotAvailable(slot *TimeSlot) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, meeting := range c.ScheduledMeetings {
		meetingSlot := &TimeSlot{StartTime: meeting.StartTime, EndTime: meeting.EndTime}
		if meetingSlot.Overlaps(slot) {
			return false
		}
	}
	return true
}

func (c *Calendar) ScheduleMeeting(meeting *Meeting) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	meetingSlot := &TimeSlot{StartTime: meeting.StartTime, EndTime: meeting.EndTime}

	// Check Conflicts
	for _, existingMeeting := range c.ScheduledMeetings {
		existingSlot := &TimeSlot{StartTime: existingMeeting.StartTime, EndTime: existingMeeting.EndTime}
		if existingSlot.Overlaps(meetingSlot) {
			return errors.New("time slot conflicts")
		}
	}

	c.ScheduledMeetings[meeting.MeetingID] = meeting
	return nil
}

// ==================== RoomBookingStrategy ====================
type RoomBookingStrategy interface {
	BookRoom(rooms []*MeetingRoom, slot *TimeSlot, participantsCount int) *MeetingRoom
}

type FCFSRoomBookingStrategy struct{}

func NewFCFSRoomBookingStrategy() *FCFSRoomBookingStrategy {
	return &FCFSRoomBookingStrategy{}
}

func (c *FCFSRoomBookingStrategy) BookRoom(rooms []*MeetingRoom, slot *TimeSlot, participantsCount int) *MeetingRoom {
	for _, room := range rooms {
		if room.Capacity >= participantsCount && room.IsAvailable(slot) {
			return room
		}
	}
	return nil
}

// ==================== MeetingScheduler ====================
type MeetingSchedular struct {
	MeetingRooms        []*MeetingRoom
	HistoryMeetings     []*Meeting
	RoomBookingStrategy RoomBookingStrategy
	mu                  sync.RWMutex
	meetingCounter      int
}

func NewMeetingSchedular(strategy RoomBookingStrategy) *MeetingSchedular {
	return &MeetingSchedular{
		MeetingRooms:        make([]*MeetingRoom, 0),
		HistoryMeetings:     make([]*Meeting, 0),
		RoomBookingStrategy: strategy,
		meetingCounter:      1,
	}
}

func (ms *MeetingSchedular) Notify(meeting MeetingObserver) {
	meeting.NotifyUsers()
}

func (ms *MeetingSchedular) AddMeetingRoom(room *MeetingRoom) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MeetingRooms = append(ms.MeetingRooms, room)
}

func (ms *MeetingSchedular) ScheduleMeeting(title, description string, startTime time.Time, duration time.Duration, host NotificationObserver, participants []NotificationObserver) (*Meeting, error) {
	slot := NewTimeSlot(startTime, duration)
	participantsCount := len(participants) + 1

	// find available room
	ms.mu.RLock()
	room := ms.RoomBookingStrategy.BookRoom(ms.MeetingRooms, slot, participantsCount)
	ms.mu.RUnlock()

	if room == nil {
		return nil, errors.New("no available room found for the given time slot")
	}

	// create meeting
	ms.mu.Lock()
	meeting := NewMeeting(ms.meetingCounter, room.MeetingRoomID, title, description, startTime, startTime.Add(duration), host)
	ms.meetingCounter++
	ms.mu.Unlock()

	// add participants
	for _, p := range participants {
		meeting.AddParticipants(p)
	}

	// book the room
	if err := room.BookMeetingRoom(meeting); err != nil {
		return nil, err
	}

	// add to history
	ms.mu.Lock()
	ms.HistoryMeetings = append(ms.HistoryMeetings, meeting)
	ms.mu.Unlock()

	return meeting, nil
}

// ---------------------------------------------------------------------------------------------

func main() {
	// new meeting schedular
	meetingSchedular := NewMeetingSchedular(NewFCFSRoomBookingStrategy())

	// Rooms
	room1 := NewMeetingRoom(1, "Alpha", 5)
	room2 := NewMeetingRoom(2, "Beta", 10)

	meetingSchedular.AddMeetingRoom(room1)
	meetingSchedular.AddMeetingRoom(room2)

	// Users
	host := NewUser(1, "Animesh", "animesh@gmail.com", "1111")
	u1 := NewUser(2, "Rahul", "rahul@gmail.com", "2222")
	u2 := NewUser(3, "Amit", "amit@gmail.com", "3333")

	start := time.Now().Add(1 * time.Hour)

	var wg sync.WaitGroup
	wg.Add(1)

	// Concurrent booking attempt
	go func() {
		defer wg.Done()
		meeting, err := meetingSchedular.ScheduleMeeting(
			"Design Discussion",
			"LLD interview prep",
			start,
			30*time.Minute,
			host,
			[]NotificationObserver{u1, u2},
		)
		if err != nil {
			fmt.Println("Booking 1 failed:", err)
			return
		}
		meetingSchedular.Notify(meeting)
	}()

	// go func() {
	// 	defer wg.Done()
	// 	_, err := meetingSchedular.ScheduleMeeting(
	// 		"Parallel Booking",
	// 		"Conflict test",
	// 		start,
	// 		30*time.Minute,
	// 		host,
	// 		[]NotificationObserver{u1},
	// 	)
	// 	if err != nil {
	// 		fmt.Println("Booking 2 failed:", err)
	// 	}
	// }()

	wg.Wait()
}
