package api

type HarryPotterCharacter struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	House         string `json:"house"`
	OriginStatus  string `json:"originStatus"`
}

var Students = []HarryPotterCharacter{
	{1, "Harry Potter", "Gryffindor", "Half-blood"},
	{2, "Hermione Granger","Gryffindor","Muggle-blood"},
	{3, "Ron Weasley","Gryffindor","Pure-blood"},
	{4, "Cedric Diggory","Hufflepuff","Half-blood"},
	{5, "Luna Lovegood", "Ravenclaw", "Pure-blood"},
	{6, "Draco Malfoy","Slytherin", "Pure-blood"},
}

func GetStudentById(id int) *HarryPotterCharacter {
	for i := range Students {
		if Students[i].ID == id {
			return &Students[i]
		}
	}
	return nil
}