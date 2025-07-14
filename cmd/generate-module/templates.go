package main

const docTemplate = `// Package {{.ModuleName}} - documentação do módulo
// Coloque aqui a descrição do que este módulo faz
package {{.ModuleName}}
`

const initTemplate = `package {{.ModuleName}}

import (
	"github.com/rafaelcoelhox/labbend/pkg/database"
)

func init() {
	// Registre aqui os modelos do GORM para migration automática
	database.RegisterModel(&{{.ModuleNameCap}}{})
}
`

const modelTemplate = `package {{.ModuleName}}

import (
	"time"
	"gorm.io/gorm"
)

// {{.ModuleNameCap}} - estrutura principal do módulo
// Adicione aqui os campos que sua entidade precisa
type {{.ModuleNameCap}} struct {
	ID        uint           ` + "`" + `gorm:"primaryKey" json:"id"` + "`" + `
	Nome      string         ` + "`" + `gorm:"not null" json:"nome"` + "`" + `
	Descricao string         ` + "`" + `json:"descricao"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"deleted_at"` + "`" + `
}

// Create{{.ModuleNameCap}}Input - input para criar nova entidade
// Adicione aqui os campos necessários para criação
type Create{{.ModuleNameCap}}Input struct {
	Nome      string ` + "`" + `json:"nome" validate:"required"` + "`" + `
	Descricao string ` + "`" + `json:"descricao"` + "`" + `
}

// Update{{.ModuleNameCap}}Input - input para atualizar entidade
// Adicione aqui os campos que podem ser atualizados
type Update{{.ModuleNameCap}}Input struct {
	Nome      *string ` + "`" + `json:"nome,omitempty"` + "`" + `
	Descricao *string ` + "`" + `json:"descricao,omitempty"` + "`" + `
}
`

const repositoryTemplate = `package {{.ModuleName}}

import (
	"context"
	"gorm.io/gorm"
)

// Repository - interface para operações de banco de dados
// Adicione aqui os métodos que precisa para acessar dados
type Repository interface {
	Create(ctx context.Context, {{.ModuleName}} *{{.ModuleNameCap}}) error
	GetByID(ctx context.Context, id uint) (*{{.ModuleNameCap}}, error)
	GetAll(ctx context.Context) ([]{{.ModuleNameCap}}, error)
	Update(ctx context.Context, id uint, input Update{{.ModuleNameCap}}Input) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Implementação básica dos métodos CRUD
// Customize conforme suas necessidades
func (r *repository) Create(ctx context.Context, {{.ModuleName}} *{{.ModuleNameCap}}) error {
	return r.db.WithContext(ctx).Create({{.ModuleName}}).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*{{.ModuleNameCap}}, error) {
	var {{.ModuleName}} {{.ModuleNameCap}}
	err := r.db.WithContext(ctx).First(&{{.ModuleName}}, id).Error
	return &{{.ModuleName}}, err
}

func (r *repository) GetAll(ctx context.Context) ([]{{.ModuleNameCap}}, error) {
	var {{.ModuleName}}s []{{.ModuleNameCap}}
	err := r.db.WithContext(ctx).Find(&{{.ModuleName}}s).Error
	return {{.ModuleName}}s, err
}

func (r *repository) Update(ctx context.Context, id uint, input Update{{.ModuleNameCap}}Input) error {
	return r.db.WithContext(ctx).Model(&{{.ModuleNameCap}}{}).Where("id = ?", id).Updates(input).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&{{.ModuleNameCap}}{}, id).Error
}
`

const serviceTemplate = `package {{.ModuleName}}

import (
	"context"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// Service - interface para lógica de negócio
// Adicione aqui os métodos de business logic
type Service interface {
	Create(ctx context.Context, input Create{{.ModuleNameCap}}Input) (*{{.ModuleNameCap}}, error)
	GetByID(ctx context.Context, id uint) (*{{.ModuleNameCap}}, error)
	GetAll(ctx context.Context) ([]{{.ModuleNameCap}}, error)
	Update(ctx context.Context, id uint, input Update{{.ModuleNameCap}}Input) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo   Repository
	logger logger.Logger
}

func NewService(repo Repository, logger logger.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Implementação dos métodos de negócio
// Adicione aqui validações, regras de negócio, etc.
func (s *service) Create(ctx context.Context, input Create{{.ModuleNameCap}}Input) (*{{.ModuleNameCap}}, error) {
	// Adicione validações aqui
	{{.ModuleName}} := &{{.ModuleNameCap}}{
		Nome:      input.Nome,
		Descricao: input.Descricao,
	}
	
	err := s.repo.Create(ctx, {{.ModuleName}})
	if err != nil {
		s.logger.Error("erro ao criar {{.ModuleName}}", "error", err)
		return nil, err
	}
	
	return {{.ModuleName}}, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*{{.ModuleNameCap}}, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]{{.ModuleNameCap}}, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, id uint, input Update{{.ModuleNameCap}}Input) error {
	// Adicione validações aqui
	return s.repo.Update(ctx, id, input)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
`

const graphqlTemplate = `package {{.ModuleName}}

import (
	"context"
	"strconv"
	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// GraphQL Types - configure os tipos GraphQL aqui
var {{.ModuleName}}Type = graphql.NewObject(graphql.ObjectConfig{
	Name: "{{.ModuleNameCap}}",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if {{.ModuleName}}, ok := p.Source.(*{{.ModuleNameCap}}); ok {
					return strconv.Itoa(int({{.ModuleName}}.ID)), nil
				}
				return nil, nil
			},
		},
		"nome": &graphql.Field{
			Type: graphql.String,
		},
		"descricao": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Queries - configure as consultas GraphQL aqui
func Queries(service Service, logger logger.Logger) graphql.Fields {
	return graphql.Fields{
		"{{.ModuleName}}": &graphql.Field{
			Type: {{.ModuleName}}Type,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := strconv.Atoi(p.Args["id"].(string))
				return service.GetByID(p.Context, uint(id))
			},
		},
		"{{.ModuleName}}s": &graphql.Field{
			Type: graphql.NewList({{.ModuleName}}Type),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return service.GetAll(p.Context)
			},
		},
	}
}

// Mutations - configure as mutações GraphQL aqui  
func Mutations(service Service, logger logger.Logger) graphql.Fields {
	return graphql.Fields{
		"create{{.ModuleNameCap}}": &graphql.Field{
			Type: {{.ModuleName}}Type,
			Args: graphql.FieldConfigArgument{
				"nome": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"descricao": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := Create{{.ModuleNameCap}}Input{
					Nome:      p.Args["nome"].(string),
					Descricao: p.Args["descricao"].(string),
				}
				return service.Create(p.Context, input)
			},
		},
		"delete{{.ModuleNameCap}}": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := strconv.Atoi(p.Args["id"].(string))
				err := service.Delete(p.Context, uint(id))
				return err == nil, err
			},
		},
	}
}
`
