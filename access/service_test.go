package access

import (
	"testing"
)

const policyFromDb = `{"data": [
{"role_id": 1, "role_code": "DEVELOPER", "resource_id": 1, "resource_code": "IDDQD", "user_access_level": 15, "resource_type_code": "ADMIN", "required_access_level": 1},
{"role_id": 2, "role_code": "ADMIN", "resource_id": 5, "resource_code": "POST_/api/test_create", "user_access_level": 7, "resource_type_code": "API", "required_access_level": 5},
{"role_id": 2, "role_code": "ADMIN", "resource_id": 2, "resource_code": "BOARD_ADMIN", "user_access_level": 1, "resource_type_code": "APP_MODULE", "required_access_level": 1},
{"role_id": 3, "role_code": "LEAD", "resource_id": 5, "resource_code": "POST_/api/test_create", "user_access_level": 7, "resource_type_code": "API", "required_access_level": 5},
{"role_id": 3, "role_code": "LEAD", "resource_id": 6, "resource_code": "GET_/api/test_get", "user_access_level": 7, "resource_type_code": "API", "required_access_level": 3},
{"role_id": 3, "role_code": "LEAD", "resource_id": 7, "resource_code": "DELETE_/api/test_delete", "user_access_level": 1, "resource_type_code": "API", "required_access_level": 9},
{"role_id": 3, "role_code": "LEAD", "resource_id": 3, "resource_code": "BOARD_REPORT", "user_access_level": 1, "resource_type_code": "APP_MODULE", "required_access_level": 1},
{"role_id": 3, "role_code": "LEAD", "resource_id": 4, "resource_code": "BOARD_WORKER", "user_access_level": 1, "resource_type_code": "APP_MODULE", "required_access_level": 1},
{"role_id": 4, "role_code": "WORKER", "resource_id": 5, "resource_code": "POST_/api/test_create", "user_access_level": 3, "resource_type_code": "API", "required_access_level": 5},
{"role_id": 4, "role_code": "WORKER", "resource_id": 6, "resource_code": "GET_/api/test_get", "user_access_level": 3, "resource_type_code": "API", "required_access_level": 3},
{"role_id": 4, "role_code": "WORKER", "resource_id": 7, "resource_code": "DELETE_/api/test_delete", "user_access_level": 1, "resource_type_code": "API", "required_access_level": 9},
{"role_id": 4, "role_code": "WORKER", "resource_id": 4, "resource_code": "BOARD_WORKER", "user_access_level": 1, "resource_type_code": "APP_MODULE", "required_access_level": 1}
]}`

func Test_NewService(t *testing.T) {
	service := NewService()

	if service.access != nil {
		t.Errorf("Ожидается, что поле access будет nil, получено %v", service.access)
	}
}

