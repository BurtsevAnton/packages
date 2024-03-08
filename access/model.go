package access

// Policies - структура для для хранения списка политик доступа.
type Policies struct {
	Data []Policy `json:"data"`
}

// Policy - структура для хранения политики доступа.
type Policy struct {
	RoleId              int    `json:"role_id"`
	RoleCode            string `json:"role_code"`
	ResourceId          int    `json:"resource_id"`
	UserAccessLevel     int    `json:"user_access_level"`
	ResourceCode        string `json:"resource_code"`
	ResourceTypeCode    string `json:"resource_type_code"`
	RequiredAccessLevel int    `json:"required_access_level"`
}

// LevelIndex - структура для хранения индекса доступных уровней.
type LevelIndex struct {
	Levels map[int]map[string]Policy
}
