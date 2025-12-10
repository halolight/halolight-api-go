package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type TeamService interface {
	List(userID string, page, limit int, search string) ([]models.Team, int64, error)
	Get(id string) (*models.Team, error)
	Create(name, description, ownerID string) (*models.Team, error)
	Update(id, name, description string) (*models.Team, error)
	Delete(id string) error
	AddMember(teamID, userID string, roleID *string) (*models.TeamMember, error)
	RemoveMember(teamID, userID string) error
	IsOwner(teamID, userID string) bool
}

type teamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) TeamService {
	return &teamService{db: db}
}

func (s *teamService) List(userID string, page, limit int, search string) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	query := s.db.Model(&models.Team{}).
		Where("owner_id = ? OR id IN (SELECT team_id FROM team_members WHERE user_id = ?)", userID, userID)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Owner").Offset(offset).Limit(limit).Order("created_at DESC").Find(&teams).Error

	return teams, total, err
}

func (s *teamService) Get(id string) (*models.Team, error) {
	var team models.Team
	err := s.db.Preload("Owner").Preload("Members.User").First(&team, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *teamService) Create(name, description, ownerID string) (*models.Team, error) {
	team := &models.Team{
		Name:        name,
		Description: &description,
		OwnerID:     ownerID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		// Add owner as member
		member := &models.TeamMember{
			TeamID: team.ID,
			UserID: ownerID,
		}
		return tx.Create(member).Error
	})

	return team, err
}

func (s *teamService) Update(id, name, description string) (*models.Team, error) {
	team, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		team.Name = name
	}
	if description != "" {
		team.Description = &description
	}

	err = s.db.Save(team).Error
	return team, err
}

func (s *teamService) Delete(id string) error {
	return s.db.Delete(&models.Team{}, "id = ?", id).Error
}

func (s *teamService) AddMember(teamID, userID string, roleID *string) (*models.TeamMember, error) {
	member := &models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		RoleID: roleID,
	}
	err := s.db.Create(member).Error
	if err != nil {
		return nil, err
	}
	// Reload with user
	s.db.Preload("User").First(member, "team_id = ? AND user_id = ?", teamID, userID)
	return member, nil
}

func (s *teamService) RemoveMember(teamID, userID string) error {
	return s.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{}).Error
}

func (s *teamService) IsOwner(teamID, userID string) bool {
	var team models.Team
	err := s.db.Select("owner_id").First(&team, "id = ?", teamID).Error
	if err != nil {
		return false
	}
	return team.OwnerID == userID
}
