package main

type Records struct {
	ID					int 		`json:"ID"`
	AccountID			int 		`json:"AccountID"`
	ContactRole			string 		`json:"ContactRole"`
	StudentCount		int 		`json:"StudentCount"`	
	AcadYear			string 		`json:"AcadYear"`
	Title				string 		`json:"Title"`
	CompanyName			string 		`json:"CompanyName"`
	CompanyPOC			string 		`json:"CompanyPOC"`
	Description			string 		`json:"Description"`
	CreationDate        string 		`json:"CreationDate"`
	IsDeleted			bool 		`json:"IsDeleted"`
}