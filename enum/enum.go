package enum

type Role int

const (
	RoleAdmin    Role = 1
	RoleDoctor   Role = 2
	RoleResearch Role = 3
	RoleNurse    Role = 4
)

type QuestionType int

const (
	QuestionTypeFreeText      QuestionType = 1
	QuestionTypeFreeNumerical QuestionType = 2
	QuestionTypeRadioGroup    QuestionType = 3
	QuestionTypeCheckbox      QuestionType = 4
)
