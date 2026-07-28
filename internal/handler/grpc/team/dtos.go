package team

type CreateTeamReq struct {
	Name         string `json:"name"          validate:"required,team-name"`
	City         string `json:"city"          validate:"required"`
	Description  string `json:"description"   validate:"required"`
	ContactNum   string `json:"contact_num"   validate:"required"`
	Email        string `json:"contact_email" validate:"required,email"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"       validate:"required"`
	FullName     string `json:"full_name"     validate:"required"`
}

type UpdateTeamDetailsReq struct {
	TeamID       string `json:"team_id"        validate:"required"`
	TeamMemberID string `json:"team_member_id" validate:"required"`
	Name         string `json:"name"           validate:"omitempty,team-name"`
	City         string `json:"city"           validate:"omitempty,min=3,max=50"`
	Description  string `json:"description"    validate:"omitempty,max=255"`
	ShortName    string `json:"short_name"     validate:"omitempty,min=3,max=50"`
	Email        string `json:"email"          validate:"omitempty,email"`
	PhoneNum     string `json:"phone_num"      validate:"omitempty,phone"`
}



