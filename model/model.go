package model

type Choice struct {
	ID   string `json:"id"`
	Gone bool   `json:"gone"`
	Come bool   `json:"come"`
}

//define a Choice CRUD struct(that cannot be exportd) that can access the Choice members and define sql queries for CRUD operations within it as its own members
//it will use the members of the Choice Struct to as arguments of the SQL code
/*
	type choiceCRUD struct{ 	//or type [c Choice] choiceCRUD struct{}
	Choice
	GET func()string{
		return fmt.Sprintf("", Choice.ID)
		}
	}
*/

//or I can simply write GET as a method of the choice struct
