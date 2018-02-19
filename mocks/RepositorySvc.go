// Code generated by mockery v1.0.0
package mocks

import context "context"
import github "github.com/google/go-github/github"

import mock "github.com/stretchr/testify/mock"

// RepositorySvc is an autogenerated mock type for the RepositorySvc type
type RepositorySvc struct {
	mock.Mock
}

// GetLatestRelease provides a mock function with given fields: ctx, owner, repo
func (_m *RepositorySvc) GetLatestRelease(ctx context.Context, owner string, repo string) (*github.RepositoryRelease, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo)

	var r0 *github.RepositoryRelease
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *github.RepositoryRelease); ok {
		r0 = rf(ctx, owner, repo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryRelease)
		}
	}

	var r1 *github.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *github.Response); ok {
		r1 = rf(ctx, owner, repo)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, owner, repo)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetReleaseByTag provides a mock function with given fields: ctx, owner, repo, tag
func (_m *RepositorySvc) GetReleaseByTag(ctx context.Context, owner string, repo string, tag string) (*github.RepositoryRelease, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, tag)

	var r0 *github.RepositoryRelease
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *github.RepositoryRelease); ok {
		r0 = rf(ctx, owner, repo, tag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryRelease)
		}
	}

	var r1 *github.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) *github.Response); ok {
		r1 = rf(ctx, owner, repo, tag)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, string) error); ok {
		r2 = rf(ctx, owner, repo, tag)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