func TestService_IsAccessAllowed(t *testing.T) {
	service := NewService()
	if err := service.SetAccessPolicies(policyFromDb); err != nil {
		t.Error("SetAccessPolicies() fail")
	}

	type args struct {
		roleId       int
		resourceCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_accessParser_1_POST",
			args: args{
				roleId:       1, // Developer
				resourceCode: "POST_/api/test_create",
			},
			want: true,
		},
		{
			name: "Test_accessParser_1_GET",
			args: args{
				roleId:       1,
				resourceCode: "GET_/api/test_get",
			},
			want: true,
		},
		{
			name: "Test_accessParser_1_DELETE",
			args: args{
				roleId:       1,
				resourceCode: "DELETE_/api/test_delete",
			},
			want: true,
		},
		{
			name: "Test_accessParser_1_RUN",
			args: args{
				roleId:       1,
				resourceCode: "SomeResource",
			},
			want: true,
		},

		{
			name: "Test_accessParser_2_POST",
			args: args{
				roleId:       2, // Admin
				resourceCode: "POST_/api/test_create",
			},
			want: true, // Разрешено
		},
		{
			name: "Test_accessParser_2_GET",
			args: args{
				roleId:       2,
				resourceCode: "GET_/api/test_get",
			},
			want: false, //Не найден уровень доступа для роли 2 и ресурса GET_/api/test_get
		},
		{
			name: "Test_accessParser_2_DELETE",
			args: args{
				roleId:       2,
				resourceCode: "DELETE_/api/test_delete",
			},
			want: false, // Не найден уровень доступа для роли 2 и ресурса DELETE_/api/test_delete
		},
		{
			name: "Test_accessParser_2_RUN",
			args: args{
				roleId:       2,
				resourceCode: "SomeResource",
			},
			want: false, // Запрещено
		},

		{
			name: "Test_accessParser_3_POST",
			args: args{
				roleId:       3, // Lead
				resourceCode: "POST_/api/test_create",
			},
			want: true, // Разрешено
		},
		{
			name: "Test_accessParser_3_GET",
			args: args{
				roleId:       3,
				resourceCode: "GET_/api/test_get",
			},
			want: true, // Разрешено
		},
		{
			name: "Test_accessParser_3_DELETE",
			args: args{
				roleId:       3,
				resourceCode: "DELETE_/api/test_delete",
			},
			want: false, // Запрещено
		},
		{
			name: "Test_accessParser_3_RUN",
			args: args{
				roleId:       3,
				resourceCode: "SomeResource",
			},
			want: false, // Запрещено // Не найден уровень доступа для роли 3 и ресурса SomeResource
		},

		{
			name: "Test_accessParser_4_POST",
			args: args{
				roleId:       4, // Worker
				resourceCode: "POST_/api/test_create",
			},
			want: false, // Запрещено
		},
		{
			name: "Test_accessParser_4_GET",
			args: args{
				roleId:       4,
				resourceCode: "GET_/api/test_get",
			},
			want: true, // Разрешено
		},
		{
			name: "Test_accessParser_4_DELETE",
			args: args{
				roleId:       4,
				resourceCode: "DELETE_/api/test_delete",
			},
			want: false, // Запрещено
		},
		{
			name: "Test_accessParser_4_RUN",
			args: args{
				roleId:       4,
				resourceCode: "SomeResource",
			},
			want: false, // Запрещено // Не найден уровень доступа для роли 4 и ресурса SomeResource
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.IsAccessAllowed(tt.args.roleId, tt.args.resourceCode); got != tt.want {
				t.Errorf("IsAccessAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SetAccessPolicies(t *testing.T) {
	service := NewService()

	err := service.SetAccessPolicies(`invalid`)
	if err == nil {
		t.Error("Ожидается, что JSON будет неправильный")
	}

	err = service.SetAccessPolicies(policyFromDb)
	if err != nil {
		t.Errorf("Неожиданная ошибка: %v", err)
	}

	if len(service.access.Levels) != 4 {
		t.Errorf("Ожидается размерность 4, получено %d", len(service.access.Levels))
	}

	if service.access.Levels[1]["IDDQD"].UserAccessLevel != 15 {
		t.Errorf("Ожидается уровень 15 для роли 1, получено %d", service.access.Levels[1]["IDDQD"].UserAccessLevel)
	}

	if service.access.Levels[2]["POST_/api/test_create"].RequiredAccessLevel != 5 {
		t.Errorf("Ожидается уровень 5 для ресурса 'POST_/api/test_create', получено %d", service.access.Levels[2]["POST_/api/test_create"].RequiredAccessLevel)
	}

	if service.access.Levels[2]["GET_/api/test_get"].UserAccessLevel != 0 {
		t.Errorf("Ожидается отсутствие уровняя доступа (0) у роли 2 для ресурсу 'GET_/api/test_get', получено %d", service.access.Levels[2]["GET_/api/test_get"].UserAccessLevel)
	}

	if service.access.Levels[4]["POST_/api/test_create"].UserAccessLevel != 3 {
		t.Errorf("Ожидается уровень 3 для роли 4, получено %d", service.access.Levels[4]["POST_/api/test_create"].UserAccessLevel)
	}

	if service.access.Levels[4]["POST_/api/test_create"].RequiredAccessLevel != 5 {
		t.Errorf("Ожидается уровень 5 для ресурса 'POST_/api/test_create', получено %d", service.access.Levels[4]["POST_/api/test_create"].RequiredAccessLevel)
	}
}
