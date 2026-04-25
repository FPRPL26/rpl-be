package entity

type Skill struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

func (s *Skill) TableName() string {
	return "skills"
}
