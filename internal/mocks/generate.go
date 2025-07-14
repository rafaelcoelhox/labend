// Package mocks provides generated mocks for testing
package mocks

//go:generate mockgen -destination=users_repository_mock.go -package=mocks -mock_names=Repository=MockUsersRepository github.com/rafaelcoelhox/labbend/internal/users Repository
//go:generate mockgen -destination=users_service_mock.go -package=mocks -mock_names=Service=MockUsersService github.com/rafaelcoelhox/labbend/internal/users Service
//go:generate mockgen -destination=challenges_repository_mock.go -package=mocks -mock_names=Repository=MockChallengesRepository github.com/rafaelcoelhox/labbend/internal/challenges Repository
//go:generate mockgen -destination=challenges_service_mock.go -package=mocks -mock_names=Service=MockChallengesService github.com/rafaelcoelhox/labbend/internal/challenges Service
//go:generate mockgen -destination=challenges_userservice_mock.go -package=mocks -mock_names=UserService=MockChallengesUserService github.com/rafaelcoelhox/labbend/internal/challenges UserService
//go:generate mockgen -destination=users_eventbus_mock.go -package=mocks -mock_names=EventBus=MockUsersEventBus github.com/rafaelcoelhox/labbend/internal/users EventBus
//go:generate mockgen -destination=challenges_eventbus_mock.go -package=mocks -mock_names=EventBus=MockChallengesEventBus github.com/rafaelcoelhox/labbend/internal/challenges EventBus
//go:generate mockgen -destination=eventbus_handler_mock.go -package=mocks -mock_names=EventHandler=MockEventHandler github.com/rafaelcoelhox/labbend/pkg/eventbus EventHandler
//go:generate mockgen -destination=logger_mock.go -package=mocks -mock_names=Logger=MockLogger github.com/rafaelcoelhox/labbend/pkg/logger Logger
//go:generate mockgen -destination=notifications_repository_mock.go -package=mocks -mock_names=Repository=MockNotificationsRepository github.com/rafaelcoelhox/labbend/internal/notifications Repository
//go:generate mockgen -destination=notifications_service_mock.go -package=mocks -mock_names=Service=MockNotificationsService github.com/rafaelcoelhox/labbend/internal/notifications Service
//go:generate mockgen -destination=notifications_eventbus_mock.go -package=mocks -mock_names=EventBus=MockNotificationsEventBus github.com/rafaelcoelhox/labbend/internal/notifications EventBus
//go:generate mockgen -destination=notifications_userservice_mock.go -package=mocks -mock_names=UserService=MockNotificationsUserService github.com/rafaelcoelhox/labbend/internal/notifications UserService
