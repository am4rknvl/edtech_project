package graph

import (
	"fmt"
	"sync"

	"github.com/am4rknvl/edtech_project/graph/model"
)

// very small concurrent in-memory store for development/testing only
var store = newMemoryStore()

type memoryStore struct {
	mu          sync.RWMutex
	subjects    map[string]*model.Subject
	courses     map[string]*model.Course
	units       map[string]*model.Unit
	lessons     map[string]*model.Lesson
	submissions map[string]*model.Submission
	nextID      int
}

func newMemoryStore() *memoryStore {
	s := &memoryStore{
		subjects:    map[string]*model.Subject{},
		courses:     map[string]*model.Course{},
		units:       map[string]*model.Unit{},
		lessons:     map[string]*model.Lesson{},
		submissions: map[string]*model.Submission{},
		nextID:      1,
	}
	// seed a subject
	subj := &model.Subject{ID: "1", Name: "Math"}
	s.subjects[subj.ID] = subj
	return s
}

func (s *memoryStore) next() string {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.mu.Unlock()
	return fmt.Sprintf("%d", id)
}

func (s *memoryStore) CreateUser(input model.SignUpInput) *model.User {
	id := s.next()
	u := &model.User{ID: id, Role: input.Role}
	return u
}

func (s *memoryStore) CreateStudentProfile(input model.CreateStudentProfileInput) *model.StudentProfile {
	id := s.next()
	p := &model.StudentProfile{ID: id, Grade: input.Grade, Age: input.Age}
	return p
}

func (s *memoryStore) ListSubjects() []*model.Subject {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Subject, 0, len(s.subjects))
	for _, v := range s.subjects {
		out = append(out, v)
	}
	return out
}

func (s *memoryStore) CreateCourse(input model.CreateCourseInput) (*model.Course, error) {
	id := s.next()
	subj := s.subjects[input.SubjectID]
	if subj == nil {
		return nil, fmt.Errorf("subject not found")
	}
	c := &model.Course{ID: id, Subject: subj, Grade: input.Grade, Title: input.Title}
	s.courses[id] = c
	return c, nil
}

func (s *memoryStore) ListCourses(grade *int, subjectID *string) []*model.Course {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*model.Course{}
	for _, c := range s.courses {
		if grade != nil && c.Grade != *grade {
			continue
		}
		if subjectID != nil && c.Subject != nil && c.Subject.ID != *subjectID {
			continue
		}
		out = append(out, c)
	}
	return out
}

func (s *memoryStore) GetCourse(id string) (*model.Course, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c := s.courses[id]
	if c == nil {
		return nil, nil
	}
	return c, nil
}

func (s *memoryStore) CreateUnit(input model.CreateUnitInput) (*model.Unit, error) {
	id := s.next()
	c := s.courses[input.CourseID]
	if c == nil {
		return nil, fmt.Errorf("course not found")
	}
	u := &model.Unit{ID: id, Course: c, Title: input.Title, Order: 0}
	s.units[id] = u
	return u, nil
}

func (s *memoryStore) ListUnits(courseID string) []*model.Unit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*model.Unit{}
	for _, u := range s.units {
		if u.Course != nil && u.Course.ID == courseID {
			out = append(out, u)
		}
	}
	return out
}

func (s *memoryStore) CreateLesson(input model.CreateLessonInput) (*model.Lesson, error) {
	id := s.next()
	u := s.units[input.UnitID]
	if u == nil {
		return nil, fmt.Errorf("unit not found")
	}
	l := &model.Lesson{ID: id, Unit: u, Title: input.Title, Order: 0, Version: 1}
	s.lessons[id] = l
	return l, nil
}

func (s *memoryStore) ListLessons(unitID string) []*model.Lesson {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*model.Lesson{}
	for _, l := range s.lessons {
		if l.Unit != nil && l.Unit.ID == unitID {
			out = append(out, l)
		}
	}
	return out
}

func (s *memoryStore) GetLesson(id string) (*model.Lesson, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lessons[id], nil
}

func (s *memoryStore) AddTextBlock(lessonID string, input model.TextBlockInput) (*model.TextBlock, error) {
	id := s.next()
	tb := &model.TextBlock{ID: id, Text: input.Text, Order: 0}
	return tb, nil
}

func (s *memoryStore) AddImageBlock(lessonID string, input model.ImageBlockInput) (*model.ImageBlock, error) {
	id := s.next()
	ib := &model.ImageBlock{ID: id, URL: input.URL, Alt: input.Alt, Order: 0}
	return ib, nil
}

func (s *memoryStore) AddAudioBlock(lessonID string, input model.AudioBlockInput) (*model.AudioBlock, error) {
	id := s.next()
	ab := &model.AudioBlock{ID: id, URL: input.URL, DurationSec: input.DurationSec, Order: 0}
	return ab, nil
}

func (s *memoryStore) AddVideoBlock(lessonID string, input model.VideoBlockInput) (*model.VideoBlock, error) {
	id := s.next()
	vb := &model.VideoBlock{ID: id, URL: input.URL, DurationSec: input.DurationSec, Order: 0}
	return vb, nil
}

func (s *memoryStore) CreateQuiz(lessonID string, questions []*model.QuestionInput) (*model.Quiz, error) {
	id := s.next()
	q := &model.Quiz{ID: id, LessonID: lessonID}
	return q, nil
}

func (s *memoryStore) SubmitForReview(lessonID string) (*model.Submission, error) {
	id := s.next()
	sub := &model.Submission{ID: id, State: "SUBMITTED"}
	s.submissions[id] = sub
	return sub, nil
}

func (s *memoryStore) ApproveSubmission(lessonID string) (*model.Submission, error) {
	// simple behavior
	for _, sub := range s.submissions {
		sub.State = "APPROVED"
		return sub, nil
	}
	return nil, fmt.Errorf("no submission")
}

func (s *memoryStore) PublishLesson(lessonID string) (*model.Lesson, error) {
	if l, ok := s.lessons[lessonID]; ok {
		l.Status = ptrString("PUBLISHED")
		return l, nil
	}
	return nil, fmt.Errorf("lesson not found")
}

func (s *memoryStore) StartLesson(lessonID string) (*model.Progress, error) {
	id := s.next()
	p := &model.Progress{ID: id}
	return p, nil
}

func (s *memoryStore) CompleteLesson(lessonID string, score *int) (*model.Progress, error) {
	id := s.next()
	p := &model.Progress{ID: id}
	return p, nil
}

func (s *memoryStore) SubmitQuiz(input model.SubmitQuizInput) (*model.Attempt, error) {
	id := s.next()
	a := &model.Attempt{ID: id}
	return a, nil
}

func (s *memoryStore) RecommendedLessons(grade int) []*model.Lesson {
	return []*model.Lesson{}
}

func (s *memoryStore) SearchLessons(q string, grade *int) []*model.Lesson {
	return []*model.Lesson{}
}

func (s *memoryStore) GetProgressByLesson(lessonID string) (*model.Progress, error) {
	return &model.Progress{ID: "1"}, nil
}

func (s *memoryStore) ListSubmissions(state *string) []*model.Submission {
	out := []*model.Submission{}
	for _, sub := range s.submissions {
		out = append(out, sub)
	}
	return out
}

func ptrString(s string) *string { return &s }
