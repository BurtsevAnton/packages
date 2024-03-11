package access

import (
	"encoding/json"
	"regexp"
	"strings"
	"sync"
)

// Service - сервис политики доступа. Содержит основное хранилище политик доступа.
type Service struct {
	access *LevelIndex
	sync.RWMutex
}

// NewService - возвращает новый экземпляр сервиса политики доступа.
func NewService() *Service {
	return &Service{
		access: nil,
	}
}

// IsAccessAllowed - проверяет доступ пользователя к ресурсу.
func (s *Service) IsAccessAllowed(roleId int, resourceCode string) bool {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.access.Levels[roleId]["IDDQD"]; ok {
		return true
	}

	if policy, ok := s.access.Levels[roleId][resourceCode]; ok {

		return policy.UserAccessLevel&policy.RequiredAccessLevel == policy.RequiredAccessLevel
	}

	// Если статичный путь не найден, то проверяем динамичные (шаблонные) пути.
	for key := range s.access.Levels[roleId] {
		keyMod := strings.ReplaceAll(key, `/`, `\/`)
		keyMod = strings.ReplaceAll(keyMod, `{%s}`, `[\w-]+`)
		r := regexp.MustCompile("^" + keyMod + "$")

		if r.MatchString(resourceCode) {
			if policy, ok := s.access.Levels[roleId][key]; ok {

				return policy.UserAccessLevel&policy.RequiredAccessLevel == policy.RequiredAccessLevel
			}
		}
	}

	return false
}

// SetAccessPolicies - записывает в структуру политику доступов на основе переданных данных.
func (s *Service) SetAccessPolicies(policiesRaw string) error {
	var policies Policies
	if err := json.Unmarshal([]byte(policiesRaw), &policies); err != nil {
		return err
	}

	var access LevelIndex
	access.Levels = make(map[int]map[string]Policy)

	for _, policy := range policies.Data {
		if _, ok := access.Levels[policy.RoleId]; !ok {
			access.Levels[policy.RoleId] = make(map[string]Policy)
		}
		access.Levels[policy.RoleId][policy.ResourceCode] = policy
	}

	s.Lock()
	s.access = &access
	s.Unlock()

	return nil
}
