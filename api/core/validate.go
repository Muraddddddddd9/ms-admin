package core

import (
	"fmt"
	"ms-admin/api/constants"
	"strings"

	"github.com/Muraddddddddd9/ms-database/models"
	"golang.org/x/crypto/bcrypt"
)

type ValidatableModel interface {
	Validate() error
}

type StudentsModel struct {
	models.StudentsModel `bson:",inline"`
}

type GroupsModel struct {
	models.GroupsModel `bson:",inline"`
}

type ObjectsGroupsModel struct {
	models.ObjectsGroupsModel `bson:",inline"`
}

type ObjectsModel struct {
	models.ObjectsModel `bson:",inline"`
}

type StatusesModel struct {
	models.StatusesModel `bson:",inline"`
}

type TeachersModel struct {
	models.TeachersModel `bson:",inline"`
}

func (s *StudentsModel) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Surname = strings.TrimSpace(s.Surname)
	s.Patronymic = strings.TrimSpace(s.Patronymic)
	s.Email = strings.TrimSpace(strings.ToLower(s.Email))
	s.Password = strings.TrimSpace(s.Password)
	s.Diplomas = []string{}
	s.IPs = []string{}

	fields := map[string]string{
		"name":       s.Name,
		"surname":    s.Surname,
		"patronymic": s.Patronymic,
		"email":      s.Email,
		"password":   s.Password,
	}

	for name, value := range fields {
		if value == "" {
			return fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	s.Password = string(bcryptPassword)

	return nil
}

func (g *GroupsModel) Validate() error {
	g.Group = strings.ToLower(strings.TrimSpace(g.Group))

	fileds := map[string]string{
		"group": g.Group,
	}

	for name, value := range fileds {
		if value == "" {
			return fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	return nil
}

func (og *ObjectsGroupsModel) Validate() error {
	return nil
}

func (o *ObjectsModel) Validate() error {
	o.Object = strings.TrimSpace(strings.ToLower(o.Object))

	fileds := map[string]string{
		"object": o.Object,
	}

	for name, value := range fileds {
		if value == "" {
			return fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	return nil
}

func (s *StatusesModel) Validate() error {
	s.Status = strings.TrimSpace(strings.ToLower(s.Status))

	fields := map[string]string{
		"status": s.Status,
	}

	for name, value := range fields {
		if value == "" {
			return fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	return nil
}

func (t *TeachersModel) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Surname = strings.TrimSpace(t.Surname)
	t.Patronymic = strings.TrimSpace(t.Patronymic)
	t.Email = strings.TrimSpace(strings.ToLower(t.Email))
	t.Password = strings.TrimSpace(t.Password)
	t.IPs = []string{}

	fields := map[string]string{
		"name":       t.Name,
		"surname":    t.Surname,
		"patronymic": t.Patronymic,
		"email":      t.Email,
		"password":   t.Password,
	}

	for name, value := range fields {
		if value == "" {
			return fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
	t.Password = string(bcryptPassword)

	return nil
}
