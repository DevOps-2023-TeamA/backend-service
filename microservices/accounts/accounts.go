package main

type Accounts struct {
	ID       		int    	`json:"ID"`
	Name     		string 	`json:"Name"`
	Username     	string 	`json:"Username"`
	Password		string 	`json:"Password"`
	Role			string 	`json:"Role"`
	CreationDate	string	`json:"CreationDate"`
	IsApproved		bool 	`json:"IsApproved"`
	IsDeleted		bool 	`json:"IsDeleted"`
}

