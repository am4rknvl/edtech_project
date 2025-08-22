package graph

// Hand-written resolver implementations. These are intentionally minimal and
// use an in-memory store defined in graph/store.go. Replace with a real DB
// and proper auth for production use.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/am4rknvl/edtech_project/graph/generated"
	"github.com/am4rknvl/edtech_project/graph/model"
)

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// --- Mutations (minimal, in-memory) ---
func (r *mutationResolver) SignUp(ctx context.Context, input model.SignUpInput) (*model.User, error) {
	u := store.CreateUser(input)
	return u, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, email *string, phone *string, otp *string, password *string) (*string, error) {
	// Return a fake token for now; integrate real auth later.
	token := "dev-token"
	return &token, nil
}

func (r *mutationResolver) CreateStudentProfile(ctx context.Context, input model.CreateStudentProfileInput) (*model.StudentProfile, error) {
	p := store.CreateStudentProfile(input)
	return p, nil
}

func (r *mutationResolver) CreateCourse(ctx context.Context, input model.CreateCourseInput) (*model.Course, error) {
	c, err := store.CreateCourse(input)
	return c, err
}

func (r *mutationResolver) CreateUnit(ctx context.Context, input model.CreateUnitInput) (*model.Unit, error) {
	u, err := store.CreateUnit(input)
	return u, err
}

func (r *mutationResolver) CreateLesson(ctx context.Context, input model.CreateLessonInput) (*model.Lesson, error) {
	l, err := store.CreateLesson(input)
	return l, err
}

func (r *mutationResolver) AddTextBlock(ctx context.Context, lessonID string, input model.TextBlockInput) (*model.TextBlock, error) {
	return store.AddTextBlock(lessonID, input)
}

func (r *mutationResolver) AddImageBlock(ctx context.Context, lessonID string, input model.ImageBlockInput) (*model.ImageBlock, error) {
	return store.AddImageBlock(lessonID, input)
}

func (r *mutationResolver) AddAudioBlock(ctx context.Context, lessonID string, input model.AudioBlockInput) (*model.AudioBlock, error) {
	return store.AddAudioBlock(lessonID, input)
}

func (r *mutationResolver) AddVideoBlock(ctx context.Context, lessonID string, input model.VideoBlockInput) (*model.VideoBlock, error) {
	return store.AddVideoBlock(lessonID, input)
}

func (r *mutationResolver) CreateQuiz(ctx context.Context, lessonID string, questions []*model.QuestionInput) (*model.Quiz, error) {
	return store.CreateQuiz(lessonID, questions)
}

func (r *mutationResolver) SubmitForReview(ctx context.Context, lessonID string) (*model.Submission, error) {
	return store.SubmitForReview(lessonID)
}

func (r *mutationResolver) ApproveSubmission(ctx context.Context, lessonID string) (*model.Submission, error) {
	return store.ApproveSubmission(lessonID)
}

func (r *mutationResolver) PublishLesson(ctx context.Context, lessonID string) (*model.Lesson, error) {
	return store.PublishLesson(lessonID)
}

func (r *mutationResolver) StartLesson(ctx context.Context, lessonID string) (*model.Progress, error) {
	return store.StartLesson(lessonID)
}

func (r *mutationResolver) CompleteLesson(ctx context.Context, lessonID string, score *int) (*model.Progress, error) {
	return store.CompleteLesson(lessonID, score)
}

func (r *mutationResolver) SubmitQuiz(ctx context.Context, input model.SubmitQuizInput) (*model.Attempt, error) {
	return store.SubmitQuiz(input)
}

func (r *mutationResolver) CreateUploadURL(ctx context.Context, input model.CreateUploadURLInput) (string, error) {
	// In prod, integrate with S3/GCS and return a signed URL. For now return a placeholder.
	return fmt.Sprintf("https://storage.example/%s/%d", input.Filename, time.Now().Unix()), nil
}

// --- Queries ---
func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	// In prod, resolve viewer from auth. Return nil in anonymous dev mode.
	return nil, nil
}

func (r *queryResolver) Subjects(ctx context.Context) ([]*model.Subject, error) {
	return store.ListSubjects(), nil
}

func (r *queryResolver) Courses(ctx context.Context, grade *int, subjectID *string) ([]*model.Course, error) {
	return store.ListCourses(grade, subjectID), nil
}

func (r *queryResolver) Course(ctx context.Context, id string) (*model.Course, error) {
	return store.GetCourse(id)
}

func (r *queryResolver) Units(ctx context.Context, courseID string) ([]*model.Unit, error) {
	return store.ListUnits(courseID), nil
}

func (r *queryResolver) Lessons(ctx context.Context, unitID string) ([]*model.Lesson, error) {
	return store.ListLessons(unitID), nil
}

func (r *queryResolver) Lesson(ctx context.Context, id string) (*model.Lesson, error) {
	return store.GetLesson(id)
}

func (r *queryResolver) RecommendedLessons(ctx context.Context, grade int) ([]*model.Lesson, error) {
	return store.RecommendedLessons(grade), nil
}

func (r *queryResolver) SearchLessons(ctx context.Context, q string, grade *int) ([]*model.Lesson, error) {
	return store.SearchLessons(q, grade), nil
}

func (r *queryResolver) MyProgress(ctx context.Context, courseID *string) ([]*model.Progress, error) {
	return []*model.Progress{}, nil
}

func (r *queryResolver) ProgressByLesson(ctx context.Context, lessonID string) (*model.Progress, error) {
	return store.GetProgressByLesson(lessonID)
}

func (r *queryResolver) MyDraftLessons(ctx context.Context) ([]*model.Lesson, error) {
	return []*model.Lesson{}, nil
}

func (r *queryResolver) Submissions(ctx context.Context, state *string) ([]*model.Submission, error) {
	return store.ListSubmissions(state), nil
}

// --- simple validation helpers ---
func ensureFound[T any](v *T, errMsg string) (*T, error) {
	if v == nil {
		return nil, errors.New(errMsg)
	}
	return v, nil
}
