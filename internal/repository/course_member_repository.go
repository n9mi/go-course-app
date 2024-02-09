package repository

import "github.com/n9mi/go-course-app/internal/entity"

type CourseMemberRepository struct {
	Repository[entity.CourseMember]
}

func NewCourseMemberRepository() *CourseMemberRepository {
	return new(CourseMemberRepository)
}
